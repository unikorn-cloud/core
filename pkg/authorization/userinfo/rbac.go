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

package userinfo

import (
	"context"

	"github.com/unikorn-cloud/core/pkg/authorization/rbac"
)

func NewScopedAuthorizer(ctx context.Context, organization string) (rbac.Authorizer, error) {
	userinfo := FromContext(ctx)

	return rbac.NewScoped(userinfo.RBAC, organization)
}

func NewUnscopedAuthorizer(ctx context.Context) (rbac.Authorizer, error) {
	userinfo := FromContext(ctx)

	return rbac.NewUnscoped(userinfo.RBAC)
}
