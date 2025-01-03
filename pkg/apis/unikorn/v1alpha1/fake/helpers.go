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

package fake

import (
	unikornv1 "github.com/unikorn-cloud/core/pkg/apis/unikorn/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (r *ManagedResource) ResourceLabels() (labels.Set, error) {
	return labels.Set(r.Labels), nil
}

func (r *ManagedResource) Paused() bool {
	return false
}

func (r *ManagedResource) StatusConditionRead(t unikornv1.ConditionType) (*unikornv1.Condition, error) {
	return unikornv1.GetCondition(r.Status.Conditions, t)
}

func (r *ManagedResource) StatusConditionWrite(t unikornv1.ConditionType, status corev1.ConditionStatus, reason unikornv1.ConditionReason, message string) {
	unikornv1.UpdateCondition(&r.Status.Conditions, t, status, reason, message)
}
