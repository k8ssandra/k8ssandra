package cleaner

import (
	"context"
	"log"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"

	cassdcapi "github.com/k8ssandra/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	managedLabel      = "app.kubernetes.io/managed-by"
	managedLabelValue = "Helm"
	instanceLabel     = "app.kubernetes.io/instance"
	nameLabel         = "app.kubernetes.io/name"
	nameLabelValue    = "k8ssandra"
	releaseAnnotation = "meta.helm.sh/release-name"
)

// Agent is a cleaner utility for resources which helm pre-delete requires
type Agent struct {
	Client    client.Client
	Namespace string
}

// New returns a new instance of cleaning agent
func New(namespace string) (*Agent, error) {
	_ = api.AddToScheme(scheme.Scheme)
	_ = cassdcapi.AddToScheme(scheme.Scheme)

	c, err := client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Fatal(err)
	}

	return &Agent{
		Client:    c,
		Namespace: namespace,
	}, nil
}

// RemoveResources deletes all the resources with finalizers or which we want an operator to trigger a deletion
func (a *Agent) RemoveResources(releaseName string) error {
	// Remove CassandraDatacenter (cass-operator should delete all the finalizers and associated resources)
	if err := a.removeCassandraDatacenter(releaseName); err != nil {
		log.Fatalf("Failed to remove Cassandra cluster(s): %v", err)
		return err
	}
	return nil
}

func (a *Agent) removeCassandraDatacenter(releaseName string) error {
	log.Printf("Removing CassandraDatacenter(s) managed in release %s from namespace %s\n", releaseName, a.Namespace)
	releaseLabels := client.MatchingLabels{
		managedLabel:  managedLabelValue,
		instanceLabel: releaseName,
		nameLabel:     nameLabelValue,
	}
	list := &cassdcapi.CassandraDatacenterList{}
	err := a.Client.List(context.Background(), list, client.InNamespace(a.Namespace), releaseLabels)
	if err != nil {
		log.Fatalf("Failed to list CassandraDatacenters in namespace: %s", a.Namespace)
		return err
	}

	for _, cassdc := range list.Items {
		if err = a.Client.Delete(context.Background(), &cassdc); err != nil && !apierrors.IsNotFound(err) && !apierrors.IsResourceExpired(err) {
			log.Fatalf("failed to delete CassandraDatacenter %s: %s",
				types.NamespacedName{Namespace: cassdc.Namespace, Name: cassdc.Name}, err)
		}
	}

	// We need to wait until the CassandraDatacenter is terminated; otherwise, cass-operator could get
	// deleted before it has a chance to clear the CassandraDatacenter's finalizer.
	return wait.PollImmediate(10*time.Second, 10*time.Minute, func() (bool, error) {
		list := &cassdcapi.CassandraDatacenterList{}
		err := a.Client.List(context.Background(), list, client.InNamespace(a.Namespace), releaseLabels)
		if err != nil {
			log.Printf("failed to list CassandraDatacenters: %s\n", err)
			return false, err
		}
		return len(list.Items) == 0, nil
	})
}
