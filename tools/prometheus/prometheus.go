package prometheus

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePrometheus(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigFile(ctx, "prometheus-yaml", &yaml.ConfigFileArgs{
		File:      "prometheus/prometheus.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
