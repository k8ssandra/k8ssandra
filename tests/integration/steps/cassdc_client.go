package steps

import (
	cassdcapi "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Client returns a controller-runtime client with cass-operator API defined
func CassDcClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	err := cassdcapi.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}
	c, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return c, nil
}
