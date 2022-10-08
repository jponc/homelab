package postgres

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePostgresDB(ctx *pulumi.Context) error {
	secret, err := yaml.NewConfigFile(ctx, "postgres-secret", &yaml.ConfigFileArgs{
		File:      "postgres/secret.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	pv, err := yaml.NewConfigFile(ctx, "postgres-pv", &yaml.ConfigFileArgs{
		File:      "postgres/persistent-volume.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	pvc, err := yaml.NewConfigFile(ctx, "postgres-pvc", &yaml.ConfigFileArgs{
		File:      "postgres/persistent-volume-claim.yaml",
		SkipAwait: false,
	}, pulumi.DependsOn([]pulumi.Resource{pv}))
	if err != nil {
		return err
	}

	deployment, err := yaml.NewConfigFile(ctx, "postgres-deployment", &yaml.ConfigFileArgs{
		File:      "postgres/deployment.yaml",
		SkipAwait: false,
	}, pulumi.DependsOn([]pulumi.Resource{secret, pv, pvc}))
	if err != nil {
		return err
	}

	_, err = yaml.NewConfigFile(ctx, "postgres-service", &yaml.ConfigFileArgs{
		File:      "postgres/service.yaml",
		SkipAwait: false,
	}, pulumi.DependsOn([]pulumi.Resource{deployment}))
	if err != nil {
		return err
	}

	return nil
}
