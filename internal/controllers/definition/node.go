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

package definition

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yndd/ndd-runtime/pkg/meta"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	LabelKeyTopologyPosition   = "topology.yndd.io/position"
	LabelKeyTopologyNodeIndex  = "topology.yndd.io/NodeIndex"
	LabelKeyTopologyPodIndex   = "topology.yndd.io/PodIndex"
	LabelKeyTopologyPlatform   = "topology.yndd.io/Platform"
	LabelKeyTopologyVendorType = "topology.yndd.io/VendorType"
	LabelKeyOrganization       = "org.yndd.io/organization"
	LabelKeyDeployment         = "org.yndd.io/deployment"
	LabelKeyAvailabilityZone   = "org.yndd.io/availabilityzone"
	LabelKeyTopology           = "org.yndd.io/topology"
)

func renderNode(drName string, cr *topov1alpha1.Definition, t *targetv1.Target) *topov1alpha1.Node { // nolint:interfacer,gocyclo
	return &topov1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.Join([]string{cr.GetName(), t.GetName()}, "."),
			Namespace: cr.Namespace,
			Labels: map[string]string{
				LabelKeyDiscoveryRule: drName,
			},
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(cr, topov1alpha1.DefinitionGroupVersionKind))},
		},
		Spec: topov1alpha1.NodeSpec{
			Properties: &topov1alpha1.NodeProperties{
				//VendorType: t.GetDiscoveryInfo().VendorType,
				Platform: t.GetDiscoveryInfo().Platform,
				//Index:
				//Position:
				// Tags://
			},
		},
	}
}

type FabricNodeInfo struct {
	Position   topov1alpha1.Position
	NodeIndex  uint32 // relative number within the position, pod
	PodIndex   uint32 // pod index
	VendorInfo *topov1alpha1.FabricTierVendorInfo
}

func renderFabricNode(cr *topov1alpha1.Definition, fabricNodeInfo *FabricNodeInfo) *topov1alpha1.Node { // nolint:interfacer,gocyclo
	var nodeName string
	labels := map[string]string{
		LabelKeyTopologyPosition:   string(fabricNodeInfo.Position),
		LabelKeyTopologyNodeIndex:  strconv.Itoa(int(fabricNodeInfo.NodeIndex)),
		LabelKeyTopologyVendorType: string(fabricNodeInfo.VendorInfo.VendorType),
		LabelKeyTopologyPlatform:   string(fabricNodeInfo.VendorInfo.Platform),
		LabelKeyOrganization:       cr.GetOrganization(),
		LabelKeyDeployment:         cr.GetDeployment(),
		LabelKeyAvailabilityZone:   cr.GetAvailabilityZone(),
		LabelKeyTopology:           cr.GetTopologyName(),
	}
	if fabricNodeInfo.Position != topov1alpha1.PositionSuperspine {
		nodeName = fmt.Sprintf("pod%d-%s%d", fabricNodeInfo.PodIndex, fabricNodeInfo.Position, fabricNodeInfo.NodeIndex)
		labels[LabelKeyTopologyPodIndex] = strconv.Itoa(int(fabricNodeInfo.PodIndex))
	} else {
		nodeName = fmt.Sprintf("%s%d", fabricNodeInfo.Position, fabricNodeInfo.NodeIndex)

	}
	return &topov1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:            strings.Join([]string{cr.GetName(), nodeName}, "."),
			Namespace:       cr.Namespace,
			Labels:          labels,
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(cr, topov1alpha1.DefinitionGroupVersionKind))},
		},
		Spec: topov1alpha1.NodeSpec{
			Properties: &topov1alpha1.NodeProperties{
				VendorType: fabricNodeInfo.VendorInfo.VendorType,
				Platform:   fabricNodeInfo.VendorInfo.Platform,
				//Index:
				Position: fabricNodeInfo.Position,
				// Tags://
			},
		},
	}
}
