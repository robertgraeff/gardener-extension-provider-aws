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

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"
)

//go:embed resources/controllerregistration.yaml
var rawDefaultControllerRegistration []byte

const (
	controllerRegistrationName = "provider-aws"
)

func applyControllerRegistration(ctx context.Context, log logr.Logger, clt client.Client, controllerRegistration *v1beta1.ControllerRegistration) error {
	log.Info("Applying controller registration")

	c := emptyControllerRegistration()
	_, err := controllerutil.CreateOrUpdate(ctx, clt, c, func() error {
		c.Spec = controllerRegistration.Spec
		return nil
	})

	return err
}

func deleteControllerRegistration(ctx context.Context, log logr.Logger, clt client.Client) error {
	log.Info("Deleting controller registration")

	c := emptyControllerRegistration()
	if err := clt.Delete(ctx, c); client.IgnoreNotFound(err) != nil {
		return err
	}

	return nil
}

func constructControllerRegistration(
	log logr.Logger,
	imports *Imports,
) (*v1beta1.ControllerRegistration, error) {
	log.Info("Constructing controller registration")

	controllerRegistration := &v1beta1.ControllerRegistration{}
	if err := yaml.Unmarshal(rawDefaultControllerRegistration, controllerRegistration); err != nil {
		return nil, err
	}

	if len(imports.ControllerRegistration.ControllerResources) > 0 {
		controllerRegistration.Spec.Resources = imports.ControllerRegistration.ControllerResources
	}

	return controllerRegistration, nil
}

func emptyControllerRegistration() *v1beta1.ControllerRegistration {
	return &v1beta1.ControllerRegistration{ObjectMeta: metav1.ObjectMeta{Name: controllerRegistrationName}}
}
