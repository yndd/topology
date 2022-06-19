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
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yndd/app-runtime/pkg/app"
	"github.com/yndd/app-runtime/pkg/appcontext"
	"github.com/yndd/app-runtime/pkg/reconciler/managed"
	"github.com/yndd/catalog"
	discoveryv1alphav1 "github.com/yndd/discovery/api/v1alpha1"
	"github.com/yndd/ndd-runtime/pkg/event"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/meta"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/shared"
	statev1alpha1 "github.com/yndd/state/apis/state/v1alpha1"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	// timers
	//reconcileTimeout = 1 * time.Minute
	//veryShortWait    = 1 * time.Second
	// errors
	errUnexpectedResource = "unexpected object"
	errGetK8sResource     = "cannot get organization resource"
)

// Setup adds a controller that reconciles infra.
func Setup(mgr ctrl.Manager, nddcopts *shared.NddControllerOptions) error {
	name := strings.Join([]string{topov1alpha1.Group, strings.ToLower(topov1alpha1.DefinitionKind)}, "/")
	//tlfn := func() topov1alpha1.Tl { return &topov1alpha1.TopologyLink{} }
	//tllfn := func() topov1alpha1.TlList { return &topov1alpha1.TopologyLinkList{} }
	//tpfn := func() topov1alpha1.Tp { return &topov1alpha1.Topology{} }
	//tlfn := func() targetv1.TgList { return &targetv1.TargetList{} }
	cat := catalog.Default
	topoFn, err := cat.GetFn(catalog.FnKey{Name: "configure_topology", Version: "latest", Vendor: targetv1.VendorTypeUnknown})
	if err != nil {
		return err
	}
	nodeFn, err := cat.GetFn(catalog.FnKey{Name: "configure_node", Version: "latest", Vendor: targetv1.VendorTypeUnknown})
	if err != nil {
		return err
	}
	configLLDPFn, err := cat.GetFn(catalog.FnKey{Name: "configure_lldp", Version: "latest", Vendor: targetv1.VendorTypeUnknown})
	if err != nil {
		return err
	}
	stateLLDPFn, err := cat.GetFn(catalog.FnKey{Name: "state_lldp", Version: "latest", Vendor: targetv1.VendorTypeUnknown})
	if err != nil {
		return err
	}
	topoResourceFn := func() resource.Managed { return &topov1alpha1.Topology{} }
	topoResourceListFn := func() resource.ManagedList { return &topov1alpha1.TopologyList{} }
	nodeResourceFn := func() resource.Managed { return &topov1alpha1.Node{} }
	nodeResourceListFn := func() resource.ManagedList { return &topov1alpha1.NodeList{} }
	stateResourceFn := func() resource.Managed { return &statev1alpha1.State{} }
	stateResourceListFn := func() resource.ManagedList { return &statev1alpha1.StateList{} }

	c := resource.ClientApplicator{
		Client:     mgr.GetClient(),
		Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
	}

	ac := appcontext.New(
		appcontext.WithClient(c),
		appcontext.WithLogging(nddcopts.Logger.WithValues("appcontext", name)),
		appcontext.WithResourceFn(topov1alpha1.TopologyGroupVersionKind.String(), topoFn, topoResourceFn, topoResourceListFn),
		appcontext.WithResourceFn(topov1alpha1.NodeGroupVersionKind.String(), nodeFn, nodeResourceFn, nodeResourceListFn),
		appcontext.WithResourceFn(topov1alpha1.NodeGroupVersionKind.String(), configLLDPFn, nodeResourceFn, nodeResourceListFn),
		appcontext.WithResourceFn(statev1alpha1.StateGroupVersionKind.String(), stateLLDPFn, stateResourceFn, stateResourceListFn),
	)

	r := managed.NewReconciler(mgr, resource.ManagedKind(topov1alpha1.DefinitionGroupVersionKind),
		managed.WithLogger(nddcopts.Logger.WithValues("controller", name)),
		managed.WithApplogic(&applogic{
			log:    nddcopts.Logger.WithValues("applogic", name),
			client: c,
			/*
				topoFn:       topoFn,
				nodeFn:       nodeFn,
				configLLDPFn: configLLDPFn,
				stateLLDPFn:  stateLLDPFn,
				m:            sync.Mutex{},
				intents:      make(map[string]*intent.Compositeintent),
			*/
			ac: ac,
		}),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	)

	targetHandler := &EnqueueRequestForAllTargets{
		client: mgr.GetClient(),
		log:    nddcopts.Logger,
		ctx:    context.Background(),
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(nddcopts.Copts).
		For(&topov1alpha1.Definition{}).
		Owns(&topov1alpha1.Definition{}).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Watches(&source.Kind{Type: &targetv1.Target{}}, targetHandler).
		Complete(r)
}

type applogic struct {
	client resource.ClientApplicator
	log    logging.Logger

	//newTargetList func() targetv1.TgList

	/*
		// Functions used in the app
		topoFn       func(in *catalog.Input) (resource.Managed, error)
		nodeFn       func(in *catalog.Input) (resource.Managed, error)
		configLLDPFn func(in *catalog.Input) (resource.Managed, error)
		stateLLDPFn  func(in *catalog.Input) (resource.Managed, error)

		m       sync.Mutex
		intents map[string]*intent.Compositeintent
		//abstractions map[string]*abstraction.Compositeabstraction
	*/

	ac appcontext.AppContext
}

func (r *applogic) Initialize(ctx context.Context, mr resource.Managed) error {
	return nil
}

func (r *applogic) Update(ctx context.Context, mr resource.Managed) (map[string]string, error) {
	if err := r.populateSchema(ctx, mr); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *applogic) FinalUpdate(ctx context.Context, mr resource.Managed) {
}

func (r *applogic) Timeout(ctx context.Context, mr resource.Managed) time.Duration {
	return 0
}

func (r *applogic) Delete(ctx context.Context, mr resource.Managed) (bool, error) {
	return true, nil
}

func (r *applogic) FinalDelete(ctx context.Context, mr resource.Managed) {
}

func (r *applogic) populateSchema(ctx context.Context, mr resource.Managed) error {
	// cast the type to the real object/resource we expect
	cr, ok := mr.(*topov1alpha1.Definition)
	if !ok {
		return errors.New(errUnexpectedResource)
	}

	log := r.log.WithValues("crName", cr.GetNamespacedName())
	log.Debug("populateSchema")

	//r.intents[crName] = intent.New(r.client, crName)
	//r.abstractions[crName] = abstraction.New(r.client, crName)

	/*
		ci := app.New(&app.IntentInfo{
			Log:       r.log,
			Name:      cr.GetName(),
			Namespace: cr.GetNamespace(),
			Client:    r.client,
		})
	*/

	ac := r.ac.Clone(mr.GetName(), mr.GetNamespace())

	// +++++ PRE-PROCESSING  +++++
	// +++++ GET RESOURCES  +++++
	// +++++ CREATE INTENT +++++

	// create a topology

	ac.AddInstance(
		topov1alpha1.TopologyGroupVersionKind.String(),
		&catalog.Input{ObjectMeta: cr.ObjectMeta},
	)

	/*
		topo, err := r.topoFn(&catalog.Input{ObjectMeta: cr.ObjectMeta})
		if err != nil {
			log.Debug("cannot create topology", "error", err)
			return err
		}

		ci.AddEntry(
			topov1alpha1.TopologyGroupVersionKind.String(),
			topo.GetName(),
			topointentv1alpha1.InitTopology(r.client),
			topo,
		)
	*/

	// +++++ BREAKDOWN  +++++
	// per discovery rule check if the discovery rule matches within the namespace
	for _, dr := range cr.Spec.Properties.DiscoveryRules {

		namespace, drName := meta.NamespacedName(dr.NamespacedName).GetNameAndNamespace()
		opts := []client.ListOption{
			client.MatchingLabels{LabelKeyDiscoveryRule: drName},
			client.InNamespace(namespace),
		}
		// get targets in the namespace based on the discovery rule
		tl := &targetv1.TargetList{}
		r.client.List(ctx, tl, opts...)

		for _, t := range tl.Items {
			log.WithValues("target", t.GetName())
			// create a node

			// +++++ DEFINE CHILD INTENT +++++
			// +++++ NODE INTENT   - VENDOR AGNOSTIC +++++
			// +++++ STATE INTENT  - VENDOR SPECIFIC +++++
			// +++++ CONFIG INTENT - VENDOR SPECIFIC +++++

			meta := app.GetMeta(
				cr,
				map[string]string{
					discoveryv1alphav1.LabelKeyDiscoveryRule: drName,
				},
				map[string]string{},
			)

			// render node
			ac.AddInstance(
				topov1alpha1.NodeGroupVersionKind.String(),
				&catalog.Input{Object: &t, ObjectMeta: meta},
			)

			/*
				n, err := r.nodeFn(&catalog.Input{Object: &t, ObjectMeta: meta})
				if err != nil {
					log.Debug("cannot create node", "error", err)
					return err

				}

				ci.AddEntry(
					topov1alpha1.NodeGroupVersionKind.String(),
					n.GetName(),
					topointentv1alpha1.InitNode(r.client),
					n,
				)
			*/

			// render config
			ac.AddInstance(
				"dummyConfig",
				&catalog.Input{Object: &t, ObjectMeta: meta},
			)
			/*
				c, err := r.configLLDPFn(&catalog.Input{Object: &t, ObjectMeta: meta})
				if err != nil {
					log.Debug("cannot create lldp config", "error", err)
					return err
				}
				if err := r.client.Apply(ctx, c); err != nil {
					return err
				}
			*/

			// render state
			ac.AddInstance(
				statev1alpha1.StateGroupVersionKind.String(),
				&catalog.Input{Object: &t, ObjectMeta: meta},
			)
			/*
				s, err := r.stateLLDPFn(&catalog.Input{Object: &t, ObjectMeta: meta})
				if err != nil {
					log.Debug("cannot create lldp state", "error", err)
					return err
				}
				if err := r.client.Apply(ctx, s); err != nil {
					return err
				}
			*/

		}
	}
	// *** UPDATE SUBSCRIBER WITH CURRENT PATHS FROM STATE CR

	// *** LOOK AT LINK STATE FROM NATS/CACHE
	// -> query NATS per target and see which links are needed (last per subject)

	// **** COLLECTED ALL based on current k8s state               *****
	// **** DELTA - REMOVE THE ONCE NO LONGER NEEDED               *****
	// **** TRANSACTION                                            *****
	if err := ac.Apply(ctx); err != nil {
		return err
	}

	return nil
}
