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
	"encoding/json"
	cdv2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/landscaper-utils/deployutils/pkg/utils"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"
)

//go:embed resources/controllerdeployment.yaml
var rawDefaultControllerDeployment []byte

const (
	controllerDeploymentName = "provider-aws"
)

func applyControllerDeployment(ctx context.Context, log logr.Logger, clt client.Client, controllerDeployment *v1beta1.ControllerDeployment) error {
	log.Info("Applying controller deployment")

	c := emptyControllerDeployment()
	_, err := controllerutil.CreateOrUpdate(ctx, clt, c, func() error {
		c.Type = controllerDeployment.Type
		c.ProviderConfig = controllerDeployment.ProviderConfig
		return nil
	})

	return err
}

func deleteControllerDeployment(ctx context.Context, log logr.Logger, clt client.Client) error {
	log.Info("Deleting controller deployment")

	c := emptyControllerDeployment()
	if err := clt.Delete(ctx, c); client.IgnoreNotFound(err) != nil {
		return err
	}

	return nil
}

func constructControllerDeployment(o *utils.Options, imports *Imports) (*v1beta1.ControllerDeployment, error) {
	const(
		resourceNameAlpine       = "alpine"
		resourceNamePause        = "pause"
		resourceNameProviderAWS  = "gardener-extension-provider-aws"
		componentNameProviderAWS = "github.com/gardener/gardener-extension-provider-aws"
	)

	// read defaults
	controllerDeployment := &v1beta1.ControllerDeployment{}
	if err := yaml.Unmarshal(rawDefaultControllerDeployment, controllerDeployment); err != nil {
		return nil, err
	}

	providerConfig := struct {
		Chart  string                 `json:"chart,omitempty"`
		Values map[string]interface{} `json:"values,omitempty"`
	}{}
	if err := json.Unmarshal(controllerDeployment.ProviderConfig.Raw, &providerConfig); err != nil {
		return nil, err
	}

	providerConfig.Chart = imports.ControllerRegistration.Chart

	cdProviderAWS, err := o.GetComponentDescriptorByName(componentNameProviderAWS)
	if err != nil {
		return nil, err
	}

	repository, tag, err := o.GetResourceOCIRepositoryAndTag(cdProviderAWS, resourceNameProviderAWS)
	if err != nil {
		return nil, err
	}

	providerConfig.Values["image"] = map[string]string{
		"repository": repository,
		"tag": tag,
	}

	images, err := getImageReferences(o, cdProviderAWS, resourceNameAlpine, resourceNamePause)
	if err != nil {
		return nil, err
	}
	providerConfig.Values["images"] = images

	if imports.ControllerRegistration.ConcurrentSyncs > 0 {
		providerConfig.Values["controllers"] = newControllersConfig(imports.ControllerRegistration.ConcurrentSyncs)
	}

	if imports.ControllerRegistration.Resources != nil {
		providerConfig.Values["resources"] = imports.ControllerRegistration.Resources
	}

	if imports.ControllerRegistration.VPA != nil {
		providerConfig.Values["vpa"] = imports.ControllerRegistration.VPA
	}

	if len(imports.ControllerRegistration.ImageVectorOverwrite) > 0 {
		providerConfig.Values["imageVectorOverwrite"] = imports.ControllerRegistration.ImageVectorOverwrite
	}

	rawConfig, err := json.Marshal(providerConfig)
	if err != nil {
		return nil, err
	}

	controllerDeployment.ProviderConfig = runtime.RawExtension{
		Raw: rawConfig,
	}

	return controllerDeployment, nil
}

func getImageReferences(o *utils.Options, cd *cdv2.ComponentDescriptor , resourceNames ...string) (map[string]string, error) {
	imageRefs := map[string]string{}

	for _, resourceName := range resourceNames {
		imageRef, err := o.GetResourceOCIImageReference(cd, resourceName)
		if err != nil {
			return nil, err
		}

		imageRefs[resourceName] = imageRef
	}

	return imageRefs, nil
}

func newControllersConfig(concurrentSyncs int) map[string]interface{} {
	return map[string]interface{}{
		"backupentry":    map[string]int{"concurrentSyncs": concurrentSyncs},
		"controlplane":   map[string]int{"concurrentSyncs": concurrentSyncs},
		"dnsrecord":      map[string]int{"concurrentSyncs": concurrentSyncs},
		"healthcheck":    map[string]int{"concurrentSyncs": concurrentSyncs},
		"infrastructure": map[string]int{"concurrentSyncs": concurrentSyncs},
		"worker":         map[string]int{"concurrentSyncs": concurrentSyncs},
	}
}

func emptyControllerDeployment() *v1beta1.ControllerDeployment {
	return &v1beta1.ControllerDeployment{ObjectMeta: metav1.ObjectMeta{Name: controllerDeploymentName}}
}
