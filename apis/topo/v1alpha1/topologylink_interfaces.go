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
	nddov1 "github.com/yndd/nddo-runtime/apis/common/v1"
	"github.com/yndd/nddo-runtime/pkg/odns"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ TlList = &TopologyLinkList{}

// +k8s:deepcopy-gen=false
type TlList interface {
	client.ObjectList

	GetLinks() []Tl
}

func (x *TopologyLinkList) GetLinks() []Tl {
	xs := make([]Tl, len(x.Items))
	for i, r := range x.Items {
		r := r // Pin range variable so we can take its address.
		xs[i] = &r
	}
	return xs
}

var _ Tl = &TopologyLink{}

// +k8s:deepcopy-gen=false
type Tl interface {
	resource.Object
	resource.Conditioned

	GetCondition(ct nddv1.ConditionKind) nddv1.Condition
	SetConditions(c ...nddv1.Condition)
	GetOrganization() string
	GetDeployment() string
	GetAvailabilityZone() string
	GetTopologyName() string
	GetLinkName() string
	GetAdminState() string
	GetDescription() string
	GetTags() map[string]string
	GetEndpoints() []*TopoTopologyLinkEndpoints
	GetEndpointANodeName() string
	GetEndpointBNodeName() string
	GetEndpointAInterfaceName() string
	GetEndpointBInterfaceName() string
	GetEndpointATag() map[string]string
	GetEndpointBTag() map[string]string
	GetEndpointATagRaw() []*nddov1.Tag
	GetEndpointBTagRaw() []*nddov1.Tag
	GetEndPointAKind() string
	GetEndPointBKind() string
	GetEndPointAGroup() string
	GetEndPointBGroup() string
	GetEndPointAMultiHoming() bool
	GetEndPointBMultiHoming() bool
	GetEndPointAMultiHomingName() string
	GetEndPointBMultiHomingName() string
	GetLagMember() bool
	GetLag() bool
	GetLacp() bool
	GetLacpFallbackA() bool
	GetLacpFallbackB() bool
	GetLagAName() string
	GetLagBName() string
	GetStatus() string
	GetNodes() []*NddrTopologyTopologyLinkStateNode
	GetStatusTagsRaw() []*nddov1.Tag
	InitializeResource() error
	SetStatus(string)
	SetReason(string)
	SetNodeEndpoint(nodeName string, ep *NddrTopologyTopologyLinkStateNodeEndpoint)
	GetNodeEndpoints() []*NddrTopologyTopologyLinkStateNode
	SetKind(s string)
	GetKind() string
	AddEndPointATag(string, string)
	AddEndPointBTag(string, string)
	DeleteEndPointATag(key string, value string)
	DeleteEndPointBTag(key string, value string)

	SetOrganization(s string)
	SetDeployment(s string)
	SetAvailabilityZone(s string)
	SetTopologyName(string)
}

// GetCondition of this Network Node.
func (x *TopologyLink) GetCondition(ct nddv1.ConditionKind) nddv1.Condition {
	return x.Status.GetCondition(ct)
}

// SetConditions of the Network Node.
func (x *TopologyLink) SetConditions(c ...nddv1.Condition) {
	x.Status.SetConditions(c...)
}

func (x *TopologyLink) GetOrganization() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetOrganization()
}

func (x *TopologyLink) GetDeployment() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetDeployment()
}

func (x *TopologyLink) GetAvailabilityZone() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetAvailabilityZone()
}

func (x *TopologyLink) GetTopologyName() string {
	return odns.Name2OdnsTopoResource(x.GetName()).GetTopologyName()
}

func (x *TopologyLink) GetLinkName() string {
	return x.GetName()
}

func (x *TopologyLink) GetAdminState() string {
	if reflect.ValueOf(x.Spec.TopologyLink.AdminState).IsZero() {
		return ""
	}
	return *x.Spec.TopologyLink.AdminState
}

func (x *TopologyLink) GetDescription() string {
	if reflect.ValueOf(x.Spec.TopologyLink.Description).IsZero() {
		return ""
	}
	return *x.Spec.TopologyLink.Description
}

func (x *TopologyLink) GetTags() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.TopologyLink.Tag).IsZero() {
		return s
	}
	for _, tag := range x.Spec.TopologyLink.Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *TopologyLink) GetStatusTagsRaw() []*nddov1.Tag {
	if x.Status.TopologyLink != nil && x.Status.TopologyLink.Tag != nil {
		return x.Status.TopologyLink.Tag
	}
	return make([]*nddov1.Tag, 0)
}

func (x *TopologyLink) GetEndpoints() []*TopoTopologyLinkEndpoints {
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return nil
	}
	return x.Spec.TopologyLink.Endpoints
}

func (x *TopologyLink) GetEndpointANodeName() string {
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return ""
	}
	return *x.Spec.TopologyLink.Endpoints[0].NodeName
}

func (x *TopologyLink) GetEndpointBNodeName() string {
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return ""
	}
	return *x.Spec.TopologyLink.Endpoints[1].NodeName
}

func (x *TopologyLink) GetEndpointAInterfaceName() string {
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return ""
	}
	return *x.Spec.TopologyLink.Endpoints[0].InterfaceName
}

func (x *TopologyLink) GetEndpointBInterfaceName() string {
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return ""
	}
	return *x.Spec.TopologyLink.Endpoints[1].InterfaceName
}

func (x *TopologyLink) GetEndpointATag() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return s
	}
	for _, tag := range x.Spec.TopologyLink.Endpoints[0].Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *TopologyLink) GetEndpointBTag() map[string]string {
	s := make(map[string]string)
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return s
	}
	for _, tag := range x.Spec.TopologyLink.Endpoints[1].Tag {
		s[*tag.Key] = *tag.Value
	}
	return s
}

func (x *TopologyLink) GetEndpointATagRaw() []*nddov1.Tag {
	s := make([]*nddov1.Tag, 0)
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return s
	}
	return x.Spec.TopologyLink.Endpoints[0].Tag
}

func (x *TopologyLink) GetEndpointBTagRaw() []*nddov1.Tag {
	s := make([]*nddov1.Tag, 0)
	if reflect.ValueOf(x.Spec.TopologyLink.Endpoints).IsZero() {
		return s
	}
	return x.Spec.TopologyLink.Endpoints[1].Tag
}

func (x *TopologyLink) GetEndPointAKind() string {
	if n, ok := x.GetEndpointATag()[KeyLinkEPKind]; ok {
		return n
	}
	// default
	return "infra"
}

func (x *TopologyLink) GetEndPointBKind() string {
	if n, ok := x.GetEndpointBTag()[KeyLinkEPKind]; ok {
		return n
	}
	// default
	return "infra"
}

func (x *TopologyLink) GetEndPointAGroup() string {
	if n, ok := x.GetEndpointATag()[KeyLinkEPGroup]; ok {
		return n
	}
	// default
	return ""
}

func (x *TopologyLink) GetEndPointBGroup() string {
	if n, ok := x.GetEndpointBTag()[KeyLinkEPGroup]; ok {
		return n
	}
	// default
	return ""
}

func (x *TopologyLink) GetEndPointAMultiHoming() bool {
	if n, ok := x.GetEndpointATag()[KeyLinkEPMultiHoming]; ok {
		return n == "true"
	}
	// default
	return false
}

func (x *TopologyLink) GetEndPointBMultiHoming() bool {
	if n, ok := x.GetEndpointBTag()[KeyLinkEPMultiHoming]; ok {
		return n == "true"
	}
	// default
	return false
}

func (x *TopologyLink) GetEndPointAMultiHomingName() string {
	if n, ok := x.GetEndpointATag()[KeyLinkEPMultiHomingName]; ok {
		return n
	}
	// default
	return ""
}

func (x *TopologyLink) GetEndPointBMultiHomingName() string {
	if n, ok := x.GetEndpointBTag()[KeyLinkEPMultiHomingName]; ok {
		return n
	}
	// default
	return ""
}

func (x *TopologyLink) GetLacpFallbackA() bool {
	if _, ok := x.GetEndpointATag()[KeyLinkEPLacpFallback]; ok {
		return x.GetTags()[KeyLinkEPLacpFallback] == "true"
	}
	// default
	return false
}

func (x *TopologyLink) GetLacpFallbackB() bool {
	if _, ok := x.GetEndpointBTag()[KeyLinkEPLacpFallback]; ok {
		return x.GetTags()[KeyLinkEPLacpFallback] == "true"
	}
	// default
	return false
}

func (x *TopologyLink) GetLagMember() bool {
	if _, ok := x.GetTags()[KeyLinkLagMember]; ok {
		return x.GetTags()[KeyLinkLagMember] == "true"
	}
	// default is false
	return false
}

func (x *TopologyLink) GetLag() bool {
	if _, ok := x.GetTags()[KeyLinkLag]; ok {
		return x.GetTags()[KeyLinkLag] == "true"
	}
	// default is false
	return false
}

func (x *TopologyLink) GetLacp() bool {
	if _, ok := x.GetTags()[KeyLinkLacp]; ok {
		return x.GetTags()[KeyLinkLacp] == "true"
	}
	// default is true
	return true
}

func (x *TopologyLink) GetLagAName() string {
	if n, ok := x.GetEndpointATag()[KeyLinkEPLagName]; ok {
		return n
	}
	return ""
}

func (x *TopologyLink) GetLagBName() string {
	if n, ok := x.GetEndpointBTag()[KeyLinkEPLagName]; ok {
		return n
	}
	return ""
}

func (x *TopologyLink) GetStatus() string {
	if x.Status.TopologyLink != nil && x.Status.TopologyLink.State != nil && x.Status.TopologyLink.State.Status != nil {
		return *x.Status.TopologyLink.State.Status
	}
	return "unknown"
}

func (x *TopologyLink) GetNodes() []*NddrTopologyTopologyLinkStateNode {
	return x.Status.TopologyLink.State.Node
}

func (x *TopologyLink) InitializeResource() error {
	eps := make([]*NddrTopologyTopologyLinkEndpoints, 0, len(x.Spec.TopologyLink.Endpoints))
	for _, ep := range x.Spec.TopologyLink.Endpoints {
		epTags := make([]*nddov1.Tag, 0, len(ep.Tag))
		for _, tag := range ep.Tag {
			epTags = append(epTags, &nddov1.Tag{
				Key:   tag.Key,
				Value: tag.Value,
			})
		}

		eps = append(eps, &NddrTopologyTopologyLinkEndpoints{
			InterfaceName: ep.InterfaceName,
			NodeName:      ep.NodeName,
			Tag:           epTags,
		})
	}

	tags := make([]*nddov1.Tag, 0, len(x.Spec.TopologyLink.Tag))
	for _, tag := range x.Spec.TopologyLink.Tag {
		tags = append(tags, &nddov1.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}

	if x.Status.TopologyLink != nil && x.Status.TopologyLink.State != nil {
		x.Status.TopologyLink.AdminState = x.Spec.TopologyLink.AdminState
		x.Status.TopologyLink.Description = x.Spec.TopologyLink.Description
		x.Status.TopologyLink.Endpoints = eps
		x.Status.TopologyLink.Tag = tags
		return nil
	}

	x.Status.TopologyLink = &NddrTopologyTopologyLink{
		Name:        x.Spec.TopologyLink.Name,
		AdminState:  x.Spec.TopologyLink.AdminState,
		Description: x.Spec.TopologyLink.Description,
		Endpoints:   eps,
		Tag:         tags,
		State: &NddrTopologyTopologyLinkState{
			Status: utils.StringPtr(""),
			Reason: utils.StringPtr(""),
			Node:   make([]*NddrTopologyTopologyLinkStateNode, 0),
			Tag:    make([]*nddov1.Tag, 0),
		},
	}
	return nil
}

func (x *TopologyLink) SetStatus(s string) {
	x.Status.TopologyLink.State.Status = &s
}

func (x *TopologyLink) SetReason(s string) {
	x.Status.TopologyLink.State.Reason = &s
}

func (x *TopologyLink) SetNodeEndpoint(nodeName string, ep *NddrTopologyTopologyLinkStateNodeEndpoint) {
	for _, node := range x.Status.TopologyLink.State.Node {
		if *node.Name == nodeName {
			for _, nodeep := range node.Endpoint {
				if *nodeep.Name == *ep.Name {
					nodeep.Lag = ep.Lag
					nodeep.LagMemberLink = ep.LagMemberLink
					nodeep.Name = ep.Name
					return
				}
			}
			node.Endpoint = append(node.Endpoint, ep)
			return
		}
	}
	// if we come here we need to create the node
	x.Status.TopologyLink.State.Node = append(x.Status.TopologyLink.State.Node, &NddrTopologyTopologyLinkStateNode{
		Name: &nodeName,
		Endpoint: []*NddrTopologyTopologyLinkStateNodeEndpoint{
			ep,
		},
	})
}

func (x *TopologyLink) GetNodeEndpoints() []*NddrTopologyTopologyLinkStateNode {
	if x.Status.TopologyLink != nil && x.Status.TopologyLink.State != nil && x.Status.TopologyLink.State.Node != nil {
		return x.Status.TopologyLink.State.Node
	}
	return make([]*NddrTopologyTopologyLinkStateNode, 0)
}

func (x *TopologyLink) SetKind(s string) {
	for _, tag := range x.Status.TopologyLink.State.Tag {
		if *tag.Key == KeyLinkKind {
			tag.Value = &s
			return
		}
	}
	x.Status.TopologyLink.State.Tag = append(x.Status.TopologyLink.State.Tag, &nddov1.Tag{
		Key:   utils.StringPtr(KeyLinkKind),
		Value: &s,
	})
}

func (x *TopologyLink) GetKind() string {
	if x.Status.TopologyLink != nil && x.Status.TopologyLink.State != nil && x.Status.TopologyLink.State.Tag != nil {
		for _, tag := range x.Status.TopologyLink.State.Tag {
			if *tag.Key == KeyLinkKind {
				return *tag.Value
			}
		}

	}
	return LinkEPKindUnknown.String()
}

func (x *TopologyLink) AddEndPointATag(key string, value string) {
	for _, tag := range x.Spec.TopologyLink.Endpoints[0].Tag {
		if *tag.Key == key {
			tag.Value = &value
			return
		}
	}
	// if not found append
	x.Spec.TopologyLink.Endpoints[0].Tag = append(x.Spec.TopologyLink.Endpoints[0].Tag,
		&nddov1.Tag{
			Key:   &key,
			Value: &value,
		})
}

func (x *TopologyLink) AddEndPointBTag(key string, value string) {
	for _, tag := range x.Spec.TopologyLink.Endpoints[1].Tag {
		if *tag.Key == key {
			tag.Value = &value
			return
		}
	}
	// if not found append
	x.Spec.TopologyLink.Endpoints[1].Tag = append(x.Spec.TopologyLink.Endpoints[1].Tag,
		&nddov1.Tag{
			Key:   &key,
			Value: &value,
		})
}

func (x *TopologyLink) DeleteEndPointATag(key string, value string) {
	found := false
	var idx int
	for i, tag := range x.Spec.TopologyLink.Endpoints[0].Tag {
		if *tag.Key == key && *tag.Value == value {
			idx = i
			found = true
		}
	}
	if found {
		x.Spec.TopologyLink.Endpoints[0].Tag = append(x.Spec.TopologyLink.Endpoints[0].Tag[:idx], x.Spec.TopologyLink.Endpoints[0].Tag[idx+1:]...)
	}
}

func (x *TopologyLink) DeleteEndPointBTag(key string, value string) {
	found := false
	var idx int
	for i, tag := range x.Spec.TopologyLink.Endpoints[1].Tag {
		if *tag.Key == key && *tag.Value == value {
			idx = i
			found = true
		}
	}
	if found {
		x.Spec.TopologyLink.Endpoints[1].Tag = append(x.Spec.TopologyLink.Endpoints[1].Tag[:idx], x.Spec.TopologyLink.Endpoints[1].Tag[idx+1:]...)
	}
}

func (x *TopologyLink) SetOrganization(s string) {
	x.Status.SetOrganization(s)
}

func (x *TopologyLink) SetDeployment(s string) {
	x.Status.SetDeployment(s)
}

func (x *TopologyLink) SetAvailabilityZone(s string) {
	x.Status.SetAvailabilityZone(s)
}

func (x *TopologyLink) SetTopologyName(s string) {
	x.Status.TopologyName = &s
}
