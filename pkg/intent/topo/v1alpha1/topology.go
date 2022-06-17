/*
Copyright 2021 Wim Henderickx.

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

//+kubebuilder:object:generate=true
package v1alpha1

import (
	"context"

	"github.com/yndd/app-runtime/pkg/app"
	"github.com/yndd/ndd-runtime/pkg/resource"
	topov1alpha1 "github.com/yndd/topology/apis/topo/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InitTopology(c resource.ClientApplicator) app.Object {
	return &topology{
		client: c,
	}
}

type topology struct {
	// k8s client
	client resource.ClientApplicator
}

func (x *topology) List(ctx context.Context, opts []client.ListOption) (resource.ManagedList, error) {
	ol := &topov1alpha1.TopologyList{}
	if err := x.client.List(ctx, ol, opts...); err != nil {
		return nil, err
	}
	return ol, nil
}
