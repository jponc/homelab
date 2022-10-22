package kubestatemetrics

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateKubeStateMetrics(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigFile(ctx, "kubestatemetrics-yaml", &yaml.ConfigFileArgs{
		File:      "kubestatemetrics/kubestatemetrics.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
