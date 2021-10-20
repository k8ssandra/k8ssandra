package unit_test

import (
	"encoding/json"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/helm"
	cassdcv1beta1 "github.com/k8ssandra/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	helmUtils "github.com/k8ssandra/k8ssandra/tests/unit/utils/helm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Verify 4x features are created in template", func() {
	var (
		helmChartPath string
		err           error
		cassDC        cassdcv1beta1.CassandraDatacenter
	)

	BeforeEach(func() {
		helmChartPath, err = filepath.Abs(ChartsPath)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		err = nil
	})

	renderTemplate := func(options *helm.Options) error {
		return helmUtils.RenderAndUnmarshall("templates/cassandra/cassdc.yaml",
			options, helmChartPath, HelmReleaseName,
			func(renderedYaml string) error {
				return helm.UnmarshalK8SYamlE(GinkgoT(), renderedYaml, &cassDC)
			})
	}

	Context("by confirming a 3x template will not render", func() {
		It("if FQL options is defined", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version":                   "3.11.11",
					"cassandra.datacenters[0].fql.enable": "true",
					"cassandra.datacenters[0].name":       "testdc",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
		It("if audit logging options is defined", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version":                             "3.11.11",
					"cassandra.datacenters[0].audit_logging.enable": "true",
					"cassandra.datacenters[0].name":                 "testdc",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
		It("if client backpressure options is defined", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": "3.11.11",
					"cassandra.datacenters[0].client_backpressure.enable": "true",
					"cassandra.datacenters[0].name":                       "testdc",
				},
			}
			Expect(renderTemplate(options)).ShouldNot(Succeed())
		})
	})

	Context("by rendering audit log when", func() {
		It("cassandra version is 4x and FQL is enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version":                   "4.0.1",
					"cassandra.datacenters[0].fql.enable": "true",
					"cassandra.datacenters[0].name":       "testdc",
				},
			}
			Expect(renderTemplate(options)).To(Succeed())

			var dcCfg map[string]interface{}
			json.Unmarshal(cassDC.Spec.Config, &dcCfg)
			cassYaml, ok := dcCfg["cassandra-yaml"]
			if !ok {
				Fail("couldn't index cassandra-yaml in dc config")
			}
			fqlOpts, ok := cassYaml.(map[string]interface{})["full_query_logging_options"]
			if !ok {
				Fail("couldn't index fql options in dc config")
			}
			log_dir, ok := fqlOpts.(map[string]interface{})["log_dir"]
			if !ok {
				Fail("couldn't index log_dir in fql options in dc config")
			}
			if log_dir != "/var/log/cassandra/fql" {
				Fail("could not retrieve correct log_dir from FQL opts")
			}
		})
		It("cassandra version is 4x and audit logging is enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version":                             "4.0.1",
					"cassandra.datacenters[0].audit_logging.enable": "true",
					"cassandra.datacenters[0].name":                 "testdc",
				},
			}
			Expect(renderTemplate(options)).To(Succeed())

			var dcCfg map[string]interface{}
			json.Unmarshal(cassDC.Spec.Config, &dcCfg)
			cassYaml, ok := dcCfg["cassandra-yaml"]
			if !ok {
				Fail("couldn't index cassandra-yaml in dc config")
			}
			auditOpts, ok := cassYaml.(map[string]interface{})["audit_logging_options"]
			if !ok {
				Fail("couldn't index audit_logging_options in dc config")
			}
			auditEnabled, ok := auditOpts.(map[string]interface{})["enabled"]
			if !ok {
				Fail("couldn't find audit_logging_options.enabled in dc config")
			}
			if auditEnabled != true {
				Fail("audit logging was not enabled and should have been")
			}
		})
		It("cassandra version is 4x and client backpressure is enabled", func() {
			options := &helm.Options{
				KubectlOptions: defaultKubeCtlOptions,
				SetValues: map[string]string{
					"cassandra.version": "4.0.1",
					"cassandra.datacenters[0].client_backpressure.native_transport_max_concurrent_requests_in_bytes_per_ip": "1",
					"cassandra.datacenters[0].client_backpressure.native_transport_max_concurrent_requests_in_bytes":        "2",
					"cassandra.datacenters[0].name": "testdc",
				},
			}

			Expect(renderTemplate(options)).To(Succeed())

			var dcCfg map[string]interface{}
			json.Unmarshal(cassDC.Spec.Config, &dcCfg)
			cassYaml, ok := dcCfg["cassandra-yaml"]
			if !ok {
				Fail("couldn't index cassandra-yaml in dc config")
			}
			clPressPerIP, ok := cassYaml.(map[string]interface{})["native_transport_max_concurrent_requests_in_bytes_per_ip"]
			if !ok {
				Fail("couldn't index native_transport_max_concurrent_requests_in_bytes_per_ip in cassandra yaml")
			}
			Expect(clPressPerIP).To(Equal(1.0)) // When read back in, the ints in the yaml seem to get interpreted as float64s.
			clPress, ok := cassYaml.(map[string]interface{})["native_transport_max_concurrent_requests_in_bytes"]
			if !ok {
				Fail("couldn't index native_transport_max_concurrent_requests_in_bytes in dc config")
			}
			Expect(clPress).To(Equal(2.0)) // When read back in, the ints in the yaml seem to get interpreted as float64s.
		})
	})
})
