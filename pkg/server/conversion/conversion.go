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
	"context"
	"errors"
	"fmt"
	"time"
	"unicode"

	unikornv1 "github.com/unikorn-cloud/core/pkg/apis/unikorn/v1alpha1"
	"github.com/unikorn-cloud/core/pkg/authorization/userinfo"
	"github.com/unikorn-cloud/core/pkg/constants"
	"github.com/unikorn-cloud/core/pkg/openapi"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
)

var (
	ErrAnnotation = errors.New("a required annotation was missing")
)

// ConvertStatusCondition translates from Kubernetes status conditions to API ones.
func ConvertStatusCondition(in *unikornv1.Condition) openapi.ResourceProvisioningStatus {
	//nolint:exhaustive
	switch in.Reason {
	case unikornv1.ConditionReasonProvisioning:
		return openapi.ResourceProvisioningStatusProvisioning
	case unikornv1.ConditionReasonProvisioned:
		return openapi.ResourceProvisioningStatusProvisioned
	case unikornv1.ConditionReasonErrored:
		return openapi.ResourceProvisioningStatusError
	case unikornv1.ConditionReasonDeprovisioning:
		return openapi.ResourceProvisioningStatusDeprovisioning
	default:
		return openapi.ResourceProvisioningStatusUnknown
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

	if v, ok := annotations[constants.CreatorAnnotation]; ok {
		out.CreatedBy = &v
	}

	if v, ok := annotations[constants.ModifierAnnotation]; ok {
		out.ModifiedBy = &v
	}

	if v, ok := annotations[constants.ModifiedTimestampAnnotation]; ok {
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			out.ModifiedTime = &t
		}
	}

	if v := in.GetDeletionTimestamp(); v != nil {
		out.DeletionTime = &v.Time
	}

	return out
}

// OrganizationScopedResourceReadMetadata extracts organization scoped metdata from a resource
// for GET APIS.
func OrganizationScopedResourceReadMetadata(in metav1.Object, status openapi.ResourceProvisioningStatus) openapi.OrganizationScopedResourceReadMetadata {
	labels := in.GetLabels()

	temp := ResourceReadMetadata(in, status)

	out := openapi.OrganizationScopedResourceReadMetadata{
		Id:                 temp.Id,
		Name:               temp.Name,
		Description:        temp.Description,
		CreatedBy:          temp.CreatedBy,
		CreationTime:       temp.CreationTime,
		ModifiedBy:         temp.ModifiedBy,
		ModifiedTime:       temp.ModifiedTime,
		ProvisioningStatus: temp.ProvisioningStatus,
		OrganizationId:     labels[constants.OrganizationLabel],
	}

	return out
}

// ProjectScopedResourceReadMetadata extracts project scoped metdata from a resource for
// GET APIs.
func ProjectScopedResourceReadMetadata(in metav1.Object, status openapi.ResourceProvisioningStatus) openapi.ProjectScopedResourceReadMetadata {
	labels := in.GetLabels()

	temp := OrganizationScopedResourceReadMetadata(in, status)

	out := openapi.ProjectScopedResourceReadMetadata{
		Id:                 temp.Id,
		Name:               temp.Name,
		Description:        temp.Description,
		CreatedBy:          temp.CreatedBy,
		CreationTime:       temp.CreationTime,
		ModifiedBy:         temp.ModifiedBy,
		ModifiedTime:       temp.ModifiedTime,
		ProvisioningStatus: temp.ProvisioningStatus,
		OrganizationId:     temp.OrganizationId,
		ProjectId:          labels[constants.ProjectLabel],
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

// ObjectMetadata implements a builder pattern.
type ObjectMetadata struct {
	metadata       *openapi.ResourceWriteMetadata
	namespace      string
	organizationID string
	projectID      string
	labels         map[string]string
}

// NewObjectMetadata requests the bare minimum to build an object metadata object.
func NewObjectMetadata(metadata *openapi.ResourceWriteMetadata, namespace string) *ObjectMetadata {
	return &ObjectMetadata{
		metadata:  metadata,
		namespace: namespace,
	}
}

// WithOrganization adds an organization for scoped resources.
func (o *ObjectMetadata) WithOrganization(id string) *ObjectMetadata {
	o.organizationID = id

	return o
}

// WithProject adds a project for scoped resources.
func (o *ObjectMetadata) WithProject(id string) *ObjectMetadata {
	o.projectID = id

	return o
}

// WithLabel allows non-generic labels to be attached to a resource.
func (o *ObjectMetadata) WithLabel(key, value string) *ObjectMetadata {
	if o.labels == nil {
		o.labels = map[string]string{}
	}

	o.labels[key] = value

	return o
}

// Get renders the object metadata ready for inclusion into a Kubernetes resource.
func (o *ObjectMetadata) Get(ctx context.Context) metav1.ObjectMeta {
	userinfo := userinfo.FromContext(ctx)

	out := metav1.ObjectMeta{
		Namespace: o.namespace,
		Name:      generateResourceID(),
		Labels: map[string]string{
			constants.NameLabel: o.metadata.Name,
		},
		Annotations: map[string]string{
			constants.CreatorAnnotation: userinfo.Subject,
		},
	}

	if o.organizationID != "" {
		out.Labels[constants.OrganizationLabel] = o.organizationID
	}

	if o.projectID != "" {
		out.Labels[constants.ProjectLabel] = o.projectID
	}

	for k, v := range o.labels {
		out.Labels[k] = v
	}

	if o.metadata.Description != nil {
		out.Annotations[constants.DescriptionAnnotation] = *o.metadata.Description
	}

	return out
}

// UpdateObjectMetadata abstracts away metadata updates.
func UpdateObjectMetadata(required, current metav1.Object, persistedAnnotations ...string) error {
	requiredAnnotations := required.GetAnnotations()
	currentAnnotations := current.GetAnnotations()

	// Persist any component specific annotations.
	for _, annotation := range persistedAnnotations {
		v, ok := currentAnnotations[annotation]
		if !ok {
			return fmt.Errorf("%w: %s", ErrAnnotation, annotation)
		}

		requiredAnnotations[annotation] = v
	}

	// When updating, the required creator is now the updater.
	requiredAnnotations[constants.ModifierAnnotation] = requiredAnnotations[constants.CreatorAnnotation]
	requiredAnnotations[constants.ModifiedTimestampAnnotation] = time.Now().UTC().Format(time.RFC3339)

	// And preserve the original creator.
	requiredAnnotations[constants.CreatorAnnotation] = currentAnnotations[constants.CreatorAnnotation]

	if v, ok := currentAnnotations[constants.CreatorAnnotation]; ok {
		requiredAnnotations[constants.CreatorAnnotation] = v
	}

	required.SetAnnotations(requiredAnnotations)

	return nil
}
