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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/yndd/app-runtime/pkg/intent"
	"github.com/yndd/app-runtime/pkg/reconciler/managed"
	"github.com/yndd/ndd-runtime/pkg/event"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/meta"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"github.com/yndd/ndd-runtime/pkg/shared"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
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

	c := resource.ClientApplicator{
		Client:     mgr.GetClient(),
		Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
	}

	r := managed.NewReconciler(mgr, resource.ManagedKind(topov1alpha1.DefinitionGroupVersionKind),
		managed.WithLogger(nddcopts.Logger.WithValues("controller", name)),
		managed.WithApplogic(&applogic{
			log:    nddcopts.Logger.WithValues("applogic", name),
			client: c,
			//newTargetList: tlfn,
			intents: make(map[string]*intent.Compositeintent),
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

	m       sync.Mutex
	intents map[string]*intent.Compositeintent
	//abstractions map[string]*abstraction.Compositeabstraction
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
	crName := cr.GetNamespacedName()

	log := r.log.WithValues("crName", crName)
	log.Debug("populateSchema")

	r.intents[crName] = intent.New(r.client, crName)
	//r.abstractions[crName] = abstraction.New(r.client, crName)

	// +++++ PRE-PROCESSING  +++++
	// +++++ GET RESOURCES  +++++
	// +++++ CREATE INTENT +++++

	// create a topology
	topo := renderTopology(cr)
	if err := r.client.Apply(ctx, topo); err != nil {
		return err
	}

	// per template create the fabric
	for _, dt := range cr.Spec.Properties.Templates {
		log.Debug("NamespacedName input", "dt.NamespacedName", dt.NamespacedName)
		name, namespace := meta.NamespacedName(dt.NamespacedName).GetNameAndNamespace()
		log.Debug("NamespacedName output", "namespace", namespace, "name", name)
		tmpl := &topov1alpha1.Template{}
		if err := r.client.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		}, tmpl); err != nil {
			// template not defined
			return err
		}
		if err := r.createFabric(ctx, cr, tmpl); err != nil {
			return err
		}
	}

	// +++++ BREAKDOWN  +++++
	// per discovery rule check if the discovery rule matches within the namespace
	for _, dr := range cr.Spec.Properties.DiscoveryRules {

		name, namespace := meta.NamespacedName(dr.NamespacedName).GetNameAndNamespace()
		opts := []client.ListOption{
			client.MatchingLabels{LabelKeyDiscoveryRule: name},
			client.InNamespace(namespace),
		}
		// get targets in the namespace based on the discovery rule
		tl := &targetv1.TargetList{}
		r.client.List(ctx, tl, opts...)

		for _, t := range tl.Items {
			// create a node

			// +++++ CREATE CHILD INTENT +++++
			// +++++ NODE INTENT   - VENDOR AGNOSTIC +++++
			// +++++ STATE INTENT  - VENDOR SPECIFIC +++++
			// +++++ CONFIG INTENT - VENDOR SPECIFIC +++++

			/*
				ci := r.intents[crName]
				ci.AddChild(t.GetName(), intenttopov1alpha1.InitNode(r.client, ci, t.GetName()))
				node := ci.GetChildData(t.GetName())
				n, ok := node.(*topov1alpha1.TopologyNodeProperties)
				if !ok {
					return errors.New("expected ygot struct")
				}
				n.AdminState = "enable"
				n.VendorType = *t.GetDiscoveryInfo().VendorType
				n.Platform = *t.GetDiscoveryInfo().Kind
			*/

			n := renderNode(dr.NamespacedName, cr, &t)
			if err := r.client.Apply(ctx, n); err != nil {
				return err
			}

			switch n.Spec.Properties.VendorType {
			case targetv1.VendorTypeNokiaSRL:
				// populate data structure
			case targetv1.VendorTypeNokiaSROS:
			default:
				return fmt.Errorf("unsupported vendor type: %s", n.Spec.Properties.VendorType)
			}

			// create a state object per vendor type

		}
	}

	// **** COLLECT ALL                                  *****
	// **** FEEDBACK TO TOP LEVEL                        *****
	// **** SUBSCRIPTION WITH HANDLER (CREATE LINK/NODE) *****
	// **** TRANSACTION                                  *****

	return nil
}

func (r *applogic) createFabric(ctx context.Context, cr *topov1alpha1.Definition, tmpl *topov1alpha1.Template) error {
	crName := cr.GetNamespacedName()
	log := r.log.WithValues("crName", crName)
	log.Debug("createFabric...")

	totalPodNum := uint32(0)
	totalTier1Num := tmpl.Spec.Properties.Fabric.Tier1.NodeNumber
	totalTier2Info := map[uint32][]*topov1alpha1.FabricTierVendorInfo{}

	// n is number of pod definitions
	for p, pod := range tmpl.Spec.Properties.Fabric.Pods {
		totalPodNum += pod.PodNumber
		// i is the number of pods in a definition
		for i := uint32(0); i < pod.PodNumber; i++ {
			tier3NodePerPod := []*topov1alpha1.FabricTierVendorInfo{}
			totalTier2Info[(uint32(p)+1)*(i+1)] = []*topov1alpha1.FabricTierVendorInfo{}
			// kind is tier 2 or tier3
			for kind, tier := range pod.Tiers {
				vendorNum := len(tier.VendorInfo)
				switch kind {
				case "tier3":
					for n := uint32(0); n < tier.NodeNumber; n++ {
						vendorIdx := n % uint32(vendorNum)
						tier3NodePerPod = append(tier3NodePerPod, tier.VendorInfo[vendorIdx])

						//nodeNbr := n + 1
						//nodeName := fmt.Sprintf("pod%d-leaf%d", totalPodNum, nodeNbr)
						//log.Debug("create fabric node", "nodeName", nodeName)
						node := renderFabricNode(cr, &FabricNodeInfo{
							Position:   topov1alpha1.PositionLeaf,
							NodeIndex:  n + 1,
							PodIndex:   (uint32(p) + 1) * (i + 1),
							VendorInfo: tier.VendorInfo[vendorIdx],
						})
						if err := r.client.Apply(ctx, node); err != nil {
							return err
						}
					}
				case "tier2":

					for n := uint32(0); n < tier.NodeNumber; n++ {
						vendorIdx := n % uint32(vendorNum)
						totalTier2Info[(uint32(p)+1)*(i+1)] = append(totalTier2Info[(uint32(p)+1)*(i+1)],
							tier.VendorInfo[vendorIdx])
						//nodeNbr := n + 1
						//nodeName := fmt.Sprintf("pod%d-spine%d", totalPodNum, nodeNbr)
						//log.Debug("create fabric node", "nodeName", nodeName)

						node := renderFabricNode(cr, &FabricNodeInfo{
							Position:   topov1alpha1.PositionSpine,
							NodeIndex:  n + 1,
							PodIndex:   (uint32(p) + 1) * (i + 1),
							VendorInfo: tier.VendorInfo[vendorIdx],
						})
						if err := r.client.Apply(ctx, node); err != nil {
							return err
						}
					}
				default:
					return fmt.Errorf("wrong kind in the template definition: %s, value: %s, expected tier2 or tier3", tmpl.GetNamespacedName(), kind)
				}
			}

			// TODO vendorInfo
			for n := range totalTier2Info[(uint32(p)+1)*(i+1)] {
				for m := range tier3NodePerPod {
					tier2NodeNbr := n + 1
					tier3NodeNbr := m + 1
					linkName := fmt.Sprintf("pod%d-spine%d-int-1/%d-leaf%d-int-1/%d", totalPodNum, tier2NodeNbr, 1+m, tier3NodeNbr, tier2NodeNbr+26)
					log.Debug("create link", "linkName", linkName)
				}
			}
		}
	}

	// superspine/backbone
	// TODO vendorInfo
	for n := uint32(0); n < totalTier1Num; n++ {
		nodeNbr := n + 1
		nodeName := fmt.Sprintf("superspine%d", nodeNbr)
		log.Debug("create fabric node", "nodeName", nodeName)

		vendorIdx := n % uint32(len(tmpl.Spec.Properties.Fabric.Tier1.VendorInfo))
		node := renderFabricNode(cr, &FabricNodeInfo{
			Position:   topov1alpha1.PositionSuperspine,
			NodeIndex:  n + 1,
			VendorInfo: tmpl.Spec.Properties.Fabric.Tier1.VendorInfo[vendorIdx],
		})
		if err := r.client.Apply(ctx, node); err != nil {
			return err
		}
		actualTier2 := 0
		for p, tier2PodInfo := range totalTier2Info {
			for m := range tier2PodInfo {
				tier1NodeNbr := n + 1
				tier2NodeNbr := m + 1
				actualTier2++
				linkName := fmt.Sprintf("superspine%d-int-1/%d-pod%d-spine%d-int-1/%d", tier1NodeNbr, actualTier2, p, tier2NodeNbr, tier1NodeNbr+26)
				log.Debug("create link", "linkName", linkName)
			}
		}
	}

	return nil
}
