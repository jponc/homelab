package main

import (
	"homelab/tools/types"

	pulumiconfig "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func NewConfig(pulumiConfig *pulumiconfig.Config) *types.Config {
	return &types.Config{
		Route53DDNSImage:                     pulumiConfig.Require("route53ddnsImage"),
		Route53DDNSAWSAccessKeyIdBase64:      pulumiConfig.RequireSecret("route53ddnsAWSAccessKeyIdBase64"),
		Route53DDNSAWSSecretAccessKeyBase64:  pulumiConfig.RequireSecret("route53ddnsAWSSecretAccessKeyBase64"),
		Route53DDNSRoute53DomainsBase64:      pulumiConfig.RequireSecret("route53ddnsRoute53DomainsBase64"),
		Route53DDNSRoute53HostedZoneIdBase64: pulumiConfig.RequireSecret("route53ddnsRoute53HostedZoneIdBase64"),

		PostgresImage:          pulumiConfig.Require("postgresImage"),
		PostgresDbBase64:       pulumiConfig.RequireSecret("postgresDbBase64"),
		PostgresPasswordBase64: pulumiConfig.RequireSecret("postgresPasswordBase64"),
		PostgresUserBase64:     pulumiConfig.RequireSecret("postgresUserBase64"),

		MetabaseImage: pulumiConfig.Require("metabaseImage"),
	}
}
