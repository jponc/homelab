package metabase

import (
	"homelab/tools/types"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"

	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
)

func CreateMetabase(ctx *pulumi.Context, config *types.Config) error {
	deploymentName := "metabase-deployment"
	appName := "metabase"
	deployment, err := appsv1.NewDeployment(ctx, deploymentName, &appsv1.DeploymentArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(deploymentName),
		},
		Spec: appsv1.DeploymentSpecArgs{
			Selector: metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String(appName),
				},
			},
			Replicas: pulumi.Int(1),
			Template: corev1.PodTemplateSpecArgs{
				Metadata: metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String(appName),
					},
				},
				Spec: corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:            pulumi.String(appName),
							Image:           pulumi.String(config.MetabaseImage),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(3000),
								},
							},
							Env: corev1.EnvVarArray{
								corev1.EnvVarArgs{
									Name:  pulumi.String("MB_DB_TYPE"),
									Value: pulumi.String("postgres"),
								},
								corev1.EnvVarArgs{
									Name: pulumi.String("MB_DB_DBNAME"),
									ValueFrom: corev1.EnvVarSourceArgs{
										SecretKeyRef: corev1.SecretKeySelectorArgs{
											Name: pulumi.String("postgres-secret"),
											Key:  pulumi.String("POSTGRES_DB"),
										},
									},
								},
								corev1.EnvVarArgs{
									Name:  pulumi.String("MB_DB_PORT"),
									Value: pulumi.String("5432"),
								},
								corev1.EnvVarArgs{
									Name: pulumi.String("MB_DB_USER"),
									ValueFrom: corev1.EnvVarSourceArgs{
										SecretKeyRef: corev1.SecretKeySelectorArgs{
											Name: pulumi.String("postgres-secret"),
											Key:  pulumi.String("POSTGRES_USER"),
										},
									},
								},
								corev1.EnvVarArgs{
									Name: pulumi.String("MB_DB_PASS"),
									ValueFrom: corev1.EnvVarSourceArgs{
										SecretKeyRef: corev1.SecretKeySelectorArgs{
											Name: pulumi.String("postgres-secret"),
											Key:  pulumi.String("POSTGRES_PASSWORD"),
										},
									},
								},
								corev1.EnvVarArgs{
									Name:  pulumi.String("MB_DB_HOST"),
									Value: pulumi.String("postgres"),
								},
							},
						},
					},
					RestartPolicy: pulumi.String("Always"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	serviceName := "metabase"
	service, err := corev1.NewService(ctx, serviceName, &corev1.ServiceArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(serviceName),
		},
		Spec: corev1.ServiceSpecArgs{
			Type: pulumi.String("ClusterIP"),
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port:       pulumi.Int(3000),
					TargetPort: pulumi.Int(3000),
				},
			},
			Selector: pulumi.StringMap{
				"app": pulumi.String(appName),
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{deployment}))
	if err != nil {
		return err
	}

	// Setup metabase yaml
	_, err = yaml.NewConfigFile(ctx, "metabase-yaml", &yaml.ConfigFileArgs{
		File:      "metabase/metabase.yaml",
		SkipAwait: false,
	}, pulumi.DependsOn([]pulumi.Resource{service}))
	if err != nil {
		return err
	}

	return nil
}
