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
	"github.com/yndd/ndd-runtime/pkg/meta"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func renderTopology(cr *topov1alpha1.TopologyDefinition) *topov1alpha1.Topology { // nolint:interfacer,gocyclo
	return &topov1alpha1.Topology{
		ObjectMeta: metav1.ObjectMeta{
			Name:            cr.GetName(),
			Namespace:       cr.Namespace,
			Labels:          map[string]string{},
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(cr, topov1alpha1.TopologyDefinitionGroupVersionKind))},
		},
		Spec: topov1alpha1.TopologySpec{
			Properties: topov1alpha1.TopologyProperties{
				AdminState: "enable",
				Defaults: &topov1alpha1.TopologyDefaults{
					NodeAttributes: &topov1alpha1.NodeAttributes{
						//Position: "infra",
					},
				},
				VendorTypeInfo: []*topov1alpha1.VendorTypeInfo{
					{
						VendorType:     "nokia-srl",
						Platform:       "7220 IXR-D2",
						NodeAttributes: &topov1alpha1.NodeAttributes{
							//Position: "infra",
						},
					},
					{
						VendorType:     "nokia-sros",
						Platform:       "7750 SR1",
						NodeAttributes: &topov1alpha1.NodeAttributes{
							//Position: "infra",
						},
					},
				},
			},
		},
	}
}
