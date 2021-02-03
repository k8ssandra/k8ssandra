package cleaner

import (
	"context"
	"time"

	cassdcapi "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	CleanerTestNamespace = "k8ssandra"
	cleanerTestRelease   = "cleanrel"
	managedName          = "dc-managed"
	notManagedName       = "dc-not-managed"
)

var _ = Describe("Cleaning CassandraDatacenters", func() {
	Specify("but only the correct ones", func() {
		testNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: CleanerTestNamespace,
			},
		}
		Expect(k8sClient.Create(context.Background(), testNamespace)).Should(Succeed())

		cassdcKey := types.NamespacedName{
			Name:      notManagedName,
			Namespace: testNamespace.Name,
		}

		By("creating CassandraDatacenter outside helm charts")
		testDc := &cassdcapi.CassandraDatacenter{
			ObjectMeta: metav1.ObjectMeta{
				Name:        cassdcKey.Name,
				Namespace:   cassdcKey.Namespace,
				Annotations: map[string]string{},
			},
			Spec: cassdcapi.CassandraDatacenterSpec{
				ClusterName:   notManagedName,
				ServerType:    "cassandra",
				ServerVersion: "3.11.7",
				Size:          1,
			},
		}
		Expect(k8sClient.Create(context.Background(), testDc)).Should(Succeed())
		Eventually(func() bool {
			result := &cassdcapi.CassandraDatacenter{}
			_ = k8sClient.Get(context.Background(), cassdcKey, result)
			return result.Spec.Size == testDc.Spec.Size
		}, timeout, interval).Should(BeTrue())

		By("creating CassandraDatacenter managed by helm charts")
		cassdcKeyManaged := types.NamespacedName{
			Name:      managedName,
			Namespace: testNamespace.Name,
		}

		testDcManaged := &cassdcapi.CassandraDatacenter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cassdcKeyManaged.Name,
				Namespace: cassdcKeyManaged.Namespace,
				Annotations: map[string]string{
					releaseAnnotation: cleanerTestRelease,
				},
				Labels: map[string]string{
					managedLabel: managedLabelValue,
				},
			},
			Spec: cassdcapi.CassandraDatacenterSpec{
				ClusterName:   managedName,
				ServerType:    "cassandra",
				ServerVersion: "3.11.9",
				Size:          3,
			},
		}

		Expect(k8sClient.Create(context.Background(), testDcManaged)).Should(Succeed())
		Eventually(func() bool {
			result := &cassdcapi.CassandraDatacenter{}
			_ = k8sClient.Get(context.Background(), cassdcKeyManaged, result)
			return result.Spec.Size == testDcManaged.Spec.Size
		}, timeout, interval).Should(BeTrue())

		Eventually(func() bool {
			result := &cassdcapi.CassandraDatacenterList{}
			_ = k8sClient.List(context.Background(), result)
			return len(result.Items) == 2
		}, timeout, interval).Should(BeTrue())

		By("running the cleaner")
		cleaner := &Agent{
			Client:    k8sClient,
			Namespace: CleanerTestNamespace,
		}

		err := cleaner.RemoveResources(cleanerTestRelease)
		Expect(err).To(BeNil())

		By("verifying that only the managed CassandraDatacenter was deleted")
		Eventually(func() bool {
			result := &cassdcapi.CassandraDatacenterList{}
			_ = k8sClient.List(context.Background(), result, client.InNamespace(CleanerTestNamespace))
			return len(result.Items) == 1
		}, timeout, interval).Should(BeTrue())

		Consistently(func() bool {
			result := &cassdcapi.CassandraDatacenter{}
			_ = k8sClient.Get(context.Background(), cassdcKey, result)
			return result.Spec.ClusterName == notManagedName
		}, 1*time.Second, interval).Should(BeTrue())

		By("checking runs without removable CassandraDatacenters does not cause an error")
		err = cleaner.RemoveResources(cleanerTestRelease)
		Expect(err).To(BeNil())
	})

	Specify("even in empty namespaces", func() {
		cleaner := &Agent{
			Client:    k8sClient,
			Namespace: CleanerTestNamespace + "notReal",
		}

		By("running removeResources multiple times")
		Consistently(func() error {
			return cleaner.RemoveResources(cleanerTestRelease + "notReal")
		}, 1*time.Second, interval).Should(Succeed())

		result := &cassdcapi.CassandraDatacenterList{}
		Expect(k8sClient.List(context.Background(), result, client.InNamespace(CleanerTestNamespace+"notReal"))).Should(Succeed())
	})
})
