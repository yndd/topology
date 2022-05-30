/*
Copyright 2021 Wim Henderickx.

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

//+kubebuilder:object:generate=true
package v1alpha1

import (
	"context"
	"strconv"
	"strings"

	nddappv1 "github.com/yndd/app-runtime/apis/common/v1"
	"github.com/yndd/app-runtime/pkg/intent"
	"github.com/yndd/app-runtime/pkg/odns"
	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/meta"
	"github.com/yndd/ndd-runtime/pkg/resource"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InitNode(c resource.ClientApplicator, p intent.Intent, name string) intent.Intent {
	newNodeList := func() topov1alpha1.TnList { return &topov1alpha1.TopologyNodeList{} }
	return &nodeIntent{
		client:      c,
		name:        name,
		parent:      p,
		properties:  &topov1alpha1.TopologyNodeProperties{},
		labels:      map[string]string{},
		newNodeList: newNodeList,
	}
}

type nodeIntent struct {
	// k8s client
	client resource.ClientApplicator
	// key
	name string
	// parent
	parent intent.Intent
	// children
	// Data
	properties  *topov1alpha1.TopologyNodeProperties
	labels      map[string]string
	newNodeList func() topov1alpha1.TnList
}

func (x *nodeIntent) GetData() interface{} {
	return x.properties
}

func (x *nodeIntent) GetLabels() map[string]string {
	return x.labels
}

func (x *nodeIntent) Deploy(ctx context.Context, mr resource.Managed, labels map[string]string) error {
	cr, err := x.buildCR(mr, x.name, labels)
	if err != nil {
		return err
	}
	return x.client.Apply(ctx, cr)
}

func (x *nodeIntent) Destroy(ctx context.Context, mr resource.Managed, labels map[string]string) error {
	cr, err := x.buildCR(mr, x.name, labels)
	if err != nil {
		return err
	}
	return x.client.Delete(ctx, cr)
}

func (x *nodeIntent) List(ctx context.Context, mr resource.Managed, resources map[string]map[string]resource.Managed) (map[string]map[string]resource.Managed, error) {
	// local CR list
	opts := []client.ListOption{
		client.MatchingLabels{nddappv1.LabelKeyOwner: odns.GetOdnsResourceKindName(mr.GetName(), strings.ToLower(mr.GetObjectKind().GroupVersionKind().Kind))},
	}
	list := x.newNodeList()
	if err := x.client.List(ctx, list, opts...); err != nil {
		return nil, err
	}

	for _, d := range list.GetNodes() {
		if _, ok := resources[d.GetObjectKind().GroupVersionKind().Kind]; !ok {
			resources[d.GetObjectKind().GroupVersionKind().Kind] = make(map[string]resource.Managed)
		}
		resources[d.GetObjectKind().GroupVersionKind().Kind][d.GetName()] = d
	}

	return resources, nil
}

func (x *nodeIntent) Validate(ctx context.Context, mr resource.Managed, resources map[string]map[string]resource.Managed) (map[string]map[string]resource.Managed, error) {
	// local CR validation
	resourceName := odns.GetOdnsResourceName(mr.GetName(), strings.ToLower(mr.GetObjectKind().GroupVersionKind().Kind),
		[]string{
			strings.ToLower(x.name)})

	if r, ok := resources[topov1alpha1.TopologyNodeKind]; ok {
		delete(r, resourceName)
	}

	return resources, nil
}

func (x *nodeIntent) Delete(ctx context.Context, mr resource.Managed, resources map[string]map[string]resource.Managed) error {
	// local CR deletion
	if res, ok := resources[topov1alpha1.TopologyNodeKind]; ok {
		for resName := range res {
			o := &topov1alpha1.TopologyNode{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resName,
					Namespace: mr.GetNamespace(),
				},
			}
			if err := x.client.Delete(ctx, o); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *nodeIntent) buildCR(mr resource.Managed, deviceName string, labels map[string]string) (*topov1alpha1.TopologyNode, error) {
	resourceName := odns.GetOdnsResourceName(mr.GetName(), strings.ToLower(mr.GetObjectKind().GroupVersionKind().Kind),
		[]string{
			//strings.ToLower(x.name),
			strings.ToLower(deviceName)})

	//labels[nddappv1.LabelKeyLifecyclePolicy] = string(mr.GetLifecyclePolicy())
	labels[nddappv1.LabelKeyOwner] = odns.GetOdnsResourceKindName(mr.GetName(), strings.ToLower(mr.GetObjectKind().GroupVersionKind().Kind))
	labels[nddappv1.LabelKeyOwnerGeneration] = strconv.Itoa(int(mr.GetGeneration()))
	labels[nddappv1.LabelKeyTarget] = deviceName
	//labels[srlv1alpha1.LabelNddaItfce] = itfceName

	namespace := mr.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	return &topov1alpha1.TopologyNode{
		ObjectMeta: metav1.ObjectMeta{
			Name:            resourceName,
			Namespace:       namespace,
			Labels:          labels,
			OwnerReferences: []metav1.OwnerReference{meta.AsController(meta.TypedReferenceTo(mr, mr.GetObjectKind().GroupVersionKind()))},
		},
		Spec: topov1alpha1.TopologyNodeSpec{
			ResourceSpec: nddv1.ResourceSpec{},
			Properties:   x.properties,
		},
	}, nil
}
