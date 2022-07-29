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
	"context"
	"fmt"
	"sync"

	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/meta"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Option can be used to manipulate Fabric config.
type Option func(Fabric)

// WithLogger specifies how the Fabric logs messages.
func WithLogger(log logging.Logger) Option {
	return func(f Fabric) {
		f.SetLogger(log)
	}
}

// WithClient specifies the fabric to use within the client.
func WithClient(c client.Client) Option {
	return func(f Fabric) {
		f.SetClient(c)
	}
}

// +k8s:deepcopy-gen=false
type Fabric interface {
	GetFabricNodes() []FabricNode
	GetFabricLinks() []FabricLink
	PrintNodes()
	PrintLinks()

	SetLogger(logger logging.Logger)
	SetClient(c client.Client)
}

func NewFabric(namespaceName string, template *topov1alpha1.FabricTemplate, opts ...Option) (Fabric, error) {
	f := &fabric{
		tier1Nodes:      make([]FabricNode, 0),
		pods:            map[uint32]*podInfo{},
		tier2tier3Links: make([]FabricLink, 0),
		tier1tier2Links: make([]FabricLink, 0),
	}

	for _, opt := range opts {
		opt(f)
	}

	// a template can have multiple template/definition references so we need to parse them
	// to build one fabric topology
	mergedTemplate, err := f.parseTemplate(template)
	if err != nil {
		return nil, err
	}

	f.log.Debug("mergedTemplate", "mergedTemplate", mergedTemplate)

	// process leaf/spine nodes
	// p is number of pod definitions
	for p, pod := range mergedTemplate.Pod {
		// i is the number of pods in a definition
		for i := uint32(0); i < pod.GetPodNumber(); i++ {
			// podIndex is pod template index * pod index within the template
			podIndex := (uint32(p) + 1) * (i + 1)

			//log.Debug("podIndex", "podIndex", podIndex)

			// tier 2 -> spines in the pod
			f.processPodNodeTier("tier2", podIndex, pod.Tier2)
			// tier 3 -> leafs in the pod
			f.processPodNodeTier("tier3", podIndex, pod.Tier3)
		}
	}

	// proces superspines
	// the superspine is equal to the amount of spines per pod and multiplied with the number in the template
	if mergedTemplate.Tier1 != nil {
		superspines := f.getSuperSPines()
		// process superspine nodes
		for n := uint32(0); n < superspines; n++ {
			for m := uint32(0); m < mergedTemplate.Tier1.NodeNumber; m++ {
				// venndor Index is used to map to the particular node based on modulo
				// if 1 vendor -> all nodes are from 1 vendor
				// if 2 vendors -> all odd nodes will be vendor A and all even nodes will be vendor B
				// if 3 vendors -> 1st is vendorA, 2nd vendor B, 3rd is vendor C
				vendorIdx := n % uint32(len(mergedTemplate.Tier1.VendorInfo))

				// PlaneIndex: m + 1 -> starts counting from 1, used when multiple nodes are used in the superspine plane
				// NodeIndex: n + 1 -> could also be called the Plane Index
				tier1Node := NewSuperspineFabricNode(m+1, n+1, mergedTemplate.Tier1.VendorInfo[vendorIdx], f.log)

				f.addNode(topov1alpha1.PositionSuperspine, tier1Node, 0)
			}
		}
	}

	// process spine-leaf links
	for _, podInfo := range f.pods {
		for n, tier2Node := range podInfo.tier2Nodes {
			tier2NodeIndex := uint32(n) + 1
			for m, tier3Node := range podInfo.tier3Nodes {
				tier3NodeIndex := uint32(m) + 1

				// validate if the uplinks per node is not greater than max uplinks
				// otherwise there is a conflict and the algorithm behind will create
				// overlapping indexes
				uplinksPerNode := tier3Node.GetUplinkPerNode()
				if uplinksPerNode > mergedTemplate.MaxUplinksTier3ToTier2 {
					return nil, fmt.Errorf("uplink per node %d can not be bigger than maxUplinksTier3ToTier2 %d",
						uplinksPerNode, mergedTemplate.MaxUplinksTier3ToTier2)
				}

				// the algorithm needs to avoid reindixing if changes happen -> introduced maxNumUplinks
				// the allocation is first allocating the uplink Index
				// u represnts the actual uplink index
				// spine Index    -> actualUplinkId + (actual leafs  * max uplinks)
				// leaf  Index    -> actualUplinkId + (actual spines * max uplinks)
				// actualUplinkId = u + 1 -> counting starts at 1
				// actual leafs   = tier3NodeIndex - 1 -> counting from 0
				// actual spines  = tier2NodeIndex - 1 -> counting from 0
				// max uplinks    = mergedTemplate.MaxUplinksTier3ToTier2
				for u := uint32(0); u < uplinksPerNode; u++ {
					epA := &Endpoint{
						Node:   tier2Node,
						IfName: tier2Node.GetInterfaceName(u + 1 + ((tier3NodeIndex - 1) * mergedTemplate.MaxUplinksTier3ToTier2)),
					}
					epB := &Endpoint{
						Node:   tier3Node,
						IfName: tier3Node.GetInterfaceNameWithPlatfromOffset(u + 1 + ((tier2NodeIndex - 1) * mergedTemplate.MaxUplinksTier3ToTier2)),
					}
					f.addLink(topov1alpha1.PositionSpine, NewFabricLink(epA, epB))
				}
			}
		}
	}

	// process superspine-spine links
	for _, tier1Node := range f.tier1Nodes {
		for p, podInfo := range f.pods {
			for m, tier2Node := range podInfo.tier2Nodes {

				// validate if the uplinks per node is not greater than max uplinks
				// otherwise there is a conflict and the algorithm behind will create
				// overlapping indexes
				uplinksPerNode := tier2Node.GetUplinkPerNode()
				if uplinksPerNode > mergedTemplate.MaxUplinksTier2ToTier1 {
					return nil, fmt.Errorf("uplink per node %d can not be bigger than maxUplinksTier2ToTier1 %d", uplinksPerNode, mergedTemplate.MaxUplinksTier2ToTier1)
				}

				// spine and superspine line up so we only create a link if there is a match
				// on the indexes
				if (m + 1) == int(tier1Node.GetNodeIndex()) {
					// the algorithm needs to avoid reindixing if changes happen -> introduced maxNumUplinks
					// the allocation is first allocating the uplink Index
					// u represnts the actual uplink index
					// superspine Index -> actualUplinkId + (actual podIndex  * max uplinks)
					// spine Index      -> actualUplinkId + (actual spines per plane * max uplinks)
					// actualUplinkId          = u + 1 -> counting starts at 1
					// actual PodIndex         = p +1
					// actual spines per plane = tier1Node.GetNodePlaneIndex() - 1
					// max uplinks             = mergedTemplate.MaxUplinksTier2ToTier1
					for u := uint32(0); u < uplinksPerNode; u++ {
						epA := &Endpoint{
							Node:   tier1Node,
							IfName: tier1Node.GetInterfaceName(u + 1 + ((p - 1) * mergedTemplate.MaxUplinksTier2ToTier1)),
						}
						epB := &Endpoint{
							Node:   tier2Node,
							IfName: tier2Node.GetInterfaceNameWithPlatfromOffset(u + 1 + ((tier1Node.GetNodePlaneIndex() - 1) * mergedTemplate.MaxUplinksTier2ToTier1)),
						}
						f.addLink(topov1alpha1.PositionSuperspine, NewFabricLink(epA, epB))
					}
				}
			}
		}
	}
	return f, nil
}

// +k8s:deepcopy-gen=false
type fabric struct {
	log             logging.Logger
	client          client.Client
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

func (f *fabric) addNode(pos topov1alpha1.Position, n FabricNode, podIndex uint32) {
	f.m.Lock()
	defer f.m.Unlock()

	// initialize the tier3/tier3 node struct per podIndex
	if pos != topov1alpha1.PositionSuperspine {
		if _, ok := f.pods[podIndex]; !ok {
			f.pods[podIndex] = &podInfo{
				tier2Nodes: make([]FabricNode, 0),
				tier3Nodes: make([]FabricNode, 0),
			}
		}
	}

	switch pos {
	case topov1alpha1.PositionLeaf:
		f.pods[podIndex].tier3Nodes = append(f.pods[podIndex].tier3Nodes, n)
	case topov1alpha1.PositionSpine:
		f.pods[podIndex].tier2Nodes = append(f.pods[podIndex].tier2Nodes, n)
	case topov1alpha1.PositionSuperspine:
		f.tier1Nodes = append(f.tier1Nodes, n)
	}
}

// getSuperSPines identifies the max number of spines in a pod
func (f *fabric) getSuperSPines() uint32 {
	var superspines uint32
	for _, podInfo := range f.pods {
		if superspines < uint32(len(podInfo.tier2Nodes)) {
			superspines = uint32(len(podInfo.tier2Nodes))
		}
	}
	return superspines
}

func (f *fabric) addLink(pos topov1alpha1.Position, l FabricLink) {
	switch pos {
	case topov1alpha1.PositionSpine:
		f.tier2tier3Links = append(f.tier2tier3Links, l)
	case topov1alpha1.PositionSuperspine:
		f.tier1tier2Links = append(f.tier1tier2Links, l)
	}
}

func (f *fabric) processPodNodeTier(tier string, podIndex uint32, tierTempl *topov1alpha1.TierTemplate) {
	vendorNum := len(tierTempl.VendorInfo)
	for n := uint32(0); n < tierTempl.NodeNumber; n++ {
		// n is the node Index within the tier
		vendorIdx := n % uint32(vendorNum)

		var fabricNode FabricNode

		if tier == "tier3" {
			// create a leaf node in the fabric
			// podIndex is the index of the pod -> counting starts from 1
			// nodeIndex (n+1) is the nodeIndex within the pod -> countng starts from 1
			fabricNode = NewLeafFabricNode(podIndex, n+1, tierTempl.UplinksPerNode, tierTempl.VendorInfo[vendorIdx], f.log)
			f.addNode(topov1alpha1.PositionLeaf, fabricNode, podIndex)

		} else {
			// create a spine node in the fabric
			// podIndex is the index of the pod -> counting starts from 1
			// nodeIndex (n+1) is the nodeIndex within the pod -> countng starts from 1
			fabricNode = NewSpineFabricNode(podIndex, n+1, tierTempl.UplinksPerNode, tierTempl.VendorInfo[vendorIdx], f.log)
			f.addNode(topov1alpha1.PositionSpine, fabricNode, podIndex)

		}
	}
}

func (f *fabric) SetLogger(log logging.Logger) {
	f.log = log
}

func (f *fabric) SetClient(c client.Client) {
	f.client = c
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

func (f *fabric) parseTemplate(template *topov1alpha1.FabricTemplate) (*topov1alpha1.FabricTemplate, error) {
	mergedTemplate := &topov1alpha1.FabricTemplate{}

	if err := template.CheckTemplate(true); err != nil {
		return nil, err
	}

	if template.HasReference() {
		f.log.Debug("parseTemplate", "hasReference", true)
		mergedTemplate.BorderLeaf = template.BorderLeaf
		mergedTemplate.Tier1 = template.Tier1
		mergedTemplate.MaxUplinksTier2ToTier1 = template.MaxUplinksTier2ToTier1
		mergedTemplate.MaxUplinksTier3ToTier2 = template.MaxUplinksTier2ToTier1
		mergedTemplate.Pod = make([]*topov1alpha1.PodTemplate, 0)
		for _, pod := range template.Pod {
			if pod.TemplateReference != nil {
				pd, err := f.getPodDefintionFromTemplate(*pod.TemplateReference)
				if err != nil {
					return nil, err
				}
				mergedTemplate.Pod = append(mergedTemplate.Pod, pd)
			}
			if pod.DefinitionReference != nil {
				name, namespace := meta.NamespacedName(*pod.DefinitionReference).GetNameAndNamespace()
				t := &topov1alpha1.Definition{}
				if err := f.client.Get(context.TODO(), types.NamespacedName{
					Namespace: namespace,
					Name:      name,
				}, t); err != nil {
					return nil, err
				}
				if len(t.Spec.Properties.Templates) != 1 {
					return nil, fmt.Errorf("definition can only have 1 template")
				}

				pd, err := f.getPodDefintionFromTemplate(t.Spec.Properties.Templates[0].NamespacedName)
				if err != nil {
					return nil, err
				}
				mergedTemplate.Pod = append(mergedTemplate.Pod, pd)

			}
		}
	} else {
		mergedTemplate = template
	}

	return mergedTemplate, nil
}

func (f *fabric) getPodDefintionFromTemplate(name string) (*topov1alpha1.PodTemplate, error) {
	name, namespace := meta.NamespacedName(name).GetNameAndNamespace()
	t := &topov1alpha1.Template{}
	if err := f.client.Get(context.TODO(), types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}, t); err != nil {
		return nil, err
	}
	if err := t.Spec.Properties.Fabric.CheckTemplate(false); err != nil {
		return nil, err
	}
	return t.Spec.Properties.Fabric.Pod[0], nil
}
