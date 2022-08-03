package crds

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	corev1 "k8s.io/api/core/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	// +kubebuilder:scaffold:imports
)

const (
	UpgradeTestNamespace = "k8ssandra-upgrade-test"
)

func TestUpgradingCRDs(t *testing.T) {
	var cfg *rest.Config
	var k8sClient client.Client
	var testEnv *envtest.Environment

	RegisterTestingT(t)
	g := NewWithT(t)

	chartNames := []string{"k8ssandra", ""}
	for _, chartName := range chartNames {
		t.Run(fmt.Sprintf("CRD upgrade for chart name: %s", chartName), func(t *testing.T) {
			By("bootstrapping test environment")
			testEnv = &envtest.Environment{}

			var err error
			cfg, err = testEnv.Start()
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(cfg).ToNot(BeNil())

			// err = cassdcapi.AddToScheme(scheme.Scheme)
			g.Expect(err).NotTo(HaveOccurred())

			// +kubebuilder:scaffold:scheme

			k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(k8sClient).ToNot(BeNil())

			testNamespace := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: UpgradeTestNamespace,
				},
			}
			g.Expect(k8sClient.Create(context.Background(), testNamespace)).Should(Succeed())

			By("creating new upgrader")
			u, err := NewWithClient(k8sClient)
			g.Expect(err).Should(Succeed())

			By("Upgrading / installing 1.0.0")
			var crds []unstructured.Unstructured
			g.Eventually(func() bool {
				crds, err = u.Upgrade(chartName, "1.0.0")
				return err == nil
			}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 5).Should(BeTrue())

			g.Expect(err).Should(Succeed())

			testOptions := envtest.CRDInstallOptions{
				PollInterval: 100 * time.Millisecond,
				MaxTime:      10 * time.Second,
			}

			unstructuredCRD := &unstructured.Unstructured{}
			cassDCCRD := &apiextensions.CustomResourceDefinition{}
			objs := []apiextensions.CustomResourceDefinition{}
			for _, crd := range crds {
				if crd.GetName() == "cassandradatacenters.cassandra.datastax.com" {
					unstructuredCRD = crd.DeepCopy()
					err = runtime.DefaultUnstructuredConverter.FromUnstructured(crd.UnstructuredContent(), cassDCCRD)
					g.Expect(err).ShouldNot(HaveOccurred())
				}
				objs = append(objs, *cassDCCRD)
			}

			envtest.WaitForCRDs(cfg, objs, testOptions)
			err = k8sClient.Get(context.TODO(), client.ObjectKey{Name: cassDCCRD.GetName()}, cassDCCRD)
			ver := cassDCCRD.GetResourceVersion()
			g.Expect(err).Should(Succeed())

			_, found, err := unstructured.NestedFieldNoCopy(unstructuredCRD.Object, "spec", "validation", "openAPIV3Schema", "properties", "spec", "properties", "configSecret")
			g.Expect(err).Should(Succeed())
			g.Expect(found).To(BeFalse())

			By("Upgrading to 1.5.1")
			crds, err = u.Upgrade(chartName, "1.5.1")
			g.Expect(err).Should(Succeed())

			objs = []apiextensions.CustomResourceDefinition{}
			for _, crd := range crds {
				if crd.GetName() == "cassandradatacenters.cassandra.datastax.com" {
					unstructuredCRD = crd.DeepCopy()
					err = runtime.DefaultUnstructuredConverter.FromUnstructured(crd.UnstructuredContent(), cassDCCRD)
					g.Expect(err).ShouldNot(HaveOccurred())
					objs = append(objs, *cassDCCRD)
				}
			}

			envtest.WaitForCRDs(cfg, objs, testOptions)
			err = k8sClient.Get(context.TODO(), client.ObjectKey{Name: cassDCCRD.GetName()}, cassDCCRD)
			g.Expect(err).Should(Succeed())
			g.Eventually(func() bool {
				newver := cassDCCRD.GetResourceVersion()
				eq := newver == ver
				println(fmt.Sprintf("equality: %t, current resourceVersion: %s, old resourceVersion: %s", eq, newver, ver))
				return eq
			}).WithTimeout(time.Minute * 10).WithPolling(time.Second * 5).Should(BeFalse())

			versionsSlice, found, err := unstructured.NestedSlice(unstructuredCRD.Object, "spec", "versions")
			g.Expect(found).To(BeTrue())
			_, found, err = unstructured.NestedFieldNoCopy(versionsSlice[0].(map[string]interface{}), "schema", "openAPIV3Schema", "properties", "spec", "properties", "configSecret")

			g.Expect(err).Should(Succeed())
			g.Expect(found).To(BeTrue())

			By("tearing down the test environment")
			gexec.KillAndWait(5 * time.Second)
			err = testEnv.Stop()
			g.Expect(err).ToNot(HaveOccurred())
		})
	}
}
