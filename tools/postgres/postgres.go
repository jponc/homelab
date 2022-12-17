package postgres

import (
	"homelab/tools/types"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePostgresDB(ctx *pulumi.Context, config *types.Config) error {
	secretName := "postgres-secret"
	secret, err := corev1.NewSecret(ctx, secretName, &corev1.SecretArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(secretName),
		},
		Type: pulumi.String("Opaque"),
		Data: pulumi.StringMap{
			"POSTGRES_DB":       config.PostgresDbBase64,
			"POSTGRES_USER":     config.PostgresUserBase64,
			"POSTGRES_PASSWORD": config.PostgresPasswordBase64,
		},
	})
	if err != nil {
		return err
	}

	pvcName := "postgres-pvc"
	pvc, err := corev1.NewPersistentVolumeClaim(ctx, pvcName, &corev1.PersistentVolumeClaimArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(pvcName),
		},
		Spec: corev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.StringArray{
				pulumi.String("ReadWriteMany"),
			},
			StorageClassName: pulumi.String("longhorn"),
			Resources: corev1.ResourceRequirementsArgs{
				Requests: pulumi.StringMap{
					"storage": pulumi.String("10Gi"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	deploymentName := "postgres-deployment"
	appName := "postgres"
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
							Image:           pulumi.String(config.PostgresImage),
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Ports: corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(5432),
								},
							},
							EnvFrom: corev1.EnvFromSourceArray{
								corev1.EnvFromSourceArgs{
									SecretRef: corev1.SecretEnvSourceArgs{
										Name: pulumi.String(secretName),
									},
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								corev1.VolumeMountArgs{
									MountPath: pulumi.String("/var/lib/postgresql/data"),
									Name:      pulumi.String("postgresdata"),
									SubPath:   pulumi.String("postgres"),
								},
							},
						},
					},
					Volumes: corev1.VolumeArray{
						corev1.VolumeArgs{
							Name: pulumi.String("postgresdata"),
							PersistentVolumeClaim: corev1.PersistentVolumeClaimVolumeSourceArgs{
								ClaimName: pulumi.String(pvcName),
							},
						},
					},
					RestartPolicy: pulumi.String("Always"),
				},
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{secret, pvc}))
	if err != nil {
		return err
	}

	serviceName := "postgres"
	_, err = corev1.NewService(ctx, serviceName, &corev1.ServiceArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name: pulumi.String(serviceName),
		},
		Spec: corev1.ServiceSpecArgs{
			Type: pulumi.String("NodePort"),
			Ports: corev1.ServicePortArray{
				corev1.ServicePortArgs{
					Port: pulumi.Int(5432),
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

	return nil
}
