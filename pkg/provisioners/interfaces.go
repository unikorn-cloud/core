/*
Copyright 2022-2024 EscherCloud.
Copyright 2024-2025 the Unikorn Authors.

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

package provisioners

//go:generate mockgen -source=interfaces.go -destination=mock/interfaces.go -package=mock

import (
	"context"

	unikornv1 "github.com/unikorn-cloud/core/pkg/apis/unikorn/v1alpha1"
	"github.com/unikorn-cloud/core/pkg/cd"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Generator is an abstraction around the sources of remote
// clusters e.g. a cluster API or vcluster Kubernetes instance.
type RemoteCluster interface {
	// ID is the unique resource identifier for this remote cluster.
	ID() *cd.ResourceIdentifier

	// Config returns the client configuration (aka parsed Kubeconfig.)
	Config(ctx context.Context) (*clientcmdapi.Config, error)
}

// Provisioner is an abstract type that allows provisioning of Kubernetes
// packages in a technology agnostic way.  For example some things may be
// installed as a raw set of resources, a YAML manifest, Helm etc.
type Provisioner interface {
	// ProvisionerName returns the provisioner name.
	ProvisionerName() string

	// Provision deploys the requested package.
	// Implementations should ensure this receiver is idempotent.
	Provision(ctx context.Context) error

	// Deprovision does any special handling of resource/component
	// removal.  In the general case, you should rely on cascading
	// deletion i.e. kill the namespace, use owner references.
	// Deprovisioners should be gating, waiting for their resources
	// to be removed before indicating success.
	Deprovision(ctx context.Context) error
}

// ManagerProvisioner top-level manager provisioners hook directly into
// the controller runtime layer, and are a little special in that they
// abstract away type specific things.
type ManagerProvisioner interface {
	Provisioner

	// Object returns a reference to the generic object type, internally
	// the provisioner will have a type specific version.
	Object() unikornv1.ManagableResourceInterface
}
