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

package conversion

import (
	"unicode"

	unikornv1 "github.com/unikorn-cloud/core/pkg/apis/unikorn/v1alpha1"
	"github.com/unikorn-cloud/core/pkg/constants"
	"github.com/unikorn-cloud/core/pkg/openapi"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

// ConvertStatusCondition translates from Kubernetes status conditions to API ones.
func ConvertStatusCondition(in *unikornv1.Condition) openapi.ResourceProvisioningStatus {
	//nolint:exhaustive
	switch in.Reason {
	case unikornv1.ConditionReasonProvisioning:
		return openapi.Provisioning
	case unikornv1.ConditionReasonProvisioned:
		return openapi.Provisioned
	case unikornv1.ConditionReasonErrored:
		return openapi.Error
	case unikornv1.ConditionReasonDeprovisioning:
		return openapi.Deprovisioning
	default:
		return openapi.Unknown
	}
}

// ResourceReadMetadata extracts generic metadata from a resource for GET APIs.
func ResourceReadMetadata(in metav1.Object, status openapi.ResourceProvisioningStatus) openapi.ResourceReadMetadata {
	labels := in.GetLabels()
	annotations := in.GetAnnotations()

	out := openapi.ResourceReadMetadata{
		Id:                 in.GetName(),
		Name:               labels[constants.NameLabel],
		CreationTime:       in.GetCreationTimestamp().Time,
		ProvisioningStatus: status,
	}

	if v, ok := annotations[constants.DescriptionAnnotation]; ok {
		out.Description = &v
	}

	if v := in.GetDeletionTimestamp(); v != nil {
		out.DeletionTime = &v.Time
	}

	return out
}

// generateResourceID creates a valid Kubernetes name from a UUID.
func generateResourceID() string {
	for {
		// NOTE: Kubernetes UUIDs are based on version 4, aka random,
		// so the first character will be a letter eventually, like
		// a 6/16 chance: tl;dr infinite loops are... improbable.
		if id := uuid.NewUUID(); unicode.IsLetter(rune(id[0])) {
			return string(id)
		}
	}
}

// ObjectMetadata creates Kubernetes object metadata from generic request metadata.
func ObjectMetadata(in *openapi.ResourceWriteMetadata, namespace string, labels map[string]string) metav1.ObjectMeta {
	out := metav1.ObjectMeta{
		Namespace: namespace,
		Name:      generateResourceID(),
		Labels: map[string]string{
			constants.NameLabel: in.Name,
		},
	}

	for k, v := range labels {
		out.Labels[k] = v
	}

	if in.Description != nil {
		out.Annotations = map[string]string{
			constants.DescriptionAnnotation: *in.Description,
		}
	}

	return out
}

// UpdateObjectMetadata abstracts away metadata updates e.g. name and description changes.
func UpdateObjectMetadata(out, in metav1.Object) {
	out.SetLabels(in.GetLabels())
	out.SetAnnotations(in.GetAnnotations())
}
