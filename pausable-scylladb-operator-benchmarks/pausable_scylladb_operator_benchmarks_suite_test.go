package pausable_scylladb_operator_benchmarks_test

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"slices"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/gocql/gocql/scyllacloud"
	g "github.com/onsi/ginkgo/v2"
	o "github.com/onsi/gomega"
	"github.com/pausing-clusters-thesis/benchmarks/framework"
	"github.com/pausing-clusters-thesis/benchmarks/naming"
	"github.com/pausing-clusters-thesis/benchmarks/utils"
	pausingv1alpha1 "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/api/pausing/v1alpha1"
	psocontrollerhelpers "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/controllerhelpers"
	psonaming "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/naming"
	proxycsinaming "github.com/pausing-clusters-thesis/proxy-csi-driver/pkg/naming"
	"github.com/scylladb/gocqlx/v2"
	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
	"github.com/scylladb/scylla-operator/pkg/genericclioptions"
	sonaming "github.com/scylladb/scylla-operator/pkg/naming"
	"github.com/scylladb/scylla-operator/pkg/scheme"
	cqlclientv1alpha1 "github.com/scylladb/scylla-operator/pkg/scylla/api/cqlclient/v1alpha1"
	soframework "github.com/scylladb/scylla-operator/test/e2e/framework"
	sotestutils "github.com/scylladb/scylla-operator/test/e2e/utils"
	scyllaclusterverification "github.com/scylladb/scylla-operator/test/e2e/utils/verification/scyllacluster"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
)

type result struct {
	ElapsedTimeMs     int64 `json:"elapsed_time_ms"`
	ApplicationTimeMs int64 `json:"application_time_ms"`
	OverheadTimeMs    int64 `json:"overhead_time_ms"`
}

var (
	scyllaProcessStartRegex = regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-][0-9]{2}:[0-9]{ 2}))\s+INFO\s+\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:,\d+)? starting ScyllaDB\.{3}`)
	scyllaServingStartRegex = regexp.MustCompile(`(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-][0-9]{2}:[0-9]{ 2}))\s+INFO\s+\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(?:,\d+)? \[shard \d+:main\] init - serving`)
)

var (
	clientConfig = genericclioptions.NewClientConfig("pausable-scylladb-operator-benchmarks")

	topologyZoneLabelValue   string
	proxyStorageClassName    string
	backendCSIDriverName     string
	nodeCount                int
	ingressClassName         string
	ingressControllerAddress string
	destDir                  string
)

var supportedBackendCSIDriverNames = []string{
	naming.GCEPersistentDiskCSIDriverName,
	naming.EBSCSIDriverName,
}

func init() {
	flag.StringVar(&topologyZoneLabelValue, "topology-zone-label-value", topologyZoneLabelValue, "The value of the topology zone label (optional).")
	flag.StringVar(&proxyStorageClassName, "proxy-storage-class-name", proxyStorageClassName, "The name of a StorageClass provisioned by the Proxy CSI Driver to be used in the test.")
	flag.StringVar(&backendCSIDriverName, "backend-csi-driver-name", backendCSIDriverName, fmt.Sprintf("The name of the backend CSI driver to test. Supported drvier names are: %v.", supportedBackendCSIDriverNames))
	flag.IntVar(&nodeCount, "nodes", nodeCount, "The number of nodes in the cluster.")
	flag.StringVar(&ingressClassName, "ingress-class-name", ingressClassName, "Name of the IngressClass to use to configure CQL backends.")
	flag.StringVar(&ingressControllerAddress, "ingress-controller-address", ingressControllerAddress, "Overrides destination address when sending testing data to applications behind ingresses.")
	flag.StringVar(&destDir, "dest-dir", destDir, "Destination directory in which results should be saved.")
}

func TestPausableScylladbOperatorBenchmarks(t *testing.T) {
	var err error

	err = Validate()
	if err != nil {
		t.Fatal(err)
	}

	err = Complete()
	if err != nil {
		t.Fatal(err)
	}

	soframework.TestContext = &soframework.TestContextType{
		RestConfigs:   []*restclient.Config{clientConfig.RestConfig},
		CleanupPolicy: soframework.CleanupPolicyAlways,
	}

	o.RegisterFailHandler(g.Fail)

	g.RunSpecs(t, "PausableScyllaDBOperatorBenchmarks Suite")
}

func Validate() error {
	var errs []error

	if len(proxyStorageClassName) == 0 {
		errs = append(errs, fmt.Errorf("proxy-storage-class-name must not be empty"))
	}

	if nodeCount <= 0 {
		errs = append(errs, fmt.Errorf("nodes must not be less or equal to zero"))
	}

	if len(ingressControllerAddress) == 0 {
		errs = append(errs, fmt.Errorf("ingress-class-name can't be empty"))
	}

	if len(ingressControllerAddress) == 0 {
		errs = append(errs, fmt.Errorf("ingress-controller-address can't be empty"))
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

	return nil
}

var _ = g.Describe("measure time to readiness", func() {
	f := framework.NewFramework("benchmark")

	type scenarioEntry struct {
		capacity        int32
		limit           int32
		resultsFileName string
	}

	g.DescribeTable("when unpausing PausableScyllaDBDatacenter", func(ctx g.SpecContext, se *scenarioEntry) {
		resultsFilePath := path.Join(destDir, se.resultsFileName)
		resultsFile, err := os.OpenFile(resultsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		o.Expect(err).NotTo(o.HaveOccurred())
		g.DeferCleanup(resultsFile.Close)

		c := f.Cluster(0)
		ns, nsClient, ok := c.DefaultNamespaceIfAny()
		o.Expect(ok).To(o.BeTrue())

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

		framework.By("Creating a ScyllaDBDatacenterPool")
		sdcp := &pausingv1alpha1.ScyllaDBDatacenterPool{
			ObjectMeta: metav1.ObjectMeta{
				Name: "basic",
			},
			Spec: pausingv1alpha1.ScyllaDBDatacenterPoolSpec{
				Template: pausingv1alpha1.ScyllaDBDatacenterTemplate{
					Spec: getScyllaDBDatacenterSpec(backendImmediateStorageClass.GetName()),
				},
				Capacity:              se.capacity,
				Limit:                 se.limit,
				ProxyStorageClassName: proxyStorageClassName,
			},
		}

		sdcp.Spec.Template.Spec.CertificateOptions = &scyllav1alpha1.CertificateOptions{
			ServingCA: &scyllav1alpha1.TLSCertificateAuthority{
				Type: scyllav1alpha1.TLSCertificateAuthorityTypeUserManaged,
				UserManagedOptions: &scyllav1alpha1.UserManagedTLSCertificateAuthorityOptions{
					SecretName: "",
				},
			},
			ClientCA: &scyllav1alpha1.TLSCertificateAuthority{
				Type: scyllav1alpha1.TLSCertificateAuthorityTypeUserManaged,
				UserManagedOptions: &scyllav1alpha1.UserManagedTLSCertificateAuthorityOptions{
					SecretName: "",
				},
			},
		}

		sdcp, err = c.PausingAdminClient().PausingV1alpha1().ScyllaDBDatacenterPools(ns.GetName()).Create(ctx, sdcp, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for ScyllaDBDatacenterPool to roll out")
		// TODO: context
		sdcp, err = psocontrollerhelpers.WaitForScyllaDBDatacenterPoolState(ctx, c.PausingAdminClient().PausingV1alpha1().ScyllaDBDatacenterPools(ns.GetName()), sdcp.GetName(), socontrollerhelpers.WaitForStateOptions{}, psocontrollerhelpers.IsScyllaDBDatacenterPoolRolledOut)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Creating a PausableScyllaDBDatacenter in an unpaused state")
		psdc := &pausingv1alpha1.PausableScyllaDBDatacenter{
			ObjectMeta: metav1.ObjectMeta{
				Name: "basic",
			},
			Spec: pausingv1alpha1.PausableScyllaDBDatacenterSpec{
				ScyllaDBDatacenterPoolName: sdcp.GetName(),
				Paused:                     ptr.To(false),
				ExposeOptions: &pausingv1alpha1.ExposeOptions{
					CQL: &scyllav1alpha1.CQLExposeOptions{
						Ingress: &scyllav1alpha1.CQLExposeIngressOptions{
							IngressClassName: ingressClassName,
						},
					},
				},
			},
		}

		psdc, err = c.PausingAdminClient().PausingV1alpha1().PausableScyllaDBDatacenters(ns.GetName()).Create(ctx, psdc, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for PausableScyllaDBDatacenter to roll out")
		// TODO: context
		psdc, err = psocontrollerhelpers.WaitForPausableScyllaDBDatacenterState(ctx, c.PausingAdminClient().PausingV1alpha1().PausableScyllaDBDatacenters(ns.GetName()), psdc.GetName(), socontrollerhelpers.WaitForStateOptions{}, psocontrollerhelpers.IsPausableScyllaDBDatacenterRolledOut)
		o.Expect(err).NotTo(o.HaveOccurred())
		// TODO: verify

		sdcc, err := c.PausingAdminClient().PausingV1alpha1().ScyllaDBDatacenterClaims(ns.GetName()).Get(ctx, psonaming.GetScyllaDBDatacenterClaimNameForPausableScyllaDBDatacenter(psdc), metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		o.Expect(sdcc.Status.ScyllaDBDatacenterName).NotTo(o.BeNil())
		o.Expect(*sdcc.Status.ScyllaDBDatacenterName).NotTo(o.BeEmpty())
		sdcName := *sdcc.Status.ScyllaDBDatacenterName

		sdc, err := c.ScyllaAdminClient().ScyllaV1alpha1().ScyllaDBDatacenters(ns.GetName()).Get(ctx, sdcName, metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(sdc.Spec.DNSDomains).To(o.HaveLen(1))
		dnsDomain := sdc.Spec.DNSDomains[0]

		connectionBundleDir, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("connection-bundle-%s-", ns.GetName()))
		o.Expect(err).NotTo(o.HaveOccurred())
		defer func() {
			err := os.RemoveAll(connectionBundleDir)
			o.Expect(err).NotTo(o.HaveOccurred())
		}()

		var groupVersioner runtime.GroupVersioner = schema.GroupVersions([]schema.GroupVersion{cqlclientv1alpha1.GroupVersion})
		decoder := scheme.Codecs.DecoderToVersion(scheme.Codecs.UniversalDeserializer(), groupVersioner)
		encoder := scheme.Codecs.EncoderForVersion(scheme.DefaultYamlSerializer, groupVersioner)

		framework.By("Injecting ingress controller address into the CQL connection bundle")
		bundleSecret, err := nsClient.KubeClient().CoreV1().Secrets(sdc.Namespace).Get(ctx, sonaming.GetScyllaClusterLocalAdminCQLConnectionConfigsName(sdc.Name), metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		cqlConnectionConfig := &cqlclientv1alpha1.CQLConnectionConfig{}
		_, _, err = decoder.Decode(bundleSecret.Data[dnsDomain], nil, cqlConnectionConfig)
		o.Expect(err).NotTo(o.HaveOccurred())

		gossipDatacenterName := sonaming.GetScyllaDBDatacenterGossipDatacenterName(sdc)
		cqlConnectionConfig.Datacenters[gossipDatacenterName].Server = ingressControllerAddress

		cqlConnectionConfigData, err := runtime.Encode(encoder, cqlConnectionConfig)
		o.Expect(err).NotTo(o.HaveOccurred())

		cqlConnectionConfigFilePath := path.Join(connectionBundleDir, fmt.Sprintf("%s.yaml", dnsDomain))
		framework.By("Saving CQL Connection Config in %s", cqlConnectionConfigFilePath)
		err = os.WriteFile(cqlConnectionConfigFilePath, cqlConnectionConfigData, 0600)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Connecting to cluster via Ingress")
		cluster, err := scyllacloud.NewCloudCluster(cqlConnectionConfigFilePath)
		o.Expect(err).NotTo(o.HaveOccurred())

		// Increase default timeout, due to additional hop on the route to host.
		cluster.Timeout = 10 * time.Second
		cluster.Logger = nopLogger{}
		cluster.Consistency = gocql.Quorum

		session, err := gocqlx.WrapSession(cluster.CreateSession())
		o.Expect(err).NotTo(o.HaveOccurred())

		// When overriding the sessions, "hosts" are only used by the data inserter to determine a replication factor.
		hosts := slices.Repeat([]string{""}, nodeCount)

		di, err := sotestutils.NewDataInserter(hosts, sotestutils.WithSession(&session))
		o.Expect(err).NotTo(o.HaveOccurred())

		scyllaclusterverification.InsertAndVerifyCQLDataUsingDataInserter(ctx, di)
		session.Close()
		di.ForceSession(nil)

		framework.By("Getting VolumeAttachments for backend PVCs")
		var backendVAs []*storagev1.VolumeAttachment
		o.Expect(sdcp.Spec.Template.Spec.Racks).To(o.HaveLen(1))
		for i := range nodeCount {
			podName := sonaming.MemberServiceName(sdcp.Spec.Template.Spec.Racks[0], sdc, i)
			pod, err := nsClient.KubeClient().CoreV1().Pods(ns.GetName()).Get(ctx, podName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(pod.Spec.NodeName).NotTo(o.BeEmpty())

			backendPVCName := psonaming.GetBackendPersistentVolumeClaimNameForPausableScyllaDBDatacenterMember(psdc.Name, sdcp.Spec.Template.Spec.Racks[0].Name, int32(i))
			backendPVC, err := nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Get(ctx, backendPVCName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendPVC.Spec.VolumeName).NotTo(o.BeEmpty())

			backendPV, err := c.KubeAdminClient().CoreV1().PersistentVolumes().Get(ctx, backendPVC.Spec.VolumeName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendPV.Spec.CSI).NotTo(o.BeNil())

			backendVAName := getAttachmentName(backendPV.Spec.CSI.VolumeHandle, backendPV.Spec.CSI.Driver, pod.Spec.NodeName)
			backendVA, err := c.KubeAdminClient().StorageV1().VolumeAttachments().Get(ctx, backendVAName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendVA.Status.Attached).To(o.BeTrue())
			backendVAs = append(backendVAs, backendVA)
		}

		framework.By("Pausing PausableScyllaDBDatacenter")
		psdc, err = c.PausingAdminClient().PausingV1alpha1().PausableScyllaDBDatacenters(ns.GetName()).Patch(
			ctx,
			psdc.Name,
			types.JSONPatchType,
			[]byte(`[{"op": "replace", "path": "/spec/paused", "value": true}]`),
			metav1.PatchOptions{},
		)
		o.Expect(err).NotTo(o.HaveOccurred())
		o.Expect(psdc.Spec.Paused).NotTo(o.BeNil())
		o.Expect(*psdc.Spec.Paused).To(o.BeTrue())

		framework.By("Waiting for PausableScyllaDBDatacenter to roll out")
		// TODO: context
		psdc, err = psocontrollerhelpers.WaitForPausableScyllaDBDatacenterState(ctx, c.PausingAdminClient().PausingV1alpha1().PausableScyllaDBDatacenters(ns.GetName()), psdc.GetName(), socontrollerhelpers.WaitForStateOptions{}, psocontrollerhelpers.IsPausableScyllaDBDatacenterRolledOut)
		o.Expect(err).NotTo(o.HaveOccurred())
		// TODO: verify

		framework.By("Wait for ScyllaDBDatacenter to be deleted")
		// TODO: context
		err = soframework.WaitForObjectDeletion(
			ctx,
			nsClient.DynamicClient(),
			scyllav1alpha1.GroupVersion.WithResource("scylladbdatacenters"),
			ns.GetName(),
			sdc.GetName(),
			ptr.To(sdc.GetUID()),
		)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for backend PVCs to be unbound from proxy PVCs")
		o.Expect(sdcp.Spec.Template.Spec.Racks).To(o.HaveLen(1))
		for i := range int32(nodeCount) {
			backendPVCName := psonaming.GetBackendPersistentVolumeClaimNameForPausableScyllaDBDatacenterMember(psdc.Name, sdcp.Spec.Template.Spec.Racks[0].Name, i)
			_, err = socontrollerhelpers.WaitForPVCState(ctx, nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()), backendPVCName, socontrollerhelpers.WaitForStateOptions{}, IsBackendPersistentVolumeClaimUnboundFromProxyPersistentVolumeClaim)
			o.Expect(err).NotTo(o.HaveOccurred())
		}

		framework.By("Waiting for backend VolumeAttachments to be deleted")
		for _, va := range backendVAs {
			err = soframework.WaitForObjectDeletion(
				ctx,
				c.DynamicAdminClient(),
				storagev1.SchemeGroupVersion.WithResource("volumeattachments"),
				"",
				va.GetName(),
				ptr.To(va.GetUID()),
			)
			o.Expect(err).NotTo(o.HaveOccurred())
		}

		framework.By("Waiting for ScyllaDBDatacenterPool to roll out")
		// TODO: context
		sdcp, err = psocontrollerhelpers.WaitForScyllaDBDatacenterPoolState(ctx, c.PausingAdminClient().PausingV1alpha1().ScyllaDBDatacenterPools(ns.GetName()), sdcp.GetName(), socontrollerhelpers.WaitForStateOptions{}, psocontrollerhelpers.IsScyllaDBDatacenterPoolRolledOut)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Connecting to the paused cluster via Ingress")
		startTime := time.Now()
		newSession, err := gocqlx.WrapSession(cluster.CreateSession())
		stopTime := time.Now()
		framework.By("Session created successfully")
		o.Expect(err).NotTo(o.HaveOccurred())

		di.ForceSession(&newSession)

		scyllaclusterverification.VerifyCQLData(ctx, di)
		newSession.Close()

		framework.By("Waiting for PausableScyllaDBDatacenter to roll out")
		// TODO: context
		psdc, err = psocontrollerhelpers.WaitForPausableScyllaDBDatacenterState(ctx, c.PausingAdminClient().PausingV1alpha1().PausableScyllaDBDatacenters(ns.GetName()), psdc.GetName(), socontrollerhelpers.WaitForStateOptions{}, psocontrollerhelpers.IsPausableScyllaDBDatacenterRolledOut)
		o.Expect(err).NotTo(o.HaveOccurred())
		// TODO: verify

		sdcc, err = c.PausingAdminClient().PausingV1alpha1().ScyllaDBDatacenterClaims(ns.GetName()).Get(ctx, psonaming.GetScyllaDBDatacenterClaimNameForPausableScyllaDBDatacenter(psdc), metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		o.Expect(sdcc.Status.ScyllaDBDatacenterName).NotTo(o.BeNil())
		o.Expect(*sdcc.Status.ScyllaDBDatacenterName).NotTo(o.BeEmpty())
		sdcName = *sdcc.Status.ScyllaDBDatacenterName
		sdc, err = c.ScyllaAdminClient().ScyllaV1alpha1().ScyllaDBDatacenters(ns.GetName()).Get(ctx, sdcName, metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		scyllaStartToServeTime := getScyllaContainerTimeToServeCQL(ctx, nsClient.KubeClient().CoreV1().Pods(ns.GetName()), sdc)
		o.Expect(scyllaStartToServeTime).NotTo(o.BeZero())

		totalTime := stopTime.Sub(startTime)
		platformTime := totalTime - scyllaStartToServeTime
		framework.Infof("Total time: %v.\nPlatform time: %v.\nApplication time: %v.\n", totalTime, platformTime, scyllaStartToServeTime)

		jsonEncoder := json.NewEncoder(resultsFile)
		res := result{
			ElapsedTimeMs:     totalTime.Milliseconds(),
			ApplicationTimeMs: scyllaStartToServeTime.Milliseconds(),
			OverheadTimeMs:    platformTime.Milliseconds(),
		}
		err = jsonEncoder.Encode(res)
		o.Expect(err).NotTo(o.HaveOccurred())
	},
		g.Entry("with pre-warmed ScyllaDBDatacenters", &scenarioEntry{
			capacity:        1,
			limit:           1,
			resultsFileName: "prewarmed",
		}),
		g.Entry("with cold-started ScyllaDBDatacenters", &scenarioEntry{
			capacity:        0,
			limit:           0,
			resultsFileName: "cold",
		}),
	)

	g.It("cold-starting with Scylla Operator", func(ctx g.SpecContext) {
		const resultsFileName = "baseline"
		resultsFilePath := path.Join(destDir, resultsFileName)
		resultsFile, err := os.OpenFile(resultsFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		o.Expect(err).NotTo(o.HaveOccurred())
		g.DeferCleanup(resultsFile.Close)

		c := f.Cluster(0)
		ns, nsClient, ok := c.DefaultNamespaceIfAny()
		o.Expect(ok).To(o.BeTrue())

		backendStorageClass, err := utils.GetImmediateStorageClassForCSIDriver(backendCSIDriverName)
		o.Expect(err).NotTo(o.HaveOccurred())

		// Switch to wait for first consumer volume binding mode.
		backendStorageClass.VolumeBindingMode = ptr.To(storagev1.VolumeBindingWaitForFirstConsumer)

		framework.By("Creating immediate StorageClass for backend CSI driver")
		backendStorageClass, err = c.KubeAdminClient().StorageV1().StorageClasses().Create(ctx, backendStorageClass, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		g.DeferCleanup(func(ctx g.SpecContext, backendImmediateStorageClass *storagev1.StorageClass) {
			framework.By("Deleting immediate StorageClass")
			err := c.KubeAdminClient().StorageV1().StorageClasses().Delete(ctx, backendImmediateStorageClass.GetName(), metav1.DeleteOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
		}, backendStorageClass)

		originalSDC := &scyllav1alpha1.ScyllaDBDatacenter{
			ObjectMeta: metav1.ObjectMeta{
				Name: "basic",
			},
			Spec: getScyllaDBDatacenterSpec(backendStorageClass.GetName()),
		}
		originalSDC.Spec.DNSDomains = []string{fmt.Sprintf("%s.public.nodes.scylladb.com", ns.GetName())}
		if originalSDC.Spec.ExposeOptions != nil {
			originalSDC.Spec.ExposeOptions = &scyllav1alpha1.ExposeOptions{}
		}
		originalSDC.Spec.ExposeOptions.CQL = &scyllav1alpha1.CQLExposeOptions{
			Ingress: &scyllav1alpha1.CQLExposeIngressOptions{
				IngressClassName: ingressClassName,
			},
		}

		sdc := originalSDC.DeepCopy()

		framework.By("Creating the ScyllaDBDatacenter")
		sdc, err = nsClient.ScyllaClient().ScyllaV1alpha1().ScyllaDBDatacenters(ns.GetName()).Create(ctx, sdc, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for ScyllaDBDatacenter to roll out")
		// TODO: context
		sdc, err = utils.WaitForScyllaDBDatacenterState(ctx, nsClient.ScyllaClient().ScyllaV1alpha1().ScyllaDBDatacenters(ns.GetName()), sdc.GetName(), socontrollerhelpers.WaitForStateOptions{}, socontrollerhelpers.IsScyllaDBDatacenterRolledOut)
		o.Expect(err).NotTo(o.HaveOccurred())

		o.Expect(sdc.Spec.DNSDomains).To(o.HaveLen(1))
		dnsDomain := sdc.Spec.DNSDomains[0]

		connectionBundleDir, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("connection-bundle-%s-", ns.GetName()))
		o.Expect(err).NotTo(o.HaveOccurred())
		defer func() {
			err := os.RemoveAll(connectionBundleDir)
			o.Expect(err).NotTo(o.HaveOccurred())
		}()

		var groupVersioner runtime.GroupVersioner = schema.GroupVersions([]schema.GroupVersion{cqlclientv1alpha1.GroupVersion})
		decoder := scheme.Codecs.DecoderToVersion(scheme.Codecs.UniversalDeserializer(), groupVersioner)
		encoder := scheme.Codecs.EncoderForVersion(scheme.DefaultYamlSerializer, groupVersioner)

		framework.By("Injecting ingress controller address into the CQL connection bundle")
		bundleSecret, err := nsClient.KubeClient().CoreV1().Secrets(sdc.Namespace).Get(ctx, sonaming.GetScyllaClusterLocalAdminCQLConnectionConfigsName(sdc.Name), metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		cqlConnectionConfig := &cqlclientv1alpha1.CQLConnectionConfig{}
		_, _, err = decoder.Decode(bundleSecret.Data[dnsDomain], nil, cqlConnectionConfig)
		o.Expect(err).NotTo(o.HaveOccurred())

		gossipDatacenterName := sonaming.GetScyllaDBDatacenterGossipDatacenterName(sdc)
		cqlConnectionConfig.Datacenters[gossipDatacenterName].Server = ingressControllerAddress

		cqlConnectionConfigData, err := runtime.Encode(encoder, cqlConnectionConfig)
		o.Expect(err).NotTo(o.HaveOccurred())

		cqlConnectionConfigFilePath := path.Join(connectionBundleDir, fmt.Sprintf("%s.yaml", dnsDomain))
		framework.By("Saving CQL Connection Config in %s", cqlConnectionConfigFilePath)
		err = os.WriteFile(cqlConnectionConfigFilePath, cqlConnectionConfigData, 0600)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Connecting to cluster via Ingress")
		cluster, err := scyllacloud.NewCloudCluster(cqlConnectionConfigFilePath)
		o.Expect(err).NotTo(o.HaveOccurred())

		// Increase default timeout, due to additional hop on the route to host.
		cluster.Timeout = 10 * time.Second
		cluster.Logger = nopLogger{}
		cluster.Consistency = gocql.Quorum

		session, err := gocqlx.WrapSession(cluster.CreateSession())
		o.Expect(err).NotTo(o.HaveOccurred())

		// When overriding the sessions, "hosts" are only used by the data inserter to determine a replication factor.
		hosts := slices.Repeat([]string{""}, nodeCount)

		di, err := sotestutils.NewDataInserter(hosts, sotestutils.WithSession(&session))
		o.Expect(err).NotTo(o.HaveOccurred())

		scyllaclusterverification.InsertAndVerifyCQLDataUsingDataInserter(ctx, di)
		session.Close()
		di.ForceSession(nil)

		framework.By("Getting VolumeAttachments for PVCs")
		var backendVAs []*storagev1.VolumeAttachment
		o.Expect(sdc.Spec.Racks).To(o.HaveLen(1))
		for i := range nodeCount {
			podName := sonaming.MemberServiceName(sdc.Spec.Racks[0], sdc, i)
			pod, err := nsClient.KubeClient().CoreV1().Pods(ns.GetName()).Get(ctx, podName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(pod.Spec.NodeName).NotTo(o.BeEmpty())

			backendPVCName := sonaming.PVCNameForPod(podName)
			backendPVC, err := nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Get(ctx, backendPVCName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendPVC.Spec.VolumeName).NotTo(o.BeEmpty())

			backendPV, err := c.KubeAdminClient().CoreV1().PersistentVolumes().Get(ctx, backendPVC.Spec.VolumeName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendPV.Spec.CSI).NotTo(o.BeNil())

			backendVAName := getAttachmentName(backendPV.Spec.CSI.VolumeHandle, backendPV.Spec.CSI.Driver, pod.Spec.NodeName)
			backendVA, err := c.KubeAdminClient().StorageV1().VolumeAttachments().Get(ctx, backendVAName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendVA.Status.Attached).To(o.BeTrue())
			backendVAs = append(backendVAs, backendVA)
		}

		framework.By("Deleting the ScyllaDBDatacenter")
		err = nsClient.ScyllaClient().ScyllaV1alpha1().ScyllaDBDatacenters(ns.GetName()).Delete(ctx, sdc.GetName(), metav1.DeleteOptions{
			Preconditions: &metav1.Preconditions{
				UID: ptr.To(sdc.UID),
			},
			PropagationPolicy: ptr.To(metav1.DeletePropagationForeground),
		})
		o.Expect(err).NotTo(o.HaveOccurred())

		// TODO: context
		err = soframework.WaitForObjectDeletion(
			ctx,
			nsClient.DynamicClient(),
			scyllav1alpha1.GroupVersion.WithResource("scylladbdatacenters"),
			ns.GetName(),
			sdc.GetName(),
			ptr.To(sdc.UID),
		)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Waiting for backend VolumeAttachments to be deleted")
		for _, va := range backendVAs {
			err = soframework.WaitForObjectDeletion(
				ctx,
				c.DynamicAdminClient(),
				storagev1.SchemeGroupVersion.WithResource("volumeattachments"),
				"",
				va.GetName(),
				ptr.To(va.GetUID()),
			)
			o.Expect(err).NotTo(o.HaveOccurred())
		}

		framework.By("Verifying PVs' presence")

		o.Expect(sdc.Spec.Racks).To(o.HaveLen(1))
		for i := range nodeCount {
			podName := sonaming.MemberServiceName(sdc.Spec.Racks[0], sdc, i)

			backendPVCName := sonaming.PVCNameForPod(podName)
			backendPVC, err := nsClient.KubeClient().CoreV1().PersistentVolumeClaims(ns.GetName()).Get(ctx, backendPVCName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendPVC.Spec.VolumeName).NotTo(o.BeEmpty())
			o.Expect(backendPVC.ObjectMeta.DeletionTimestamp).To(o.BeNil())

			backendPV, err := c.KubeAdminClient().CoreV1().PersistentVolumes().Get(ctx, backendPVC.Spec.VolumeName, metav1.GetOptions{})
			o.Expect(err).NotTo(o.HaveOccurred())
			o.Expect(backendPV.ObjectMeta.DeletionTimestamp).To(o.BeNil())
		}

		sdc = originalSDC.DeepCopy()

		framework.By("Redeploying the ScyllaDBDatacenter")
		startTime := time.Now()
		sdc, err = nsClient.ScyllaClient().ScyllaV1alpha1().ScyllaDBDatacenters(ns.GetName()).Create(ctx, sdc, metav1.CreateOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		sdc, err = utils.WaitForScyllaDBDatacenterState(ctx, nsClient.ScyllaClient().ScyllaV1alpha1().ScyllaDBDatacenters(ns.GetName()), sdc.GetName(), socontrollerhelpers.WaitForStateOptions{}, isScyllaDBDatacenterAvailable)
		stopTime := time.Now()
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Injecting ingress controller address into the CQL connection bundle")
		bundleSecret, err = nsClient.KubeClient().CoreV1().Secrets(sdc.Namespace).Get(ctx, sonaming.GetScyllaClusterLocalAdminCQLConnectionConfigsName(sdc.Name), metav1.GetOptions{})
		o.Expect(err).NotTo(o.HaveOccurred())

		cqlConnectionConfig = &cqlclientv1alpha1.CQLConnectionConfig{}
		_, _, err = decoder.Decode(bundleSecret.Data[dnsDomain], nil, cqlConnectionConfig)
		o.Expect(err).NotTo(o.HaveOccurred())

		gossipDatacenterName = sonaming.GetScyllaDBDatacenterGossipDatacenterName(sdc)
		cqlConnectionConfig.Datacenters[gossipDatacenterName].Server = ingressControllerAddress

		cqlConnectionConfigData, err = runtime.Encode(encoder, cqlConnectionConfig)
		o.Expect(err).NotTo(o.HaveOccurred())

		cqlConnectionConfigFilePath = path.Join(connectionBundleDir, fmt.Sprintf("%s.yaml", dnsDomain))
		framework.By("Saving CQL Connection Config in %s", cqlConnectionConfigFilePath)
		err = os.WriteFile(cqlConnectionConfigFilePath, cqlConnectionConfigData, 0600)
		o.Expect(err).NotTo(o.HaveOccurred())

		framework.By("Connecting to cluster via Ingress")
		cluster, err = scyllacloud.NewCloudCluster(cqlConnectionConfigFilePath)
		o.Expect(err).NotTo(o.HaveOccurred())

		// Increase default timeout, due to additional hop on the route to host.
		cluster.Timeout = 10 * time.Second
		cluster.Logger = nopLogger{}
		cluster.Consistency = gocql.Quorum

		newSession, err := gocqlx.WrapSession(cluster.CreateSession())
		o.Expect(err).NotTo(o.HaveOccurred())

		di.ForceSession(&newSession)

		scyllaclusterverification.VerifyCQLData(ctx, di)
		newSession.Close()

		scyllaStartToServeTime := getScyllaContainerTimeToServeCQL(ctx, nsClient.KubeClient().CoreV1().Pods(ns.GetName()), sdc)
		o.Expect(scyllaStartToServeTime).NotTo(o.BeZero())

		totalTime := stopTime.Sub(startTime)
		platformTime := totalTime - scyllaStartToServeTime
		framework.Infof("Total time: %v.\nPlatform time: %v.\nApplication time: %v.\n", totalTime, platformTime, scyllaStartToServeTime)

		jsonEncoder := json.NewEncoder(resultsFile)
		res := result{
			ElapsedTimeMs:     totalTime.Milliseconds(),
			ApplicationTimeMs: scyllaStartToServeTime.Milliseconds(),
			OverheadTimeMs:    platformTime.Milliseconds(),
		}
		err = jsonEncoder.Encode(res)
		o.Expect(err).NotTo(o.HaveOccurred())
	})
})

type nopLogger struct{}

func (n nopLogger) Print(_ ...interface{}) {}

func (n nopLogger) Printf(_ string, _ ...interface{}) {}

func (n nopLogger) Println(_ ...interface{}) {}

var _ gocql.StdLogger = (*nopLogger)(nil)

func IsBackendPersistentVolumeClaimUnboundFromProxyPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) (bool, error) {
	if psocontrollerhelpers.HasAnnotation(pvc, proxycsinaming.DelayedStorageProxyPersistentVolumeClaimRefAnnotation) {
		return false, nil
	}

	return true, nil
}

func getScyllaDBDatacenterSpec(backendImmediateStorageClassName string) scyllav1alpha1.ScyllaDBDatacenterSpec {
	sdcSpec := scyllav1alpha1.ScyllaDBDatacenterSpec{
		ClusterName:    "basic",
		DatacenterName: ptr.To("us-east-1"),
		ScyllaDB: scyllav1alpha1.ScyllaDB{
			Image:               "docker.io/scylladb/scylla:6.2.3@sha256:a9d904089abe9a4f8b5b893ebb5b5bf8b5a1bd0dc6658921cf05f89d3712289c",
			EnableDeveloperMode: ptr.To(false),
		},
		ScyllaDBManagerAgent: &scyllav1alpha1.ScyllaDBManagerAgent{
			Image: ptr.To("docker.io/scylladb/scylla-manager-agent:3.4.1@sha256:392ce6d3971ae077cc58b3cd2c7da1e9572f9f76223dfd5e11445c32e7ab0396"),
		},
		ExposeOptions: &scyllav1alpha1.ExposeOptions{
			NodeService: &scyllav1alpha1.NodeServiceTemplate{
				Type: scyllav1alpha1.NodeServiceTypeHeadless,
			},
			BroadcastOptions: &scyllav1alpha1.NodeBroadcastOptions{
				Nodes: scyllav1alpha1.BroadcastOptions{
					Type: scyllav1alpha1.BroadcastAddressTypePodIP,
				},
				Clients: scyllav1alpha1.BroadcastOptions{
					Type: scyllav1alpha1.BroadcastAddressTypePodIP,
				},
			},
		},
		RackTemplate: &scyllav1alpha1.RackTemplate{
			Nodes: ptr.To(int32(nodeCount)),
			ScyllaDB: &scyllav1alpha1.ScyllaDBTemplate{
				Resources: &corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("4Gi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("1"),
						corev1.ResourceMemory: resource.MustParse("4Gi"),
					},
				},
				Storage: &scyllav1alpha1.StorageOptions{
					Capacity:         "10Gi",
					StorageClassName: ptr.To(backendImmediateStorageClassName),
				},
			},
			ScyllaDBManagerAgent: &scyllav1alpha1.ScyllaDBManagerAgentTemplate{
				Resources: &corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("10m"),
						corev1.ResourceMemory: resource.MustParse("100Mi"),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("10m"),
						corev1.ResourceMemory: resource.MustParse("100Mi"),
					},
				},
			},
		},
		Racks: []scyllav1alpha1.RackSpec{
			{
				Name: "us-east-1a",
				RackTemplate: scyllav1alpha1.RackTemplate{
					Placement: &scyllav1alpha1.Placement{
						NodeAffinity: &corev1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
								NodeSelectorTerms: []corev1.NodeSelectorTerm{
									{
										MatchExpressions: []corev1.NodeSelectorRequirement{
											{
												Key:      "scylla.scylladb.com/node-type",
												Operator: corev1.NodeSelectorOpIn,
												Values:   []string{"scylla"},
											},
										},
									},
								},
							},
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "scylla-operator.scylladb.com/dedicated",
								Operator: corev1.TolerationOpEqual,
								Value:    "scyllaclusters",
								Effect:   corev1.TaintEffectNoSchedule,
							},
						},
					},
				},
			},
		},
		MinReadySeconds: ptr.To[int32](0),
		ReadinessGates: []corev1.PodReadinessGate{
			{
				ConditionType: psonaming.IngressControllerScyllaDBMemberPodConditionType,
			},
		},
	}

	if len(topologyZoneLabelValue) > 0 {
		sdcSpec.Racks[0].Placement.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions = append(
			sdcSpec.Racks[0].Placement.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions,
			corev1.NodeSelectorRequirement{
				Key:      corev1.LabelTopologyZone,
				Operator: corev1.NodeSelectorOpIn,
				Values:   []string{topologyZoneLabelValue},
			})
	}

	return sdcSpec
}

func getAttachmentName(volumeHandle string, driverName string, nodeName string) string {
	result := sha256.Sum256([]byte(fmt.Sprintf("%s%s%s", volumeHandle, driverName, nodeName)))
	return fmt.Sprintf("csi-%x", result)
}

func getScyllaContainerTimeToServeCQL(ctx context.Context, podClient corev1client.PodInterface, sdc *scyllav1alpha1.ScyllaDBDatacenter) time.Duration {
	stsName := sonaming.StatefulSetNameForRack(sdc.Spec.Racks[0], sdc)
	rackNodes, err := socontrollerhelpers.GetRackNodeCount(sdc, sdc.Spec.Racks[0].Name)
	o.Expect(err).NotTo(o.HaveOccurred())

	scyllaTimeToServeTimes := make([]time.Duration, nodeCount)
	for i := int32(0); i < *rackNodes; i++ {
		podName := fmt.Sprintf("%s-%d", stsName, i)

		var scyllaProcessStartTime, scyllaServingStartTime time.Time

		var logs io.ReadCloser
		logs, err = podClient.GetLogs(podName, &corev1.PodLogOptions{
			Container:  sonaming.ScyllaContainerName,
			Timestamps: true,
		}).Stream(ctx)
		o.Expect(err).NotTo(o.HaveOccurred())

		scanner := bufio.NewScanner(logs)
		for scanner.Scan() {
			line := scanner.Text()
			if matches := scyllaProcessStartRegex.FindStringSubmatch(line); matches != nil && len(matches) > 1 {
				scyllaProcessStartTime, err = time.Parse(time.RFC3339Nano, matches[1])
				o.Expect(err).NotTo(o.HaveOccurred())
			} else if matches = scyllaServingStartRegex.FindStringSubmatch(line); matches != nil && len(matches) > 1 {
				scyllaServingStartTime, err = time.Parse(time.RFC3339Nano, matches[1])
				o.Expect(err).NotTo(o.HaveOccurred())
			}
		}

		o.Expect(scyllaProcessStartTime.IsZero()).To(o.BeFalse())
		o.Expect(scyllaServingStartTime.IsZero()).To(o.BeFalse())
		scyllaTimeToServeTimes[i] = scyllaServingStartTime.Sub(scyllaProcessStartTime)
	}

	return slices.Max(scyllaTimeToServeTimes)
}

func isScyllaDBDatacenterAvailable(sdc *scyllav1alpha1.ScyllaDBDatacenter) (bool, error) {
	return psocontrollerhelpers.IsScyllaDBDatacenterAvailable(sdc), nil
}
