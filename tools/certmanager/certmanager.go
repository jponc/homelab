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

	return nil
}
