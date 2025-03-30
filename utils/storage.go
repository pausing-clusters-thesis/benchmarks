package utils

import (
	"fmt"

	"github.com/pausing-clusters-thesis/benchmarks/naming"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func GetImmediateStorageClassForCSIDriver(csiDriverName string) (*storagev1.StorageClass, error) {
	switch csiDriverName {
	case naming.LocalCSIDriverName:
		return &storagev1.StorageClass{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "scylladb-local-xfs-immediate-",
			},
			Provisioner: naming.LocalCSIDriverName,
			Parameters: map[string]string{
				"csi.storage.k8s.io/fstype": "xfs",
			},
			ReclaimPolicy:     ptr.To(corev1.PersistentVolumeReclaimDelete),
			VolumeBindingMode: ptr.To(storagev1.VolumeBindingImmediate),
		}, nil

	case naming.GCEPersistentDiskCSIDriverName:
		return &storagev1.StorageClass{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "premium-rwo-immediate-",
			},
			Provisioner: naming.GCEPersistentDiskCSIDriverName,
			Parameters: map[string]string{
				"csi.storage.k8s.io/fstype": "xfs",
				"type":                      "pd-ssd",
			},
			ReclaimPolicy:     ptr.To(corev1.PersistentVolumeReclaimDelete),
			VolumeBindingMode: ptr.To(storagev1.VolumeBindingImmediate),
		}, nil

	case naming.EBSCSIDriverName:
		return &storagev1.StorageClass{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "gp3-immediate-",
			},
			Provisioner: naming.EBSCSIDriverName,
			Parameters: map[string]string{
				"csi.storage.k8s.io/fstype": "xfs",
				"type":                      "gp3",
			},
			ReclaimPolicy:     ptr.To(corev1.PersistentVolumeReclaimDelete),
			VolumeBindingMode: ptr.To(storagev1.VolumeBindingImmediate),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported CSI Driver name: %q", csiDriverName)

	}
}
