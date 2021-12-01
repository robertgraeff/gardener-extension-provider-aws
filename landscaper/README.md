# Deployment of Component "gardener-extension-provider-aws"

This directory contains the Landscaper based deployment of component `gardener-extension-provider-aws`.
The deployment is defined in a [blueprint](blueprint/blueprint.yaml) with a single [container deploy item](blueprint/...).
In particular, we do not use an aggregated blueprint with sub-installations for pre and post processing steps.

The deploy code is implemented in Go and executed by the container deployer.
It creates a `ControllerDeployment`, `ControllerRegistration`, and a `CloudProfile` on the target cluster.

### Default configuration

The default values for the three resources are contained in the following files:
- [controllerdeployment.yaml](./pkg/aws/resources/controllerdeployment.yaml).
- [controllerregistration.yaml](./pkg/aws/resources/controllerregistration.yaml)
- [cloudprofile.yaml](./pkg/aws/resources/cloudprofile.yaml)

Some of them can be overwritten and new values can be set by corresponding fields in the imports.

### Machine Images

From [landscaper-utils](https://github.com/gardener/landscaper-utils) we use:
- the [function to compute machine images](https://github.com/gardener/landscaper-utils/blob/master/machineimages/pkg/machineimages/machine_images.go) 
- the [machine images schema](https://github.com/gardener/landscaper-utils/blob/master/machineimages/.landscaper/machine-images/schemata/machine-images.json)
to type some imports.

### Image Vector

Landscaper offers a templating function 
[generateImageOverwrite](https://github.com/gardener/landscaper/blob/master/docs/usage/TemplateExecutors.md#template-executors) 
for the computation of an image vector. The image vector for the provider-aws component is computed during the 
templating of the [deploy item](./blueprint/deploy-execution-container.yaml):

```yaml
imageVectorOverwrite:
{{- generateImageOverwrite .cd .imports.lssComponentDescriptor | toYaml | nindent 8 }}
```

Function `generateImageOverwrite` gets as input the component descriptor of the provider-aws component (`.cd`)
and a component descriptor list with one element, namely the lss component descriptor (`.imports.lssComponentDescriptor`).
This component descriptor list comes from the import parameter `lssComponentDescriptor` of the blueprint.
After the templating, the resulting image vector is contained in field `.spec.config.importValues.imageVectorOverwrite`
of the deploy item. The Go code executed by the container deployer can access these import values.

### Utilities

- utilities 
  - reading imports
  - reading image references from the component descriptor,
  - creating a kubernetes client from a Target import parameter,


## Installing and running the deploy program locally

For test purposes, one can install and run the deploy code for the provider-aws component from the command line.
In this case the Landscaper is not involved. 
On the other hand, if one deploys with Landscaper, the same code is executed by the container deployer.

We assume that variable `REPO_ROOT` contains the path to the root directory of the current git repository.
Install the deploy command for the provider-aws component:

```shell
cd ${REPO_ROOT}/landscaper
make install-deployer
```

Start a deployment of the provider-aws component with the following command.
Note that the imports.yaml and component-descriptor.yaml files contain only example data that might be adjusted.
In particular, you must maintain the kubeconfig of the target cluster in the imports.yaml.

```shell
landscaper-extension-provider-aws \
  --operation RECONCILE \
  --imports-path "${REPO_ROOT}/landscaper/pkg/example/imports.yaml" \
  --exports-path "${REPO_ROOT}/landscaper/pkg/example/exports.yaml" \
  --component-desciptor-path "${REPO_ROOT}/landscaper/pkg/example/component-descriptor.yaml"
```

Instead of the command arguments, you can set the following environment variables: 

```text
OPERATION=RECONCILE
IMPORTS_PATH=${REPO_ROOT}/landscaper/pkg/example/imports.yaml
EXPORTS_PATH=${REPO_ROOT}/landscaper/pkg/example/exports.yaml
COMPONENT_DESCRIPTOR_PATH=${REPO_ROOT}/landscaper/pkg/example/component-descriptor.yaml
```

In order to uninstall the provider-aws component, set `OPERATION=DELETE`.

Use `${REPO_ROOT}/landscaper/cmd/main.go` as entry point if you want to debug the deploy code.


## Upload Component Descriptor

To manually build and push the component descriptor to the oci registry, 
`cd` into the repository root directory and run:

```shell
make cnudie-cd-build-push
```


## To Do

- Controller deployment
  - How can we access the chart that must be included in the controller deployment? Template it into the deploy item.
  - Image vector from lib, not from import

- Apline and pause images
  - Add alpine and pause image to resources.yaml. Check that name, repo and tag are the same as in the dev system.
  - Add them to the example comopnent descriptor
  - Remove the `_images.tpl` file from the chart, so that the images are not added twice to the component descriptor.
    Adjust the `_helpers.tpl` file. 
    We can also adjust the caller of the helper functions in the chart.
  
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

- Tests
