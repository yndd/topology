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

package link

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/yndd/app-runtime/pkg/odns"
	"github.com/yndd/app-runtime/pkg/reconciler/managed"
	"github.com/yndd/ndd-runtime/pkg/event"
	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-runtime/pkg/resource"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/source"

	nddv1 "github.com/yndd/ndd-runtime/apis/common/v1"
	"github.com/yndd/ndd-runtime/pkg/shared"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
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
func Setup(mgr ctrl.Manager, nddcopts *shared.NddControllerOptions) error {
	name := "nddo/" + strings.ToLower(topov1alpha1.LinkGroupKind)
	//tlfn := func() topov1alpha1.Tl { return &topov1alpha1.TopologyLink{} }
	//tllfn := func() topov1alpha1.TlList { return &topov1alpha1.TopologyLinkList{} }
	//tpfn := func() topov1alpha1.Tp { return &topov1alpha1.Topology{} }

	c := resource.ClientApplicator{
		Client:     mgr.GetClient(),
		Applicator: resource.NewAPIPatchingApplicator(mgr.GetClient()),
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(topov1alpha1.LinkGroupVersionKind),
		managed.WithLogger(nddcopts.Logger.WithValues("controller", name)),
		managed.WithApplogic(&application{
			client: c,
			hooks:  NewHook(c, nddcopts.Logger.WithValues("nodehook", name)),
			log:    nddcopts.Logger.WithValues("applogic", name),
			//newTopology:     tpfn,
			//newTopologyLink: tlfn,
			//handler:         nddcopts.Handler,
		}),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
	)

	topologyHandler := &EnqueueRequestForAllTopologies{
		client: mgr.GetClient(),
		log:    nddcopts.Logger,
		ctx:    context.Background(),
		//handler:         nddcopts.Handler,
		//newTopoLinkList: tllfn,
	}

	topologyLinkHandler := &EnqueueRequestForAllTopologyLinks{
		client: mgr.GetClient(),
		log:    nddcopts.Logger,
		ctx:    context.Background(),
		//handler:         nddcopts.Handler,
		//newTopoLinkList: tllfn,
	}

	topologyNodeHandler := &EnqueueRequestForAllTopologyNodes{
		client: mgr.GetClient(),
		log:    nddcopts.Logger,
		ctx:    context.Background(),
		//handler:         nddcopts.Handler,
		//newTopoLinkList: tllfn,
	}

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(nddcopts.Copts).
		For(&topov1alpha1.Link{}).
		Owns(&topov1alpha1.Link{}).
		WithEventFilter(resource.IgnoreUpdateWithoutGenerationChangePredicate()).
		Watches(&source.Kind{Type: &topov1alpha1.Topology{}}, topologyHandler).
		Watches(&source.Kind{Type: &topov1alpha1.Node{}}, topologyNodeHandler).
		Watches(&source.Kind{Type: &topov1alpha1.Link{}}, topologyLinkHandler).
		Complete(r)
}

type application struct {
	client resource.ClientApplicator
	log    logging.Logger

	//newTopology     func() topov1alpha1.Tp
	//newTopologyLink func() topov1alpha1.Tl

	//handler handler.Handler
	hooks Hooks
}

func getCrName(cr *topov1alpha1.Link) string {
	return strings.Join([]string{cr.GetNamespace(), cr.GetName()}, ".")
}

func (r *application) Initialize(ctx context.Context, mr resource.Managed) error {
	/*
		cr, ok := mg.(*topov1alpha1.TopologyLink)
		if !ok {
			return errors.New(errUnexpectedResource)
		}

		if err := cr.InitializeResource(); err != nil {
			r.log.Debug("Cannot initialize", "error", err)
			return err
		}
	*/

	return nil
}

func (r *application) Update(ctx context.Context, mr resource.Managed) (map[string]string, error) {
	cr, ok := mr.(*topov1alpha1.Link)
	if !ok {
		return nil, errors.New(errUnexpectedResource)
	}

	return r.handleAppLogic(ctx, cr)

	//return nil, nil
}

func (r *application) FinalUpdate(ctx context.Context, mr resource.Managed) {
}

func (r *application) Timeout(ctx context.Context, mr resource.Managed) time.Duration {
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

func (r *application) Delete(ctx context.Context, mr resource.Managed) (bool, error) {
	/*
		cr, ok := mg.(*topov1alpha1.TopologyLink)
		if !ok {
			return true, errors.New(errUnexpectedResource)
		}
		if cr.GetLagMember() {
			topologyName := cr.GetTopologyName()
			logicalLink, err := r.hooks.Get(ctx, cr, topologyName)
			if err == nil {
				r.log.Debug("logical link exists", "Logical Link", logicalLink.GetName())
				//for the multi-homed case we need to delete the tags of the member links
				// that match the mh name
				if err := r.hooks.DeleteApply(ctx, cr, logicalLink); err != nil {
					r.log.Debug("Cannot delete tags of a logical link", "error", err)
					return true, err
				}
			}
		}
	*/
	return true, nil
}

func (r *application) FinalDelete(ctx context.Context, mr resource.Managed) {
	/*
		cr, ok := mg.(*topov1alpha1.TopologyLink)
		if !ok {
			return
		}
		crName := getCrName(cr)
		r.handler.Delete(crName)
	*/
}

func (r *application) handleAppLogic(ctx context.Context, cr *topov1alpha1.Link) (map[string]string, error) {
	log := r.log.WithValues("function", "handleAppLogic", "crname", cr.GetName())
	log.Debug("handleAppLogic")

	// initialize speedy
	//crName := getCrName(cr)
	//r.handler.Init(crName)

	// get the topo name which is the full name w/o the link info
	fullTopoName := odns.GetParentResourceName(cr.GetName())

	//topo := r.newTopology()
	topo := &topov1alpha1.Topology{}
	if err := r.client.Get(ctx, types.NamespacedName{
		Namespace: cr.GetNamespace(),
		Name:      fullTopoName}, topo); err != nil {
		// can happen when the resource is not found
		//cr.SetStatus("down")
		//cr.SetReason("topology not found")
		return nil, errors.Wrap(err, "topology not found")
	}
	if topo.GetCondition(nddv1.ConditionKindReady).Status != corev1.ConditionTrue {
		//cr.SetStatus("down")
		//cr.SetReason("topology not found or ready")
		return nil, errors.New("topology not ready")
	}

	// topology found and ready

	if err := r.handleStatus(ctx, cr, topo); err != nil {
		return nil, err
	}

	_, err := r.parseLink(ctx, cr, fullTopoName)
	if err != nil {
		return nil, err
	}

	return make(map[string]string), nil
}

func (r *application) handleStatus(ctx context.Context, cr *topov1alpha1.Link, topo *topov1alpha1.Topology) error {
	// topology found

	/*
		if topo.GetStatus() == "down" {
			cr.SetStatus("down")
			cr.SetReason("parent status down")
		} else {
			if cr.GetAdminState() == "disable" {
				cr.SetStatus("down")
				cr.SetReason("admin disable")
			} else {
				cr.SetStatus("up")
				cr.SetReason("")
			}
		}
	*/
	cr.SetOrganization(cr.GetOrganization())
	cr.SetDeployment(cr.GetDeployment())
	cr.SetAvailabilityZone(cr.GetAvailabilityZone())
	return nil
}

func (r *application) parseLink(ctx context.Context, cr *topov1alpha1.Link, fullTopoName string) (*string, error) {
	// parse link

	// validates if the nodes if the links are present in the k8s api are not
	// if an error occurs during validation an error is returned
	msg, err := r.validateNodes(ctx, cr)
	if err != nil {
		return msg, err
	}
	if msg != nil {
		return nil, fmt.Errorf("%s", *msg)
	}

	/*
		// for infra links we set the kind at the link level using the information from the spec
		if cr.GetEndPointAKind() == topov1alpha1.EndpointKindInfra && cr.GetEndPointBKind() == topov1alpha1.EndpointKindInfra {
			//TODO
			//cr.SetKind(topov1alpha1.LinkEPKindInfra.String())
		}

		if cr.GetLag() {
			// this is a logical link (single homes or multihomed), we dont need to process it since the member links take care
			// of crud operation
			cr.SetOrganization(cr.GetOrganization())
			cr.SetDeployment(cr.GetDeployment())
			cr.SetAvailabilityZone(cr.GetAvailabilityZone())
			//cr.SetTopologyName(cr.GetTopologyName())
			return nil, nil
		}

		// check if the link is part of a lag
		if cr.GetLagMember() {
			logicalLink, err := r.hooks.Get(ctx, cr, fullTopoName)
			if err != nil {
				if resource.IgnoreNotFound(err) != nil {
					return nil, err
				}
				if err := r.hooks.Create(ctx, cr, fullTopoName); err != nil {
					return nil, err
				}
				r.log.Debug("logical link created")
				cr.SetOrganization(cr.GetOrganization())
				cr.SetDeployment(cr.GetDeployment())
				cr.SetAvailabilityZone(cr.GetAvailabilityZone())
				//cr.SetTopologyName(cr.GetTopologyName())
				return nil, nil

			}
			r.log.Debug("logical link exists", "Logical Link", logicalLink.GetName())

			// for the multi-homed case we need to add the tags of the other member links
			// that match the mh name
			if err := r.hooks.Apply(ctx, cr, logicalLink); err != nil {
				return nil, err
			}

		}
		cr.SetOrganization(cr.GetOrganization())
		cr.SetDeployment(cr.GetDeployment())
		cr.SetAvailabilityZone(cr.GetAvailabilityZone())
		//cr.SetTopologyName(cr.GetTopologyName())
	*/
	return nil, nil
}

func (r *application) validateNodes(ctx context.Context, cr *topov1alpha1.Link) (*string, error) {
	/*
		for i := 0; i <= 1; i++ {
			var multihoming bool
			var nodeName string
			var tags map[string]string
			lag := cr.GetLag()
			switch i {
			case 0:
				nodeName = cr.GetEndpointANodeName()
				multihoming = cr.GetEndPointAMultiHoming()
				tags = cr.GetEndpointATag()
			case 1:
				nodeName = cr.GetEndpointBNodeName()
				multihoming = cr.GetEndPointBMultiHoming()
				tags = cr.GetEndpointBTag()
			}

			// lag are logical links which are created based on member links
			// for singlehomed logical links if the node no longer exists, we delete the sh-logical-link
			// for multi-homed logical links if a member node no longer exists, we delete the tags related to the node
			// for multi-homed logical links of all member nodes no longer exist, we delete the mh-logical link
			if lag {
				if multihoming {
					// node validation happens through the endpoint tags
					// a nodetag has a prefix of node:
					found := false
					for k, v := range tags {
						if strings.Contains(k, topov1alpha1.NodePrefix) {
							nodeName := strings.TrimPrefix(k, topov1alpha1.NodePrefix+":")

							fullNodeName := strings.Join([]string{odns.GetParentResourceName(cr.GetName()), nodeName}, ".")
							node := &topov1alpha1.Node{}
							if err := r.client.Get(ctx, types.NamespacedName{
								Namespace: cr.GetNamespace(),
								Name:      fullNodeName}, node); err != nil {
								if resource.IgnoreNotFound(err) != nil {
									return nil, err
								}
								r.log.Debug("mh-ep logical-link:: member node not found, delete the ep node tags", "nodeName", nodeName)
								// node no longer exists, we can delete the node tags from the logocal element
								if err := r.hooks.DeleteApplyNode(ctx, cr, 0, k, v); err != nil {
									return nil, err
								}
							} else {
								found = true
							}
						}
					}
					if !found {
						// when none of the mh nodes are found we can delete the logical link
						if err := r.hooks.Delete(ctx, cr); err != nil {
							return nil, err
						}
						r.log.Debug("mh-ep logical-link: none of the member nodes wwere found, delete the logical-link")
						return nil, nil
					}
				} else {
					fullNodeName := strings.Join([]string{odns.GetParentResourceName(cr.GetName()), nodeName}, ".")
					node := &topov1alpha1.Node{}
					if err := r.client.Get(ctx, types.NamespacedName{
						Namespace: cr.GetNamespace(),
						Name:      fullNodeName}, node); err != nil {
						if resource.IgnoreNotFound(err) != nil {
							return nil, err
						}
						r.log.Debug("sh-ep logical-link: node not found, delete the logical-link", "nodeName", nodeName)
						// node no longer exists, we can delete the logical element
						if err := r.hooks.Delete(ctx, cr); err != nil {
							return nil, err
						}
						// when delete is successfull we finish/return
						return nil, nil
					}
				}
			} else {
				// individual links
				fullNodeName := strings.Join([]string{odns.GetParentResourceName(cr.GetName()), nodeName}, ".")
				node := &topov1alpha1.Node{}
				if err := r.client.Get(ctx, types.NamespacedName{
					Namespace: cr.GetNamespace(),
					Name:      fullNodeName}, node); err != nil {
					r.log.Debug("individual link: node not found", "nodeName", nodeName)
					//cr.SetStatus("down")
					//cr.SetReason(fmt.Sprintf("node %d not found", i))
					return utils.StringPtr(fmt.Sprintf("node %d not found", i)), nil
				}
			}
		}
	*/
	return nil, nil
}
