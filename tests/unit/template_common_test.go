package unit_test

import (
	"fmt"
	"path/filepath"

	cassdcv1beta1 "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify common templates", func() {
	var (
		helmChartPath string
		cassdc        *cassdcv1beta1.CassandraDatacenter
	)

	BeforeEach(func() {
		path, err := filepath.Abs(chartsPath)
		Expect(err).To(BeNil())
		helmChartPath = path
		cassdc = &cassdcv1beta1.CassandraDatacenter{}
	})

	renderTemplate := func(options *helm.Options) error {
		renderedOutput, err := helm.RenderTemplateE(
			GinkgoT(), options, helmChartPath, helmReleaseName,
			[]string{"templates/cassandra/cassdc.yaml"},
		)

		if err == nil {
			err = helm.UnmarshalK8SYamlE(GinkgoT(), renderedOutput, cassdc)
		}

		return err
	}

	Context("by rendering it with options", func() {
		It("using only default options", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
			}

			Expect(renderTemplate(options)).To(Succeed())