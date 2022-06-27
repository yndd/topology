package catalog

import (
	"strings"

	"github.com/yndd/catalog"
	_ "github.com/yndd/catalog/vendors/all"
	"github.com/yndd/ndd-runtime/pkg/resource"
	targetv1 "github.com/yndd/target/apis/target/v1"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	catalog.RegisterEntries(catalog.Default, Entries)
}

var Entries = map[catalog.Key]catalog.Entry{
	{
		Name:    "configure_definition",
		Version: "latest",
		//Vendor:    targetv1.VendorTypeUnknown,
		//Platform:  "",
		//SwVersion: "",
	}: {
		RenderFn: ConfigureDefinition,
		ResourceFn: func() resource.Managed {
			return &topov1alpha1.Definition{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &topov1alpha1.DefinitionList{}
		},
		MergeFn: nil, // TODO
	},
	{
		Name:    "configure_template",
		Version: "latest",
		//Vendor:    targetv1.VendorTypeUnknown,
		//Platform:  "",
		//SwVersion: "",
	}: {
		RenderFn: ConfigureTemplate,
		ResourceFn: func() resource.Managed {
			return &topov1alpha1.Template{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &topov1alpha1.TemplateList{}
		},
		MergeFn: nil, // TODO
	},
	{
		Name:    "configure_topology",
		Version: "latest",
		//Vendor:    targetv1.VendorTypeUnknown,
		//Platform:  "",
		//SwVersion: "",
	}: {
		RenderFn: ConfigureTopology,
		ResourceFn: func() resource.Managed {
			return &topov1alpha1.Topology{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &topov1alpha1.TopologyList{}
		},
		MergeFn: nil, // TODO
	},
	{
		Name:    "configure_node",
		Version: "latest",
		//Vendor:    targetv1.VendorTypeUnknown,
		//Platform:  "",
		//SwVersion: "",
	}: {
		RenderFn: ConfigureNode,
		ResourceFn: func() resource.Managed {
			return &topov1alpha1.Node{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &topov1alpha1.NodeList{}
		},
		MergeFn: nil, // TODO
	},
	{
		Name:    "configure_link",
		Version: "latest",
		//Vendor:    targetv1.VendorTypeUnknown,
		//Platform:  "",
		//SwVersion: "",
	}: {
		RenderFn: ConfigureLink,
		ResourceFn: func() resource.Managed {
			return &topov1alpha1.Link{}
		},
		ResourceListFn: func() resource.ManagedList {
			return &topov1alpha1.LinkList{}
		},
		MergeFn: nil, // TODO
	},
}

func ConfigureDefinition(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	return &topov1alpha1.Definition{}, nil
}

func ConfigureTemplate(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	return &topov1alpha1.Template{}, nil
}

func ConfigureTopology(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	return &topov1alpha1.Topology{
		ObjectMeta: in.ObjectMeta,
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

func ConfigureNode(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
	t, err := in.GetTarget()
	if err != nil {
		return nil, err
	}
	in.ObjectMeta.Name = strings.Join([]string{in.ObjectMeta.Name, t.GetName()}, ".")
	return &topov1alpha1.Node{
		ObjectMeta: in.ObjectMeta,
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

func ConfigureLink(key catalog.Key, in *catalog.Input) (resource.Managed, error) {
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
