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

package topologynode

import (
	"context"

	//ndddvrv1 "github.com/yndd/ndd-core/apis/dvr/v1"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/nddo-runtime/pkg/odns"
	topov1alpha1 "github.com/yndd/nddr-topo-registry/apis/topo/v1alpha1"
	"github.com/yndd/nddr-topo-registry/internal/handler"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type adder interface {
	Add(item interface{})
}

type EnqueueRequestForAllTopologies struct {
	client client.Client
	log    logging.Logger
	ctx    context.Context

	handler handler.Handler

	newTopoNodeList func() topov1alpha1.TnList
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTopologies) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTopologies) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.ObjectOld, q)
	e.add(evt.ObjectNew, q)
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTopologies) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

// Create enqueues a request for all infrastructures which pertains to the topology.
func (e *EnqueueRequestForAllTopologies) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	e.add(evt.Object, q)
}

func (e *EnqueueRequestForAllTopologies) add(obj runtime.Object, queue adder) {
	dd, ok := obj.(*topov1alpha1.Topology)
	if !ok {
		return
	}
	log := e.log.WithValues("function", "watch topologies", "name", dd.GetName())
	log.Debug("topologynode handleEvent")

	d := e.newTopoNodeList()
	if err := e.client.List(e.ctx, d); err != nil {
		return
	}

	watchDnsName, _ := odns.Name2OdnsTopo(dd.GetName()).GetFullOdaName()

	for _, toponode := range d.GetNodes() {
		// only enqueue if the topology name match
		//if toponode.GetTopologyName() == dd.GetName() {
		nodeDnsName, _ := odns.Name2OdnsTopo(toponode.GetName()).GetFullOdaName()
		if nodeDnsName == watchDnsName {

			crName := getCrName(toponode)
			e.handler.ResetSpeedy(crName)

			queue.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Namespace: toponode.GetNamespace(),
				Name:      toponode.GetName()}})
		}
	}
}
