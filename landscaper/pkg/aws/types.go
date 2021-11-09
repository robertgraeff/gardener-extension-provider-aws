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
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/gardener/landscaper-utils/machineimages/pkg/machineimages"
	lsv1alpha1 "github.com/gardener/landscaper/apis/core/v1alpha1"
)

// Imports
type Imports struct {
	Cluster                lsv1alpha1.Target      `json:"cluster" yaml:"cluster"`
	CloudProfile           CloudProfile           `json:"cloudProfile,omitempty"`
	ControllerDeployment   ControllerDeployment   `json:"controllerDeployment,omitempty"`
	ControllerRegistration ControllerRegistration `json:"controllerRegistration,omitempty"`
}

type CloudProfile struct {
	KubernetesVersions      []v1beta1.ExpirableVersion         `json:"kubernetesVersions,omitempty"`
	MachineImages           []machineimages.MachineImage       `json:"machineImages,omitempty"`
	MachineImagesLs         []machineimages.MachineImage       `json:"machineImagesLs,omitempty"`
	MachineImagesProviderLs []machineimages.MachineImage       `json:"machineImagesProviderLs,omitempty"`
	IncludeFilters          []machineimages.OsImagesFilterKind `json:"includeFilters,omitempty"`
	ExcludeFilters          []machineimages.OsImagesFilterKind `json:"excludeFilters,omitempty"`
	DisableMachineImages    []string                           `json:"disableMachineImages,omitempty"`
	Regions                 []v1beta1.Region                   `json:"regions,omitempty"`
}

type ControllerDeployment struct {
	// base64 encoded string of the tarred and gzipped gardener-extension-provider-aws chart
	Chart           string      `json:"chart,omitempty"`
	ConcurrentSyncs int         `json:"concurrentSyncs,omitempty"`
	Resources       interface{} `json:"resources,omitempty"`
	VPA             interface{} `json:"vpa,omitempty"`
}

type ControllerRegistration struct {
	ControllerResources []v1beta1.ControllerResource `json:"controllerResources,omitempty"`
}
