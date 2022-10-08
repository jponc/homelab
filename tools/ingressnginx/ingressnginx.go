package ingressnginx

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateIngressNginx(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigFile(ctx, "ingress-nginx", &yaml.ConfigFileArgs{
		File:      "ingressnginx/ingressnginx.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
