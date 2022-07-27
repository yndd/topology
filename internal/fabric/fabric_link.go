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

// Package v1alpha1 contains API Schema definitions for the topo v1alpha1 API group
package fabric

import (
	"fmt"
	"strings"
)

// +k8s:deepcopy-gen=false
type FabricLink interface {
	GetName() string
	GetEndpointA() *Endpoint
	GetEndpointB() *Endpoint
}

func NewFabricLink(epA *Endpoint, epB *Endpoint) FabricLink {
	linkName := fmt.Sprintf("%s-%s-%s-%s", epA.Node.GetNodeName(), epA.IfName, epB.Node.GetNodeName(), epB.IfName)
	
	return &fabricLink{
		name: linkName,
		epA:  epA,
		epB:  epB,
	}
}

// +k8s:deepcopy-gen=false
type fabricLink struct {
	name string
	epA  *Endpoint
	epB  *Endpoint
}

func (n *fabricLink) AddInterfaceName(idx uint32) {
}

// +k8s:deepcopy-gen=false
type Endpoint struct {
	Node   FabricNode
	IfName string
}

func (l *fabricLink) GetName() string {
	return strings.ReplaceAll(l.name, "/", "-")
}

func (l *fabricLink) GetEndpointA() *Endpoint {
	return l.epA
}

func (l *fabricLink) GetEndpointB() *Endpoint {
	return l.epB
}
