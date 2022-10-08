package appsnamespace

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateAppsNamespace(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigFile(ctx, "appsnamespace-yaml", &yaml.ConfigFileArgs{
		File:      "appsnamespace/namespace.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
