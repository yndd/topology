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

// TopologyDefinitionSpec struct
type TopologySpec struct {
	nddv1.ResourceSpec `json:",inline"`
	// Properties define the properties of the Topology
	Properties *TopologyProperties `json:"properties,omitempty"`
}

// A TopologyStatus represents the observed state of a Topology.
type TopologyStatus struct {
	nddv1.ResourceStatus `json:",inline"`
}

// TopologyProperties struct
type TopologyProperties struct {
	Defaults       *TopologyDefaults `json:"defaults,omitempty"`
	VendorTypeInfo []*NodeProperties `json:"vendorTypeInfo,omitempty"`
}

// TopologySpecDefaults struct
type TopologyDefaults struct {
	NodeProperties *NodeProperties   `json:",inline"`
	Tag            map[string]string `json:"tag,omitempty"`
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
// +kubebuilder:resource:categories={yndd,topo}
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
	TopologyKind             = reflect.TypeOf(Topology{}).Name()
	TopologyGroupKind        = schema.GroupKind{Group: Group, Kind: TopologyKind}.String()
	TopologyKindAPIVersion   = TopologyKind + "." + GroupVersion.String()
	TopologyGroupVersionKind = GroupVersion.WithKind(TopologyKind)
)
