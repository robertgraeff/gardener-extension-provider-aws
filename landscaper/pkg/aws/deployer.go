// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/landscaper-utils/deployutils/pkg/utils"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Process(ctx context.Context, o *utils.Options) error {
	o.Log.Info("Processing component provider-aws")

	imports := &Imports{}
	err := o.ReadImports(imports)
	if err != nil {
		return err
	}

	o.Log.Info("Creating kubernetes client based on imported target")
	scheme := runtime.NewScheme()
	utilruntime.Must(v1beta1.AddToScheme(scheme))
	clientOptions := client.Options{Scheme: scheme}
	clt, err := utils.NewClientFromTarget(imports.Cluster, clientOptions)
	if err != nil {
		return err
	}

	switch o.Operation {
	case utils.OperationReconcile:
		return deploy(ctx, o.Log, clt, imports)

	case utils.OperationDelete:
		return undeploy(ctx, o.Log, clt)

	default:
		return fmt.Errorf("unknown operation: %q", o.Operation)
	}
}

func deploy(ctx context.Context, log logr.Logger, clt client.Client, imports *Imports) error {
	controllerDeployment, err := constructControllerDeployment(imports)
	if err != nil {
		return err
	}

	controllerRegistration, err := constructControllerRegistration(log, imports)
	if err != nil {
		return err
	}

	cloudProfile, err := constructCloudProfile(ctx, log, imports)
	if err != nil {
		return err
	}

	if err := applyControllerDeployment(ctx, log, clt, controllerDeployment); err != nil {
		return err
	}

	if err := applyControllerRegistration(ctx, log, clt, controllerRegistration); err != nil {
		return err
	}

	if err := applyCloudProfile(ctx, log, clt, cloudProfile); err != nil {
		return err
	}

	return nil
}

func undeploy(ctx context.Context, log logr.Logger, clt client.Client) error {
	if err := deleteCloudProfile(ctx, log, clt); err != nil {
		return err
	}

	if err := deleteControllerRegistration(ctx, log, clt); err != nil {
		return err
	}

	if err := deleteControllerDeployment(ctx, log, clt); err != nil {
		return err
	}

	return nil
}
