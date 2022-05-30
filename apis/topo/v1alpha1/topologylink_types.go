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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TopologyLinkProperties struct
type TopologyLinkProperties struct {
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState *string `json:"admin-state,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Description *string                  `json:"description,omitempty"`
	Endpoints   []*TopologyLinkEndpoints `json:"endpoints,omitempty"`
	Name        *string                  `json:"name,omitempty"`
	Tag         []*nddv1.Tag             `json:"tag,omitempty"`
}

// TopologyLinkEndpoints struct
type TopologyLinkEndpoints struct {
	// kubebuilder:validation:MinLength=3
	// kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`int-([1-9](\d){0,1}(/[abcd])?(/[1-9](\d){0,1})?/(([1-9](\d){0,1})|(1[0-1]\d)|(12[0-8])))|`
	InterfaceName *string      `json:"interface-name"`
	NodeName      *string      `json:"node-name"`
	Tag           []*nddv1.Tag `json:"tag,omitempty"`
}

// TopologyLinkSpec struct
type TopologyLinkSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	// Properties define the properties of the Topology
	Properties TopologyLinkProperties `json:"properties,omitempty"`
}

// A TopologyLinkStatus represents the observed state of a TopologyLink.
type TopologyLinkStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	//TopologyName            *string               `json:"topology-name,omitempty"`
	//Topology                *NddrTopologyTopology `json:"topology,omitempty"`
}

// +kubebuilder:object:root=true

// TopoTopologyLink is the Schema for the TopologyLink API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda[?(@.key=='organization')].value"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda[?(@.key=='deployment')].value"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda[?(@.key=='availability-zone')].value"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.topology-name"
// +kubebuilder:printcolumn:name="LAG",type="string",JSONPath=".spec.link.tag[?(@.key=='lag')].value"
// +kubebuilder:printcolumn:name="MEMBER",type="string",JSONPath=".spec.link.tag[?(@.key=='lag-member')].value"
// +kubebuilder:printcolumn:name="NODE-EPA",type="string",JSONPath=".spec.link.endpoints[0].node-name"
// +kubebuilder:printcolumn:name="ITFCE-EPA",type="string",JSONPath=".spec.link.endpoints[0].interface-name"
// +kubebuilder:printcolumn:name="MH-EPA",type="string",JSONPath=".spec.link.endpoints[0].tag[?(@.key=='multihoming')].value"
// +kubebuilder:printcolumn:name="NODE-EPB",type="string",JSONPath=".spec.link.endpoints[1].node-name"
// +kubebuilder:printcolumn:name="ITFCE-EPB",type="string",JSONPath=".spec.link.endpoints[1].interface-name"
// +kubebuilder:printcolumn:name="MH-EPB",type="string",JSONPath=".spec.link.endpoints[1].tag[?(@.key=='multihoming')].value"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type TopologyLink struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologyLinkSpec   `json:"spec,omitempty"`
	Status TopologyLinkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TopoTopologyLinkList contains a list of TopologyLinks
type TopologyLinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopologyLink `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopologyLink{}, &TopologyLinkList{})
}

// TopologyLink type metadata.
var (
	TopologyLinkKind             = reflect.TypeOf(TopologyLink{}).Name()
	TopologyLinkGroupKind        = schema.GroupKind{Group: Group, Kind: TopologyLinkKind}.String()
	TopologyLinkKindAPIVersion   = TopologyLinkKind + "." + GroupVersion.String()
	TopologyLinkGroupVersionKind = GroupVersion.WithKind(TopologyLinkKind)
)
