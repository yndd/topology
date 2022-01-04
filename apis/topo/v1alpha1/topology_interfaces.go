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
	"reflect"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/utils"
	"github.com/yndd/nddo-runtime/pkg/odr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

	GetOrganizationName() string
	GetDeploymentName() string
	GetTopologyName() string
	GetAdminState() string
	GetDescription() string
	GetDefaultsTags() map[string]string
	GetKinds() []*TopoTopologyKind
	GetKindNames() []string
	GetKindTagsByName(string) map[string]string
	GetPlatformByKindName(string) string
	GetPlatformFromDefaults() string
	InitializeResource() error

	SetStatus(string)
	SetReason(string)
	GetStatus() string
	SetOrganizationName(string)
	SetDeploymentName(string)
	SetTopologyName(string)
}

// GetCondition of this Network Node.
func (x *Topology) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *Topology) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *Topology) GetOrganizationName() string {
	return odr.GetOrganizationName(x.GetNamespace())
}

func (x *Topology) GetDeploymentName() string {
	return odr.GetDeploymentName(x.GetNamespace())
}

func (x *Topology) GetTopologyName() string {
	return x.GetName()
}

func (x *Topology) GetAdminState() string {
	if reflect.ValueOf(x.Spec.Topology.AdminState).IsZero() {
		return ""
	}
	return *x.Spec.Topology.AdminState
}

func (x *Topology) GetDescription() string {
	if reflect.ValueOf(x.Spec.Topology.Description).IsZero() {
		return ""
	}
	return *x.Spec.Topology.Description
}

func (x *Topology) GetDefaultsTags() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.Topology.Defaults).IsZero() ||
		reflect.ValueOf(x.Spec.Topology.Defaults.Tag).IsZero() {
		return s
	}
	for _, tag := range x.Spec.Topology.Defaults.Tag {
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

func (x *Topology) GetKinds() []*TopoTopologyKind {
	if reflect.ValueOf(x.Spec.Topology.Kind).IsZero() {
		return nil
	}
	return x.Spec.Topology.Kind
}

func (x *Topology) GetKindNames() []string {
	s := make([]string, 0)
	if reflect.ValueOf(x.Spec.Topology.Kind).IsZero() {
		return s
	}
	for _, kind := range x.Spec.Topology.Kind {
		s = append(s, *kind.Name)
	}
	return s
}

func (x *Topology) GetKindTagsByName(name string) map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.Topology.Kind).IsZero() {
		return s
	}
	for _, kind := range x.Spec.Topology.Kind {
		if name == *kind.Name {
			for _, tag := range kind.Tag {
				s[*tag.Key] = *tag.Value
			}
		}
	}
	return s
}

func (x *Topology) GetPlatformByKindName(name string) string {
	if reflect.ValueOf(x.Spec.Topology.Kind).IsZero() {
		return ""
	}
	for _, kind := range x.Spec.Topology.Kind {
		if name == *kind.Name {
			for _, tag := range kind.Tag {
				if *tag.Key == KeyNodePlatform {
					return *tag.Value
				}
			}
		}
	}
	return ""
}

func (x *Topology) InitializeResource() error {
	if x.Status.Topology != nil && x.Status.Topology.State != nil {
		// pool was already initialiazed
		// copy the spec, but not the state
		x.Status.Topology.AdminState = x.Spec.Topology.AdminState
		x.Status.Topology.Description = x.Spec.Topology.Description
		return nil
	}

	x.Status.Topology = &NddrTopologyTopology{
		AdminState:  x.Spec.Topology.AdminState,
		Description: x.Spec.Topology.Description,
		State: &NddrTopologyTopologyState{
			Status: utils.StringPtr(""),
			Reason: utils.StringPtr(""),
		},
	}
	return nil
}

func (x *Topology) SetStatus(s string) {
	x.Status.Topology.State.Status = &s
}

func (x *Topology) SetReason(s string) {
	x.Status.Topology.State.Reason = &s
}

func (x *Topology) GetStatus() string {
	if x.Status.Topology != nil && x.Status.Topology.State != nil && x.Status.Topology.State.Status != nil {
		return *x.Status.Topology.State.Status
	}
	return "unknown"
}

func (x *Topology) SetOrganizationName(s string) {
	x.Status.OrganizationName = &s
}

func (x *Topology) SetDeploymentName(s string) {
	x.Status.DeploymentName = &s
}

func (x *Topology) SetTopologyName(s string) {
	x.Status.TopologyName = &s
}
