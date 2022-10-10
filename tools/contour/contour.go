package contour

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateContour(ctx *pulumi.Context) error {
	_, err := yaml.NewConfigFile(ctx, "contour", &yaml.ConfigFileArgs{
		File:      "contour/contour.yaml",
		SkipAwait: false,
	})
	if err != nil {
		return err
	}

	return nil
}
