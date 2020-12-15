package tests_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tests Suite")
}

var _ = BeforeSuite(func(done Done) {
	helmChartPath, err := filepath.Abs("../../charts/k8ssandra-cluster")
	Expect(err).To(BeNil())
	Expect(helmChartPath).NotTo(BeNil())
	close(done)
})
