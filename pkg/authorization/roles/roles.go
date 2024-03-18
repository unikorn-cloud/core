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
