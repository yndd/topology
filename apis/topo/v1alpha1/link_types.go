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
	//targetv1alpha1pb "github.com/yndd/topology/gen/go/apis/topo/v1alpha1"
)

// TopologyLinkSpec struct
type LinkSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	// Properties define the properties of the Topology
	Properties *LinkProperties `json:"properties,omitempty"`
}

// A LinkStatus represents the observed state of a Link.
type LinkStatus struct {
	nddv1.ResourceStatus `json:",inline"`
}

// LinkProperties struct
type LinkProperties struct {
	Endpoints []*Endpoints       `json:"endpoints,omitempty"`
	LagMember bool               `json:"lagMember,omitempty"`
	Lacp      bool               `json:"lacp,omitempty"`
	Lag       bool               `json:"lag,omitempty"`
	Kind      LinkKindProperties `json:"kind,omitempty"`
	Tag       map[string]string  `json:"tag,omitempty"`
}

// LinkEndpoints struct
type Endpoints struct {
	// kubebuilder:validation:MinLength=3
	// kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`int-([1-9](\d){0,1}(/[abcd])?(/[1-9](\d){0,1})?/(([1-9](\d){0,1})|(1[0-1]\d)|(12[0-8])))|`
	InterfaceName   string                 `json:"interfaceName"`
	NodeName        string                 `json:"nodeName"`
	Kind            EndpointKindProperties `json:"kind,omitempty"`
	LacpFallback    bool                   `json:"lacpFallback,omitempty"`
	LagName         string                 `json:"lagName,omitempty"`
	EndpointGroup   string                 `json:"endpointGroup,omitempty"`
	MultiHoming     bool                   `json:"multiHoming,omitempty"`
	MultiHomingName string                 `json:"multiHomingName,omitempty"`
	Tag             map[string]string      `json:"tag,omitempty"`
}

type LinkKindProperties string

// LinkPropertiesKind enums.
const (
	LinkKindUnknown LinkKindProperties = "unknown"
	LinkKindInfra   LinkKindProperties = "infra"
	LinkKindLoop    LinkKindProperties = "loop"
)

type EndpointKindProperties string

// LinkPropertiesKind enums.
const (
	EndpointKindUnknown  EndpointKindProperties = "unknown"
	EndpointKindInfra    EndpointKindProperties = "infra"
	EndpointKindLoop     EndpointKindProperties = "loop"
	EndpointKindExternal EndpointKindProperties = "external"
	EndpointKindOob      EndpointKindProperties = "oob"
)

// +kubebuilder:object:root=true

// Link is the Schema for the Link API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda.organization"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda.deployment"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda.availabilityZone"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.oda.resourceName"
// +kubebuilder:printcolumn:name="LAG",type="boolean",JSONPath=".spec.properties.lag"
// +kubebuilder:printcolumn:name="MEMBER",type="boolean",JSONPath=".spec.properties.lagmember"
// +kubebuilder:printcolumn:name="NODE-EPA",type="string",JSONPath=".spec.properties.endpoints[0].nodeName"
// +kubebuilder:printcolumn:name="ITFCE-EPA",type="string",JSONPath=".spec.properties.endpoints[0].interfaceName"
// +kubebuilder:printcolumn:name="MH-EPA",type="string",JSONPath=".spec.properties.endpoints[0].multiHomingName"
// +kubebuilder:printcolumn:name="NODE-EPB",type="string",JSONPath=".spec.properties.endpoints[1].nodeName"
// +kubebuilder:printcolumn:name="ITFCE-EPB",type="string",JSONPath=".spec.properties.endpoints[1].interfaceName"
// +kubebuilder:printcolumn:name="MH-EPB",type="string",JSONPath=".spec.properties.endpoints[1].multiHomingName"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type Link struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LinkSpec   `json:"spec,omitempty"`
	Status LinkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LinkList contains a list of Links
type LinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Link `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Link{}, &LinkList{})
}

// Link type metadata.
var (
	LinkKind             = reflect.TypeOf(Link{}).Name()
	LinkGroupKind        = schema.GroupKind{Group: Group, Kind: LinkKind}.String()
	LinkKindAPIVersion   = LinkKind + "." + GroupVersion.String()
	LinkGroupVersionKind = GroupVersion.WithKind(LinkKind)
)
