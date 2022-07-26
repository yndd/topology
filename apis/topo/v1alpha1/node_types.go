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
)

// NodeSpec struct
type NodeSpec struct {
	nddv1.ResourceSpec `json:",inline"`
	// Properties define the properties of the Topology
	Properties *NodeProperties `json:"properties,omitempty"`
}

// A NodeStatus represents the observed state of a node.
type NodeStatus struct {
	nddv1.ResourceStatus `json:",inline"`
}

// NodeProperties struct
type NodeProperties struct {
	VendorType        targetv1.VendorType `json:"vendorType,omitempty"`
	Platform          string              `json:"platform,omitempty"`
	Position          Position            `json:"position,omitempty"`
	MacAddress        string              `json:"macAddress,omitempty"`
	SerialNumber      string              `json:"serialNumber,omitempty"`
	ExpectedSWVersion string              `json:"expectedSwVersion,omitempty"`
	MgmtIPAddress     string              `json:"mgmtIPAddress,omitempty"`
	Tag               map[string]string   `json:"tag,omitempty"`
}

type Position string

// Position enums.
const (
	PositionUnknown    Position = "unknown"
	PositionLeaf       Position = "leaf"
	PositionSpine      Position = "spine"
	PositionSuperspine Position = "superspine"
	PositionDcgw       Position = "dcgw"
	PositionWan        Position = "wan"
	PositionCpe        Position = "cpe"
	PositionServer     Position = "server"
	PositionInfra      Position = "infra"
)

// +kubebuilder:object:root=true

// Node is the Schema for the Node API
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="SYNC",type="string",JSONPath=".status.conditions[?(@.kind=='Synced')].status"
// +kubebuilder:printcolumn:name="STATUS",type="string",JSONPath=".status.conditions[?(@.kind=='Ready')].status"
// +kubebuilder:printcolumn:name="ORG",type="string",JSONPath=".status.oda.organization"
// +kubebuilder:printcolumn:name="DEP",type="string",JSONPath=".status.oda.deployment"
// +kubebuilder:printcolumn:name="AZ",type="string",JSONPath=".status.oda.availabilityZone"
// +kubebuilder:printcolumn:name="TOPO",type="string",JSONPath=".status.oda.resourceName"
// +kubebuilder:printcolumn:name="KIND",type="string",JSONPath=".spec.properties.vendorType"
// +kubebuilder:printcolumn:name="PLATFORM",type="string",JSONPath="..spec.properties.platform"
// +kubebuilder:printcolumn:name="POSITION",type="string",JSONPath="..spec.properties.position"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:categories={yndd,topo}
type Node struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeSpec   `json:"spec,omitempty"`
	Status NodeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NodeList contains a list of Nodes
type NodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Node `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Node{}, &NodeList{})
}

// Node type metadata.
var (
	NodeKind             = reflect.TypeOf(Node{}).Name()
	NodeGroupKind        = schema.GroupKind{Group: Group, Kind: NodeKind}.String()
	NodeKindAPIVersion   = NodeKind + "." + GroupVersion.String()
	NodeGroupVersionKind = GroupVersion.WithKind(NodeKind)
)
