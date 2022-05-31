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
)

/*
var _ TpList = &TopologyList{}

// +k8s:deepcopy-gen=false
type TpList interface {
	client.ObjectList

	GetTopologies() []Tp
}

func (x *TopologyList) GetTopologies() []Tp {
	xs := make([]Tp, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ Tp = &Topology{}

// +k8s:deepcopy-gen=false
type Tp interface {
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
	GetAdminState() string
	GetDescription() string
	GetDefaultsTags() map[string]string
	GetVendorTypeInfo() []*VendorTypeInfo
	GetKindVendorTypes() []string
	GetKindTagsByVendorType(string) map[string]string
	GetPlatformByVendorType(string) string
	GetPlatformFromDefaults() string
	InitializeResource() error

	SetOrganization(s string)
	SetDeployment(s string)
	SetAvailabilityZone(s string)
}
*/
// GetCondition of this Network Node.
func (x *Topology) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *Topology) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *Topology) SetHealthConditions(c nddv1.HealthConditionedStatus) {
	x.Status.Health = c
}

func (x *Topology) GetDeletionPolicy() nddv1.DeletionPolicy {
	return x.Spec.Lifecycle.DeletionPolicy
}

func (x *Topology) SetDeletionPolicy(c nddv1.DeletionPolicy) {
	x.Spec.Lifecycle.DeletionPolicy = c
}

func (x *Topology) GetDeploymentPolicy() nddv1.DeploymentPolicy {
	return x.Spec.Lifecycle.DeploymentPolicy
}

func (x *Topology) SetDeploymentPolicy(c nddv1.DeploymentPolicy) {
	x.Spec.Lifecycle.DeploymentPolicy = c
}

func (x *Topology) GetTargetReference() *nddv1.Reference {
	return x.Spec.TargetReference
}

func (x *Topology) SetTargetReference(p *nddv1.Reference) {
	x.Spec.TargetReference = p
}

func (x *Topology) GetRootPaths() []string {
	return x.Status.RootPaths
}

func (x *Topology) SetRootPaths(rootPaths []string) {
	x.Status.RootPaths = rootPaths
}

func (x *Topology) GetOrganization() string {
	return odns.Name2OdnsTopo(x.GetName()).GetOrganization()
}

func (x *Topology) GetDeployment() string {
	return odns.Name2OdnsTopo(x.GetName()).GetDeployment()
}

func (x *Topology) GetAvailabilityZone() string {
	return odns.Name2OdnsTopo(x.GetName()).GetAvailabilityZone()
}

func (x *Topology) GetTopologyName() string {
	return odns.Name2OdnsTopo(x.GetName()).GetTopologyName()
}

/*

func (x *Topology) GetDefaultsTags() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.Properties.Defaults).IsZero() ||
		reflect.ValueOf(x.Spec.Properties.Defaults.Tag).IsZero() {
		return s
	}
	for _, tag := range x.Spec.Properties.Defaults.Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *Topology) GetPlatformFromDefaults() string {
	if p, ok := x.GetDefaultsTags()[KeyNodePlatform]; ok {
		return p
	}
	return ""
}

func (x *Topology) GetVendorTypeInfo() []*VendorTypeInfo {
	if reflect.ValueOf(x.Spec.Properties.VendorTypeInfo).IsZero() {
		return nil
	}
	return x.Spec.Properties.VendorTypeInfo
}

func (x *Topology) GetKindVendorTypes() []string {
	s := make([]string, 0)
	if reflect.ValueOf(x.Spec.Properties.VendorTypeInfo).IsZero() {
		return s
	}
	for _, vendorTypeInfo := range x.Spec.Properties.VendorTypeInfo {
		s = append(s, vendorTypeInfo.VendorType)
	}
	return s
}

func (x *Topology) GetKindTagsByVendorType(vendorType string) map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.Properties.VendorTypeInfo).IsZero() {
		return s
	}
	for _, vendorTypeInfo := range x.Spec.Properties.VendorTypeInfo {
		if vendorType == vendorTypeInfo.VendorType {
			for _, tag := range vendorTypeInfo.Tag {
				s[*tag.Key] = *tag.Value
			}
		}
	}
	return s
}

func (x *Topology) GetPlatformByVendorType(vendorType string) string {
	if reflect.ValueOf(x.Spec.Properties.VendorTypeInfo).IsZero() {
		return ""
	}
	for _, vendorTypeInfo := range x.Spec.Properties.VendorTypeInfo {
		if vendorType == vendorTypeInfo.VendorType {
			return vendorTypeInfo.Platform
		}
	}
	return ""
}

func (x *Topology) InitializeResource() error {
	return nil
}

func (x *Topology) SetOrganization(s string) {
	x.Status.SetOrganization(s)
}

func (x *Topology) SetDeployment(s string) {
	x.Status.SetDeployment(s)
}

func (x *Topology) SetAvailabilityZone(s string) {
	x.Status.SetAvailabilityZone(s)
}
*/
