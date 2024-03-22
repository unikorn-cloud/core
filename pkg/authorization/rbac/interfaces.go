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

// Authorizer defines an interface for authorization.
type Authorizer interface {
	// AllowedByRole allows access based on role, typically used for admin only interfaces.
	AllowedByRole(role roles.Role) error
	// AllowedByGroup allows access based on groups assigned to a restricted resource.
	AllowedByGroup(groupIDs []string) error
	// AllowedByGroupRole allows access baseed on groups and the role, typically for write access.
	AllowedByGroupRole(groupIDs []string, role roles.Role) error
}
