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

package controllers

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/yndd/ndd-target-runtime/pkg/shared"
	"github.com/yndd/topology/internal/controllers/definition"
	"github.com/yndd/topology/internal/controllers/link"
	"github.com/yndd/topology/internal/controllers/node"
	"github.com/yndd/topology/internal/controllers/topology"
)

// Setup package controllers.
func Setup(mgr ctrl.Manager, nddcopts *shared.NddControllerOptions) error {
	for _, setup := range []func(ctrl.Manager, *shared.NddControllerOptions) error{
		definition.Setup,
		topology.Setup,
		link.Setup,
		node.Setup,
	} {
		if err := setup(mgr, nddcopts); err != nil {
			return err
		}
	}

	return nil
}
