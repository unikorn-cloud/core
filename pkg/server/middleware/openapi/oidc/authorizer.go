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
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/spf13/pflag"
	"golang.org/x/oauth2"

	"github.com/unikorn-cloud/core/pkg/server/errors"
)

type Options struct {
	// issuer is used to perform OIDC discovery and verify access tokens
	// using the JWKS endpoint.
	issuer string

	// issuerCA is the root CA of the identity endpoint.
	issuerCA []byte
}

func (o *Options) AddFlags(f *pflag.FlagSet) {
	f.StringVar(&o.issuer, "oidc-issuer", "", "OIDC issuer URL to use for token validation.")
	f.BytesBase64Var(&o.issuerCA, "oidc-issuer-ca", nil, "base64 OIDC endpoint CA certificate.")
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

// authorizeOAuth2 checks APIs that require and oauth2 bearer token.
func (a *Authorizer) authorizeOAuth2(r *http.Request) (string, *oidc.UserInfo, error) {
	authorizationScheme, rawToken, err := getHTTPAuthenticationScheme(r)
	if err != nil {
		return "", nil, err
	}

	if !strings.EqualFold(authorizationScheme, "bearer") {
		return "", nil, errors.OAuth2InvalidRequest("authorization scheme not allowed").WithValues("scheme", authorizationScheme)
	}

	// Handle non-public CA certiifcates used in development.
	ctx := r.Context()

	if a.options.issuerCA != nil {
		certPool := x509.NewCertPool()

		if ok := certPool.AppendCertsFromPEM(a.options.issuerCA); !ok {
			return "", nil, errors.OAuth2InvalidRequest("failed to parse oidc issuer CA cert")
		}

		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:    certPool,
					MinVersion: tls.VersionTLS13,
				},
			},
		}

		ctx = oidc.ClientContext(ctx, client)
	}

	// Perform userinfo call against the identity service that will validate the token
	// and also return some information about the user.
	provider, err := oidc.NewProvider(ctx, a.options.issuer)
	if err != nil {
		return "", nil, errors.OAuth2ServerError("oidc service discovery failed").WithError(err)
	}

	token := &oauth2.Token{
		AccessToken: rawToken,
		TokenType:   authorizationScheme,
	}

	userinfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return "", nil, err
	}

	return rawToken, userinfo, nil
}

// Authorize checks the request against the OpenAPI security scheme.
func (a *Authorizer) Authorize(authentication *openapi3filter.AuthenticationInput) (string, *oidc.UserInfo, error) {
	if authentication.SecurityScheme.Type == "oauth2" {
		return a.authorizeOAuth2(authentication.RequestValidationInput.Request)
	}

	return "", nil, errors.OAuth2InvalidRequest("authorization scheme unsupported").WithValues("scheme", authentication.SecurityScheme.Type)
}
