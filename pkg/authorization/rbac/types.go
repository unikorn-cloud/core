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
	"github.com/unikorn-cloud/core/pkg/authorization/constants"
)

// GroupPermissions are privilege grants for a project.
type GroupPermissions struct {
	// ID is the unique, immutable project identifier.
	ID string `json:"id"`
	// Roles are the privileges a user has for the group.
	Roles []string `json:"roles"`
}

// ProjectPermissions define projects the user hass access to
// and the roles that granted those permissions.
type ProjectPermissions struct {
	// Name is the project name.
	Name string `json:"name"`
	// Roles are the privileges a user has within the project.
	Roles []string `json:"roles"`
}

// OrganizationPermissions are privilege grants for an organization.
type OrganizationPermissions struct {
	// Name is the name of the organization.
	Name string `json:"name"`
	// Groups are any groups the user belongs to in an organization.
	// These define access control lists.
	Groups []GroupPermissions `json:"groups,omitempty"`
	// Projects are any projects the user belongs to in an organization
	// via group inclusion.  These define scoping rules when accessing
	// resources.
	Projects []ProjectPermissions `json:"projects,omitempty"`
}

// Permissions are privilege grants for the entire system.
type Permissions struct {
	// IsSuperAdmin HAS SUPER COW POWERS!!!
	IsSuperAdmin bool `json:"isSuperAdmin,omitempty"`
	// Organizations are any organizations the user has access to.
	Organizations []OrganizationPermissions `json:"organizations,omitempty"`
}

// Scope maps a named API scope to a set of permissions.
type Scope struct {
	// Name is the name of the scope.
	Name string `json:"name"`
	// Permissions is the set of permissions allowed for that scope.
	Permissions []constants.Permission `json:"permissions"`
}

// ACL maps scopes to permissions.
type ACL struct {
	IsSuperAdmin bool `json:"isSuperAdmin,omitempty"`
	// Scopes is the set of scoped APIs the role can access.
	Scopes []*Scope `json:"scopes,omitempty"`
}
