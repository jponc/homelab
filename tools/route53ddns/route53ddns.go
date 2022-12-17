package route53ddns

import (
	"homelab/tools/types"

	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/batch/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateRoute53DDNS(ctx *pulumi.Context, config *types.Config) error {
	secretName := "route53ddns-secret"
	secret, err := corev1.NewSecret(ctx, secretName, &corev1.SecretArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(secretName),
		},
		Type: pulumi.String("Opaque"),
		Data: pulumi.StringMap{
			"AWS_ACCESS_KEY_ID":      config.Route53DDNSAWSAccessKeyIdBase64,
			"AWS_SECRET_ACCESS_KEY":  config.Route53DDNSAWSSecretAccessKeyBase64,
			"ROUTE53_DOMAINS":        config.Route53DDNSRoute53DomainsBase64,
			"ROUTE53_HOSTED_ZONE_ID": config.Route53DDNSRoute53HostedZoneIdBase64,
		},
	})
	if err != nil {
		return err
	}

	crontName := "route53ddns-cronjob"
	_, err = batchv1.NewCronJob(ctx, crontName, &batchv1.CronJobArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(crontName),
		},
		Spec: batchv1.CronJobSpecArgs{
			Schedule: pulumi.String("*/10 * * * *"), // Every 10 minutes
			JobTemplate: batchv1.JobTemplateSpecArgs{
				Spec: batchv1.JobSpecArgs{
					Template: corev1.PodTemplateSpecArgs{
						Spec: corev1.PodSpecArgs{
							Containers: corev1.ContainerArray{
								corev1.ContainerArgs{
									Name:            pulumi.String(crontName),
									Image:           pulumi.String(config.Route53DDNSImage),
									ImagePullPolicy: pulumi.String("IfNotPresent"),
									EnvFrom: corev1.EnvFromSourceArray{
										corev1.EnvFromSourceArgs{
											SecretRef: corev1.SecretEnvSourceArgs{
												Name: pulumi.String(secretName),
											},
										},
									},
								},
							},
							RestartPolicy: pulumi.String("OnFailure"),
						},
					},
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{secret}))
	if err != nil {
		return err
	}

	return nil
}
