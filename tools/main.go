package main

import (
	"homelab/tools/certmanager"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		err := certmanager.CreateCertManager(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
