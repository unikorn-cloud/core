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

package claims

import (
	"github.com/go-jose/go-jose/v3/jwt"

	"github.com/unikorn-cloud/core/pkg/authorization/oauth2/scope"
)

// Role defines the role a user has within the scope of a group.
// +kubebuilder:validation:Enum=superAdmin;admin;user;reader
type Role string

const (
	// SuperAdmin users can do anything, anywhere, and should be
	// restricted to platform operators only.
	SuperAdmin = "superAdmin"
	// Admin users can do anything within an organization.
	Admin Role = "admin"
	// Users can do anything within allowed projects.
	User Role = "user"
	// Readers have read-only access within allowed projects.
	Reader Role = "reader"
)

// Group records RBAC data in the claims.
type Group struct {
	// ID is the immutable group ID.
	ID string `json:"id"`
	// Roles are a list of roles the group possesses.
	Roles []Role `json:"roles,omitempty"`
}

// UnikornClaims contains all application specific claims in a single
// top-level claim that won't clash with the ones defined by IETF.
type UnikornClaims struct {
	// Organization is the top level organization the user belongs to.
	Organization string `json:"org"`

	// Groups is a list of groups and roles the token has access to.
	// Resources should be scoped to some group/groups that the resource
	// server can filter based on the access token.  Then it can determine
	// what operations are allowed based on the roles assigned to those
	// groups.
	Groups []Group `json:"groups,omitempty"`
}

// Claims is an application specific set of claims.
// TODO: this technically isn't conformant to oauth2 in that we don't specify
// the client_id claim, and there are probably others.
type Claims struct {
	jwt.Claims `json:",inline"`

	// Scope is the oauth2 scope of the token.
	Scope scope.Scope `json:"scope,omitempty"`

	// Unikorn claims are application specific extensions.
	Unikorn *UnikornClaims `json:"unikorn,omitempty"`
}
