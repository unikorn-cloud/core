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
	"errors"
	"fmt"

	"github.com/unikorn-cloud/core/pkg/authorization/roles"
)

var (
	ErrPermissionDenied = errors.New("access denied")
)

// SuperAdminAuthorizer allows access to everything.
type SuperAdminAuthorizer struct{}

func (a *SuperAdminAuthorizer) Allow(_ string, _ roles.Permission) error {
	return nil
}

// OrganizationAuthorizer is scoped to a specific organization.
type OrganizationAuthorizer struct {
	permissions *OrganizationPermissions
}

type PermissionMap map[roles.Permission]interface{}

type ScopedPermissionMap map[string]PermissionMap

func (a *OrganizationAuthorizer) Allow(scope string, permission roles.Permission) error {
	roleManager := roles.New()

	scopedPermissions := ScopedPermissionMap{}

	// Build up a set of scopes and permissions based on group membership and the roles
	// associated with that.
	for _, group := range a.permissions.Groups {
		for _, r := range group.Roles {
			rolePermissions, err := roleManager.GetRole(r)
			if err != nil {
				return err
			}

			for roleScope, perms := range rolePermissions.Permissions {
				if _, ok := scopedPermissions[roleScope]; !ok {
					scopedPermissions[roleScope] = PermissionMap{}
				}

				for _, perm := range perms {
					scopedPermissions[roleScope][perm] = nil
				}
			}
		}
	}

	s, ok := scopedPermissions[scope]
	if !ok {
		return fmt.Errorf("%w: not permitted to access the %v scope", ErrPermissionDenied, scope)
	}

	if _, ok := s[permission]; !ok {
		return fmt.Errorf("%w: not permitted to %v within the %v scope", ErrPermissionDenied, permission, scope)
	}

	return nil
}

func New(permissions *Permissions, organizationName string) (Authorizer, error) {
	if permissions == nil {
		return nil, fmt.Errorf("%w: user has no RBAC information", ErrPermissionDenied)
	}

	if permissions.IsSuperAdmin {
		return &SuperAdminAuthorizer{}, nil
	}

	organization, err := permissions.LookupOrganization(organizationName)
	if err != nil {
		return nil, err
	}

	authorizer := &OrganizationAuthorizer{
		permissions: organization,
	}

	return authorizer, nil
}
