package main

import (
	"homelab/tools/certmanager"
	"homelab/tools/grafana"
	"homelab/tools/kubestatemetrics"
	"homelab/tools/postgres"
	"homelab/tools/prometheus"
	"homelab/tools/pvcs"
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

		_, monitoringPVC, err := pvcs.CreatePVCs(ctx)
		if err != nil {
			return err
		}

		err = postgres.CreatePostgresDB(ctx, config)
		if err != nil {
			return err
		}

		err = prometheus.CreatePrometheus(ctx, monitoringPVC)
		if err != nil {
			return err
		}

		err = grafana.CreateGrafana(ctx, monitoringPVC)
		if err != nil {
			return err
		}

		err = kubestatemetrics.CreateKubeStateMetrics(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
