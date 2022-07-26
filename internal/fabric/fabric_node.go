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

	"github.com/yndd/ndd-runtime/pkg/logging"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
)

// +k8s:deepcopy-gen=false
type FabricNode interface {
	GetNodeName() string
	GetPosition() topov1alpha1.Position
	GetNodeIndex() uint32
	GetPodIndex() uint32
	GetInterfaceName(idx uint32) string
	GetInterfaceNameWithPlatfromOffset(idx uint32) string
	GetVendorType() targetv1.VendorType
	GetPlatform() string
}

func NewLeafFabricNode(podIndex, nodeIndex uint32, vendorInfo *topov1alpha1.FabricTierVendorInfo, log logging.Logger) FabricNode {
	return &fabricNode{
		log:        log,
		position:   topov1alpha1.PositionLeaf,
		podIndex:   podIndex,
		nodeIndex:  nodeIndex,
		vendorInfo: vendorInfo,
	}
}

func NewSpineFabricNode(podIndex, nodeIndex uint32, vendorInfo *topov1alpha1.FabricTierVendorInfo, log logging.Logger) FabricNode {
	return &fabricNode{
		log:        log,
		position:   topov1alpha1.PositionSpine,
		podIndex:   podIndex,
		nodeIndex:  nodeIndex,
		vendorInfo: vendorInfo,
	}
}

func NewSuperspineFabricNode(nodeIndex uint32, vendorInfo *topov1alpha1.FabricTierVendorInfo, log logging.Logger) FabricNode {
	return &fabricNode{
		log:        log,
		position:   topov1alpha1.PositionSuperspine,
		nodeIndex:  nodeIndex,
		vendorInfo: vendorInfo,
	}
}

// +k8s:deepcopy-gen=false
type fabricNode struct {
	log        logging.Logger
	position   topov1alpha1.Position
	nodeIndex  uint32 // relative number within the position, pod
	podIndex   uint32 // pod index
	vendorInfo *topov1alpha1.FabricTierVendorInfo
}

func (n *fabricNode) GetInterfaceName(idx uint32) string {
	return fmt.Sprintf("int-1/%d", idx)
}

func (n *fabricNode) GetInterfaceNameWithPlatfromOffset(idx uint32) string {
	n.log.Debug("GetInterfaceNameWithPlatformOffset",
		"idx", idx,
		"nodeName", n.GetNodeName(),
		"podIndex", n.GetPodIndex(),
		"vendorType", n.GetVendorType(),
		"platform", n.GetPlatform(),
		"position", n.GetPosition(),
	)

	n.log.Debug("GetInterfaceNameWithPlatformOffset",
		"vendorType", n.GetVendorType(),
		"vendorType", targetv1.VendorTypeNokiaSRL,
	)
	var actualIndex uint32
	switch n.GetVendorType() {
	case targetv1.VendorTypeNokiaSRL:
		n.log.Debug("GetInterfaceNameWithPlatformOffset", "vendorType", targetv1.VendorTypeNokiaSRL)
		switch n.GetPosition() {
		case topov1alpha1.PositionLeaf:
			n.log.Debug("GetInterfaceNameWithPlatformOffset", "position", targetv1.VendorTypeNokiaSRL)
			switch n.GetPlatform() {
			case "IXR-D3":
				n.log.Debug("GetInterfaceNameWithPlatformOffset", "platform", "IXR-D3")
				actualIndex = idx + 26
			case "IXR-D2":
				n.log.Debug("GetInterfaceNameWithPlatformOffset", "platform", "IXR-D2")
				actualIndex = idx + 48
			}
		case topov1alpha1.PositionSpine:
			switch n.GetPlatform() {
			case "IXR-D3":
				n.log.Debug("GetInterfaceNameWithPlatformOffset", "platform", "IXR-D3")
				actualIndex = idx + 24
			}
		}
	case targetv1.VendorTypeNokiaSROS:
		// TODO
	}
	n.log.Debug("GetInterfaceNameWithPlatformOffset",
		"actualIndex", actualIndex,
		"nodeName", n.GetNodeName(),
		"podIndex", n.GetPodIndex(),
		"vendorType", n.GetVendorType(),
		"platform", n.GetPlatform(),
		"position", n.GetPosition(),
	)
	return fmt.Sprintf("int-1/%d", actualIndex)
}

func (n *fabricNode) GetPosition() topov1alpha1.Position {
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
	if n.GetPosition() != topov1alpha1.PositionSuperspine {
		return fmt.Sprintf("pod%d-%s%d", n.podIndex, n.position, n.nodeIndex)

	} else {
		return fmt.Sprintf("%s%d", n.position, n.nodeIndex)
	}
}
