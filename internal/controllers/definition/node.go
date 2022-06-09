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

	"github.com/yndd/ndd-runtime/pkg/meta"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
