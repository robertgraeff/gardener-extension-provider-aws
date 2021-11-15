module github.com/gardener/gardener-extension-provider-aws/landscaper

go 1.16

require (
	github.com/gardener/gardener v1.34.1
	github.com/gardener/gardener-extension-provider-aws v1.29.1
	//github.com/gardener/landscaper-utils v0.1.0 // indirect
	github.com/gardener/landscaper-utils/deployutils   v0.0.0-20211115155816-a36c4149b98c
	github.com/gardener/landscaper-utils/machineimages v0.0.0-20211115155816-a36c4149b98c
	github.com/gardener/landscaper/apis v0.15.1
	github.com/go-logr/logr v0.4.0
	github.com/spf13/cobra v1.2.1
	k8s.io/apimachinery v0.22.3
	sigs.k8s.io/controller-runtime v0.10.2
	sigs.k8s.io/yaml v1.3.0
)

replace (
	github.com/gardener/gardener-resource-manager/api => github.com/gardener/gardener-resource-manager/api v0.25.0
	k8s.io/client-go => k8s.io/client-go v0.21.2
)
