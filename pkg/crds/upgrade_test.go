package crds

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	corev1 "k8s.io/api/core/v1"
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

	var cassdcCRD *unstructured.Unstructured

	By("creating new upgrader")
	u, err := NewWithClient(k8sClient)
	g.Expect(err).Should(Succeed())

	By("Upgrading / installing 1.0.0")
	crds, err := u.Upgrade("1.0.0")
	g.Expect(err).Should(Succeed())

	testOptions := envtest.CRDInstallOptions{
		PollInterval: 100 * time.Millisecond,
		MaxTime:      10 * time.Second,
	}

	var objs []runtime.Object
	for _, crd := range crds {
		if crd.GetName() == "cassandradatacenters.cassandra.datastax.com" {
			cassdcCRD = crd.DeepCopy()
		}
		objs = append(objs, &crd)
	}

	envtest.WaitForCRDs(cfg, objs, testOptions)
	err = k8sClient.Get(context.TODO(), client.ObjectKey{Name: cassdcCRD.GetName()}, cassdcCRD)
	ver := cassdcCRD.GetResourceVersion()
	g.Expect(err).Should(Succeed())

	_, found, err := unstructured.NestedFieldNoCopy(cassdcCRD.Object, "spec", "validation", "openAPIV3Schema", "properties", "spec", "properties", "configSecret")
	g.Expect(err).Should(Succeed())
	g.Expect(found).To(BeFalse())

	By("Upgrading to 1.1.0")
	crds, err = u.Upgrade("1.1.0")
	g.Expect(err).Should(Succeed())

	objs = []runtime.Object{}
	for _, crd := range crds {
		objs = append(objs, &crd)
	}

	envtest.WaitForCRDs(cfg, objs, testOptions)
	err = k8sClient.Get(context.TODO(), client.ObjectKey{Name: cassdcCRD.GetName()}, cassdcCRD)
	g.Expect(err).Should(Succeed())
	g.Expect(cassdcCRD.GetResourceVersion()).ToNot(Equal(ver))
	ver = cassdcCRD.GetResourceVersion()

	By("Upgrading to 1.2.0-20210514022645-da7547a5")
	crds, err = u.Upgrade("1.2.0-20210514022645-da7547a5")
	g.Expect(err).Should(Succeed())

	objs = []runtime.Object{}
	for _, crd := range crds {
		objs = append(objs, &crd)
	}

	envtest.WaitForCRDs(cfg, objs, testOptions)
	err = k8sClient.Get(context.TODO(), client.ObjectKey{Name: cassdcCRD.GetName()}, cassdcCRD)
	g.Expect(err).Should(Succeed())
	g.Expect(cassdcCRD.GetResourceVersion()).ToNot(Equal(ver))
	ver = cassdcCRD.GetResourceVersion()

	// 1.2.0 shipped with v1beta1 CRD version
	g.Expect("v1beta1", cassdcCRD.GetAPIVersion())

	_, found, err = unstructured.NestedFieldNoCopy(cassdcCRD.Object, "spec", "validation", "openAPIV3Schema", "properties", "spec", "properties", "configSecret")
	g.Expect(err).Should(Succeed())
	g.Expect(found).To(BeTrue())

	// TODO We want to update to the newest one here
	By("Upgrading to current head")
	crds, err = u.Upgrade("1.4.0-20211018063502-085ddce7")
	g.Expect(err).Should(Succeed())

	objs = []runtime.Object{}
	for _, crd := range crds {
		objs = append(objs, &crd)
	}

	envtest.WaitForCRDs(cfg, objs, testOptions)
	err = k8sClient.Get(context.TODO(), client.ObjectKey{Name: cassdcCRD.GetName()}, cassdcCRD)
	g.Expect(err).Should(Succeed())
	g.Expect(cassdcCRD.GetResourceVersion()).ToNot(Equal(ver))

	// Newest version should have v1 instead of v1beta1
	g.Expect("v1", cassdcCRD.GetAPIVersion())

	By("tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
	err = testEnv.Stop()
	g.Expect(err).ToNot(HaveOccurred())
}
