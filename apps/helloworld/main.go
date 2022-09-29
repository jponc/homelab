package main

import (
	_ "embed"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed index.html
var indexHtml string

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Creates configmap data which has index html file
		configMap, err := corev1.NewConfigMap(ctx, "configmap", &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("helloworld"),
			},
			Data: pulumi.StringMap{
				"index.html": pulumi.String(indexHtml),
			},
		})
		if err != nil {
			return err
		}

		// Runs the entire yaml for helloworld app
		_, err = yaml.NewConfigFile(ctx, "helloworld-yaml", &yaml.ConfigFileArgs{
			File:      "helloworld.yaml",
			SkipAwait: false,
		}, pulumi.DependsOn([]pulumi.Resource{configMap}))
		if err != nil {
			return err
		}

		return nil
	})
}
