package framework

import (
	"context"

	o "github.com/onsi/gomega"
	pausingclient "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned"
	soframework "github.com/scylladb/scylla-operator/test/e2e/framework"
	corev1 "k8s.io/api/core/v1"
)

type Cluster struct {
	soframework.ClusterInterface
}

func (c *Cluster) PausingAdminClient() *pausingclient.Clientset {
	cs, err := pausingclient.NewForConfig(c.AdminClientConfig())
	o.Expect(err).NotTo(o.HaveOccurred())
	return cs
}

func (c *Cluster) PausingClient() *pausingclient.Clientset {
	cs, err := pausingclient.NewForConfig(c.AdminClientConfig())
	o.Expect(err).NotTo(o.HaveOccurred())
	return cs
}

func (c *Cluster) CreateUserNamespace(ctx context.Context) (*corev1.Namespace, Client) {
	ns, nsClient := c.ClusterInterface.CreateUserNamespace(ctx)
	return ns, Client{
		Client: nsClient,
	}
}
