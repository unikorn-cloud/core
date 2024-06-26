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

func (p *Permissions) LookupOrganization(organization string) (*OrganizationPermissions, error) {
	for i, o := range p.Organizations {
		if o.Name == organization {
			return &p.Organizations[i], nil
		}
	}

	return nil, ErrPermissionDenied
}

func (p *OrganizationPermissions) HasRole(name string) bool {
	if p == nil {
		return false
	}

	for _, group := range p.Groups {
		for _, role := range group.Roles {
			if role == name {
				return true
			}
		}
	}

	return false
}

func (a *ACL) GetScope(name string) *Scope {
	for _, scope := range a.Scopes {
		if scope.Name == name {
			return scope
		}
	}

	return nil
}
