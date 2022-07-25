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
package v1alpha1

import (
	"fmt"

	targetv1 "github.com/yndd/target/apis/target/v1"
)

// +k8s:deepcopy-gen=false
type FabricNode interface {
	GetNodeName() string
	GetPosition() Position
	GetNodeIndex() uint32
	GetPodIndex() uint32
	GetInterfaceName(idx uint32) string
	GetInterfaceNameWithPlatformOffset(idx uint32) string
	GetVendorType() targetv1.VendorType
	GetPlatform() string
}

func NewLeafFabricNode(podIndex, nodeIndex uint32, vendorInfo *FabricTierVendorInfo) FabricNode {
	return &fabricNode{
		position:   PositionLeaf,
		podIndex:   podIndex,
		nodeIndex:  nodeIndex,
		vendorInfo: vendorInfo,
	}
}

func NewSpineFabricNode(podIndex, nodeIndex uint32, vendorInfo *FabricTierVendorInfo) FabricNode {
	return &fabricNode{
		position:   PositionSpine,
		podIndex:   podIndex,
		nodeIndex:  nodeIndex,
		vendorInfo: vendorInfo,
	}
}

func NewSuperspineFabricNode(nodeIndex uint32, vendorInfo *FabricTierVendorInfo) FabricNode {
	return &fabricNode{
		position:   PositionSuperspine,
		nodeIndex:  nodeIndex,
		vendorInfo: vendorInfo,
	}
}

// +k8s:deepcopy-gen=false
type fabricNode struct {
	position   Position
	nodeIndex  uint32 // relative number within the position, pod
	podIndex   uint32 // pod index
	vendorInfo *FabricTierVendorInfo
}

func (n *fabricNode) GetInterfaceName(idx uint32) string {
	return fmt.Sprintf("int-1/%d", idx)
}

func (n *fabricNode) GetInterfaceNameWithPlatformOffset(idx uint32) string {
	var actualIndex uint32
	switch n.vendorInfo.VendorType {
	case targetv1.VendorTypeNokiaSRL:
		switch n.position {
		case PositionLeaf:
			switch n.vendorInfo.Platform {
			case "IXR-D3":
				actualIndex = idx + 26
			case "IXR-D2":
				actualIndex = idx + 48
			}
		case PositionSpine:
			switch n.vendorInfo.Platform {
			case "IXR-D3":
				actualIndex = idx + 24
			}
		default:
			// not expected
		}
	case targetv1.VendorTypeNokiaSROS:
		// TODO
	default:

	}
	return fmt.Sprintf("int-1/%d", actualIndex)
}

func (n *fabricNode) GetPosition() Position {
	return n.position
}

func (n *fabricNode) GetNodeIndex() uint32 {
	return n.nodeIndex
}

func (n *fabricNode) GetPodIndex() uint32 {
	return n.podIndex
}

func (n *fabricNode) GetVendorType() targetv1.VendorType {
	return n.vendorInfo.VendorType
}

func (n *fabricNode) GetPlatform() string {
	return n.vendorInfo.Platform
}

func (n *fabricNode) GetNodeName() string {
	if n.GetPosition() != PositionSuperspine {
		return fmt.Sprintf("pod%d-%s%d", n.podIndex, n.position, n.nodeIndex)

	} else {
		return fmt.Sprintf("%s%d", n.position, n.nodeIndex)
	}
}
