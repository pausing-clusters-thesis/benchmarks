package utils

import (
	"context"

	scyllav1alpha1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1alpha1"
	scyllav1alpha1client "github.com/scylladb/scylla-operator/pkg/client/scylla/clientset/versioned/typed/scylla/v1alpha1"
	socontrollerhelpers "github.com/scylladb/scylla-operator/pkg/controllerhelpers"
)

func WaitForScyllaDBDatacenterState(ctx context.Context, client scyllav1alpha1client.ScyllaDBDatacenterInterface, name string, options socontrollerhelpers.WaitForStateOptions, condition func(sdc *scyllav1alpha1.ScyllaDBDatacenter) (bool, error), additionalConditions ...func(sdc *scyllav1alpha1.ScyllaDBDatacenter) (bool, error)) (*scyllav1alpha1.ScyllaDBDatacenter, error) {
	return socontrollerhelpers.WaitForObjectState[*scyllav1alpha1.ScyllaDBDatacenter, *scyllav1alpha1.ScyllaDBDatacenterList](ctx, client, name, options, condition, additionalConditions...)
}
