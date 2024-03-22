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
	"slices"

	"github.com/unikorn-cloud/core/pkg/authorization/roles"
)

var (
	ErrPermissionDenied = errors.New("access denied")
)

type AdminAuthorizer struct{}

func (a *AdminAuthorizer) AllowedByRole(_ roles.Role) error {
	return nil
}

func (a *AdminAuthorizer) AllowedByGroup(_ []string) error {
	return nil
}

func (a *AdminAuthorizer) AllowedByGroupRole(_ []string, _ roles.Role) error {
	return nil
}

type OrganizationAuthorizer struct {
	permissions *OrganizationPermissions
}

func (a *OrganizationAuthorizer) AllowedByRole(role roles.Role) error {
	return ErrPermissionDenied
}

func (a *OrganizationAuthorizer) AllowedByGroup(groups []string) error {
	for _, group := range a.permissions.Groups {
		if slices.Contains(groups, group.ID) {
			return nil
		}
	}

	return ErrPermissionDenied
}

func (a *OrganizationAuthorizer) AllowedByGroupRole(groups []string, role roles.Role) error {
	for _, group := range a.permissions.Groups {
		if slices.Contains(groups, group.ID) && slices.Contains(group.Roles, role) {
			return nil
		}
	}

	return ErrPermissionDenied
}

func NewAuthorizer(permissions *Permissions, organizationName string) (Authorizer, error) {
	if permissions == nil {
		return nil, fmt.Errorf("%w: user has no RBAC information", ErrPermissionDenied)
	}

	if permissions.IsSuperAdmin {
		return &AdminAuthorizer{}, nil
	}

	organization, err := permissions.LookupOrganization(organizationName)
	if err != nil {
		return nil, err
	}

	if organization.IsAdmin {
		return &AdminAuthorizer{}, nil
	}

	return &OrganizationAuthorizer{
		permissions: organization,
	}, nil
}
