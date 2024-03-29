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
	"github.com/unikorn-cloud/core/pkg/authorization/roles"
)

// GroupPermissions are privilege grants for a project.
type GroupPermissions struct {
	// ID is the unique, immutable project identifier.
	ID string `json:"id"`
	// Roles are the privileges a user has for the group.
	Roles []roles.Role `json:"roles"`
}

// OrganizationPermissions are privilege grants for an organization.
type OrganizationPermissions struct {
	// Name is the name of the organization.
	Name string `json:"name"`
	// Groups are any groups the user belongs to in an organization.
	Groups []GroupPermissions `json:"groups,omitempty"`
}

// Permissions are privilege grants for the entire system.
type Permissions struct {
	// IsSuperAdmin HAS SUPER COW POWERS!!!
	IsSuperAdmin bool `json:"isSuperAdmin,omitempty"`
	// Organizations are any organizations the user has access to.
	Organizations []OrganizationPermissions `json:"organizations,omitempty"`
}
