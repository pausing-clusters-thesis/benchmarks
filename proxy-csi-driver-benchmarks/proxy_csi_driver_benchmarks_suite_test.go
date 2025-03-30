package proxy_csi_driver_benchmarks_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"
	"github.com/pausing-clusters-thesis/benchmarks/naming"
	"github.com/pausing-clusters-thesis/benchmarks/utils"
	proxycsinaming "github.com/pausing-clusters-thesis/proxy-csi-driver/pkg/naming"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/genericclioptions"
	"github.com/scylladb/scylla-operator/pkg/helpers/slices"
	"github.com/scylladb/scylla-operator/test/e2e/framework"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type result struct {
	ElapsedTimeMs int64 `json:"elapsed_time_ms"`
}

const (
	prewarmTimeout = 1 * time.Minute
)

var (
	clientConfig = genericclioptions.NewClientConfig("proxy-csi-driver-benchmarks")

	imagePullPolicyString = string(corev1.PullAlways)

	proxyStorageClassName string
	backendCSIDriverName  string
	imagePullPolicy       corev1.PullPolicy
	destDir               string
)

var supportedImagePullPolicyStrings = []string{
	string(corev1.PullNever),
	string(corev1.PullAlways),
}

var supportedBackendCSIDriverNames = []string{
	naming.LocalCSIDriverName,
	naming.GCEPersistentDiskCSIDriverName,
	naming.EBSCSIDriverName,
}

func init() {
	flag.StringVar(&proxyStorageClassName, "proxy-storage-class-name", proxyStorageClassName, "The name of a StorageClass provisioned by the Proxy CSI Driver to be used in the test.")
	flag.StringVar(&backendCSIDriverName, "backend-csi-driver-name", backendCSIDriverName, fmt.Sprintf("The name of the backend CSI driver to test. Supported drvier names are: %v.", supportedBackendCSIDriverNames))
	flag.StringVar(&imagePullPolicyString, "image-pull-policy", imagePullPolicyString, fmt.Sprintf("ImagePullPolicy to use for containers. Supported policies are: %v. In case of PullNever, the image has to pre-pulled.", supportedImagePullPolicyStrings))
	flag.StringVar(&destDir, "dest-dir", destDir, "Destination directory in which results should be saved.")
}

func TestProxyCsiDriverBenchmarks(t *testing.T) {
	var err error

	err = Validate()
	if err != nil {
		t.Fatal(err)
	}

	err = Complete()
	if err != nil {
		t.Fatal(err)
	}

	framework.TestContext = &framework.TestContextType{
		RestConfigs:   []*restclient.Config{clientConfig.RestConfig},
		CleanupPolicy: framework.CleanupPolicyAlways,
	}

	o.RegisterFailHandler(g.Fail)

	g.RunSpecs(t, "ProxyCSIDriverBenchmarks Suite")
}

func Validate() error {
	var errs []error

	if len(proxyStorageClassName) == 0 {
		errs = append(errs, fmt.Errorf("proxy-storage-class-name must not be empty"))
	}

	if len(backendCSIDriverName) == 0 {
		errs = append(errs, fmt.Errorf("backend-csi-driver-name must not be empty"))
	} else if !slices.ContainsItem(supportedBackendCSIDriverNames, backendCSIDriverName) {
		errs = append(errs, fmt.Errorf("backend-csi-driver-name %q is not supported. Supported names are: %v", backendCSIDriverName, supportedBackendCSIDriverNames))
	}

	if len(imagePullPolicyString) == 0 {
		errs = append(errs, fmt.Errorf("image-pull-policy must not be empty"))
	} else if !slices.ContainsItem(supportedImagePullPolicyStrings, imagePullPolicyString) {
		errs = append(errs, fmt.Errorf("image-pull-policy %q is not supported. Supported policies are: %v", imagePullPolicyString, supportedImagePullPolicyStrings))
	}

	if len(destDir) > 0 {
		fi, err := os.Stat(destDir)
		if err != nil {
			errs = append(errs, fmt.Errorf("can't stat dest-dir %q: %w", destDir, err))
		} else if !fi.IsDir() {
			errs = append(errs, fmt.Errorf("dest-dir %q must be a directory", destDir))
		}
	} else {
		errs = append(errs, fmt.Errorf("dest-dir can't be empty"))
	}

	return errors.Join(errs...)
}

func Complete() error {
	var err error

	err = clientConfig.Complete()
	if err != nil {
		return err
	}

	imagePullPolicy = corev1.PullPolicy(imagePullPolicyString)

	return nil
}

var _ = g.Describe("measure time to readiness", func() {
	f := framework.NewFramework("benchmark")

	g.It("creating pod from scratch (baseline)", func(ctx g.SpecContext) {
		const resultsFileName = "baseline"
		resultsFilePath := path.Join(destDir, resultsFileName)
		resultsFile, err := os.OpenFile(resultsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		o.Expect(err).NotTo(o.HaveOccurred())
		g.DeferCleanup(resultsFile.Close)

		c := f.Cluster(0)
		ns, nsClient := c.CreateUserNamespace(ctx)

		backendImmediateStorageClass, err := utils.GetImmediateStorageClassForCSIDriver(backendCSIDriverName)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Creating immediate StorageClass for backend CSI driver")
		backendImmediateStorageClass, err = c.KubeAdminClient().StorageV1().StorageClasses().Create(ctx, backendImmediateStorageClass, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.DeferCleanup(func(ctx g.SpecContext, backendImmediateStorageClass *storagev1.StorageClass) {
			framework.By("Deleting immediate StorageClass")
			err := c.KubeAdminClient().StorageV1().StorageClasses().Delete(ctx, backendImmediateStorageClass.GetName(), metav1.DeleteOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
		}, backendImmediateStorageClass)

		framework.By("Creating backend PVC")
		backendPVC := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "backend-pvc",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("100Mi"),
					},
				},
				StorageClassName: ptr.To(backendImmediateStorageClass.GetName()),
				VolumeMode:       ptr.To(corev1.PersistentVolumeFilesystem),
			},
		}

		backendPVC, err = nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Create(ctx, backendPVC, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for backend PVC to be bound")
		backendPVCBindingCtx, backendPVCBindingCtxCancel := context.WithTimeout(ctx, 2*time.Minute)
		defer backendPVCBindingCtxCancel()
		backendPVC, err = waitForPersistentVolumeClaimState(backendPVCBindingCtx, f.KubeAdminClient().CoreV1().PersistentVolumeClaims(ns.GetName()), backendPVC.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPersistentVolumeClaimBound)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(backendPVC.Spec.VolumeName).NotTo(o.BeEmpty())

		backendPV, err := c.KubeAdminClient().CoreV1().PersistentVolumes().Get(ctx, backendPVC.Spec.VolumeName, metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		pod := getPodTemplate(ns.GetName(), backendPVC.GetName())
		pod.Spec.Containers[0].Command = []string{
			"bin/sh",
			"-euEo",
			"pipefail",
			"-c",
			strings.TrimSpace(`
trap 'kill $( jobs -p ); exit 0' TERM

touch /data/test

sleep infinity &
wait $!
`),
		}

		// Set Pod's NodeSelector to match backend PV's NodeSelector.
		if backendPV.Spec.NodeAffinity != nil && backendPV.Spec.NodeAffinity.Required != nil {
			pod.Spec.Affinity = &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: backendPV.Spec.NodeAffinity.Required.DeepCopy(),
				},
			}
		}

		startTime := time.Now()

		pod, err = nsClient.KubeClient().CoreV1().Pods(ns.GetName()).Create(ctx, pod, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		pod, err = socontrollerhelpers.WaitForPodState(ctx, f.KubeAdminClient().CoreV1().Pods(ns.GetName()), pod.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPodReady)
		stopTime := time.Now()
		o.Expect(err).NotTo(o.HaveOccurred())

		podReadyCondition, err := getPodCondition(pod, corev1.PodReady)
		o.Expect(err).NotTo(o.HaveOccurred())

		elapsedTime := stopTime.Sub(startTime)
		framework.Infof("Time elapsed: %v, stop time %v, pod ready condition timestamp %v", elapsedTime, stopTime, podReadyCondition.LastTransitionTime)

		encoder := json.NewEncoder(resultsFile)
		res := result{ElapsedTimeMs: elapsedTime.Milliseconds()}
		err = encoder.Encode(res)
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("pre-warming pod with proxy volume", func(ctx g.SpecContext) {
		const resultsFileName = "busywait"
		resultsFilePath := path.Join(destDir, resultsFileName)
		resultsFile, err := os.OpenFile(resultsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		o.Expect(err).NotTo(o.HaveOccurred())
		g.DeferCleanup(resultsFile.Close)

		c := f.Cluster(0)
		ns, nsClient := c.CreateUserNamespace(ctx)

		backendImmediateStorageClass, err := utils.GetImmediateStorageClassForCSIDriver(backendCSIDriverName)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Creating immediate StorageClass for backend CSI driver")
		backendImmediateStorageClass, err = c.KubeAdminClient().StorageV1().StorageClasses().Create(ctx, backendImmediateStorageClass, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.DeferCleanup(func(ctx g.SpecContext, backendImmediateStorageClass *storagev1.StorageClass) {
			framework.By("Deleting immediate StorageClass")
			err := c.KubeAdminClient().StorageV1().StorageClasses().Delete(ctx, backendImmediateStorageClass.GetName(), metav1.DeleteOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
		}, backendImmediateStorageClass)

		framework.By("Creating backend PVC")
		backendPVC := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "backend-pvc",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("100Mi"),
					},
				},
				StorageClassName: ptr.To(backendImmediateStorageClass.GetName()),
				VolumeMode:       ptr.To(corev1.PersistentVolumeFilesystem),
			},
		}

		backendPVC, err = nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Create(ctx, backendPVC, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for backend PVC to be bound")
		backendPVCBindingCtx, backendPVCBindingCtxCancel := context.WithTimeout(ctx, 2*time.Minute)
		defer backendPVCBindingCtxCancel()
		backendPVC, err = waitForPersistentVolumeClaimState(backendPVCBindingCtx, f.KubeAdminClient().CoreV1().PersistentVolumeClaims(ns.GetName()), backendPVC.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPersistentVolumeClaimBound)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(backendPVC.Spec.VolumeName).NotTo(o.BeEmpty())

		backendPV, err := c.KubeAdminClient().CoreV1().PersistentVolumes().Get(ctx, backendPVC.Spec.VolumeName, metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Creating proxy PVC")
		proxyPVC := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: "proxy-pvc",
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						// Doesn't matter, only a proxy directory is created.
						corev1.ResourceStorage: resource.MustParse("1Mi"),
					},
				},
				StorageClassName: ptr.To(proxyStorageClassName),
				VolumeMode:       ptr.To(corev1.PersistentVolumeFilesystem),
			},
		}

		proxyPVC, err = nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Create(ctx, proxyPVC, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		pod := getPodTemplate(ns.GetName(), proxyPVC.GetName())
		pod.Spec.Containers[0].Command = []string{
			"bin/sh",
			"-euEo",
			"pipefail",
			"-c",
			strings.TrimSpace(`
trap 'kill $( jobs -p ); exit 0' TERM

while true; do
	touch /data/test && break
done

sleep infinity &
wait $!
`),
		}

		// Set Pod's NodeSelector to match backend PV's NodeSelector.
		if backendPV.Spec.NodeAffinity != nil && backendPV.Spec.NodeAffinity.Required != nil {
			pod.Spec.Affinity = &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: backendPV.Spec.NodeAffinity.Required.DeepCopy(),
				},
			}
		}

		framework.By("Creating pre-warmed Pod")
		pod, err = nsClient.KubeClient().CoreV1().Pods(ns.GetName()).Create(ctx, pod, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for Pod to be running")
		podRunningCtx, podRunningCtxCancel := context.WithTimeout(ctx, prewarmTimeout)
		defer podRunningCtxCancel()
		pod, err = socontrollerhelpers.WaitForPodState(podRunningCtx, nsClient.KubeClient().CoreV1().Pods(ns.GetName()), pod.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPodRunning)
		o.Expect(err).NotTo(o.HaveOccurred())

		startTime := time.Now()

		annotateProxyPVCWithBackendPVCRef(ctx, nsClient.KubeClient().CoreV1(), pod, proxyPVC, backendPVC.GetName())

		pod, err = socontrollerhelpers.WaitForPodState(ctx, f.KubeAdminClient().CoreV1().Pods(ns.GetName()), pod.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPodReady)
		stopTime := time.Now()
		o.Expect(err).NotTo(o.HaveOccurred())

		podReadyCondition, err := getPodCondition(pod, corev1.PodReady)
		o.Expect(err).NotTo(o.HaveOccurred())

		elapsedTime := stopTime.Sub(startTime)
		framework.Infof("Time elapsed: %v, stop time %v, pod ready condition timestamp %v", elapsedTime, stopTime, podReadyCondition.LastTransitionTime)

		encoder := json.NewEncoder(resultsFile)
		res := result{ElapsedTimeMs: elapsedTime.Milliseconds()}
		err = encoder.Encode(res)
		o.Expect(err).NotTo(o.HaveOccurred())
	})

	g.It("pre-warms pod with proxy volume and suggested wait mechanism", func(ctx g.SpecContext) {
		const resultsFileName = "sidecar"
		resultsFilePath := path.Join(destDir, resultsFileName)
		resultsFile, err := os.OpenFile(resultsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		o.Expect(err).NotTo(o.HaveOccurred())
		g.DeferCleanup(resultsFile.Close)

		c := f.Cluster(0)
		ns, nsClient := c.CreateUserNamespace(ctx)

		backendImmediateStorageClass, err := utils.GetImmediateStorageClassForCSIDriver(backendCSIDriverName)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Creating immediate StorageClass for backend CSI driver")
		backendImmediateStorageClass, err = c.KubeAdminClient().StorageV1().StorageClasses().Create(ctx, backendImmediateStorageClass, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.DeferCleanup(func(ctx g.SpecContext, backendImmediateStorageClass *storagev1.StorageClass) {
			framework.By("Deleting immediate StorageClass")
			err := c.KubeAdminClient().StorageV1().StorageClasses().Delete(ctx, backendImmediateStorageClass.GetName(), metav1.DeleteOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
		}, backendImmediateStorageClass)

		framework.By("Creating a ServiceAccount")
		podSA, err := nsClient.KubeClient().CoreV1().ServiceAccounts(ns.Name).Create(ctx, &corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: "basic",
			},
		}, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		podRole, err := nsClient.KubeClient().RbacV1().Roles(ns.Name).Create(ctx, &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pods",
			},
			Rules: []rbacv1.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"pods"},
					Verbs:     []string{"get", "list", "watch"},
				},
			},
		}, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		_, err = nsClient.KubeClient().RbacV1().RoleBindings(ns.Name).Create(ctx, &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pods-role",
			},
			Subjects: []rbacv1.Subject{
				{
					APIGroup:  corev1.GroupName,
					Kind:      rbacv1.ServiceAccountKind,
					Namespace: podSA.Namespace,
					Name:      podSA.Name,
				},
			},
			RoleRef: rbacv1.RoleRef{
				APIGroup: rbacv1.GroupName,
				Kind:     "Role",
				Name:     podRole.Name,
			},
		}, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Creating backend PVC")
		backendPVC := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "backend-pvc",
				Namespace:    ns.GetName(),
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("100Mi"),
					},
				},
				StorageClassName: ptr.To(backendImmediateStorageClass.GetName()),
				VolumeMode:       ptr.To(corev1.PersistentVolumeFilesystem),
			},
		}

		backendPVC, err = nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Create(ctx, backendPVC, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for backend PVC to be bound")
		backendPVCBindingCtx, backendPVCBindingCtxCancel := context.WithTimeout(ctx, 2*time.Minute)
		defer backendPVCBindingCtxCancel()
		backendPVC, err = waitForPersistentVolumeClaimState(backendPVCBindingCtx, f.KubeAdminClient().CoreV1().PersistentVolumeClaims(ns.GetName()), backendPVC.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPersistentVolumeClaimBound)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(backendPVC.Spec.VolumeName).NotTo(o.BeEmpty())

		backendPV, err := c.KubeAdminClient().CoreV1().PersistentVolumes().Get(ctx, backendPVC.Spec.VolumeName, metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Creating proxy PVC")
		proxyPVC := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "proxy-pvc-",
				Namespace:    ns.GetName(),
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.VolumeResourceRequirements{
					Requests: corev1.ResourceList{
						// Doesn't matter, only a proxy directory is created.
						corev1.ResourceStorage: resource.MustParse("1Mi"),
					},
				},
				StorageClassName: ptr.To(proxyStorageClassName),
				VolumeMode:       ptr.To(corev1.PersistentVolumeFilesystem),
			},
		}

		proxyPVC, err = nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Create(ctx, proxyPVC, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		pod := getPodTemplate(ns.GetName(), proxyPVC.GetName())
		pod.Spec.ServiceAccountName = podSA.Name
		pod.Spec.Containers[0].Command = []string{
			"bin/sh",
			"-euEo",
			"pipefail",
			"-c",
			strings.TrimSpace(`
trap 'kill $( jobs -p ); exit 0' TERM

while true; do
	test -f "/var/lib/shared/backend-volume-mounting.done" && break
done
touch /data/test

sleep infinity &
wait $!
`),
		}

		// Set Pod's NodeSelector to match backend PV's NodeSelector.
		if backendPV.Spec.NodeAffinity != nil && backendPV.Spec.NodeAffinity.Required != nil {
			pod.Spec.Affinity = &corev1.Affinity{
				NodeAffinity: &corev1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: backendPV.Spec.NodeAffinity.Required.DeepCopy(),
				},
			}
		}

		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, corev1.VolumeMount{
			Name:      "shared",
			ReadOnly:  true,
			MountPath: "/var/lib/shared",
		})

		pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
			Name:  "wait",
			Image: "docker.io/rzetelskik/proxy-csi-driver:latest@sha256:7f22416a68afc8b16abd88d3cc5f9bfb399e83310946f118e9e7933a060276fa",
			Args: []string{
				"wait",
				"--pod-name=$(POD_NAME)",
				"--volume=data",
				"--signal-file-path=/var/lib/shared/backend-volume-mounting.done",
			},
			Env: []corev1.EnvVar{
				{
					Name: "POD_NAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "shared",
					MountPath: "/var/lib/shared",
				},
			},
			ImagePullPolicy: imagePullPolicy,
		})
		pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
			Name: "shared",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})

		framework.By("Creating pre-warmed Pod")
		pod, err = nsClient.KubeClient().CoreV1().Pods(ns.GetName()).Create(ctx, pod, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for Pod to be running")
		podRunningCtx, podRunningCtxCancel := context.WithTimeout(ctx, prewarmTimeout)
		defer podRunningCtxCancel()
		pod, err = socontrollerhelpers.WaitForPodState(podRunningCtx, nsClient.KubeClient().CoreV1().Pods(ns.GetName()), pod.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPodRunning)
		o.Expect(err).NotTo(o.HaveOccurred())

		startTime := time.Now()

		annotateProxyPVCWithBackendPVCRef(ctx, nsClient.KubeClient().CoreV1(), pod, proxyPVC, backendPVC.GetName())

		pod, err = socontrollerhelpers.WaitForPodState(ctx, f.KubeAdminClient().CoreV1().Pods(ns.GetName()), pod.GetName(), socontrollerhelpers.WaitForStateOptions{}, isPodReady)
		stopTime := time.Now()
		o.Expect(err).NotTo(o.HaveOccurred())

		podReadyCondition, err := getPodCondition(pod, corev1.PodReady)
		o.Expect(err).NotTo(o.HaveOccurred())

		elapsedTime := stopTime.Sub(startTime)
		framework.Infof("Time elapsed: %v, stop time %v, pod ready condition timestamp %v", elapsedTime, stopTime, podReadyCondition.LastTransitionTime)

		encoder := json.NewEncoder(resultsFile)
		res := result{ElapsedTimeMs: elapsedTime.Milliseconds()}
		err = encoder.Encode(res)
		o.Expect(err).NotTo(o.HaveOccurred())
	})
})

func waitForPersistentVolumeClaimState(ctx context.Context, client corev1client.PersistentVolumeClaimInterface, name string, options socontrollerhelpers.WaitForStateOptions, condition func(*corev1.PersistentVolumeClaim) (bool, error), additionalConditions ...func(*corev1.PersistentVolumeClaim) (bool, error)) (*corev1.PersistentVolumeClaim, error) {
	return socontrollerhelpers.WaitForObjectState[*corev1.PersistentVolumeClaim, *corev1.PersistentVolumeClaimList](ctx, client, name, options, condition, additionalConditions...)
}

func isPersistentVolumeClaimBound(claim *corev1.PersistentVolumeClaim) (bool, error) {
	return len(claim.Spec.VolumeName) != 0, nil
}

func getPodTemplate(namespace string, dataPVCName string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "data",
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: dataPVCName,
							ReadOnly:  false,
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:  "sleep",
					Image: "docker.io/library/busybox:latest@sha256:498a000f370d8c37927118ed80afe8adc38d1edcbfc071627d17b25c88efcab0",
					// Exit 1 to enforce setting the command explicitly.
					Command: []string{
						"exit",
						"1",
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:             "data",
							ReadOnly:         false,
							MountPath:        "/data",
							MountPropagation: ptr.To(corev1.MountPropagationHostToContainer),
						},
					},
					ImagePullPolicy: imagePullPolicy,
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							Exec: &corev1.ExecAction{
								Command: []string{
									"test",
									"-f",
									"/data/test",
								},
							},
						},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      300,
						PeriodSeconds:       1,
						SuccessThreshold:    1,
						FailureThreshold:    300,
					},
				},
			},
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    ptr.To[int64](65534),
				RunAsGroup:   ptr.To[int64](65534),
				RunAsNonRoot: ptr.To(true),
				FSGroup:      ptr.To[int64](65534),
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
			//Affinity: &corev1.Affinity{
			//	NodeAffinity: &corev1.NodeAffinity{
			//		RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
			//			NodeSelectorTerms: []corev1.NodeSelectorTerm{
			//				{
			//					MatchExpressions: []corev1.NodeSelectorRequirement{
			//						{
			//							Key:      "scylla.scylladb.com/node-type",
			//							Operator: corev1.NodeSelectorOpIn,
			//							Values:   []string{"scylla"},
			//						},
			//					},
			//				},
			//			},
			//		},
			//	},
			//},
			Tolerations: []corev1.Toleration{
				{
					Key:      "scylla-operator.scylladb.com/dedicated",
					Operator: corev1.TolerationOpEqual,
					Value:    "scyllaclusters",
					Effect:   corev1.TaintEffectNoSchedule,
				},
			},
		},
	}
}

func getPodCondition(pod *corev1.Pod, conditionType corev1.PodConditionType) (*corev1.PodCondition, error) {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == conditionType {
			return &cond, nil
		}
	}

	return nil, fmt.Errorf("pod condition of type %q not found", conditionType)
}

func isPodReady(pod *corev1.Pod) (bool, error) {
	return isPodStatusConditionPresentAndTrue(pod.Status.Conditions, corev1.PodReady), nil
}

func isPodStatusConditionPresentAndTrue(conditions []corev1.PodCondition, conditionType corev1.PodConditionType) bool {
	return isPodStatusConditionPresentAndEqual(conditions, conditionType, corev1.ConditionTrue)
}

func isPodStatusConditionPresentAndEqual(conditions []corev1.PodCondition, conditionType corev1.PodConditionType, status corev1.ConditionStatus) bool {
	for _, condition := range conditions {
		if condition.Type != conditionType {
			continue
		}

		return condition.Status == status
	}

	return false
}

func isPodRunning(pod *corev1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case corev1.PodRunning:
		return true, nil
	case corev1.PodFailed, corev1.PodSucceeded:
		return false, fmt.Errorf("pod unexpectedly ran to completion")
	default:
		return false, nil
	}
}

func annotateProxyPVCWithBackendPVCRef(ctx context.Context, client corev1client.CoreV1Interface, pod *corev1.Pod, proxyPVC *corev1.PersistentVolumeClaim, backendPVCName string) {
	framework.By("Annotating proxy PVC to mount backend storage")
	proxyPVC, err := client.PersistentVolumeClaims(proxyPVC.GetNamespace()).Patch(ctx, proxyPVC.GetName(), types.MergePatchType, []byte(fmt.Sprintf(`{"metadata": {"annotations": {%q: %q} } }`, proxycsinaming.DelayedStorageBackendPersistentVolumeClaimRefAnnotation, backendPVCName)), metav1.PatchOptions{})
	o.Expect(err).NotTo(o.HaveOccurred())
}
