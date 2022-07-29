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
	"fmt"
	"strings"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	//"github.com/yndd/ndd_runtime/pkg/resource"
	//"sigs.k8s.io/controller-runtime/pkg/client"
)

/*
var _ TdList = &DefinitionList{}

// +k8s:deepcopy-gen=false
type TdList interface {
	client.ObjectList

	GetDefinitions() []Td
}

func (x *DefinitionList) GetDefinitions() []Td {
	xs := make([]Td, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ Td = &Definition{}

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

func (x *Template) GetNumPods() uint32 {
	if x.Spec.Properties.Fabric.Pod == nil {
		return 0
	}
	numPod := uint32(0)
	for _, p := range x.Spec.Properties.Fabric.Pod {
		numPod += *p.PodNumber
	}
	return numPod
}

func (x *FabricTemplate) CheckTemplate(master bool) error {
	if x.Pod == nil {
		return nil
	}
	// for a non master template we expect only a single pod definition
	if !master && len(x.Pod) != 1 {
		return fmt.Errorf("a child template can only have 1 pod defined")
	}
	for _, p := range x.Pod {
		if err := p.CheckPodTemplate(master); err != nil {
			return err
		}
	}
	return nil
}

func (x *FabricTemplate) HasDefinitionReference() bool {
	if x.Pod == nil {
		return false
	}
	for _, p := range x.Pod {
		if p.HasDefinitionReference() {
			return true
		}
	}
	return false
}

func (x *FabricTemplate) HasTemplateReference() bool {
	if x.Pod == nil {
		return false
	}
	for _, p := range x.Pod {
		if p.HasTemplateReference() {
			return true
		}
	}
	return false
}

func (x *FabricTemplate) HasReference() bool {
	return x.HasDefinitionReference() || x.HasTemplateReference()
}

func (x *FabricTemplate) HasPodDefinition() bool {
	if x.Pod == nil {
		return false
	}
	for _, p := range x.Pod {
		if p.Tier2 != nil {
			return true
		}
		if p.Tier3 != nil {
			return true
		}
	}
	return false
}

func (x *FabricTemplate) HasTier1() bool {
	return x.Tier1 == nil
}

func (x *FabricTemplate) HasBorderLeaf() bool {
	return x.BorderLeaf == nil
}

func (x *PodTemplate) CheckPodTemplate(master bool) error {
	// check mix of native definition
	if x.Tier2 != nil || x.Tier3 != nil {
		if x.TemplateReference != nil || x.DefinitionReference != nil {
			// this i not allowed
			return fmt.Errorf("podTemplate error: native pod definition can not be mixed with template/definition references")
		}
	}
	if master {
		// master template
		if x.HasReference() && x.PodNumber != nil {
			return fmt.Errorf("a template with a reference cannot define the pod number")
		}
		if !x.HasReference() && x.PodNumber == nil {
			return fmt.Errorf("a pod template w/o references should have a podNumber defined")
		}
	} else {
		// this is a child template
		if x.HasReference() {
			return fmt.Errorf("a child template cannot have another child template")
		}
		if x.PodNumber != nil && *x.PodNumber != 1 {
			return fmt.Errorf("a child reference can only define 1 pod instance")
		}
	}
	return nil
}

func (x *PodTemplate) HasReference() bool {
	return x.HasTemplateReference() || x.HasDefinitionReference()
}

func (x *PodTemplate) HasTemplateReference() bool {
	return x.TemplateReference != nil
}

func (x *PodTemplate) HasDefinitionReference() bool {
	return x.TemplateReference != nil
}

func (x *PodTemplate) GetPodNumber() uint32 {
	if x.PodNumber == nil {
		return 1
	}
	return *x.PodNumber
}
