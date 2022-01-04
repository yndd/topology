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
	"strconv"
	"strings"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/utils"
	nddov1 "github.com/yndd/nddo-runtime/apis/common/v1"
	"github.com/yndd/nddo-runtime/pkg/odr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ TnList = &TopologyNodeList{}

// +k8s:deepcopy-gen=false
type TnList interface {
	client.ObjectList

	GetNodes() []Tn
}

func (x *TopologyNodeList) GetNodes() []Tn {
	xs := make([]Tn, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ Tn = &TopologyNode{}

// +k8s:deepcopy-gen=false
type Tn interface {
	resource.Object
	resource.Conditioned

	GetOrganizationName() string
	GetDeploymentName() string
	GetTopologyName() string
	GetNodeName() string
	GetKindName() string
	GetAdminState() string
	GetDescription() string
	GetTags() map[string]string
	GetStateTags() map[string]string
	GetPlatform() string
	GetPosition() string
	GetNodeIndex() uint32
	GetStatus() string
	InitializeResource() error
	SetStatus(string)
	SetReason(string)
	SetPlatform(string)
	SetPosition(string)
	SetNodeEndpoint(ep *NddrTopologyTopologyLinkStateNodeEndpoint)
	GetNodeEndpoints() []*NddrTopologyTopologyNodeStateEndpoint
	SetOrganizationName(string)
	SetDeploymentName(string)
	SetTopologyName(string)
}

// GetCondition of this Network Node.
func (x *TopologyNode) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *TopologyNode) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *TopologyNode) GetOrganizationName() string {
	return odr.GetOrganizationName(x.GetNamespace())
}

func (x *TopologyNode) GetDeploymentName() string {
	return odr.GetDeploymentName(x.GetNamespace())
}

func (x *TopologyNode) GetTopologyName() string {
	split := strings.Split(x.GetName(), ".")
	if len(split) > 1 {
		return split[0]
	}
	return ""
}

func (x *TopologyNode) GetNodeName() string {
	split := strings.Split(x.GetName(), ".")
	if len(split) > 1 {
		return split[1]
	}
	return ""
}

func (x *TopologyNode) GetKindName() string {
	if reflect.ValueOf(x.Spec.TopologyNode.KindName).IsZero() {
		return ""
	}
	return *x.Spec.TopologyNode.KindName
}

func (x *TopologyNode) GetAdminState() string {
	if reflect.ValueOf(x.Spec.TopologyNode.AdminState).IsZero() {
		return ""
	}
	return *x.Spec.TopologyNode.AdminState
}

func (x *TopologyNode) GetDescription() string {
	if reflect.ValueOf(x.Spec.TopologyNode.Description).IsZero() {
		return ""
	}
	return *x.Spec.TopologyNode.Description
}

func (x *TopologyNode) GetTags() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.TopologyNode.Tag).IsZero() {
		return s
	}
	for _, tag := range x.Spec.TopologyNode.Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *TopologyNode) GetStatus() string {
	if x.Status.TopologyNode != nil && x.Status.TopologyNode.State != nil && x.Status.TopologyNode.State.Status != nil {
		return *x.Status.TopologyNode.State.Status
	}
	return "unknown"
}

func (x *TopologyNode) GetStateTags() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Status.TopologyNode.State.Tag).IsZero() {
		return s
	}
	for _, tag := range x.Status.TopologyNode.State.Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *TopologyNode) GetPlatform() string {
	if t, ok := x.GetStateTags()[KeyNodePlatform]; ok {
		return t
	}
	return ""
}

func (x *TopologyNode) GetPosition() string {
	if t, ok := x.GetTags()[KeyNodePosition]; ok {
		return t
	}
	return ""
}

func (x *TopologyNode) GetNodeIndex() uint32 {
	if t, ok := x.GetTags()[KeyNodeIndex]; ok {
		if i, err := strconv.Atoi(t); err == nil {
			return uint32(i)
		}
		return MaxUint32
	}
	return MaxUint32
}

func (x *TopologyNode) InitializeResource() error {
	tags := make([]*nddov1.Tag, 0, len(x.Spec.TopologyNode.Tag))
	for _, tag := range x.Spec.TopologyNode.Tag {
		tags = append(tags, &nddov1.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	if x.Status.TopologyNode != nil && x.Status.TopologyNode.State != nil {
		// pool was already initialiazed
		x.Status.TopologyNode.AdminState = x.Spec.TopologyNode.AdminState
		x.Status.TopologyNode.Description = x.Spec.TopologyNode.Description
		x.Status.TopologyNode.KindName = x.Spec.TopologyNode.KindName
		x.Status.TopologyNode.Tag = tags
		return nil
	}

	x.Status.TopologyNode = &NddrTopologyTopologyNode{
		AdminState:  x.Spec.TopologyNode.AdminState,
		Description: x.Spec.TopologyNode.Description,
		KindName:    x.Spec.TopologyNode.KindName,
		Tag:         tags,
		State: &NddrTopologyTopologyNodeState{
			Status:   utils.StringPtr(""),
			Reason:   utils.StringPtr(""),
			Endpoint: make([]*NddrTopologyTopologyNodeStateEndpoint, 0),
			Tag:      make([]*nddov1.Tag, 0),
		},
	}
	return nil
}

func (x *TopologyNode) SetStatus(s string) {
	x.Status.TopologyNode.State.Status = &s
}

func (x *TopologyNode) SetReason(s string) {
	x.Status.TopologyNode.State.Reason = &s
}

func (x *TopologyNode) SetPlatform(s string) {
	for _, tag := range x.Status.TopologyNode.State.Tag {
		if *tag.Key == KeyNodePlatform {
			tag.Value = &s
			return
		}
	}
	x.Status.TopologyNode.State.Tag = append(x.Status.TopologyNode.State.Tag, &nddov1.Tag{
		Key:   utils.StringPtr(KeyNodePlatform),
		Value: &s,
	})
}

func (x *TopologyNode) SetPosition(s string) {
	for _, tag := range x.Spec.TopologyNode.Tag {
		if *tag.Key == KeyNodePosition {
			tag.Value = &s
			return
		}
	}
	x.Spec.TopologyNode.Tag = append(x.Spec.TopologyNode.Tag, &nddov1.Tag{
		Key:   utils.StringPtr(KeyNodePosition),
		Value: &s,
	})
}

func (x *TopologyNode) SetNodeEndpoint(ep *NddrTopologyTopologyLinkStateNodeEndpoint) {
	if x.Status.TopologyNode.State.Endpoint == nil {
		x.Status.TopologyNode.State.Endpoint = make([]*NddrTopologyTopologyNodeStateEndpoint, 0)
	}
	for _, nodeep := range x.Status.TopologyNode.State.Endpoint {
		if *nodeep.Name == *ep.Name {
			// endpoint exists, so we update the information
			nodeep = &NddrTopologyTopologyNodeStateEndpoint{
				Name:       ep.Name,
				Lag:        ep.Lag,
				LagSubLink: ep.LagMemberLink,
			}
			return
		}
	}
	// new endpoint
	x.Status.TopologyNode.State.Endpoint = append(x.Status.TopologyNode.State.Endpoint,
		&NddrTopologyTopologyNodeStateEndpoint{
			Name:       ep.Name,
			Lag:        ep.Lag,
			LagSubLink: ep.LagMemberLink,
		})
}

func (x *TopologyNode) GetNodeEndpoints() []*NddrTopologyTopologyNodeStateEndpoint {
	if x.Status.TopologyNode != nil && x.Status.TopologyNode.State != nil && x.Status.TopologyNode.State.Endpoint != nil {
		return x.Status.TopologyNode.State.Endpoint
	}
	return make([]*NddrTopologyTopologyNodeStateEndpoint, 0)
}

func (x *TopologyNode) SetOrganizationName(s string) {
	x.Status.OrganizationName = &s
}

func (x *TopologyNode) SetDeploymentName(s string) {
	x.Status.DeploymentName = &s
}

func (x *TopologyNode) SetTopologyName(s string) {
	x.Status.TopologyName = &s
}
