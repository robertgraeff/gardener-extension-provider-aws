#!/usr/bin/env bash
#
# Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# SPDX-License-Identifier: Apache-2.0

set -e

SOURCE_PATH="$(dirname $0)/.."
source "${SOURCE_PATH}/hack/environment.sh"

TMP_DIR="$(mktemp -d)"
INSTALLATION_PATH="${TMP_DIR}/installation.yaml"

REGISTRY=$(get_cd_registry)
COMPONENT_NAME=$(get_cd_component_name)

cat << EOF > ${INSTALLATION_PATH}
apiVersion: landscaper.gardener.cloud/v1alpha1
kind: Installation
metadata:
  name: provider-aws
spec:
  componentDescriptor:
    ref:
      repositoryContext:
        type: ociRegistry
        baseUrl: ${REGISTRY}
      componentName: ${COMPONENT_NAME}
      version: ${EFFECTIVE_VERSION}

  blueprint:
    ref:
      resourceName: blueprint

  imports:
    targets:
      - name: cluster
        target: "#cluster"

    componentDescriptors:
      - name: lssComponentDescriptor
        list:
        - ref:
            componentName: github.wdf.sap.corp/kubernetes/landscape-setup
            version: 0.2510.0
            repositoryContext:
              type: ociRegistry
              baseUrl: eu.gcr.io/sap-se-gcr-k8s-private/cnudie/gardener/development

  importDataMappings:
    cloudProfile:
      machineImages:
        - name: test1
          versions:
            - version: v1
              region: r1
            - version: v2
              region: r2
        - name: test2
          versions:
            - version: v3
              region: r3

    kubernetesVersions:
      - classification: preview
        version: 1.22.2
      - classification: supported
        version: 1.21.5

    controllerRegistration:
      concurrentSyncs: 50
      resources:
        requests:
          cpu: 100m
          memory: 512Mi
        limits:
          cpu: 1000m
          memory: 1Gi
      vpa:
        enabled: true
        resourcePolicy:
          minAllowed:
            cpu: 50m
            memory: 256Mi
        updatePolicy:
          updateMode: "Auto"
EOF

echo "Installation stored at ${INSTALLATION_PATH}"