# Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

EXTENSION_PREFIX                       := gardener-extension
NAME                                   := provider-aws
LANDSCAPER_NAME                        := landscaper-extension-provider-aws
PROVIDER_AWS_DEPLOYER_IMAGE_REPOSITORY := eu.gcr.io/gardener-project/landscaper-gardener-deployers/provider-aws-deployer
LANDSCAPER_DIR                         := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
REPO_ROOT                              := $(LANDSCAPER_DIR)/..
HACK_DIR                               := $(REPO_ROOT)/hack
GARDENER_HACK_DIR                      := $(REPO_ROOT)/vendor/github.com/gardener/gardener/hack
VERSION                                := $(shell cat "$(REPO_ROOT)/VERSION")
EFFECTIVE_VERSION                      := $(VERSION)-$(shell git rev-parse HEAD)


.PHONY: install
install:
	CGO_ENABLED=0 GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) GO111MODULE=on \
	LD_FLAGS="-w -X github.com/gardener/$(EXTENSION_PREFIX)-$(NAME)/landscaper/pkg/version.Version=$(EFFECTIVE_VERSION)" \
	go install -mod=vendor $(LANDSCAPER_DIR)/cmd/...


.PHONY: docker-images
docker-images:
	@echo "Building docker images for version $(EFFECTIVE_VERSION)"
	@docker build --build-arg EFFECTIVE_VERSION=$(EFFECTIVE_VERSION) -t $(PROVIDER_AWS_DEPLOYER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION) -f Dockerfile .


.PHONY: docker-push
docker-push:
	@echo "Pushing docker images for version $(EFFECTIVE_VERSION) to registry $(REGISTRY)"
	@if ! docker images $(PROVIDER_AWS_DEPLOYER_IMAGE_REPOSITORY) | awk '{ print $$2 }' | grep -q -F $(EFFECTIVE_VERSION); then echo "$(PROVIDER_AWS_DEPLOYER_IMAGE_REPOSITORY) version $(EFFECTIVE_VERSION) is not yet built. Please run 'make docker-images'"; false; fi
	@docker push $(PROVIDER_AWS_DEPLOYER_IMAGE_REPOSITORY):$(EFFECTIVE_VERSION)
