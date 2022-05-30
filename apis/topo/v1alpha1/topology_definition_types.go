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

// TopologyDefinitionProperties define the properties of the TopologyDefinition
type TopologyDefinitionProperties struct {
	Templates      []*Template      `json:"templates,omitempty"`
	DiscoveryRules []*DiscoveryRule `json:"discovery-rules,omitempty"`
}

type Template struct {
	Rule `json:",inline"`
}

type DiscoveryRule struct {
	Rule `json:",inline"`
}

type Rule struct {
	Name string `json:"name"`
	// +kubebuilder:default=false
	DigitalTwin bool `json:"digital-twin,omitempty"`
}

// TopologyDefinitionSpec struct
type TopologyDefinitionSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	// Properties define the properties of the TopologyDefinition
	Properties TopologyDefinitionProperties `json:"properties,omitempty"`
}

// A TopologyDefinitionStatus represents the observed state of a TopologyDefinition.
type TopologyDefinitionStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	//TopologyName            *string               `json:"topology-name,omitempty"`
	//Topology                *NddrTopologyTopology `json:"topology,omitempty"`
}

// +kubebuilder:object:root=true

// TopologyDefinition is the Schema for the Topology API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda[?(@.key=='organization')].value"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda[?(@.key=='deployment')].value"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda[?(@.key=='availability-zone')].value"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.topology-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type TopologyDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologyDefinitionSpec   `json:"spec,omitempty"`
	Status TopologyDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TopologyDefinitionList contains a list of TopologyDefinitions
type TopologyDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopologyDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopologyDefinition{}, &TopologyDefinitionList{})
}

// TopologyDefinition type metadata.
var (
	TopologyDefinitionKind             = reflect.TypeOf(TopologyDefinition{}).Name()
	TopologyDefinitionGroupKind        = schema.GroupKind{Group: Group, Kind: TopologyDefinitionKind}.String()
	TopologyDefinitionKindAPIVersion   = TopologyDefinitionKind + "." + GroupVersion.String()
	TopologyDefinitionGroupVersionKind = GroupVersion.WithKind(TopologyDefinitionKind)
)
