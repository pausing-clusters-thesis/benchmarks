package framework

import (
	soframework "github.com/scylladb/scylla-operator/test/e2e/framework"
)

type Framework struct {
	soframework.Framework
}

func NewFramework(namePrefix string) *Framework {
	return &Framework{
		Framework: *soframework.NewFramework(namePrefix),
	}
}

func (f *Framework) Cluster(idx int) *Cluster {
	return &Cluster{
		ClusterInterface: f.Framework.Cluster(idx),
	}
}
