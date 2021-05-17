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

func TestAFunctionality(t *testing.T) {
	var cfg *rest.Config
	var k8sClient client.Client
	var testEnv *envtest.Environment

	RegisterTestingT(t)
	g := NewWithT(t)

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		// CRDDirectoryPaths: []string{filepath.Join("..", "..", "charts", "cass-operator", "crds")},
	}

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
	crds, err := u.Upgrade("1.0.0")
	g.Expect(err).Should(Succeed())

	testOptions := envtest.CRDInstallOptions{
		PollInterval: 100 * time.Millisecond,
		MaxTime:      10 * time.Second,
	}

	var objs []runtime.Object
	for _, crd := range crds {
		objs = append(objs, &crd)
	}

	envtest.WaitForCRDs(cfg, objs, testOptions)

	By("Upgrading to 1.1.0")
	crds, err = u.Upgrade("1.1.0")
	g.Expect(err).Should(Succeed())

	objs = []runtime.Object{}
	for _, crd := range crds {
		objs = append(objs, &crd)
	}

	envtest.WaitForCRDs(cfg, objs, testOptions)

	By("Upgrading to 1.2.0-20210514022645-da7547a5")
	crds, err = u.Upgrade("1.2.0-20210514022645-da7547a5")
	g.Expect(err).Should(Succeed())

	objs = []runtime.Object{}
	for _, crd := range crds {
		objs = append(objs, &crd)
	}

	envtest.WaitForCRDs(cfg, objs, testOptions)

	// TODO Check that CassandraDatacenter has new property

	By("tearing down the test environment")
	gexec.KillAndWait(5 * time.Second)
	err = testEnv.Stop()
	g.Expect(err).ToNot(HaveOccurred())
}
