/*
Copyright 2021 Wim Henderickx.

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

//+kubebuilder:object:generate=true
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
)

// Condition Kinds.
const (
	// A ConditionKindAllocationReady indicates whether the allocation is ready.
	ConditionKindReady nddv1.ConditionKind = "Ready"
)

// ConditionReasons a package is or is not installed.
const (
	ConditionReasonReady        nddv1.ConditionReason = "Ready"
	ConditionReasonNotReady     nddv1.ConditionReason = "NotReady"
	ConditionReasonAllocating   nddv1.ConditionReason = "Allocating"
	ConditionReasonDeAllocating nddv1.ConditionReason = "DeAllocating"
)

// Ready indicates that the resource is ready.
func Ready() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonReady,
	}
}

// NotReady indicates that the resource is not ready.
func NotReady() nddv1.Condition {
	return nddv1.Condition{
		Kind:               ConditionKindReady,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
		Reason:             ConditionReasonNotReady,
	}
}
