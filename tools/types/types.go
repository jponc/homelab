package types

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

type Config struct {
	Route53DDNSImage                     string
	Route53DDNSAWSAccessKeyIdBase64      pulumi.StringOutput
	Route53DDNSAWSSecretAccessKeyBase64  pulumi.StringOutput
	Route53DDNSRoute53DomainsBase64      pulumi.StringOutput
	Route53DDNSRoute53HostedZoneIdBase64 pulumi.StringOutput
}
