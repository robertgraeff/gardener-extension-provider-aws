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

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	mi "github.com/gardener/landscaper-utils/machineimages/pkg/machineimages"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/yaml"

	"github.com/gardener/gardener-extension-provider-aws/pkg/apis/aws/v1alpha1"
)

//go:embed resources/cloudprofile.yaml
var rawDefaultCloudProfile []byte

//go:embed resources/os_image_config.yaml
var rawMachineImagesProvider []byte

const (
	cloudProfileName = "aws"
)

func applyCloudProfile(ctx context.Context, log logr.Logger, clt client.Client, cloudProfile *v1beta1.CloudProfile) error {
	log.Info("Applying cloud profile")

	c := emptyCloudProfile()
	_, err := controllerutil.CreateOrUpdate(ctx, clt, c, func() error {
		c.Spec = cloudProfile.Spec
		return nil
	})

	return err
}

func deleteCloudProfile(ctx context.Context, log logr.Logger, clt client.Client) error {
	log.Info("Deleting cloud profile")

	c := emptyCloudProfile()
	if err := clt.Delete(ctx, c); client.IgnoreNotFound(err) != nil {
		return err
	}

	return nil
}

func constructCloudProfile(ctx context.Context, log logr.Logger, imports *Imports) (*v1beta1.CloudProfile, error) {
	log.Info("Constructing cloud profile")

	cloudProfile := &v1beta1.CloudProfile{}
	if err := yaml.Unmarshal(rawDefaultCloudProfile, cloudProfile); err != nil {
		return nil, err
	}

	machineImages, providerConfig, err := getMachineImages(ctx, log, imports)
	if err != nil {
		return nil, err
	}

	cloudProfile.Spec.MachineImages = machineImages
	cloudProfile.Spec.ProviderConfig = providerConfig
	cloudProfile.Spec.Kubernetes = v1beta1.KubernetesSettings{Versions: imports.CloudProfile.KubernetesVersions}
	if len(imports.CloudProfile.Regions) > 0 {
		cloudProfile.Spec.Regions = imports.CloudProfile.Regions
	}

	return cloudProfile, nil
}

func getMachineImages(ctx context.Context, log logr.Logger, imports *Imports) (
	[]v1beta1.MachineImage,
	*runtime.RawExtension,
	error,
) {
	providerSpecificMachineImages, err := getProviderSpecificMachineImageConfig()
	if err != nil {
		return nil, nil, err
	}

	machineImages, err := mi.ComputeMachineImages(
		ctx,
		log,
		imports.CloudProfile.MachineImages,
		imports.CloudProfile.MachineImagesLs,
		providerSpecificMachineImages,
		imports.CloudProfile.MachineImagesProviderLs,
		imports.CloudProfile.DisableMachineImages,
		imports.CloudProfile.IncludeFilters,
		imports.CloudProfile.ExcludeFilters,
	)
	if err != nil {
		return nil, nil, err
	}

	convertedMachineImages, err := convertMachineImages(machineImages)
	if err != nil {
		return nil, nil, err
	}

	providerConfig, err := constructProviderConfig(machineImages)
	if err != nil {
		return nil, nil, err
	}

	return convertedMachineImages, providerConfig, nil
}

func getProviderSpecificMachineImageConfig() ([]mi.MachineImage, error) {
	machineImages := []mi.MachineImage{}
	err := yaml.Unmarshal(rawMachineImagesProvider, &machineImages)
	if err != nil {
		return nil, err
	}
	return machineImages, nil
}

// convertMachineImages converts []machineimages.MachineImage into []v1beta1.MachineImage.
func convertMachineImages(machineImages []mi.MachineImage) ([]v1beta1.MachineImage, error) {
	data, err := json.Marshal(machineImages)
	if err != nil {
		return nil, err
	}

	result := []v1beta1.MachineImage{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func constructProviderConfig(machineImages []mi.MachineImage) (*runtime.RawExtension, error) {
	awsMachineImages, err := convertToAwsMachineImages(machineImages)
	if err != nil {
		return nil, err
	}

	cloudProfileConfig := &v1alpha1.CloudProfileConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CloudProfileConfig",
			APIVersion: "aws.provider.extensions.gardener.cloud/v1alpha1",
		},
		MachineImages: awsMachineImages,
	}

	cloudProfileConfigJSON, err := json.Marshal(cloudProfileConfig)
	if err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: cloudProfileConfigJSON,
	}, nil
}

// convertToAwsMachineImages converts []machineimages.MachineImage into aws specific []v1alpha1.MachineImage.
func convertToAwsMachineImages(machineImages []mi.MachineImage) ([]v1alpha1.MachineImages, error) {
	data, err := json.Marshal(machineImages)
	if err != nil {
		return nil, err
	}

	result := []v1alpha1.MachineImages{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func emptyCloudProfile() *v1beta1.CloudProfile {
	return &v1beta1.CloudProfile{ObjectMeta: metav1.ObjectMeta{Name: cloudProfileName}}
}
