package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	namespace   = "k8ssandra"
	releaseName string
)

var _ = Describe("Install the cluster", func() {
	// TODO This assumes the kind was started with correctly setup cluster. Could we automate it here to provider proper e2e?

	BeforeSuite(func() {
		kubectlOptions := k8s.NewKubectlOptions("", "", "default")
		options := &helm.Options{KubectlOptions: kubectlOptions}

		// Add traefik repo and update repos
		helm.RunHelmCommandAndGetOutputE(GinkgoT(), options, "repo", "add", "traefik", "https://helm.traefik.io/traefik")
		helm.RunHelmCommandAndGetOutputE(GinkgoT(), options, "repo", "update")

		// Deploy traefik
		// helm install traefik traefik/traefik -n traefik --create-namespace -f docs/content/en/docs/topics/ingress/traefik/kind-deployment/traefik.values.yaml
		valuesPath, _ := filepath.Abs("../../docs/content/en/docs/topics/ingress/traefik/kind-deployment/traefik.values.yaml")
		helm.RunHelmCommandAndGetOutputE(GinkgoT(), options, "install", "traefik", "traefik/traefik", "-n", "traefik", "--create-namespace", "-f", valuesPath)
	})

	BeforeEach(func() {
		kubectlOptions := k8s.NewKubectlOptions("", "", "default")
		k8s.CreateNamespace(GinkgoT(), kubectlOptions, "k8ssandra")
	})

	AfterEach(func() {
		kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
		options := &helm.Options{KubectlOptions: kubectlOptions}

		// Verify everything was actually deleted.. I don't think Service is
		helm.Delete(GinkgoT(), options, releaseName, true)

		// Delete all the resources from namespace
		kubectlOptions = k8s.NewKubectlOptions("", "", "default")
		k8s.DeleteNamespace(GinkgoT(), kubectlOptions, namespace)
	})

	Context("deploying", func() {
		It("with default options", func() {
			clusterChartPath, err := filepath.Abs("../../charts/k8ssandra")
			Expect(err).To(BeNil())

			//kubectlOptions := k8s.NewKubectlOptions("", "", "default")
			//k8s.CreateNamespace(GinkgoT(), kubectlOptions, namespace)
			kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
			options := &helm.Options{
				// Enable traefik to allow redirections for testing
				SetValues: map[string]string{
					"ingress.traefik.enabled":                    "true",
					"ingress.traefik.monitoring.grafana.host":    "grafana.localhost",
					"ingress.traefik.monitoring.prometheus.host": "prometheus.localhost",
					"reaper.ingress.enabled":                     "true",
					"reaper.ingress.host":                        "repair.localhost",
				},
				KubectlOptions: k8s.NewKubectlOptions("", "", namespace),
			}

			releaseName = fmt.Sprintf(
				"demo-%s", strings.ToLower(random.UniqueId()))
			helm.Install(GinkgoT(), options, clusterChartPath, releaseName)

			// k8s module has no select by label
			// We could also use kubectl wait --for=condition=Ready pod -l name=cass-operator
			clientset, err := k8s.GetKubernetesClientFromOptionsE(GinkgoT(), kubectlOptions)
			pods, _ := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: "name=cass-operator"})

			Expect(len(pods.Items)).To(Equal(1))
			k8s.WaitUntilPodAvailable(GinkgoT(), kubectlOptions, pods.Items[0].Name, 50, 500*time.Millisecond)

			// Wait for CassandraDatacenter to be ready..
			k8s.RunKubectl(GinkgoT(), kubectlOptions, "wait", "--for=condition=Ready", "cassandradatacenter/dc1", "--timeout=300s")

			// Verify all the correct services are there
			// Grafana
			service := k8s.GetService(GinkgoT(), kubectlOptions, "grafana-service")
			Expect(service).ToNot(BeNil())

			// Prometheus
			service = k8s.GetService(GinkgoT(), kubectlOptions, "prometheus-operated")
			Expect(service).ToNot(BeNil())

			// Cassandra
			service = k8s.GetService(GinkgoT(), kubectlOptions, "k8ssandra-dc1-service")
			Expect(service).ToNot(BeNil())
			service = k8s.GetService(GinkgoT(), kubectlOptions, "k8ssandra-dc1-all-pods-service")
			Expect(service).ToNot(BeNil())
			service = k8s.GetService(GinkgoT(), kubectlOptions, "k8ssandra-seed-service")
			Expect(service).ToNot(BeNil())

			// Reaper
			// Medusa (if enabled)

			// Verify traefik is ready
			kubectlOptions = k8s.NewKubectlOptions("", "", "traefik")
			k8s.RunKubectl(GinkgoT(), kubectlOptions, "wait", "--for=condition=Ready", "pod", "-l", "app.kubernetes.io/name=traefik", "--timeout=300s")
			kubectlOptions = k8s.NewKubectlOptions("", "", namespace)

			// TODO Verify with cqlsh that Cassandra is working properly

			// Verify that prometheus is polling the Cassandra instance correctly
			// Wait for the Prometheus to be ready
			k8s.RunKubectl(GinkgoT(), kubectlOptions, "wait", "--for=condition=Ready", "pod", "-l", "app=prometheus", "--timeout=300s")

			// Poll Prometheus targets and check that it contains our cluster
			res, err := http.Get("http://prometheus.localhost:8080/api/v1/targets")
			Expect(err).To(BeNil())
			reply, err := ioutil.ReadAll(res.Body)
			Expect(err).To(BeNil())

			promReply := make(map[string]interface{})
			err = json.Unmarshal(reply, &promReply)
			Expect(err).To(BeNil())

			data := promReply["data"].(map[string]interface{})
			activeTargets := data["activeTargets"].([]interface{})
			first := activeTargets[0].(map[string]interface{})
			labels := first["labels"].(map[string]interface{})
			health := first["health"].(string)
			cluster := labels["cassandra_datastax_com_cluster"].(string)
			Expect(cluster).To(Equal("k8ssandra"))
			Expect(health).To(Equal("up"))

			// TODO More advanced, ensure repair is scheduled, take backup.. restore backup..
		})
	})
})
