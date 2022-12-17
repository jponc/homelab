package main

import (
	"homelab/tools/certmanager"
	"homelab/tools/monitoring"
	"homelab/tools/postgres"
	"homelab/tools/route53ddns"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pulumiconfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		pulumiConfig := pulumiconfig.New(ctx, "")
		config := NewConfig(pulumiConfig)

		err := certmanager.CreateCertManager(ctx)
		if err != nil {
			return err
		}

		err = route53ddns.CreateRoute53DDNS(ctx, config)
		if err != nil {
			return err
		}

		err = postgres.CreatePostgresDB(ctx, config)
		if err != nil {
			return err
		}

		err = monitoring.CreateMonitoring(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
