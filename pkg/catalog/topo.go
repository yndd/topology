package catalog

import (
	"github.com/yndd/catalog"
	_ "github.com/yndd/catalog/vendors/all"
	"github.com/yndd/ndd-runtime/pkg/resource"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	catalog.RegisterFns(catalog.Default, Fns)
}

var Fns = map[catalog.FnKey]catalog.Fn{
	{
		Name:      "configure_definition",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeUnknown,
		Platform:  "",
		SwVersion: "",
	}: ConfigureDefinition,
	{
		Name:      "configure_template",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeUnknown,
		Platform:  "",
		SwVersion: "",
	}: ConfigureTemplate,
	{
		Name:      "configure_topology",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeUnknown,
		Platform:  "",
		SwVersion: "",
	}: ConfigureTopology,
	{
		Name:      "configure_node",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeUnknown,
		Platform:  "",
		SwVersion: "",
	}: ConfigureNode,
	{
		Name:      "configure_link",
		Version:   "latest",
		Vendor:    targetv1.VendorTypeUnknown,
		Platform:  "",
		SwVersion: "",
	}: ConfigureLink,
}

func ConfigureDefinition(in *catalog.Input) (resource.Managed, error) {
	return &topov1alpha1.Definition{}, nil
}

func ConfigureTemplate(in *catalog.Input) (resource.Managed, error) {
	return &topov1alpha1.Template{}, nil
}

func ConfigureTopology(in *catalog.Input) (resource.Managed, error) {
	return &topov1alpha1.Topology{
		ObjectMeta: metav1.ObjectMeta{
			Name:      in.ObjectMeta.Name,
			Namespace: in.ObjectMeta.Namespace,
		},
		Spec: topov1alpha1.TopologySpec{
			Properties: &topov1alpha1.TopologyProperties{
				Defaults: &topov1alpha1.TopologyDefaults{
					NodeProperties: &topov1alpha1.NodeProperties{
						Position: topov1alpha1.PositionInfra,
					},
				},
				VendorTypeInfo: []*topov1alpha1.NodeProperties{
					{
						VendorType: targetv1.VendorTypeNokiaSRL,
						Platform:   "7220 IXR-D2",
						Position:   topov1alpha1.PositionInfra,
					},
					{
						VendorType: targetv1.VendorTypeNokiaSROS,
						Platform:   "7750 SR1",
						Position:   topov1alpha1.PositionInfra,
					},
				},
			},
		},
	}, nil
}

func ConfigureNode(in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
	return &topov1alpha1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:      in.ObjectMeta.Name, // how to do cr and target
			Namespace: in.ObjectMeta.Namespace,
		},
		Spec: topov1alpha1.NodeSpec{
			Properties: &topov1alpha1.NodeProperties{
				VendorType: t.GetDiscoveryInfo().VendorType,
				Platform:   t.GetDiscoveryInfo().Platform,
				//Index:
				//Position:
				// Tags://
			},
		},
	}, nil

}

func ConfigureLink(in *catalog.Input) (resource.Managed, error) {
	return &topov1alpha1.Link{
		ObjectMeta: metav1.ObjectMeta{
			Name:      in.ObjectMeta.Name, // how to do cr and target
			Namespace: in.ObjectMeta.Namespace,
		},
		Spec: topov1alpha1.LinkSpec{
			Properties: &topov1alpha1.LinkProperties{},
		},
	}, nil
}
