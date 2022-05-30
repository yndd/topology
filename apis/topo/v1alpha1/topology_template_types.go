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

// TopologyTemplateProperties define the properties of the TopologyTemplate
type TopologyTemplateProperties struct {
}

// TopologyTemplateSpec struct
type TopologyTemplateSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	// Properties define the properties of the TopologyTemplate
	Properties TopologyTemplateProperties `json:"properties,omitempty"`
}

// A TopologyTemplateStatus represents the observed state of a TopologyTemplate.
type TopologyTemplateStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	//TopologyName            *string               `json:"topology-name,omitempty"`
	//Topology                *NddrTopologyTopology `json:"topology,omitempty"`
}

// +kubebuilder:object:root=true

// TopologyTemplate is the Schema for the TopologyTemplate API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda[?(@.key=='organization')].value"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda[?(@.key=='deployment')].value"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda[?(@.key=='availability-zone')].value"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.topology-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type TopologyTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopologyTemplateSpec   `json:"spec,omitempty"`
	Status TopologyTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TopologyTemplateList contains a list of TopologyTemplates
type TopologyTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TopologyTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TopologyTemplate{}, &TopologyTemplateList{})
}

// TopologyTemplate type metadata.
var (
	TopologyTemplateKind             = reflect.TypeOf(TopologyTemplate{}).Name()
	TopologyTemplateGroupKind        = schema.GroupKind{Group: Group, Kind: TopologyTemplateKind}.String()
	TopologyTemplateKindAPIVersion   = TopologyTemplateKind + "." + GroupVersion.String()
	TopologyTemplateGroupVersionKind = GroupVersion.WithKind(TopologyTemplateKind)
)
