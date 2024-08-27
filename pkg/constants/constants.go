/*
Copyright 2022-2024 EscherCloud.
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

package constants

import (
	"time"
)

const (
	// This is the default version in the Makefile.
	DeveloperVersion = "0.0.0"

	// NameLabel is attached to every resource to give it a mutable display
	// name.  While the character set is limited to [0-9A-Za-z_-.] it is at least
	// indexed in etcd which gives us another string to our bow.
	NameLabel = "unikorn-cloud.org/name"

	// DescriptionAnnotation is optionally attached to a resource to allow
	// an unconstriained and verbose description about the resource.
	DescriptionAnnotation = "unikorn-cloud.org/description"

	// CreatorAnnotation is optionally attached to a resource to show who
	// created it.
	CreatorAnnotation = "unikorn-cloud.org/creator"

	// ModifierAnnotation is optionally attached to a resource to show who
	// last modified it.
	ModifierAnnotation = "unikorn-cloud.org/modifier"

	// ModifiedTimestampAnnotation augments Kubernetes metadata to provide mtime
	// like functionality.
	ModifiedTimestampAnnotation = "unikorn-cloud.org/modifiedTimestamp"

	// KindLabel is used to match a resource that may be owned by a particular kind.
	// For example, projects and cluster managers are modelled on namespaces.  For CPs
	// you have to select based on project and CP name, because of name reuse, but
	// this raises the problem that selecting a project's namespace will match multiple
	// so this provides a concrete type associated with each resource.
	KindLabel = "unikorn-cloud.org/kind"

	// KindLabelValueOrganization is used to denote a resource belongs to this type.
	KindLabelValueOrganization = "organization"

	// KindLabelValueProject is used to denote a resource belongs to this type.
	KindLabelValueProject = "project"

	// KindLabelValueClusterManager is used to denote a resource belongs to this type.
	KindLabelValueClusterManager = "clustermanager"

	// KindLabelValueKubernetesCluster is used to denote a resource belongs to this type.
	KindLabelValueKubernetesCluster = "kubernetescluster"

	// KindLabelValueBaremetalCluster is used to denote a resource belongs to this type.
	KindLabelValueBaremetalCluster = "baremetalcluster"

	// OrganizationLabel is a label applied to namespaces to indicate it is under
	// control of this tool.  Useful for label selection.
	OrganizationLabel = "unikorn-cloud.org/organization"

	// ProjectLabel is a label applied to namespaces to indicate it is under
	// control of this tool.  Useful for label selection.
	ProjectLabel = "unikorn-cloud.org/project"

	// KubernetesClusterLabel is applied to resources to indicate it belongs
	// to a specific cluster.
	KubernetesClusterLabel = "unikorn-cloud.org/kubernetescluster"

	// BaremetalClusterLabel is applied to resources to indicate it belongs
	// to a specific cluster.
	BaremetalClusterLabel = "unikorn-cloud.org/baremetalcluster"

	// ApplicationLabel is applied to ArgoCD applications to differentiate
	// between them.
	ApplicationLabel = "unikorn-cloud.org/application"

	// ApplicationIDLabel is used to lookup applications based on their ID.
	ApplicationIDLabel = "unikorn-cloud.org/application-id"

	// ConfigurationHashAnnotation is used where application owners refuse to
	// poll configuration updates and we (and all other users) are forced into
	// manually restarting services based on a Deployment/DaemonSet changing.
	ConfigurationHashAnnotation = "unikorn-cloud.org/config-hash"

	// IdentityAnnotation tells you the cloud identity (in the context of
	// the region controller) that a resource owns.
	IdentityAnnotation = "unikorn-cloud.org/identity-id"

	// PhysicalNetworkAnnotation tells you the physical network (in the
	// context of a region controller) that a recource owns.
	PhysicalNetworkAnnotation = "unikorn-cloud.org/physical-network-id"

	// Finalizer is applied to resources that need to be deleted manually
	// and do other complex logic.
	Finalizer = "unikorn"

	// DefaultYieldTimeout allows N seconds for a provisioner to do its thing
	// and report a healthy status before yielding and giving someone else
	// a go.
	DefaultYieldTimeout = 10 * time.Second
)

// LabelPriorities assigns a priority to the labels for sorting.  Most things
// use the labels to uniquely identify a resource.  For example, when we create
// a remote cluster in ArgoCD we use a tuple of project, cluster manager and optionally
// the cluster.  This gives a unique identifier given projects and cluster managers
// provide a namespace abstraction, and a deterministic one as the order is defined.
// This function is required because labels are given as a map, and thus are
// no-deterministically ordered when iterating in go.
func LabelPriorities() []string {
	return []string{
		KubernetesClusterLabel,
		ProjectLabel,
		OrganizationLabel,
	}
}
