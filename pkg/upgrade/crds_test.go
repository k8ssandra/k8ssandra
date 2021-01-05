package upgrade

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	UpgradeTestNamespace = "k8ssandra-upgrade"
)

var _ = Describe("Upgrading CRDs", func() {
	Specify("verify by updating from 0.12.0 to 0.15.0", func() {
		testNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: UpgradeTestNamespace,
			},
		}
		Expect(k8sClient.Create(context.Background(), testNamespace)).Should(Succeed())

		By("creating new upgrader")
		u, err := NewWithClient(k8sClient)
		Expect(err).Should(Succeed())

		By("Upgrading / installing 0.12.0")
		Expect(u.Upgrade("0.12.0")).Should(Succeed())

		By("Upgrading to 0.15.0")
		Expect(u.Upgrade("0.15.0")).Should(Succeed())

		By("verifying that we have updated CRDs installed")
		// gvk, err := k8sClient.RESTMapper().ResourceFor()
		// rMapping, err := k8sClient.RESTMapper().RESTMapping()
		// Expect(err).Should(Succeed())

	})
})
