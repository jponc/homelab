package certmanager

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateCertManager(ctx *pulumi.Context) error {
	// Create CertManager from Yaml
	_, err := yaml.NewConfigFile(ctx, "certmanager", &yaml.ConfigFileArgs{
		File:      "certmanager/cert-manager.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	// Add letsencrypt clusterissuer
	_, err = yaml.NewConfigFile(ctx, "letsencrypt-clusterissuer", &yaml.ConfigFileArgs{
		File:      "certmanager/letsencrypt-clusterissuer.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	// Add jponc.dev wildcard cert
	_, err = yaml.NewConfigFile(ctx, "jponc-dev-certificate", &yaml.ConfigFileArgs{
		File:      "certmanager/jponc-dev-certificate.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
