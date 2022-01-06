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
	nddov1 "github.com/yndd/nddo-runtime/apis/common/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Topology struct
type TopoTopology struct {
	// +kubebuilder:validation:Enum=`disable`;`enable`
	// +kubebuilder:default:="enable"
	AdminState *string               `json:"admin-state,omitempty"`
	Defaults   *TopoTopologyDefaults `json:"defaults,omitempty"`
	// kubebuilder:validation:MinLength=1
	// kubebuilder:validation:MaxLength=255
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="[A-Za-z0-9 !@#$^&()|+=`~.,'/_:;?-]*"
	Description *string             `json:"description,omitempty"`
	Kind        []*TopoTopologyKind `json:"kind,omitempty"`
	Name        *string             `json:"name,omitempty"`
}

// TopoTopologyDefaults struct
type TopoTopologyDefaults struct {
	Tag []*nddov1.Tag `json:"tag,omitempty"`
}

// TopologyKind struct
type TopoTopologyKind struct {
	Name *string       `json:"name"`
	Tag  []*nddov1.Tag `json:"tag,omitempty"`
}

// A TopologySpec defines the desired state of a Topology.
type TopologySpec struct {
	//nddov1.OdaInfo `json:",inline"`
	Topology *TopoTopology `json:"topology,omitempty"`
}

// A TopologyStatus represents the observed state of a Topology.
type TopologyStatus struct {
	nddv1.ConditionedStatus `json:",inline"`
	nddov1.OdaInfo          `json:",inline"`
	TopologyName            *string               `json:"topology-name,omitempty"`
	Topology                *NddrTopologyTopology `json:"topology,omitempty"`
}

// +kubebuilder:object:root=true

// Topology is the Schema for the Topology API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda[?(@.key=='organization')].value"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda[?(@.key=='deployment')].value"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda[?(@.key=='availability-zone')].value"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.topology-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
type Topology struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologySpec   `json:"spec,omitempty"`
	Status TopologyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TopologyList contains a list of Topologies
type TopologyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Topology `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Topology{}, &TopologyList{})
}

// Topology type metadata.
var (
	TopologyKindKind         = reflect.TypeOf(Topology{}).Name()
	TopologyGroupKind        = schema.GroupKind{Group: Group, Kind: TopologyKindKind}.String()
	TopologyKindAPIVersion   = TopologyKindKind + "." + GroupVersion.String()
	TopologyGroupVersionKind = GroupVersion.WithKind(TopologyKindKind)
)
