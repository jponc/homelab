package main

import (
	"homelab/tools/appsnamespace"
	"homelab/tools/certmanager"
	"homelab/tools/grafana"
	"homelab/tools/kubestatemetrics"
	"homelab/tools/postgres"
	"homelab/tools/prometheus"
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

		err = postgres.CreatePostgresDB(ctx)
		if err != nil {
			return err
		}

		err = appsnamespace.CreateAppsNamespace(ctx)
		if err != nil {
			return err
		}

		err = prometheus.CreatePrometheus(ctx)
		if err != nil {
			return err
		}

		err = grafana.CreateGrafana(ctx)
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
