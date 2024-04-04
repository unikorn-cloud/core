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

	"github.com/unikorn-cloud/core/pkg/authorization/constants"
)

type ACLGetter interface {
	// Get returns an ACL from whatever source the interface reprensents.
	Get(ctx context.Context) (*ACL, error)
}

// Authorizer defines an interface for authorization.
type Authorizer interface {
	// Allow allows access based on API scope and required permissions.
	Allow(ctx context.Context, scope string, permission constants.Permission) error
}
