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
	targetv1 "github.com/yndd/target/apis/target/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	//targetv1alpha1pb "github.com/yndd/topology/gen/go/apis/topo/v1alpha1"
)

type SupportServers struct {
	DnsServers []*string `json:"dnsServers,omitempty"`
	NtPServers []*string `json:"netServers,omitempty"`
}

type TemplateSubnet struct {
	IPSubnet       string `json:"ipSubnet,omitempty"`
	SupportServers `json:"inline"`
}

type Fabric struct {
	// superspine
	Tier1 *FabricTier  `json:"tier1,omitempty"`
	Pods  []*FabricPod `json:"pod,omitempty"`
}

type FabricPod struct {
	PodNumber uint32        `json:"num,omitempty"`
	Tiers     []*FabricTier `json:"tiers,omitempty"`
}

type FabricTier struct {
	// tier3, tier2
	Kind string `json:"kind,omitempty"`
	// list to support multiple vendors in a tier - typically criss-cross
	VendorInfo []*FabricTierVendorInfo `json:"vendorInfo,omitempty"`
	// number of nodes in the tier
	NodeNumber uint32 `json:"num,omitempty"`
	// oversubscription ratio
	Oversubscription string `json:"oversubscription,omitempty"`
}

type FabricTierVendorInfo struct {
	Platform   string              `json:"platform,omitempty"`
	VendorType targetv1.VendorType `json:"vendorType,omitempty"`
}

// TemplateProperties define the properties of the Template
type TemplateProperties struct {
	SupportServers `json:"inline"`
	Subnet         *TemplateSubnet `json:"subnet,omitempty"`
	Fabric         *Fabric         `json:"fabric,omitempty"`
}

// TemplateSpec struct
type TemplateSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	// Properties define the properties of the Template
	Properties TemplateProperties `json:"properties,omitempty"`
}

// A TemplateStatus represents the observed state of a Template.
type TemplateStatus struct {
	nddv1.ResourceStatus `json:",inline"`
	//TopologyName            *string               `json:"topology-name,omitempty"`
	//Topology                *NddrTopologyTopology `json:"topology,omitempty"`
}

// +kubebuilder:object:root=true

// Template is the Schema for the Template API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda[?(@.key=='organization')].value"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda[?(@.key=='deployment')].value"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda[?(@.key=='availability-zone')].value"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.topology-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type Template struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateSpec   `json:"spec,omitempty"`
	Status TemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TemplateList contains a list of Templates
type TemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Template `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Template{}, &TemplateList{})
}

// Template type metadata.
var (
	TemplateKind             = reflect.TypeOf(Template{}).Name()
	TemplateGroupKind        = schema.GroupKind{Group: Group, Kind: TemplateKind}.String()
	TemplateKindAPIVersion   = TemplateKind + "." + GroupVersion.String()
	TemplateGroupVersionKind = GroupVersion.WithKind(TemplateKind)
)
