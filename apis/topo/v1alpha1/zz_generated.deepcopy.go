//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/yndd/ndd-runtime/apis/common/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DiscoveryRule) DeepCopyInto(out *DiscoveryRule) {
	*out = *in
	out.Rule = in.Rule
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DiscoveryRule.
func (in *DiscoveryRule) DeepCopy() *DiscoveryRule {
	if in == nil {
		return nil
	}
	out := new(DiscoveryRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeAttributes) DeepCopyInto(out *NodeAttributes) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeAttributes.
func (in *NodeAttributes) DeepCopy() *NodeAttributes {
	if in == nil {
		return nil
	}
	out := new(NodeAttributes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Rule) DeepCopyInto(out *Rule) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Rule.
func (in *Rule) DeepCopy() *Rule {
	if in == nil {
		return nil
	}
	out := new(Rule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Template) DeepCopyInto(out *Template) {
	*out = *in
	out.Rule = in.Rule
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Template.
func (in *Template) DeepCopy() *Template {
	if in == nil {
		return nil
	}
	out := new(Template)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Topology) DeepCopyInto(out *Topology) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Topology.
func (in *Topology) DeepCopy() *Topology {
	if in == nil {
		return nil
	}
	out := new(Topology)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Topology) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyDefaults) DeepCopyInto(out *TopologyDefaults) {
	*out = *in
	if in.NodeAttributes != nil {
		in, out := &in.NodeAttributes, &out.NodeAttributes
		*out = new(NodeAttributes)
		**out = **in
	}
	if in.Tag != nil {
		in, out := &in.Tag, &out.Tag
		*out = make([]*v1.Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(v1.Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyDefaults.
func (in *TopologyDefaults) DeepCopy() *TopologyDefaults {
	if in == nil {
		return nil
	}
	out := new(TopologyDefaults)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyDefinition) DeepCopyInto(out *TopologyDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyDefinition.
func (in *TopologyDefinition) DeepCopy() *TopologyDefinition {
	if in == nil {
		return nil
	}
	out := new(TopologyDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyDefinitionList) DeepCopyInto(out *TopologyDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TopologyDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyDefinitionList.
func (in *TopologyDefinitionList) DeepCopy() *TopologyDefinitionList {
	if in == nil {
		return nil
	}
	out := new(TopologyDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyDefinitionProperties) DeepCopyInto(out *TopologyDefinitionProperties) {
	*out = *in
	if in.Templates != nil {
		in, out := &in.Templates, &out.Templates
		*out = make([]*Template, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Template)
				**out = **in
			}
		}
	}
	if in.DiscoveryRules != nil {
		in, out := &in.DiscoveryRules, &out.DiscoveryRules
		*out = make([]*DiscoveryRule, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(DiscoveryRule)
				**out = **in
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyDefinitionProperties.
func (in *TopologyDefinitionProperties) DeepCopy() *TopologyDefinitionProperties {
	if in == nil {
		return nil
	}
	out := new(TopologyDefinitionProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyDefinitionSpec) DeepCopyInto(out *TopologyDefinitionSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.Properties.DeepCopyInto(&out.Properties)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyDefinitionSpec.
func (in *TopologyDefinitionSpec) DeepCopy() *TopologyDefinitionSpec {
	if in == nil {
		return nil
	}
	out := new(TopologyDefinitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyDefinitionStatus) DeepCopyInto(out *TopologyDefinitionStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyDefinitionStatus.
func (in *TopologyDefinitionStatus) DeepCopy() *TopologyDefinitionStatus {
	if in == nil {
		return nil
	}
	out := new(TopologyDefinitionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyLink) DeepCopyInto(out *TopologyLink) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyLink.
func (in *TopologyLink) DeepCopy() *TopologyLink {
	if in == nil {
		return nil
	}
	out := new(TopologyLink)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyLink) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyLinkEndpoints) DeepCopyInto(out *TopologyLinkEndpoints) {
	*out = *in
	if in.InterfaceName != nil {
		in, out := &in.InterfaceName, &out.InterfaceName
		*out = new(string)
		**out = **in
	}
	if in.NodeName != nil {
		in, out := &in.NodeName, &out.NodeName
		*out = new(string)
		**out = **in
	}
	if in.Tag != nil {
		in, out := &in.Tag, &out.Tag
		*out = make([]*v1.Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(v1.Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyLinkEndpoints.
func (in *TopologyLinkEndpoints) DeepCopy() *TopologyLinkEndpoints {
	if in == nil {
		return nil
	}
	out := new(TopologyLinkEndpoints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyLinkList) DeepCopyInto(out *TopologyLinkList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TopologyLink, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyLinkList.
func (in *TopologyLinkList) DeepCopy() *TopologyLinkList {
	if in == nil {
		return nil
	}
	out := new(TopologyLinkList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyLinkList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyLinkProperties) DeepCopyInto(out *TopologyLinkProperties) {
	*out = *in
	if in.AdminState != nil {
		in, out := &in.AdminState, &out.AdminState
		*out = new(string)
		**out = **in
	}
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]*TopologyLinkEndpoints, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(TopologyLinkEndpoints)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Tag != nil {
		in, out := &in.Tag, &out.Tag
		*out = make([]*v1.Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(v1.Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyLinkProperties.
func (in *TopologyLinkProperties) DeepCopy() *TopologyLinkProperties {
	if in == nil {
		return nil
	}
	out := new(TopologyLinkProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyLinkSpec) DeepCopyInto(out *TopologyLinkSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.Properties.DeepCopyInto(&out.Properties)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyLinkSpec.
func (in *TopologyLinkSpec) DeepCopy() *TopologyLinkSpec {
	if in == nil {
		return nil
	}
	out := new(TopologyLinkSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyLinkStatus) DeepCopyInto(out *TopologyLinkStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyLinkStatus.
func (in *TopologyLinkStatus) DeepCopy() *TopologyLinkStatus {
	if in == nil {
		return nil
	}
	out := new(TopologyLinkStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyList) DeepCopyInto(out *TopologyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Topology, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyList.
func (in *TopologyList) DeepCopy() *TopologyList {
	if in == nil {
		return nil
	}
	out := new(TopologyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyNode) DeepCopyInto(out *TopologyNode) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyNode.
func (in *TopologyNode) DeepCopy() *TopologyNode {
	if in == nil {
		return nil
	}
	out := new(TopologyNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyNode) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyNodeList) DeepCopyInto(out *TopologyNodeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TopologyNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyNodeList.
func (in *TopologyNodeList) DeepCopy() *TopologyNodeList {
	if in == nil {
		return nil
	}
	out := new(TopologyNodeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyNodeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyNodeProperties) DeepCopyInto(out *TopologyNodeProperties) {
	*out = *in
	if in.Tag != nil {
		in, out := &in.Tag, &out.Tag
		*out = make([]*v1.Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(v1.Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyNodeProperties.
func (in *TopologyNodeProperties) DeepCopy() *TopologyNodeProperties {
	if in == nil {
		return nil
	}
	out := new(TopologyNodeProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyNodeSpec) DeepCopyInto(out *TopologyNodeSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	if in.Properties != nil {
		in, out := &in.Properties, &out.Properties
		*out = new(TopologyNodeProperties)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyNodeSpec.
func (in *TopologyNodeSpec) DeepCopy() *TopologyNodeSpec {
	if in == nil {
		return nil
	}
	out := new(TopologyNodeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyNodeStatus) DeepCopyInto(out *TopologyNodeStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyNodeStatus.
func (in *TopologyNodeStatus) DeepCopy() *TopologyNodeStatus {
	if in == nil {
		return nil
	}
	out := new(TopologyNodeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyProperties) DeepCopyInto(out *TopologyProperties) {
	*out = *in
	if in.Defaults != nil {
		in, out := &in.Defaults, &out.Defaults
		*out = new(TopologyDefaults)
		(*in).DeepCopyInto(*out)
	}
	if in.VendorTypeInfo != nil {
		in, out := &in.VendorTypeInfo, &out.VendorTypeInfo
		*out = make([]*VendorTypeInfo, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(VendorTypeInfo)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyProperties.
func (in *TopologyProperties) DeepCopy() *TopologyProperties {
	if in == nil {
		return nil
	}
	out := new(TopologyProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologySpec) DeepCopyInto(out *TopologySpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.Properties.DeepCopyInto(&out.Properties)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologySpec.
func (in *TopologySpec) DeepCopy() *TopologySpec {
	if in == nil {
		return nil
	}
	out := new(TopologySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyStatus) DeepCopyInto(out *TopologyStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyStatus.
func (in *TopologyStatus) DeepCopy() *TopologyStatus {
	if in == nil {
		return nil
	}
	out := new(TopologyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyTemplate) DeepCopyInto(out *TopologyTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyTemplate.
func (in *TopologyTemplate) DeepCopy() *TopologyTemplate {
	if in == nil {
		return nil
	}
	out := new(TopologyTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyTemplateList) DeepCopyInto(out *TopologyTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]TopologyTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyTemplateList.
func (in *TopologyTemplateList) DeepCopy() *TopologyTemplateList {
	if in == nil {
		return nil
	}
	out := new(TopologyTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *TopologyTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyTemplateProperties) DeepCopyInto(out *TopologyTemplateProperties) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyTemplateProperties.
func (in *TopologyTemplateProperties) DeepCopy() *TopologyTemplateProperties {
	if in == nil {
		return nil
	}
	out := new(TopologyTemplateProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyTemplateSpec) DeepCopyInto(out *TopologyTemplateSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	out.Properties = in.Properties
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyTemplateSpec.
func (in *TopologyTemplateSpec) DeepCopy() *TopologyTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(TopologyTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TopologyTemplateStatus) DeepCopyInto(out *TopologyTemplateStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TopologyTemplateStatus.
func (in *TopologyTemplateStatus) DeepCopy() *TopologyTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(TopologyTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VendorTypeInfo) DeepCopyInto(out *VendorTypeInfo) {
	*out = *in
	if in.NodeAttributes != nil {
		in, out := &in.NodeAttributes, &out.NodeAttributes
		*out = new(NodeAttributes)
		**out = **in
	}
	if in.Tag != nil {
		in, out := &in.Tag, &out.Tag
		*out = make([]*v1.Tag, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(v1.Tag)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VendorTypeInfo.
func (in *VendorTypeInfo) DeepCopy() *VendorTypeInfo {
	if in == nil {
		return nil
	}
	out := new(VendorTypeInfo)
	in.DeepCopyInto(out)
	return out
}
