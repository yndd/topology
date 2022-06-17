/*
Copyright 2021 NDD.

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

package v1alpha1

import (
	"strings"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	//"github.com/yndd/ndd_runtime/pkg/resource"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
var _ TdList = &TemplateList{}

// +k8s:deepcopy-gen=false
type TdList interface {
	client.ObjectList

	GetTemplates() []Td
}

func (x *TemplateList) GetTemplates() []Td {
	xs := make([]Td, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ Td = &Template{}

// +k8s:deepcopy-gen=false
type Td interface {
	resource.Object
	resource.Conditioned

	GetCondition(ct nddv1.ConditionKind) nddv1.Condition
	SetConditions(c ...nddv1.Condition)

	SetHealthConditions(c nddv1.HealthConditionedStatus)

	GetDeletionPolicy() nddv1.DeletionPolicy
	SetDeletionPolicy(p nddv1.DeletionPolicy)
	GetDeploymentPolicy() nddv1.DeploymentPolicy
	SetDeploymentPolicy(p nddv1.DeploymentPolicy)

	GetTargetReference() *nddv1.Reference
	SetTargetReference(p *nddv1.Reference)

	GetRootPaths() []string
	SetRootPaths(rootPaths []string)

	GetNamespacedName() string
}
*/
// GetCondition of this Network Node.
func (x *Template) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *Template) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *Template) SetHealthConditions(c nddv1.HealthConditionedStatus) {
	x.Status.Health = c
}

func (x *Template) GetDeletionPolicy() nddv1.DeletionPolicy {
	return x.Spec.Lifecycle.DeletionPolicy
}

func (x *Template) SetDeletionPolicy(c nddv1.DeletionPolicy) {
	x.Spec.Lifecycle.DeletionPolicy = c
}

func (x *Template) GetDeploymentPolicy() nddv1.DeploymentPolicy {
	return x.Spec.Lifecycle.DeploymentPolicy
}

func (x *Template) SetDeploymentPolicy(c nddv1.DeploymentPolicy) {
	x.Spec.Lifecycle.DeploymentPolicy = c
}

func (x *Template) GetTargetReference() *nddv1.Reference {
	return x.Spec.TargetReference
}

func (x *Template) SetTargetReference(p *nddv1.Reference) {
	x.Spec.TargetReference = p
}

func (x *Template) GetRootPaths() []string {
	return x.Status.RootPaths
}

func (x *Template) SetRootPaths(rootPaths []string) {
	x.Status.RootPaths = rootPaths
}

func (x *Template) GetNamespacedName() string {
	return strings.Join([]string{x.Namespace, x.Name}, "/")
}
