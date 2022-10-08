package main

import (
	"homelab/tools/certmanager"
	"homelab/tools/route53ddns"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := certmanager.CreateCertManager(ctx)
		if err != nil {
			return err
		}

		err = route53ddns.CreateRoute53DDNS(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
