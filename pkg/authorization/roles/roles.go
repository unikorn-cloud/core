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

package roles

import (
	"errors"
	"fmt"
)

var (
	ErrRoleNotFound = errors.New("requested role not found")
)

// Role defines the role a user has within the Scope of a group.
// +kubebuilder:validation:Enum=superAdmin;admin;user;reader
type Role string

const (
	// SuperAdmin users can do anything, anywhere, and should be
	// restricted to platform operators only.
	SuperAdmin Role = "superAdmin"
	// Admin users can do anything within an organization.
	Admin Role = "admin"
	// Users can do anything within allowed projects.
	User Role = "user"
	// Readers have read-only access within allowed projects.
	Reader Role = "reader"
)

type Permission string

const (
	Create Permission = "create"
	Read   Permission = "read"
	Update Permission = "update"
	Delete Permission = "delete"
)

type Permissions []Permission

type ScopedPermissionsMap map[string]Permissions

type RoleInfo struct {
	// Permissions defined what the role can do.
	Permissions ScopedPermissionsMap
}

type RoleManager struct {
	// mapper provides a quick lookup of role.
	mapper map[Role]*RoleInfo
}

func New() *RoleManager {
	reader := &RoleInfo{
		Permissions: ScopedPermissionsMap{
			"organization": Permissions{
				Read,
			},
			"project": Permissions{
				Read,
			},
			"infrastructure": Permissions{
				Read,
			},
		},
	}

	user := &RoleInfo{
		Permissions: ScopedPermissionsMap{
			"organization": Permissions{
				Read,
			},
			"project": Permissions{
				Read,
			},
			"infrastructure": Permissions{
				Create,
				Read,
				Update,
				Delete,
			},
		},
	}

	// admin can do most things in an organization, except create and
	// delete them as creation will require a separate verification and billing
	// flow, deletion is too damned dangerous.
	admin := &RoleInfo{
		Permissions: ScopedPermissionsMap{
			"organization": Permissions{
				Read,
				Update,
			},
			"groups": Permissions{
				Create,
				Read,
				Update,
				Delete,
			},
			"oauth2provider:public": Permissions{
				Read,
			},
			"oauth2provider:private": Permissions{
				Create,
				Read,
				Update,
				Delete,
			},
			"project": Permissions{
				Create,
				Read,
				Update,
				Delete,
			},
			"infrastructure": Permissions{
				Create,
				Read,
				Update,
				Delete,
			},
		},
	}

	return &RoleManager{
		mapper: map[Role]*RoleInfo{
			Reader: reader,
			User:   user,
			Admin:  admin,
		},
	}
}

func (m *RoleManager) GetRole(role Role) (*RoleInfo, error) {
	roleInfo, ok := m.mapper[role]
	if !ok {
		return nil, fmt.Errorf("%w: %v", ErrRoleNotFound, role)
	}

	return roleInfo, nil
}
