package steps

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cassdcapi "github.com/datastax/cass-operator/operator/pkg/apis/cassandra/v1beta1"
)

const (
	retryTimeout   = 15 * time.Minute
	retryInterval  = 30 * time.Second
	releaseName    = "k8ssandra"
	datacenterName = "dc1"
)

var (
	Info    = Color("\033[1;33m%s\033[0m")
	Outline = Color("\033[1;34m%s\033[0m")
	Step    = Color("\033[1;36m%s\033[0m")
	Success = Color("\033[1;32m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Println(Info(fmt.Sprintf("%s took %s", name, elapsed)))
}

func RunShellCommand(command *exec.Cmd) error {
	err := command.Run()
	return err
}

func RunShellCommandAndGetOutput(command *exec.Cmd) string {
	var outb bytes.Buffer
	command.Stdout = &outb
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}

	return string(outb.String())
}

// Find returns true if val exists in the slice array, false otherwise
func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func g(t *testing.T) *WithT {
	return NewWithT(t)
}

func getKubectlOptions(namespace string) *k8s.KubectlOptions {
	return k8s.NewKubectlOptions("", "", namespace)
}

func deployCluster(t *testing.T, namespace, customValues string, helmValues map[string]string, upgrade bool) {
	clusterChartPath, err := filepath.Abs("../../charts/k8ssandra")
	g(t).Expect(err).To(BeNil())

	customChartPath, err := filepath.Abs("charts/" + customValues)
	g(t).Expect(err).To(BeNil())

	if os.Getenv("K8SSANDRA_CASSANDRA_VERSION") != "" {
		log.Println(Info(fmt.Sprintf("Using Cassandra version %s", os.Getenv("K8SSANDRA_CASSANDRA_VERSION"))))
		helmValues["cassandra.version"] = os.Getenv("K8SSANDRA_CASSANDRA_VERSION")
	}

	helmOptions := &helm.Options{
		// Enable traefik to allow redirections for testing
		SetValues:      helmValues,
		KubectlOptions: k8s.NewKubectlOptions("", "", namespace),
		ValuesFiles:    []string{customChartPath},
	}

	defer timeTrack(time.Now(), "Installing and starting k8ssandra")
	if upgrade {
		err = helm.UpgradeE(t, helmOptions, clusterChartPath, releaseName)
	} else {
		err = helm.InstallE(t, helmOptions, clusterChartPath, releaseName)
	}
	g(t).Expect(err).To(BeNil(), "Failed installing k8ssandra with Helm")
	// Wait for cass-operator pod to be ready
	g(t).Eventually(func() bool {
		return PodWithLabelIsReady(t, namespace, "app.kubernetes.io/name=cass-operator")
	}, retryTimeout, retryInterval).Should(BeTrue())

	// Wait for CassandraDatacenter to be udpating..
	WaitForCassDcToBeUpdating(t, namespace)

	// Wait for CassandraDatacenter to be ready..
	WaitForCassDcToBeReady(t, namespace)
}

func DeployClusterWithValues(t *testing.T, namespace, options, customValues string, nodes int, upgrade bool) {
	log.Printf("Deploying a cluster with %s options using the %s values", options, customValues)

	helmValues := map[string]string{}
	if options == "default" {
		helmValues = map[string]string{
			"reaper.ingress.host": "repair.127.0.0.1.nip.io",
		}
	}
	if options == "minio" {
		serviceName := MinioServiceName(t)
		helmValues = map[string]string{
			"medusa.storage_properties.host": fmt.Sprintf("%s.minio.svc.cluster.local", serviceName),
		}
	}
	helmValues["cassandra.datacenters[0].size"] = strconv.Itoa(nodes)
	helmValues["cassandra.datacenters[0].name"] = datacenterName
	deployCluster(t, namespace, customValues, helmValues, upgrade)
}

func WaitForCassDcToBeUpdating(t *testing.T, namespace string) {
	waitForCassDcToBe(t, namespace, cassdcapi.ProgressUpdating)
}

func WaitForCassDcToBeReady(t *testing.T, namespace string) {
	waitForCassDcToBe(t, namespace, cassdcapi.ProgressReady)
}

func waitForCassDcToBe(t *testing.T, namespace string, progress cassdcapi.ProgressState) {
	cassdcKey := types.NamespacedName{
		Name:      datacenterName,
		Namespace: namespace,
	}

	k8sClient, err := CassDcClient()
	g(t).Expect(err).To(BeNil(), "Couldn't instantiate controller-runtime client with cassdc API")

	g(t).Eventually(func() bool {
		log.Printf("Checking cassandradatacenter %s state in namespace %s...", cassdcKey.Name, cassdcKey.Namespace)
		cassdc := &cassdcapi.CassandraDatacenter{}
		err := k8sClient.Get(context.Background(), cassdcKey, cassdc)
		if err != nil {
			t.Logf("Failed getting cassdc: %s", err.Error())
			return false
		}
		return cassdc.Status.CassandraOperatorProgress == progress &&
			(cassdc.GetConditionStatus(cassdcapi.DatacenterReady) == v1.ConditionTrue || progress == cassdcapi.ProgressUpdating)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func resourceWithLabelIsPresent(t *testing.T, namespace, resourceType, label string) bool {
	switch resourceType {
	case "pod":
		return CountPodsWithLabel(t, namespace, label) == 1
	case "service":
		services := getServicesWithLabel(t, namespace, label)
		if len(services.Items) == 1 {
			return true
		}
	default:
		log.Printf("Unsupported resource type for presence check: %s", resourceType)
		t.FailNow()
	}
	return false
}

func getPodsWithLabel(t *testing.T, namespace, label string) *v1.PodList {
	clientset, _ := k8s.GetKubernetesClientFromOptionsE(t, getKubectlOptions(namespace))
	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("Failed listing pods with label %s", label))
	return pods
}

func getServicesWithLabel(t *testing.T, namespace, label string) *v1.ServiceList {
	clientset, _ := k8s.GetKubernetesClientFromOptionsE(t, getKubectlOptions(namespace))
	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("Failed listing services with label %s", label))
	return services
}

func CountPodsWithLabel(t *testing.T, namespace, label string) int {
	pods := getPodsWithLabel(t, namespace, label)
	return len(pods.Items)
}

func PodWithLabelIsReady(t *testing.T, namespace string, label string) bool {
	g(t).Eventually(func() bool {
		return resourceWithLabelIsPresent(t, namespace, "pod", label)
	}, retryTimeout, retryInterval).Should(BeTrue())

	pods := getPodsWithLabel(t, namespace, label)
	if len(pods.Items) == 1 {
		return strings.ToLower(string(pods.Items[0].Status.Phase)) == "running"
	}
	return false
}

func WaitForPodWithLabelToBeReady(t *testing.T, namespace, label string) {
	g(t).Eventually(func() bool {
		return PodWithLabelIsReady(t, namespace, label)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func DeployMinioAndCreateBucket(t *testing.T, bucketName string) {
	helmOptions := &helm.Options{
		KubectlOptions: getKubectlOptions("default"),
	}
	helm.RunHelmCommandAndGetOutputE(t, helmOptions, "repo", "add", "minio", "https://helm.min.io/")

	helm.RunHelmCommandAndGetOutputE(t, helmOptions, "install",
		"--set", fmt.Sprintf("accessKey=minio_key,secretKey=minio_secret,defaultBucket.enabled=true,defaultBucket.name=%s", bucketName),
		"minio", "minio/minio", "-n", "minio", "--create-namespace")
}

func MinioServiceName(t *testing.T) string {
	minioService, err := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions("minio"), "get", "services", "-l", "app=minio", "-o", "jsonpath={.items[0].metadata.name}")
	g(t).Expect(err).To(BeNil())

	log.Printf("Minio service: %s", minioService)
	return minioService
}

func CheckResourceWithLabelIsPresent(t *testing.T, namespace, resourceType, label string) {
	g(t).Eventually(func() bool {
		return resourceWithLabelIsPresent(t, namespace, resourceType, label)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func checkResourcePresence(t *testing.T, namespace, resourceType, name string) {
	switch resourceType {
	case "service":
		k8s.GetService(t, getKubectlOptions(namespace), name)
	default:
		t.Logf("Unsupported resource type: %s", resourceType)
		t.FailNow()
	}
}

func CheckClusterExpectedResources(t *testing.T, namespace string) {
	checkResourcePresence(t, namespace, "service", fmt.Sprintf("%s-%s-all-pods-service", releaseName, datacenterName))
	checkResourcePresence(t, namespace, "service", fmt.Sprintf("%s-%s-service", releaseName, datacenterName))
	checkResourcePresence(t, namespace, "service", fmt.Sprintf("%s-seed-service", releaseName))
}

func CheckK8sClusterIsReachable(t *testing.T) {
	output, err := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions("default"), "get", "namespace", "default", "-o", "jsonpath={.metadata.name}")
	g(t).Expect(err).To(BeNil())
	g(t).Expect(output).Should(Equal("default"))
}

func CheckNamespaceWasCreated(t *testing.T, namespace string) {
	g(t).Expect(namespaceIsAbsent(t, namespace)).To(BeFalse())
}

func CheckSecretIsPresent(t *testing.T, namespace, secret string) {
	_, err := k8s.GetSecretE(t, getKubectlOptions(namespace), secret)
	g(t).Expect(err).To(BeNil())
}

func CheckNamespaceIsAbsent(t *testing.T, namespace string) {
	g(t).Eventually(func() bool {
		return namespaceIsAbsent(t, namespace) || namespaceIsTerminating(t, namespace)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func namespaceIsAbsent(t *testing.T, namespace string) bool {
	namespaceObject, _ := k8s.GetNamespaceE(t, getKubectlOptions("default"), namespace)
	return namespaceObject.Name != namespace
}

func namespaceIsTerminating(t *testing.T, namespace string) bool {
	namespaceObject, _ := k8s.GetNamespaceE(t, getKubectlOptions("default"), namespace)

	return namespaceObject.Status.Phase == v1.NamespaceTerminating
}

func CreateNamespace(t *testing.T) string {
	namespace := func() string {
		if os.Getenv("K8SSANDRA_NS") == "" {
			// Generate a namespace name if the env variable doesn't exist
			return fmt.Sprintf("k8ssandra%s", time.Now().Format("2006010215040507"))
		} else {
			return os.Getenv("K8SSANDRA_NS")
		}
	}()
	if namespaceIsAbsent(t, namespace) {
		log.Println(fmt.Sprintf("Creating namespace %s", namespace))
		k8s.CreateNamespace(t, getKubectlOptions("default"), namespace)
	} else {
		log.Println(Outline(fmt.Sprintf("Namespace %s already exists", namespace)))
	}
	return namespace
}

func getCassDcClient(t *testing.T) client.Client {
	client, err := CassDcClient()
	g(t).Expect(err).To(BeNil(), "Couldn't instantiate controller-runtime client with cassdc API")
	return client
}

func getCassandraDatacenter(t *testing.T, key types.NamespacedName) (*cassdcapi.CassandraDatacenter, error) {
	cassdc := &cassdcapi.CassandraDatacenter{}
	err := getCassDcClient(t).Get(context.Background(), key, cassdc)
	return cassdc, err
}

func WaitForCassandraDatacenterDeletion(t *testing.T, namespace string) {
	dcKey := types.NamespacedName{Namespace: namespace, Name: datacenterName}
	// Wait cassandradatacenter object to be actually deleted
	g(t).Eventually(func() bool {
		_, err := getCassandraDatacenter(t, dcKey)
		return apierrors.IsNotFound(err)
	}, retryTimeout, retryInterval).Should(BeTrue(), "cassandradatacenter object wasn't deleted within timeout")
}

func UninstallTraefikHelmRelease(t *testing.T, traefikNamespace string) {
	err := RunShellCommand(exec.Command("helm", "uninstall", "traefik", "-n", traefikNamespace))
	if err != nil {
		t.Logf("Failed uninstalling Traefik Helm release: %s", err.Error())
	}
}

func UninstallMinioHelmRelease(t *testing.T, minioNamespace string) {
	if !namespaceIsAbsent(t, minioNamespace) {
		err := RunShellCommand(exec.Command("helm", "uninstall", "minio", "-n", minioNamespace))
		if err != nil {
			t.Logf("Failed uninstalling Minio Helm release: %s", err.Error())
		}
	}
}

func UninstallK8ssandraHelmRelease(t *testing.T, namespace string) {
	err := RunShellCommand(exec.Command("helm", "uninstall", releaseName, "-n", namespace))
	if err != nil {
		t.Logf("Failed uninstalling K8ssandra Helm release: %s", err.Error())
	}
}

func DeleteNamespace(t *testing.T, namespace string) {
	if !namespaceIsAbsent(t, namespace) {
		err := k8s.DeleteNamespaceE(t, getKubectlOptions("default"), namespace)
		if err != nil {
			t.Logf("Failed deleting namespace %s: %s", namespace, err.Error())
		}
	}
}

func InstallTraefik(t *testing.T) {
	kubectlOptions := getKubectlOptions("default")
	_, err := k8s.GetNamespaceE(t, kubectlOptions, "traefik")
	if err != nil {
		// Namespace doesn't exist yet, let's create it
		options := &helm.Options{KubectlOptions: kubectlOptions}

		// Add traefik repo and update repos
		helm.RunHelmCommandAndGetOutputE(t, options, "repo", "add", "traefik", "https://helm.traefik.io/traefik")
		helm.RunHelmCommandAndGetOutputE(t, options, "repo", "update")

		// Deploy traefik
		// helm install traefik traefik/traefik -n traefik --create-namespace -f docs/content/en/tasks/connect/ingress/kind-deployment/traefik.values.yaml
		valuesPath, _ := filepath.Abs("../../docs/content/en/tasks/connect/ingress/kind-deployment/traefik.values.yaml")
		_, err = helm.RunHelmCommandAndGetOutputE(t, options, "install", "traefik", "traefik/traefik", "-n", "traefik", "--create-namespace", "-f", valuesPath)
		g(t).Expect(err).To(BeNil())
	}
}

type credentials struct {
	username string
	password string
}

func ExtractUsernamePassword(t *testing.T, secretName, namespace string) credentials {
	secret := k8s.GetSecret(t, getKubectlOptions(namespace), secretName)
	username := secret.Data["username"]
	password := secret.Data["password"]
	creds := credentials{string(username), string(password)}
	return creds
}

func runCassandraQueryAndGetOutput(t *testing.T, namespace, query string) string {
	cqlCredentials := ExtractUsernamePassword(t, "k8ssandra-superuser", namespace)
	// Get reaper service
	output, _ := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions(namespace), "exec", "-it", fmt.Sprintf("%s-%s-default-sts-0", releaseName, datacenterName), "--", "/opt/cassandra/bin/cqlsh", "--username", cqlCredentials.username, "--password", cqlCredentials.password, "-e", query)
	return output
}

func CheckKeyspaceExists(t *testing.T, namespace, keyspace string) {
	keyspaces := runCassandraQueryAndGetOutput(t, namespace, "describe keyspaces")

	// Check that the keyspace exists in the list of keyspaces
	g(t).Expect(keyspaces).Should(ContainSubstring(keyspace))
}

func WaitForReaperPod(t *testing.T, namespace string) {
	WaitForPodWithLabelToBeReady(t, namespace, "app.kubernetes.io/managed-by=reaper-operator")
}

func CheckRowCountInTable(t *testing.T, nbRows int, namespace, tableName, keyspaceName string) {
	output := runCassandraQueryAndGetOutput(t, namespace, fmt.Sprintf("SELECT id FROM %s.%s", keyspaceName, tableName))

	// Check that we have the right number of rows
	g(t).Expect(output).Should(ContainSubstring(fmt.Sprintf("(%d rows)", nbRows)))
}

func CreateCassandraTable(t *testing.T, namespace, tableName, keyspaceName string) {
	runCassandraQueryAndGetOutput(t, namespace, fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s with replication = {'class':'SimpleStrategy', 'replication_factor':1};", keyspaceName))
	runCassandraQueryAndGetOutput(t, namespace, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s(id timeuuid PRIMARY KEY, val text);", keyspaceName, tableName))
	runCassandraQueryAndGetOutput(t, namespace, fmt.Sprintf("TRUNCATE %s.%s;", keyspaceName, tableName))
}

func LoadRowsInTable(t *testing.T, nbRows int, namespace, tableName, keyspaceName string) {
	for i := 0; i < nbRows; i++ {
		runCassandraQueryAndGetOutput(t, namespace, fmt.Sprintf("INSERT INTO %s.%s(id,val) values(now(), '%d');", keyspaceName, tableName, i))
	}
}
