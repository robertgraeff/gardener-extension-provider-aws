# Landscaper Based Deployment of Component "provider-aws"

Consumes landscaper-utils lib for:
- machine image computation
- utilities 
  - `options` object to read imports, write exports
  - read a resource item from the component descriptor,
  - and function to create a kubernetes client from a target,

Will consume the image vector function

## Debug Configuration 

Enter the kubeconfig of the target cluster in $REPO_ROOT/landscaper/pkg/example/imports.yaml

Start $REPO_ROOT/landscaper/cmd/main.go with the following environment variables (adjust the paths!):

```text
OPERATION=RECONCILE
IMPORTS_PATH=.../gardener-extension-provider-aws/landscaper/pkg/example/imports.yaml
EXPORTS_PATH=.../gardener-extension-provider-aws/landscaper/pkg/example/exports.yaml
COMPONENT_DESCRIPTOR_PATH=.../gardener-extension-provider-aws/landscaper/pkg/example/component-descriptor.yaml
```

## To Do

- Controller deployment
  - Implementation of function `constructControllerDeployment`
  - Can we access the chart that must be included in the controller deployment?
  - Image vector from lib, not from import

- Utils 
  - Function `getResolve` to read an entry from the component descriptor

- The cloudprofile, controllerdeployment, and controllerregistration,
  do they have a namespace ("garden") or not?

- Blueprint
  - Container DeployItem 
  - The imports in the blueprint must correspond with the go structures in landscaper/pkg/aws/imports.go

- Simplify logger in landscaper-utils/deployutils

- Dev Process
  - Include chart
