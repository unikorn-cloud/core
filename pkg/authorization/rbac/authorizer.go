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
	"errors"
	"fmt"
	"slices"

	"github.com/unikorn-cloud/core/pkg/authorization/constants"
)

var (
	ErrPermissionDenied = errors.New("access denied")
)

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
