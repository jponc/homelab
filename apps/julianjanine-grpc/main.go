package main

import (
	"homelab/apps/julianjanine-grpc/frontend"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Setup Kubernetes resources
		_, err := yaml.NewConfigFile(ctx, "julianjanine-grpc-yaml", &yaml.ConfigFileArgs{
			File:      "julianjanine-grpc.yaml",
			SkipAwait: false,
		})
		if err != nil {
			return err
		}

		// Configure gRPC service environment variables
		_, err = yaml.NewConfigFile(ctx, "secret", &yaml.ConfigFileArgs{
			File:      "secret.yaml",
			SkipAwait: false,
		})
		if err != nil {
			return err
		}

		// Create frontend cloudfront distribution
		err = frontend.CreateFrontend(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
