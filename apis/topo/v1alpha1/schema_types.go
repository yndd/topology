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
	nddov1 "github.com/yndd/nddo-runtime/apis/common/v1"
)

// NddrTopology struct
type NddrTopology struct {
	Topology []*NddrTopologyTopology `json:"topology,omitempty"`
}

// NddrTopologyTopology struct
type NddrTopologyTopology struct {
	AdminState  *string                       `json:"admin-state,omitempty"`
	Defaults    *NddrTopologyTopologyDefaults `json:"defaults,omitempty"`
	Description *string                       `json:"description,omitempty"`
	Kind        []*NddrTopologyTopologyKind   `json:"kind,omitempty"`
	Link        []*NddrTopologyTopologyLink   `json:"link,omitempty"`
	Name        *string                       `json:"name,omitempty"`
	Node        []*NddrTopologyTopologyNode   `json:"node,omitempty"`
	State       *NddrTopologyTopologyState    `json:"state,omitempty"`
}

// NddrTopologyTopologyDefaults struct
type NddrTopologyTopologyDefaults struct {
	Tag []*nddov1.Tag `json:"tag,omitempty"`
}

// NddrTopologyTopologyKind struct
type NddrTopologyTopologyKind struct {
	Name *string       `json:"name"`
	Tag  []*nddov1.Tag `json:"tag,omitempty"`
}

// NddrTopologyTopologyLink struct
type NddrTopologyTopologyLink struct {
	AdminState  *string                              `json:"admin-state,omitempty"`
	Description *string                              `json:"description,omitempty"`
	Endpoints   []*NddrTopologyTopologyLinkEndpoints `json:"endpoints,omitempty"`
	Name        *string                              `json:"name,omitempty"`
	State       *NddrTopologyTopologyLinkState       `json:"state,omitempty"`
	Tag         []*nddov1.Tag                        `json:"tag,omitempty"`
}

// NddrTopologyTopologyLinkEndpoints struct
type NddrTopologyTopologyLinkEndpoints struct {
	InterfaceName *string       `json:"interface-name"`
	NodeName      *string       `json:"node-name"`
	Tag           []*nddov1.Tag `json:"tag,omitempty"`
}

// NddrTopologyTopologyLinkState struct
type NddrTopologyTopologyLinkState struct {
	Reason *string                              `json:"reason,omitempty"`
	Status *string                              `json:"status,omitempty"`
	Node   []*NddrTopologyTopologyLinkStateNode `json:"node,omitempty"`
	Tag    []*nddov1.Tag                        `json:"tag,omitempty"`
}

// NddrTopologyTopologyLinkStateNode struct
type NddrTopologyTopologyLinkStateNode struct {
	Name     *string                                      `json:"name,omitempty"`
	Endpoint []*NddrTopologyTopologyLinkStateNodeEndpoint `json:"endpoint,omitempty"`
}

// NddrTopologyTopologyLinkStateNodeEndpoint struct
type NddrTopologyTopologyLinkStateNodeEndpoint struct {
	Lag           *bool   `json:"lag,omitempty"`
	LagMemberLink *bool   `json:"lag-member-link,omitempty"`
	Name          *string `json:"name,omitempty"`
}

// NddrTopologyTopologyNode struct
type NddrTopologyTopologyNode struct {
	AdminState  *string                        `json:"admin-state,omitempty"`
	Description *string                        `json:"description,omitempty"`
	KindName    *string                        `json:"kind-name,omitempty"`
	Name        *string                        `json:"name,omitempty"`
	State       *NddrTopologyTopologyNodeState `json:"state,omitempty"`
	Tag         []*nddov1.Tag                  `json:"tag,omitempty"`
}

// NddrTopologyTopologyNodeState struct
type NddrTopologyTopologyNodeState struct {
	Endpoint   []*NddrTopologyTopologyNodeStateEndpoint `json:"endpoint,omitempty"`
	LastUpdate *string                                  `json:"last-update,omitempty"`
	Reason     *string                                  `json:"reason,omitempty"`
	Status     *string                                  `json:"status,omitempty"`
	Tag        []*nddov1.Tag                            `json:"tag,omitempty"`
}

// NddrTopologyTopologyNodeStateEndpoint struct
type NddrTopologyTopologyNodeStateEndpoint struct {
	Lag        *bool   `json:"lag,omitempty"`
	LagSubLink *bool   `json:"lag-sub-link,omitempty"`
	Name       *string `json:"name"`
}

// NddrTopologyTopologyState struct
type NddrTopologyTopologyState struct {
	Reason *string       `json:"reason,omitempty"`
	Status *string       `json:"status,omitempty"`
	Tag    []*nddov1.Tag `json:"tag,omitempty"`
}

// Root is the root of the schema
type Root struct {
	TopoNddrTopology *NddrTopology `json:"Nddr-topology,omitempty"`
}
