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
	"github.com/yndd/app-runtime/pkg/odns"
	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
)

/*
var _ TnList = &NodeList{}

// +k8s:deepcopy-gen=false
type TnList interface {
	client.ObjectList

	GetNodes() []Tn
}

func (x *NodeList) GetNodes() []Tn {
	xs := make([]Tn, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ Tn = &Node{}

// +k8s:deepcopy-gen=false
type Tn interface {
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

	GetOrganization() string
	GetDeployment() string
	GetAvailabilityZone() string
	GetTopologyName() string
	GetNodeName() string
	GetVendorType() string
	GetAdminState() string
	GetDescription() string
	GetTags() map[string]string
	//GetStateTags() map[string]string
	//GetPlatform() string
	GetPosition() string
	GetNodeIndex() uint32
	InitializeResource() error

	//SetStatus(string)
	//SetReason(string)
	//SetPlatform(string)
	//SetPosition(string)
	//SetNodeEndpoint(ep *NddrTopologyTopologyLinkStateNodeEndpoint)
	//GetNodeEndpoints() []*NddrTopologyNodeStateEndpoint
	SetOrganization(s string)
	SetDeployment(s string)
	SetAvailabilityZone(s string)
	//SetTopologyName(string)
}
*/

// GetCondition of this Network Node.
func (x *Node) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *Node) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *Node) SetHealthConditions(c nddv1.HealthConditionedStatus) {
	x.Status.Health = c
}

func (x *Node) GetDeletionPolicy() nddv1.DeletionPolicy {
	return x.Spec.Lifecycle.DeletionPolicy
}

func (x *Node) SetDeletionPolicy(c nddv1.DeletionPolicy) {
	x.Spec.Lifecycle.DeletionPolicy = c
}

func (x *Node) GetDeploymentPolicy() nddv1.DeploymentPolicy {
	return x.Spec.Lifecycle.DeploymentPolicy
}

func (x *Node) SetDeploymentPolicy(c nddv1.DeploymentPolicy) {
	x.Spec.Lifecycle.DeploymentPolicy = c
}

func (x *Node) GetTargetReference() *nddv1.Reference {
	return x.Spec.TargetReference
}

func (x *Node) SetTargetReference(p *nddv1.Reference) {
	x.Spec.TargetReference = p
}

func (x *Node) GetRootPaths() []string {
	return x.Status.RootPaths
}

func (x *Node) SetRootPaths(rootPaths []string) {
	x.Status.RootPaths = rootPaths
}

func (x *Node) GetOrganization() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetOrganization()
}

func (x *Node) GetDeployment() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetDeployment()
}

func (x *Node) GetAvailabilityZone() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetAvailabilityZone()
}

func (x *Node) GetTopologyName() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetTopologyName()
}

/*

func (x *Node) GetNodeName() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetResourceName()
}

func (x *Node) GetVendorType() string {
	if reflect.ValueOf(x.Spec.Properties.VendorType).IsZero() {
		return ""
	}
	return x.Spec.Properties.VendorType
}

func (x *Node) GetAdminState() string {
	if reflect.ValueOf(x.Spec.Properties.AdminState).IsZero() {
		return ""
	}
	return x.Spec.Properties.AdminState
}

func (x *Node) GetDescription() string {
	if reflect.ValueOf(x.Spec.Properties.Description).IsZero() {
		return ""
	}
	return x.Spec.Properties.Description
}

func (x *Node) GetTags() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.Properties.Tag).IsZero() {
		return s
	}
	for _, tag := range x.Spec.Properties.Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *Node) GetPosition() string {
	if t, ok := x.GetTags()[KeyNodePosition]; ok {
		return t
	}
	return ""
}

func (x *Node) GetNodeIndex() uint32 {
	if t, ok := x.GetTags()[KeyNodeIndex]; ok {
		if i, err := strconv.Atoi(t); err == nil {
			return uint32(i)
		}
		return MaxUint32
	}
	return MaxUint32
}

func (x *Node) InitializeResource() error {
	return nil
}

func (x *Node) SetOrganization(s string) {
	x.Status.SetOrganization(s)
}

func (x *Node) SetDeployment(s string) {
	x.Status.SetDeployment(s)
}

func (x *Node) SetAvailabilityZone(s string) {
	x.Status.SetAvailabilityZone(s)
}
*/

func (x *NodeList) GetItems() []resource.Managed {
	rl := []resource.Managed{}
	for _, l := range x.Items {
		rl = append(rl, &l)
	}
	return rl
}