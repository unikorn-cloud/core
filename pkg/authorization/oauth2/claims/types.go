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

// Claims is an application specific set of claims.
// TODO: this technically isn't conformant to oauth2 in that we don't specify
// the client_id claim, and there are probably others.
type Claims struct {
	jwt.Claims `json:",inline"`

	// Organization is the top level organization the user belongs to.
	Organization string `json:"org"`

	// Scope is the oauth2 scope of the token.
	Scope scope.Scope `json:"scope,omitempty"`
}
