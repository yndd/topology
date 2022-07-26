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
	"sync"

	"github.com/yndd/ndd-runtime/pkg/logging"
)

// +k8s:deepcopy-gen=false
type Fabric interface {
	GetFabricNodes() []FabricNode
	GetFabricLinks() []FabricLink
	PrintNodes()
	PrintLinks()
}

func NewFabric(namespaceName string, template *FabricTemplate, log logging.Logger) (Fabric, error) {
	f := &fabric{
		log:             log,
		tier1Nodes:      make([]FabricNode, 0),
		pods:            map[uint32]*podInfo{},
		tier2tier3Links: make([]FabricLink, 0),
		tier1tier2Links: make([]FabricLink, 0),
	}

	// process superspine nodes
	for n := uint32(0); n < template.Tier1.NodeNumber; n++ {
		vendorIdx := n % uint32(len(template.Tier1.VendorInfo))
		tier1NodeIndex := n + 1
		tier1Node := NewSuperspineFabricNode(tier1NodeIndex, template.Tier1.VendorInfo[vendorIdx], f.log)

		f.addNode(PositionSuperspine, tier1Node, 0)
	}

	// process leaf/spine nodes
	// p is number of pod definitions
	for p, pod := range template.Pods {
		// i is the number of pods in a definition
		for i := uint32(0); i < pod.PodNumber; i++ {
			podIndex := (uint32(p) + 1) * (i + 1)

			log.Debug("podIndex", "podIndex", podIndex)

			// kind is tier 2 or tier3
			for kind, tier := range pod.Tiers {
				vendorNum := len(tier.VendorInfo)
				if kind != "tier3" && kind != "tier2" {
					return nil, fmt.Errorf("wrong kind in the template definition: %s, value: %s, expected tier2 or tier3", namespaceName, kind)
				}
				for n := uint32(0); n < tier.NodeNumber; n++ {
					vendorIdx := n % uint32(vendorNum)

					var fabricNode FabricNode

					if kind == "tier3" {
						fabricNode = NewLeafFabricNode(podIndex, n+1, tier.VendorInfo[vendorIdx], f.log)
						f.addNode(PositionLeaf, fabricNode, podIndex)

					} else {
						fabricNode = NewSpineFabricNode(podIndex, n+1, tier.VendorInfo[vendorIdx], f.log)
						f.addNode(PositionSpine, fabricNode, podIndex)
					}
				}
			}
		}
	}

	// process spine-leaf links
	for _, podInfo := range f.pods {
		for n, tier2Node := range podInfo.tier2Nodes {
			tier2NodeIndex := uint32(n) + 1
			for m, tier3Node := range podInfo.tier3Nodes {
				tier3NodeIndex := uint32(m) + 1

				epA := &Endpoint{
					Node:   tier2Node,
					IfName: tier2Node.GetInterfaceName(tier3NodeIndex),
				}
				epB := &Endpoint{
					Node:   tier3Node,
					IfName: tier3Node.GetInterfaceNameWithPlatformOffset(tier2NodeIndex),
				}
				f.addLink(PositionSpine, NewFabricLink(epA, epB))
			}
		}
	}

	// process superspine-spine links
	for n, tier1Node := range f.tier1Nodes {
		tier1NodeIndex := uint32(n) + 1
		// we need to get all the spine in ll the pods
		actualTier2 := uint32(0)
		for _, podInfo := range f.pods {
			for _, tier2Node := range podInfo.tier2Nodes {
				actualTier2++
				// this represents the total network index for the Spine
				//tier2NodeIndex := uint32(m) + 1 + (p * uint32(len(podInfo.tier2Nodes)))

				epA := &Endpoint{
					Node:   tier1Node,
					IfName: tier1Node.GetInterfaceName(actualTier2),
				}
				epB := &Endpoint{
					Node:   tier2Node,
					IfName: tier2Node.GetInterfaceNameWithPlatformOffset(tier1NodeIndex),
				}
				f.addLink(PositionSuperspine, NewFabricLink(epA, epB))

			}
		}
	}
	return f, nil
}

// +k8s:deepcopy-gen=false
type fabric struct {
	log             logging.Logger
	m               sync.Mutex
	tier1Nodes      []FabricNode
	pods            map[uint32]*podInfo
	tier2tier3Links []FabricLink
	tier1tier2Links []FabricLink
}

type podInfo struct {
	tier2Nodes []FabricNode // fabric nodes are stored per podIndex
	tier3Nodes []FabricNode // fabric nodes are stored per podIndex
}

func (f *fabric) addNode(pos Position, n FabricNode, podIndex uint32) {
	f.m.Lock()
	defer f.m.Unlock()

	// initialize the tier3/tier3 node struct per podIndex
	if _, ok := f.pods[podIndex]; !ok {
		f.pods[podIndex] = &podInfo{
			tier2Nodes: make([]FabricNode, 0),
			tier3Nodes: make([]FabricNode, 0),
		}
	}

	switch pos {
	case PositionLeaf:
		f.pods[podIndex].tier3Nodes = append(f.pods[podIndex].tier3Nodes, n)
	case PositionSpine:
		f.pods[podIndex].tier2Nodes = append(f.pods[podIndex].tier2Nodes, n)
	case PositionSuperspine:
		f.tier1Nodes = append(f.tier1Nodes, n)
	}
}

/*
func (f *fabric) getPodIndexes() uint32 {
	return uint32(len(f.tier2tier3Links))
}
*/

/*
func (f *fabric) getNodesPerPodIndex(pos Position, podIndex uint32) []FabricNode {
	f.m.Lock()
	defer f.m.Unlock()

	switch pos {
	case PositionLeaf:
		return f.tier3Nodes[podIndex]
	case PositionSpine:
		return f.tier2Nodes[podIndex]
	default:
		return nil
	}
}
*/

func (f *fabric) addLink(pos Position, l FabricLink) {
	switch pos {
	case PositionSpine:
		f.tier2tier3Links = append(f.tier2tier3Links, l)
	case PositionSuperspine:
		f.tier1tier2Links = append(f.tier1tier2Links, l)
	}
}

func (f *fabric) GetFabricNodes() []FabricNode {
	fn := make([]FabricNode, 0)
	fn = append(fn, f.tier1Nodes...)

	f.log.Debug("tier2Nodes", "length", len(f.pods))

	for _, podInfo := range f.pods {
		fn = append(fn, podInfo.tier2Nodes...)
		fn = append(fn, podInfo.tier3Nodes...)
	}
	return fn
}

func (f *fabric) GetFabricLinks() []FabricLink {
	fl := make([]FabricLink, 0, len(f.tier1tier2Links)+len(f.tier2tier3Links))
	fl = append(fl, f.tier1tier2Links...)
	fl = append(fl, f.tier2tier3Links...)
	return fl
}

func (f *fabric) PrintNodes() {
	for _, node := range f.tier1Nodes {
		f.log.Debug("tier1 node",
			"nodeName", node.GetNodeName(),
			"podIndex", node.GetPodIndex(),
			"vendorType", node.GetVendorType(),
			"platform", node.GetPlatform(),
			"position", node.GetPosition(),
		)
	}

	for _, podInfo := range f.pods {
		for _, node := range podInfo.tier2Nodes {
			f.log.Debug("tier2 node",
				"nodeName", node.GetNodeName(),
				"podIndex", node.GetPodIndex(),
				"vendorType", node.GetVendorType(),
				"platform", node.GetPlatform(),
				"position", node.GetPosition(),
			)
		}
		for _, node := range podInfo.tier3Nodes {
			f.log.Debug("tier3 node",
				"nodeName", node.GetNodeName(),
				"podIndex", node.GetPodIndex(),
				"vendorType", node.GetVendorType(),
				"platform", node.GetPlatform(),
				"position", node.GetPosition(),
			)
		}
	}
}

func (f *fabric) PrintLinks() {
	for _, link := range f.tier1tier2Links {
		f.log.Debug("link tier1tier2",
			"ep A nodeName", link.GetEndpointA().Node.GetNodeName(),
			"ep A podIndex", link.GetEndpointA().Node.GetPodIndex(),
			"ep A ifName", link.GetEndpointA().IfName,
			"ep B nodeName", link.GetEndpointB().Node.GetNodeName(),
			"ep B podIndex", link.GetEndpointB().Node.GetPodIndex(),
			"ep B ifName", link.GetEndpointB().IfName,
		)
	}
	for _, link := range f.tier2tier3Links {
		f.log.Debug("link tier2tier3",
			"ep A nodeName", link.GetEndpointA().Node.GetNodeName(),
			"ep A podIndex", link.GetEndpointA().Node.GetPodIndex(),
			"ep A ifName", link.GetEndpointA().IfName,
			"ep B nodeName", link.GetEndpointB().Node.GetNodeName(),
			"ep B podIndex", link.GetEndpointB().Node.GetPodIndex(),
			"ep B ifName", link.GetEndpointB().IfName,
		)
	}
}
