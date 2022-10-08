package route53ddns

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateRoute53DDNS(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigFile(ctx, "route53ddns-cronjob", &yaml.ConfigFileArgs{
		File:      "route53ddns/cronjob.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
