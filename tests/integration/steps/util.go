package steps

import (
	"fmt"
	cassdcapi "github.com/k8ssandra/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var testClient client.Client

// InitTestClient initializes a controller-runtime client. This is a no-op if
// it has already been called. It should be called prior to any test execution.
func InitTestClient() error {
	err := cassdcapi.AddToScheme(scheme.Scheme)
	if err != nil {
		return err
	}
	config, err := ctrl.GetConfig()
	if err != nil {
		return fmt.Errorf("cannot contact Kind test cluster: %w", err)
	}
	testClient, err = client.New(config, client.Options{Scheme: scheme.Scheme})

	return err
}
