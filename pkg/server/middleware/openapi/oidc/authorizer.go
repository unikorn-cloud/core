/*
Copyright 2022-2024 EscherCloud.
Copyright 2024 the Unikorn Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oidc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strconv"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/oauth2"

	"github.com/unikorn-cloud/core/pkg/authorization/userinfo"
	"github.com/unikorn-cloud/core/pkg/server/errors"
)

type Options struct {
	// Issuer is used to perform OIDC discovery and verify access tokens
	// using the JWKS endpoint.
	Issuer string

	// IssuerCA is the root CA of the identity endpoint.
	IssuerCA []byte
}

func (o *Options) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&o.Issuer, "oidc-issuer", "", "OIDC issuer URL to use for token validation.")
	f.BytesBase64Var(&o.IssuerCA, "oidc-issuer-ca", nil, "base64 OIDC endpoint CA certificate.")
}

// Authorizer provides OpenAPI based authorization middleware.
type Authorizer struct {
	options *Options
}

// NewAuthorizer returns a new authorizer with required parameters.
func NewAuthorizer(options *Options) *Authorizer {
	return &Authorizer{
		options: options,
	}
}

// getHTTPAuthenticationScheme grabs the scheme and token from the HTTP
// Authorization header.
func getHTTPAuthenticationScheme(r *http.Request) (string, string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", "", errors.OAuth2InvalidRequest("authorization header missing")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", "", errors.OAuth2InvalidRequest("authorization header malformed")
	}

	return parts[0], parts[1], nil
}

type propagationFunc func(r *http.Request)

type propagatingTransport struct {
	base http.Transport
	f    propagationFunc
}

func newPropagatingTransport(ctx context.Context) *propagatingTransport {
	return &propagatingTransport{
		f: func(r *http.Request) {
			otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
		},
	}
}

func (t *propagatingTransport) Clone() *propagatingTransport {
	return &propagatingTransport{
		f: t.f,
	}
}

func (t *propagatingTransport) CloseIdleConnections() {
	t.base.CloseIdleConnections()
}

func (t *propagatingTransport) RegisterProtocol(scheme string, rt http.RoundTripper) {
	t.base.RegisterProtocol(scheme, rt)
}

func (t *propagatingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.f(req)

	return t.base.RoundTrip(req)
}

// oidcErrorIsUnauthorized tries to convert the error returned by the OIDC library
// into a proper status code, as it doesn't wrap anything useful.
// The error looks like "{code} {text code}: {body}".
func oidcErrorIsUnauthorized(err error) bool {
	// Does it look like it contains the colon?
	fields := strings.Split(err.Error(), ":")
	if len(fields) < 2 {
		return false
	}

	// What about a number followed by a string?
	fields = strings.Split(fields[0], " ")
	if len(fields) < 2 {
		return false
	}

	code, err := strconv.Atoi(fields[0])
	if err != nil {
		return false
	}

	// Is the number a 403?
	return code == http.StatusUnauthorized
}

// authorizeOAuth2 checks APIs that require and oauth2 bearer token.
func (a *Authorizer) authorizeOAuth2(r *http.Request) (string, *userinfo.UserInfo, error) {
	authorizationScheme, rawToken, err := getHTTPAuthenticationScheme(r)
	if err != nil {
		return "", nil, err
	}

	if !strings.EqualFold(authorizationScheme, "bearer") {
		return "", nil, errors.OAuth2InvalidRequest("authorization scheme not allowed").WithValues("scheme", authorizationScheme)
	}

	// Handle non-public CA certiifcates used in development.
	ctx := r.Context()

	transport := newPropagatingTransport(ctx)

	if a.options.IssuerCA != nil {
		certPool := x509.NewCertPool()

		if ok := certPool.AppendCertsFromPEM(a.options.IssuerCA); !ok {
			return "", nil, errors.OAuth2InvalidRequest("failed to parse oidc issuer CA cert")
		}

		transport.base.TLSClientConfig = &tls.Config{
			RootCAs:    certPool,
			MinVersion: tls.VersionTLS13,
		}
	}

	client := &http.Client{
		Transport: transport,
	}

	ctx = oidc.ClientContext(ctx, client)

	// Perform userinfo call against the identity service that will validate the token
	// and also return some information about the user.
	provider, err := oidc.NewProvider(ctx, a.options.Issuer)
	if err != nil {
		return "", nil, errors.OAuth2ServerError("oidc service discovery failed").WithError(err)
	}

	token := &oauth2.Token{
		AccessToken: rawToken,
		TokenType:   authorizationScheme,
	}

	ui, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		if oidcErrorIsUnauthorized(err) {
			return "", nil, errors.OAuth2AccessDenied("token validation failed").WithError(err)
		}

		return "", nil, err
	}

	uiInternal := &userinfo.UserInfo{}

	if err := ui.Claims(uiInternal); err != nil {
		return "", nil, errors.OAuth2ServerError("failed to extrac user information").WithError(err)
	}

	return rawToken, uiInternal, nil
}

// Authorize checks the request against the OpenAPI security scheme.
func (a *Authorizer) Authorize(authentication *openapi3filter.AuthenticationInput) (string, *userinfo.UserInfo, error) {
	if authentication.SecurityScheme.Type == "oauth2" {
		return a.authorizeOAuth2(authentication.RequestValidationInput.Request)
	}

	return "", nil, errors.OAuth2InvalidRequest("authorization scheme unsupported").WithValues("scheme", authentication.SecurityScheme.Type)
}
