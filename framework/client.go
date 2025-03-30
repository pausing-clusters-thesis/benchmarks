package framework

import (
	o "github.com/onsi/gomega"
	pausingclient "github.com/pausing-clusters-thesis/pausable-scylladb-operator/pkg/client/pausing/clientset/versioned"
	soframework "github.com/scylladb/scylla-operator/test/e2e/framework"
)

type Client struct {
	soframework.Client
}

var _ soframework.GenericClientInterface = (*Client)(nil)
var _ soframework.ClientInterface = (*Client)(nil)

func (c *Client) PausingClient() *pausingclient.Clientset {
	cs, err := pausingclient.NewForConfig(c.ClientConfig())
	o.Expect(err).NotTo(o.HaveOccurred())
	return cs
}
