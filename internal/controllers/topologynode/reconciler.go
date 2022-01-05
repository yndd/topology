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
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yndd/ndd-runtime/pkg/event"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/nddo-runtime/pkg/reconciler/managed"
	"github.com/yndd/nddo-runtime/pkg/resource"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/source"

	topov1alpha1 "github.com/yndd/nddr-topo-registry/apis/topo/v1alpha1"
	"github.com/yndd/nddr-topo-registry/internal/handler"
	"github.com/yndd/nddr-topo-registry/internal/shared"
	corev1 "k8s.io/api/core/v1"
)

const (
	// timers
	reconcileTimeout = 1 * time.Minute
	veryShortWait    = 1 * time.Second
	// errors
	errUnexpectedResource = "unexpected organization object"
	errGetK8sResource     = "cannot get organization resource"
)

// Setup adds a controller that reconciles infra.
func Setup(mgr ctrl.Manager, o controller.Options, nddcopts *shared.NddControllerOptions) error {
	name := "nddo/" + strings.ToLower(topov1alpha1.TopologyNodeGroupKind)
	tnfn := func() topov1alpha1.Tn { return &topov1alpha1.TopologyNode{} }
	tnlfn := func() topov1alpha1.TnList { return &topov1alpha1.TopologyNodeList{} }
	tpfn := func() topov1alpha1.Tp { return &topov1alpha1.Topology{} }

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(topov1alpha1.TopologyNodeGroupVersionKind),
		managed.WithLogger(nddcopts.Logger.WithValues("controller", name)),
		managed.WithApplication(&application{
			client: resource.ClientApplicator{
				Client:     mgr.GetClient(),
				Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
			},
			log:             nddcopts.Logger.WithValues("applogic", name),
			newTopology:     tpfn,
			newTopologyNode: tnfn,
			handler:         nddcopts.Handler,
		}),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	)

	topologyHandler := &EnqueueRequestForAllTopologies{
		client:          mgr.GetClient(),
		log:             nddcopts.Logger,
		ctx:             context.Background(),
		handler:         nddcopts.Handler,
		newTopoNodeList: tnlfn,
	}

	topologyLinkHandler := &EnqueueRequestForAllTopologyLinks{
		client:          mgr.GetClient(),
		log:             nddcopts.Logger,
		ctx:             context.Background(),
		handler:         nddcopts.Handler,
		newTopoNodeList: tnlfn,
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o).
		For(&topov1alpha1.TopologyNode{}).
		Owns(&topov1alpha1.TopologyNode{}).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Watches(&source.Kind{Type: &topov1alpha1.Topology{}}, topologyHandler).
		Watches(&source.Kind{Type: &topov1alpha1.TopologyLink{}}, topologyLinkHandler).
		Complete(r)
}

type application struct {
	client resource.ClientApplicator
	log    logging.Logger

	newTopology     func() topov1alpha1.Tp
	newTopologyNode func() topov1alpha1.Tn

	handler handler.Handler
}

func getCrName(cr topov1alpha1.Tn) string {
	return strings.Join([]string{cr.GetNamespace(), cr.GetName()}, ".")
}

func (r *application) Initialize(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*topov1alpha1.TopologyNode)
	if !ok {
		return errors.New(errUnexpectedResource)
	}

	if err := cr.InitializeResource(); err != nil {
		r.log.Debug("Cannot initialize", "error", err)
		return err
	}

	return nil
}

func (r *application) Update(ctx context.Context, mg resource.Managed) (map[string]string, error) {
	cr, ok := mg.(*topov1alpha1.TopologyNode)
	if !ok {
		return nil, errors.New(errUnexpectedResource)
	}

	return r.handleAppLogic(ctx, cr)
}

func (r *application) FinalUpdate(ctx context.Context, mg resource.Managed) {
}

func (r *application) Timeout(ctx context.Context, mg resource.Managed) time.Duration {
	/*
		cr, _ := mg.(*orgv1alpha1.Organization)
		crName := getCrName(cr)
		speedy := r.handler.GetSpeedy(crName)
		if speedy <= 2 {
			r.handler.IncrementSpeedy(crName)
			r.log.Debug("Speedy incr", "number", r.handler.GetSpeedy(crName))
			switch speedy {
			case 0:
				return veryShortWait
			case 1, 2:
				return shortWait
			}

		}
	*/
	return reconcileTimeout
}

func (r *application) Delete(ctx context.Context, mg resource.Managed) (bool, error) {
	return true, nil
}

func (r *application) FinalDelete(ctx context.Context, mg resource.Managed) {
	cr, ok := mg.(*topov1alpha1.TopologyNode)
	if !ok {
		return
	}
	crName := getCrName(cr)
	r.handler.Delete(crName)
}

func (r *application) handleAppLogic(ctx context.Context, cr topov1alpha1.Tn) (map[string]string, error) {
	log := r.log.WithValues("function", "handleAppLogic", "crname", cr.GetName())
	log.Debug("handleAppLogic")

	// initialize speedy
	crName := getCrName(cr)
	r.handler.Init(crName)

	// get the topo
	topo := r.newTopology()
	if err := r.client.Get(ctx, types.NamespacedName{
		Namespace: cr.GetNamespace(),
		Name:      cr.GetTopologyName()}, topo); err != nil {
		// can happen when the resource is not found
		cr.SetStatus("down")
		cr.SetReason("topology not found")
		return nil, errors.Wrap(err, "topology not found")
	}
	if topo.GetCondition(topov1alpha1.ConditionKindReady).Status != corev1.ConditionTrue {
		cr.SetStatus("down")
		cr.SetReason("topology not found or ready")
		return nil, errors.New("topology not ready")
	}

	// topology found

	if err := r.handleStatus(ctx, cr, topo); err != nil {
		return nil, err
	}

	if err := r.setPlatform(ctx, cr, topo); err != nil {
		return nil, err
	}

	cr.SetOrganization(cr.GetOrganization())
	cr.SetDeployment(cr.GetDeployment())
	cr.SetAvailabilityZone(cr.GetAvailabilityZone())
	cr.SetTopologyName(cr.GetTopologyName())

	return make(map[string]string), nil
}

func (r *application) handleStatus(ctx context.Context, cr topov1alpha1.Tn, topo topov1alpha1.Tp) error {
	if cr.GetAdminState() == "disable" {
		cr.SetStatus("down")
		cr.SetReason("admin disabled")
	} else {
		cr.SetStatus("up")
		cr.SetReason("")
	}
	return nil
}

func (r *application) setPlatform(ctx context.Context, cr topov1alpha1.Tn, topo topov1alpha1.Tp) error {
	r.log.Debug("Setflatform", "platform", cr.GetPlatform())
	if cr.GetPlatform() == "" && cr.GetPosition() != topov1alpha1.NodePositionServer.String() {
		// platform is not defined at node level
		p := topo.GetPlatformByKindName(cr.GetKindName())
		if p != "" {
			cr.SetPlatform(p)
			return nil
		}
		p = topo.GetPlatformFromDefaults()
		if p != "" {
			cr.SetPlatform(p)
			return nil
		}
		// platform is not defined we use the global default
		cr.SetPlatform("ixrd2")
		return nil

	}
	// all good since the platform is already set
	return nil
}
