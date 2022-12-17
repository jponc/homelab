package prometheus

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePrometheus(ctx *pulumi.Context, monitoringPVC *corev1.PersistentVolumeClaim) error {
	_, err := yaml.NewConfigFile(ctx, "prometheus-yaml", &yaml.ConfigFileArgs{
		File:      "prometheus/prometheus.yaml",
		SkipAwait: false,
	}, pulumi.DependsOn([]pulumi.Resource{monitoringPVC}))
	if err != nil {
		return err
	}

	return nil
}
