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

package link

import (

	//nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"

	"strings"

	"github.com/yndd/app-runtime/pkg/odns"
	"github.com/yndd/ndd-runtime/pkg/meta"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	shPrefix    = "logical-sh-link"
	mhPrefix    = "logical-mh-link"
	labelPrefix = "nddo-infra"
)

func buildLogicalTopologyLink(cr *topov1alpha1.Link) *topov1alpha1.Link {
	// the name of the logical link is set based on multi-homing or single homing
	// sh-lag:              <org>.<depl>.<topo>.<logical-sh-link>.<node-name-epA>.<lag-name-epA>.<node-name-epB><lag-name-epB>
	// mh-lag-A - sh-lag-B: <org>.<depl>.<topo>.<logical-mh-link>.<multihoming-A-name>.<node-name-epB><lag-name-epB>
	// sh-lag-A - mh-lag-B: <org>.<depl>.<topo>.<logical-mh-link>.<node-name-epB><lag-name-epB>.<multihoming-A-name>
	// mh-lag-A - mh-lag-B: <org>.<depl>.<topo>.<logical-mh-link>.<multihoming-A-name>.<multihoming-B-name>

	var (
		name           string
		nodeNameA      string
		nodeNameB      string
		interfaceNameA string
		interfaceNameB string
		epATags        map[string]string
		epBTags        map[string]string
	)
	/*
		//fmt.Printf("buildLogicalTopologyLink: %s mhA: %t, mhB: %t\n", cr.GetName(), cr.GetEndPointAMultiHoming(), cr.GetEndPointBMultiHoming())
		if cr.GetEndPointAMultiHoming() || cr.GetEndPointBMultiHoming() {
			if cr.GetEndPointAMultiHoming() {
				// multihomed endpoint A
				name = strings.Join([]string{mhPrefix, cr.GetEndPointAMultiHomingName()}, "-")
				nodeNameA = ""
				interfaceNameA = cr.GetEndPointAMultiHomingName()
				//tagNodeName := strings.Join([]string{topov1alpha1.NodePrefix, cr.GetEndpointANodeName()}, ":")

					epATags = cr.GetEndpointATagRaw()
					epATags = append(epATags, []*nddov1.Tag{
						{Key: utils.StringPtr(topov1alpha1.KeyLinkEPMultiHoming), Value: utils.StringPtr("true")},
						{Key: utils.StringPtr(tagNodeName), Value: utils.StringPtr(cr.GetLagAName())},
					}...)

			} else {
				name = strings.Join([]string{name, cr.GetEndpointANodeName()}, "-")
				nodeNameA = cr.GetEndpointANodeName()
				interfaceNameA = cr.GetLagAName()
				//epATags = cr.GetEndpointATagRaw()
			}
			if cr.GetEndPointBMultiHoming() {
				// multihomed endpoint B
				name = strings.Join([]string{mhPrefix, cr.GetEndPointBMultiHomingName()}, "-")
				nodeNameB = ""
				interfaceNameB = cr.GetEndPointAMultiHomingName()
				//tagNodeName := strings.Join([]string{topov1alpha1.NodePrefix, cr.GetEndpointBNodeName()}, ":")
					epBTags = cr.GetEndpointBTagRaw()
					epBTags = append(epBTags, []*nddov1.Tag{
						{Key: utils.StringPtr(topov1alpha1.KeyLinkEPMultiHoming), Value: utils.StringPtr("true")},
						{Key: utils.StringPtr(tagNodeName), Value: utils.StringPtr(cr.GetLagBName())},
					}...)
			} else {
				name = strings.Join([]string{name, cr.GetEndpointBNodeName()}, "-")
				nodeNameB = cr.GetEndpointBNodeName()
				interfaceNameB = cr.GetLagBName()
				//epBTags = cr.GetEndpointBTagRaw()
			}
			// prepend the topologyname

		} else {
			name = strings.Join([]string{shPrefix, cr.GetEndpointANodeName(), cr.GetLagAName(), cr.GetEndpointBNodeName(), cr.GetLagBName()}, "-")
			nodeNameA = cr.GetEndpointANodeName()
			interfaceNameA = cr.GetLagAName()
			//epATags = cr.GetEndpointATagRaw()
			nodeNameB = cr.GetEndpointBNodeName()
			interfaceNameB = cr.GetLagBName()
			//epBTags = cr.GetEndpointBTagRaw()
		}
	*/

	// prepend the parent logic link
	name = strings.Join([]string{odns.GetParentResourceName(cr.GetName()), name}, ".")

	//ndda := nddov1.NewOdaInfo()
	//ndda.SetOrganization(cr.GetOrganization())
	//ndda.SetDeployment(cr.GetDeployment())
	//ndda.SetAvailabilityZone(cr.GetAvailabilityZone())

	//fmt.Printf("buildLogicalTopologyLink: name: %s nodeA: %s, nodeB: %s, itfcA: %s, itfceB: %s\n", name, nodeNameA, nodeNameB, interfaceNameA, interfaceNameB)
	//fmt.Printf("buildLogicalTopologyLink: epAtags: %v, epBtags: %v\n", cr.GetEndpointATag(), cr.GetEndpointBTag())
	return &topov1alpha1.Link{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Namespace:       cr.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(cr, topov1alpha1.LinkGroupVersionKind))},
		},
		Spec: topov1alpha1.LinkSpec{
			Properties: &topov1alpha1.LinkProperties{
				Endpoints: []*topov1alpha1.Endpoints{
					{
						NodeName:      nodeNameA,
						InterfaceName: interfaceNameA,
						Tag:           epATags,
					},
					{
						NodeName:      nodeNameB,
						InterfaceName: interfaceNameB,
						Tag:           epBTags,
					},
				},
				Tag: map[string]string{"lag": "true"},
				/*
					Tag: []*nddv1.Tag{
						{Key: utils.StringPtr("lag"), Value: utils.StringPtr("true")},
					},
				*/
			},
			//TopologyName: utils.StringPtr(topologyName),

		},
		/*
			Status: topov1alpha1.TopologyLinkStatus{
				TopologyLink: &topov1alpha1.NddrTopologyTopologyLink{
					Tag: cr.GetStatusTagsRaw(),
				},
			},
		*/
	}
	//l.Spec.Oda = ndda.Oda
	//return l
}

func updateLogicalTopologyLink(cr *topov1alpha1.Link, mhtl *topov1alpha1.Link) *topov1alpha1.Link {
	/*
		if cr.GetEndPointAMultiHoming() && mhtl.GetEndPointAMultiHoming() && (cr.GetEndPointAMultiHomingName() == mhtl.GetEndPointAMultiHomingName()) {
			nodeNameA := cr.GetEndpointANodeName()
			interfaceNameA := cr.GetLagAName()

			//fmt.Printf("updateLogicalTopologyLink: nodename: %s, itfcename: %s\n", nodeNameA, interfaceNameA)

			found := false
			for _, tag := range mhtl.GetEndpointATagRaw() {
				if *tag.Key == nodeNameA && *tag.Value == interfaceNameA {
					found = true
					break
				}
			}
			tagNodeName := strings.Join([]string{topov1alpha1.NodePrefix, nodeNameA}, ":")
			if !found {
				mhtl.AddEndPointATag(tagNodeName, interfaceNameA)
			}
		}
		if cr.GetEndPointBMultiHoming() && mhtl.GetEndPointBMultiHoming() && (cr.GetEndPointBMultiHomingName() == mhtl.GetEndPointBMultiHomingName()) {
			nodeNameB := cr.GetEndpointBNodeName()
			interfaceNameB := cr.GetLagBName()

			found := false
			for _, tag := range mhtl.GetEndpointBTagRaw() {
				if *tag.Key == nodeNameB && *tag.Value == interfaceNameB {
					found = true
					break
				}
			}
			tagNodeName := strings.Join([]string{topov1alpha1.NodePrefix, nodeNameB}, ":")
			if !found {
				mhtl.AddEndPointBTag(tagNodeName, interfaceNameB)
			}
		}
	*/
	return mhtl
}

func updateDeleteLogicalTopologyLink(cr *topov1alpha1.Link, mhtl *topov1alpha1.Link) *topov1alpha1.Link {
	/*
		if cr.GetEndPointAMultiHoming() && mhtl.GetEndPointAMultiHoming() && (cr.GetEndPointAMultiHomingName() == mhtl.GetEndPointAMultiHomingName()) {
			interfaceNameA := cr.GetLagAName()
			tagNodeName := strings.Join([]string{topov1alpha1.NodePrefix, cr.GetEndpointANodeName()}, ":")

			//fmt.Printf("updateLogicalTopologyLink: nodename: %s, itfcename: %s\n", nodeNameA, interfaceNameA)

			mhtl.DeleteEndPointATag(tagNodeName, interfaceNameA)
		}
		if cr.GetEndPointBMultiHoming() && mhtl.GetEndPointBMultiHoming() && (cr.GetEndPointBMultiHomingName() == mhtl.GetEndPointBMultiHomingName()) {
			interfaceNameB := cr.GetLagBName()
			tagNodeName := strings.Join([]string{topov1alpha1.NodePrefix, cr.GetEndpointBNodeName()}, ":")

			mhtl.DeleteEndPointBTag(tagNodeName, interfaceNameB)
		}
	*/
	return mhtl
}

func updateDeleteLogicalTopologyLinkNodeEndpoint(cr *topov1alpha1.Link, i int, nodeName, interfaceName string) *topov1alpha1.Link {
	/*
		switch i {
		case 0:
			cr.DeleteEndPointATag(nodeName, interfaceName)
		case 1:
			cr.DeleteEndPointATag(nodeName, interfaceName)
		}
	*/
	return cr
}
