package pvcs

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreatePVCs(ctx *pulumi.Context) (*corev1.PersistentVolumeClaim, *corev1.PersistentVolumeClaim, error) {
	defaultPVC, err := corev1.NewPersistentVolumeClaim(ctx, "default-pvc", &corev1.PersistentVolumeClaimArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String("default-pvc"),
			Namespace: pulumi.String("default"),
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
		return nil, nil, err
	}

	monitoringPVC, err := corev1.NewPersistentVolumeClaim(ctx, "monitoring-pvc", &corev1.PersistentVolumeClaimArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:      pulumi.String("monitoring-pvc"),
			Namespace: pulumi.String("monitoring"),
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
		return nil, nil, err
	}

	return defaultPVC, monitoringPVC, nil
}
