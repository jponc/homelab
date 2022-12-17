package monitoring

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateMonitoring(ctx *pulumi.Context) error {
	namespace, err := corev1.NewNamespace(ctx, "monitoring", &corev1.NamespaceArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String("monitoring"),
		},
	})
	if err != nil {
		return err
	}

	pvc, err := corev1.NewPersistentVolumeClaim(ctx, "monitoring-pvc", &corev1.PersistentVolumeClaimArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String("monitoring-pvc"),
			Namespace: pulumi.String("monitoring"),
		},
		Spec: corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{
				pulumi.String("ReadWriteMany"),
			},
			StorageClassName: pulumi.String("longhorn"),
			Resources: corev1.ResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": pulumi.String("30Gi"),
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{namespace}))
	if err != nil {
		return err
	}

	_, err = yaml.NewConfigFile(ctx, "prometheus-yaml", &yaml.ConfigFileArgs{
		File:      "monitoring/prometheus.yaml",
		SkipAwait: false,
	}, pulumi.DependsOn([]pulumi.Resource{pvc}))
	if err != nil {
		return err
	}

	_, err = yaml.NewConfigFile(ctx, "grafana-yaml", &yaml.ConfigFileArgs{
		File:      "monitoring/grafana.yaml",
		SkipAwait: false,
	}, pulumi.DependsOn([]pulumi.Resource{pvc}))
	if err != nil {
		return err
	}

	_, err = yaml.NewConfigFile(ctx, "kubestatemetrics-yaml", &yaml.ConfigFileArgs{
		File:      "kubestatemetrics/kubestatemetrics.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
