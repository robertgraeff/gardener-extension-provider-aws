# Deployment of Component "gardener-extension-provider-aws"

This directory contains the Landscaper based deployment of component `gardener-extension-provider-aws`.
The deployment is defined in a [blueprint](blueprint/blueprint.yaml) with a single [container deploy item](blueprint/...).
In particular, we do not use an aggregated blueprint with sub-installations for pre and post processing steps.

The deploy code that is executed by the container deployer is implemented in Go.
It creates a `ControllerDeployment`, `ControllerRegistration`, and a `CloudProfile`.

Consumes landscaper-utils lib for:
- machine image computation
- image vector computation (todo)
- utilities 
  - reading imports
  - reading image references from the component descriptor,
  - creating a kubernetes client from a Target import parameter,


## Debug Configuration 

Enter the kubeconfig of the target cluster in $REPO_ROOT/landscaper/pkg/example/imports.yaml

Start $REPO_ROOT/landscaper/cmd/main.go with the following environment variables (adjust the paths!):

```text
OPERATION=RECONCILE
IMPORTS_PATH=.../gardener-extension-provider-aws/landscaper/pkg/example/imports.yaml
EXPORTS_PATH=.../gardener-extension-provider-aws/landscaper/pkg/example/exports.yaml
COMPONENT_DESCRIPTOR_PATH=.../gardener-extension-provider-aws/landscaper/pkg/example/component-descriptor.yaml
```

## Upload Component Descriptor

To manually build and push the component descriptor to the oci registry, 
`cd` into the repository root directory and run:

```shell
make cnudie-cd-build-push
```

## To Do

- Controller deployment
  - Implementation of function `constructControllerDeployment`
  - How can we access the chart that must be included in the controller deployment? Template it into the deploy item.
  - Image vector from lib, not from import

- Image vector
  - where should the charts/images.yaml file go to?
  - honour component descriptor?
  - use lss?

- The cloudprofile, controllerdeployment, and controllerregistration,
  do they have a namespace ("garden") or not?

- Blueprint
  - Container DeployItem with a generic image using an executable in the blueprint directory
  - Adjust imports of the blueprint (for example machine images)
  - The imports in the blueprint must correspond with the go structures in landscaper/pkg/aws/imports.go

- Simplify logger in landscaper-utils/deployutils

- Dev Process
  - Include chart as base64 string in the example imports.yaml
  - Build the deploy code and put the executable into the blueprint directory
