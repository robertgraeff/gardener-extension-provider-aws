#!/usr/bin/env bash
#
# Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
#
# SPDX-License-Identifier: Apache-2.0

set -e

SOURCE_PATH="$(dirname $0)/.."

TMP_DIR="$(mktemp -d)"
INSTALLATION_PATH="${TMP_DIR}/installation.yaml"

REGISTRY="eu.gcr.io/sap-se-gcr-k8s-private/cnudie/gardener/development"
COMPONENT_NAME="github.com/gardener/gardener-extension-provider-aws"

cat << EOF > ${INSTALLATION_PATH}
apiVersion: landscaper.gardener.cloud/v1alpha1
kind: Installation
metadata:
  name: provider-aws
  namespace: test
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
        target: "#rh-cluster"

    componentDescriptors:
    - name: "lssComponentDescriptor"
      list:
      - ref:
          componentName: github.wdf.sap.corp/kubernetes/landscape-setup
          version: 0.2583.0
          repositoryContext:
            type: ociRegistry
            baseUrl: ${REGISTRY}

  importDataMappings:
    cloudProfile:
      kubernetesVersions: []

    machineImages:
    - name: gardenlinux
      versions:
        - classification: preview
          cri:
            - name: docker
            - containerRuntimes:
                - type: gvisor
              name: containerd
          version: 576.0.0
        - classification: preview
          cri:
            - name: docker
            - containerRuntimes:
                - type: gvisor
              name: containerd
          version: 318.9.0
        - classification: supported
          cri:
            - name: docker
            - containerRuntimes:
                - type: gvisor
              name: containerd
          version: 318.8.0
        - classification: deprecated
          cri:
            - name: docker
            - containerRuntimes:
                - type: gvisor
              name: containerd
          expirationDate: '2022-01-15T23:59:59Z'
          version: 184.0.0
    - name: suse-chost
      versions:
        - classification: supported
          cri:
            - name: docker
          version: 15.2.20211025
        - cri:
            - name: docker
          version: 15.2.20211025-gen2
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-03-31T23:59:59Z'
          version: 15.2.20210913
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-03-31T23:59:59Z'
          version: 15.2.20210913-gen2
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-01-31T23:59:59Z'
          version: 15.2.20210722
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-01-31T23:59:59Z'
          version: 15.2.20210722-gen2
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-01-15T23:59:59Z'
          version: 15.2.20210610
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-01-15T23:59:59Z'
          version: 15.2.20210610-gen2
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-01-08T23:59:59Z'
          version: 15.2.20210405
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-01-08T23:59:59Z'
          version: 15.2.20210405-gen2
    - name: ubuntu
      versions:
        - classification: supported
          cri:
            - name: docker
            - containerRuntimes:
                - type: gvisor
              name: containerd
          version: 18.4.20210415
        - classification: deprecated
          cri:
            - name: docker
            - containerRuntimes:
                - type: gvisor
              name: containerd
          expirationDate: '2021-11-30T23:59:59Z'
          version: 18.4.20200228
    - name: memoryone-chost
      versions:
        - classification: preview
          cri:
            - name: docker
          version: 10.5.2450-19
        - classification: preview
          cri:
            - name: docker
          version: 10.5.2450-17
        - classification: preview
          cri:
            - name: docker
          version: 10.5.2450-16
        - classification: supported
          cri:
            - name: docker
          version: 10.5.2450-13
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2021-11-30T23:59:59Z'
          version: 10.5.2450-9
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2021-11-30T23:59:59Z'
          version: 10.5.2450-5
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2021-11-30T23:59:59Z'
          version: 10.5.2450-4
        - classification: deprecated
          cri:
            - name: docker
          expirationDate: '2022-02-28T23:59:59Z'
          version: 10.5.2400-5

    machineImagesLs:
    - name: flatcar
      versions:
        - version: 2765.2.6
          classification: supported
          cri:
            - name: containerd
              containerRuntimes:
                - type: gvisor
            - name: docker

    machineImagesProviderLs:
    - name: flatcar
      versions:
        - version: 2765.2.6
          regions:
            - name: ap-east-1
              ami: ami-0df249fad4e423ae7
            - name: ap-northeast-1
              ami: ami-028697e8df1ff071b
            - name: ap-northeast-2
              ami: ami-0c5b8f2d07d21da16
            - name: ap-south-1
              ami: ami-055b64c22dbcd61b0
            - name: ap-southeast-1
              ami: ami-045357ea038a43fe7
            - name: ap-southeast-2
              ami: ami-05df81d055054698d
            - name: ca-central-1
              ami: ami-0f639872bfcb49738
            - name: eu-central-1
              ami: ami-055acc5a6e9587b44
            - name: eu-north-1
              ami: ami-04f64f11f4dacda92
            - name: eu-west-1
              ami: ami-019f09de46e4e3f88
            - name: eu-west-2
              ami: ami-0097d8b6241e9cf76
            - name: eu-west-3
              ami: ami-05b8b0131fbb39283
            - name: me-south-1
              ami: ami-0737c661a0881fd94
            - name: sa-east-1
              ami: ami-00bc3ae33287bc81b
            - name: us-east-1
              ami: ami-0fd66875fa1ef8395
            - name: us-east-2
              ami: ami-02eb704ee029f6b9e
            - name: us-west-1
              ami: ami-053fb35697f85574d
            - name: us-west-2
              ami: ami-019657181ea76e880

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
