/*
Copyright 2022 NDD.

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
	"context"

	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
)

const (
	LabelKeyDiscoveryRule = "discovery.yndd.io/discovery-rule"
	LabelKeyVendorType    = "discovery.yndd.io/vendor-type"
)

type adder interface {
	Add(item interface{})
}
type EnqueueRequestForAllTargets struct {
	client client.Client
	log    logging.Logger
	ctx    context.Context
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTargets) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTargets) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.ObjectOld, q)
	e.add(evt.ObjectNew, q)
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTargets) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTargets) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

func (e *EnqueueRequestForAllTargets) add(obj runtime.Object, queue adder) {
	cr, ok := obj.(*targetv1.Target)
	if !ok {
		return
	}
	log := e.log.WithValues("event handler", "Target", "namespace", cr.GetNamespace(), "name", cr.GetName())
	log.Debug("handleEvent")

	//r, ok := cr.GetLabels()[discovery.LabelKeyDiscoveryRule]
	r, ok := cr.GetLabels()[LabelKeyDiscoveryRule]
	if !ok {
		log.Debug("target without discovery rule")
		return
	}

	tdl := &topov1alpha1.DefinitionList{}
	if err := e.client.List(e.ctx, tdl); err != nil {
		log.Debug("cannot get topology definition list", "error", err)
		return
	}

	for _, td := range tdl.Items {
		if td.Spec.Properties.DiscoveryRules != nil {
			for _, dr := range td.Spec.Properties.DiscoveryRules {
				namespace, name := meta.NamespacedName(dr.NamespacedName).GetNameAndNamespace()
				if namespace == cr.GetNamespace() && name == r {
					queue.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Namespace: td.GetNamespace(),
						Name:      td.GetName()}})
					// a target can only belong to 1 topology
					// a discovery rule cannot be mapped to multiple topologies
					return
				}
			}
		}
	}
}
