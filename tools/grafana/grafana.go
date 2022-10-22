package grafana

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateGrafana(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigFile(ctx, "grafana-yaml", &yaml.ConfigFileArgs{
		File:      "grafana/grafana.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
