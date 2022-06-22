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
	cat := catalog.Default
	runPodKey := catalog.Key{Name: "run_pod", Version: "latest"}
	topoKey := catalog.Key{Name: "configure_topology", Version: "latest"}
	nodeKey := catalog.Key{Name: "configure_node", Version: "latest"}
	configLLDPKey := catalog.Key{Name: "configure_lldp", Version: "latest"}
	stateLLDPKey := catalog.Key{Name: "state_lldp", Version: "latest"}

	runPodEntry, err := cat.Get(runPodKey)
	if err != nil {
		return err
	}
	topoFnEntry, err := cat.Get(topoKey)
	if err != nil {
		return err
	}
	nodeFnEntry, err := cat.Get(nodeKey)
	if err != nil {
		return err
	}
	configLLDPFnEntry, err := cat.Get(configLLDPKey)
	if err != nil {
		return err
	}
	stateLLDPFnEntry, err := cat.Get(stateLLDPKey)
	if err != nil {
		return err
	}

	c := resource.ClientApplicator{
		Client:     mgr.GetClient(),
		Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
	}

	ac := appcontext.New(
		appcontext.WithClient(c),
		appcontext.WithLogging(nddcopts.Logger.WithValues("appcontext", name)),
		appcontext.WithResourceFn(topov1alpha1.TopologyGroupVersionKind.String(), appcontext.GvkTypeSpecific, topoKey, topoFnEntry.RenderRn, topoFnEntry.ResourceFn, topoFnEntry.ResourceListFn),
		appcontext.WithResourceFn(topov1alpha1.NodeGroupVersionKind.String(), appcontext.GvkTypeSpecific, nodeKey, nodeFnEntry.RenderRn, nodeFnEntry.ResourceFn, nodeFnEntry.ResourceListFn),
		appcontext.WithResourceFn("Config.config.yndd.io.v1alpha1", appcontext.GvkTypeGeneric, configLLDPKey, configLLDPFnEntry.RenderRn, nil, nil),
		appcontext.WithResourceFn(statev1alpha1.StateGroupVersionKind.String(), appcontext.GvkTypeSpecific, stateLLDPKey, stateLLDPFnEntry.RenderRn, stateLLDPFnEntry.ResourceFn, stateLLDPFnEntry.ResourceListFn),
		appcontext.WithResourceFn("Run.run.yndd.io.v1alpha1", appcontext.GvkTypeGeneric, runPodKey, runPodEntry.RenderRn, nil, nil),
	)

	r := managed.NewReconciler(mgr, resource.ManagedKind(topov1alpha1.DefinitionGroupVersionKind),
		managed.WithLogger(nddcopts.Logger.WithValues("controller", name)),
		managed.WithApplogic(&applogic{
			log:    nddcopts.Logger.WithValues("applogic", name),
			client: c,
			ac:     ac,
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

	ac := r.ac.WithNameSpaceName(mr.GetName(), mr.GetNamespace())

	// render a topology
	ac.AddInstance(
		topov1alpha1.TopologyGroupVersionKind.String(),
		&catalog.Input{ObjectMeta: cr.ObjectMeta},
	)

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
			// initialize meta with the context info
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

			// render config
			ac.AddInstance(
				"dummyConfig",
				&catalog.Input{Object: &t, ObjectMeta: meta},
			)

			// render state
			ac.AddInstance(
				statev1alpha1.StateGroupVersionKind.String(),
				&catalog.Input{Object: &t, ObjectMeta: meta},
			)

		}
	}
	// *** UPDATE SUBSCRIBER WITH CURRENT PATHS FROM STATE CR

	// *** LOOK AT LINK STATE FROM NATS/CACHE
	// -> query NATS per target and see which links are needed (last per subject)

	// perform diff and perform a transaction to the system
	if err := ac.Apply(ctx); err != nil {
		return err
	}

	return nil
}
