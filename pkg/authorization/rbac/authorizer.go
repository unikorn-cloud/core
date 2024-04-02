/*
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

package rbac

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/unikorn-cloud/core/pkg/authorization/accesstoken"
	"github.com/unikorn-cloud/core/pkg/authorization/constants"
)

var (
	ErrPermissionDenied = errors.New("access denied")

	ErrRequestError = errors.New("request error")

	ErrCertError = errors.New("certificate error")
)

// IdentityACLGetter grabs an ACL for the user from the identity API.
// Used for any non-identity API.
type IdentityACLGetter struct {
	host         string
	organization string
	ca           []byte
}

func NewIdentityACLGetter(host, organization string) *IdentityACLGetter {
	return &IdentityACLGetter{
		host:         host,
		organization: organization,
	}
}

func (a *IdentityACLGetter) WithCA(ca []byte) *IdentityACLGetter {
	a.ca = ca

	return a
}

func (a *IdentityACLGetter) Get(ctx context.Context) (*ACL, error) {
	client := &http.Client{}

	// Handle things like let's encrypt staging.
	if a.ca != nil {
		certPool := x509.NewCertPool()

		if ok := certPool.AppendCertsFromPEM(a.ca); !ok {
			return nil, fmt.Errorf("%w: unable to add CA certificate", ErrCertError)
		}

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:    certPool,
					MinVersion: tls.VersionTLS13,
				},
			},
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/organizations/%s/acl", a.host, a.organization), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+accesstoken.FromContext(ctx))
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code not as expected", ErrRequestError)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	acl := &ACL{}

	if err := json.Unmarshal(body, &acl); err != nil {
		return nil, err
	}

	return acl, nil
}

// SuperAdminAuthorizer allows access to everything.
type SuperAdminAuthorizer struct{}

func (a *SuperAdminAuthorizer) Allow(_ context.Context, _ string, _ constants.Permission) error {
	return nil
}

// BaseAuthorizer is scoped to a specific organization.
type BaseAuthorizer struct {
	acl *ACL
}

func (a *BaseAuthorizer) Allow(ctx context.Context, scope string, permission constants.Permission) error {
	aclScope := a.acl.GetScope(scope)
	if aclScope == nil {
		return fmt.Errorf("%w: not permitted access to the %v scope", ErrPermissionDenied, scope)
	}

	if !slices.Contains(aclScope.Permissions, permission) {
		return fmt.Errorf("%w: not permitted %v access within the %v scope", ErrPermissionDenied, permission, scope)
	}

	return nil
}

func New(ctx context.Context, getter ACLGetter) (Authorizer, error) {
	acl, err := getter.Get(ctx)
	if err != nil {
		return nil, err
	}

	if acl.IsSuperAdmin {
		return &SuperAdminAuthorizer{}, nil
	}

	authorizer := &BaseAuthorizer{
		acl: acl,
	}

	return authorizer, nil
}
