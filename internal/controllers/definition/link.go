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
	"strings"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/meta"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func renderFabricLink(cr *topov1alpha1.Definition, link topov1alpha1.FabricLink) *topov1alpha1.Link { // nolint:interfacer,gocyclo
	labels := map[string]string{
		LabelKeyOrganization:     cr.GetOrganization(),
		LabelKeyDeployment:       cr.GetDeployment(),
		LabelKeyAvailabilityZone: cr.GetAvailabilityZone(),
		LabelKeyTopology:         cr.GetTopologyName(),
	}
	return &topov1alpha1.Link{
		ObjectMeta: metav1.ObjectMeta{
			Name:            strings.Join([]string{cr.GetName(), link.GetName()}, "."),
			Namespace:       cr.Namespace,
			Labels:          labels,
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(cr, topov1alpha1.DefinitionGroupVersionKind))},
		},
		Spec: topov1alpha1.LinkSpec{
			Properties: &topov1alpha1.LinkProperties{
				Kind: topov1alpha1.LinkKindInfra,
				Endpoints: []*topov1alpha1.Endpoints{
					{
						InterfaceName: link.GetEndpointA().IfName,
						NodeName:      link.GetEndpointA().Node.GetNodeName(),
						Kind:          topov1alpha1.EndpointKindInfra,
					},
					{
						InterfaceName: link.GetEndpointB().IfName,
						NodeName:      link.GetEndpointB().Node.GetNodeName(),
						Kind:          topov1alpha1.EndpointKindInfra,
					},
				},
			},
		},
		Status: topov1alpha1.LinkStatus{
			ResourceStatus: nddv1.ResourceStatus{
				OdaInfo: nddv1.OdaInfo{
					Oda: map[string]string{
						string(nddv1.OdaKindOrganization):    cr.GetOrganization(),
						string(nddv1.OdaKindAvailabiityZone): cr.GetAvailabilityZone(),
						string(nddv1.OdaKindDeployment):      cr.GetDeployment(),
						string(nddv1.OdaKindResourceName):    cr.GetTopologyName(),
					},
				},
			},
		},
	}
}
