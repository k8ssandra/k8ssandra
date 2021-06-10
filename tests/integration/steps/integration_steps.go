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
	cassdcapi "github.com/k8ssandra/cass-operator/operator/pkg/apis/cassandra/v1beta1"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func installK8ssandraHelmRepo(t *testing.T) {
	helm.RunHelmCommandAndGetOutputE(t, &helm.Options{}, "repo", "add", "k8ssandra", "https://helm.k8ssandra.io/stable")
	helm.RunHelmCommandAndGetOutputE(t, &helm.Options{}, "repo", "update")
}

func deployCluster(t *testing.T, namespace, customValues string, helmValues map[string]string, upgrade bool, useLocalCharts bool, version string) {
	clusterChartPath, err := filepath.Abs("../../charts/k8ssandra")
	g(t).Expect(err).To(BeNil())

	if !useLocalCharts {
		installK8ssandraHelmRepo(t)
		clusterChartPath = "k8ssandra/k8ssandra"
	}

	customChartPath, err := filepath.Abs("charts/" + customValues)
	g(t).Expect(err).To(BeNil())

	if os.Getenv("K8SSANDRA_CASSANDRA_VERSION") != "" {
		if !useLocalCharts && strings.HasPrefix(os.Getenv("K8SSANDRA_CASSANDRA_VERSION"), "3.11") {
			// We should not set the 3.11 version when using the stable repo as we may end up with a patch version that's not available
			log.Println(Info("Using the default 3.11 Cassandra version available for this K8ssandra release"))
		} else {
			log.Println(Info(fmt.Sprintf("Using Cassandra version %s", os.Getenv("K8SSANDRA_CASSANDRA_VERSION"))))
			helmValues["cassandra.version"] = os.Getenv("K8SSANDRA_CASSANDRA_VERSION")
		}
	}

	helmOptions := &helm.Options{
		// Enable traefik to allow redirections for testing
		SetValues:      helmValues,
		KubectlOptions: k8s.NewKubectlOptions("", "", namespace),
		ValuesFiles:    []string{customChartPath},
	}

	if version != "" && version != "latest" {
		g(t).Expect(useLocalCharts).To(BeFalse(), "K8ssandra version can only be passed when using Helm repo based installs, not local charts.")
		helmOptions.Version = version
	}

	defer timeTrack(time.Now(), "Installing and starting k8ssandra")
	if upgrade {
		initResourceVersion := cassDcResourceVersion(t, namespace)
		err = helm.UpgradeE(t, helmOptions, clusterChartPath, releaseName)
		waitForCassDcUpgrade(t, namespace, initResourceVersion)
	} else {
		err = helm.InstallE(t, helmOptions, clusterChartPath, releaseName)
	}
	g(t).Expect(err).To(BeNil(), "Failed installing k8ssandra with Helm")
	// Wait for cass-operator pod to be ready
	labels := map[string]string{"app.kubernetes.io/name": "cass-operator"}
	g(t).Eventually(func() bool {
		return PodWithLabelsIsReady(t, namespace, labels)
	}, retryTimeout, retryInterval).Should(BeTrue())

	// Wait for CassandraDatacenter to be ready..
	WaitForCassDcToBeReady(t, namespace)
}

func DeployClusterWithValues(t *testing.T, namespace, options, customValues string, nodes int, upgrade bool, useLocalCharts bool, version string) {
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

	if options == "s3" && os.Getenv("K8SSANDRA_MEDUSA_BUCKET_NAME") != "" {
		helmValues = map[string]string{
			"medusa.bucketName":                os.Getenv("K8SSANDRA_MEDUSA_BUCKET_NAME"),
			"medusa.storage_properties.region": os.Getenv("K8SSANDRA_MEDUSA_BUCKET_REGION"),
		}
	}

	helmValues["cassandra.datacenters[0].size"] = strconv.Itoa(nodes)
	helmValues["cassandra.datacenters[0].name"] = datacenterName
	deployCluster(t, namespace, customValues, helmValues, upgrade, useLocalCharts, version)
}

func WaitForCassDcToBeUpdating(t *testing.T, namespace string) {
	waitForCassDcToBe(t, namespace, cassdcapi.ProgressUpdating)
}

// RestartStargate scales the Stargate deployment down to zero and then scales it
// back up to the prior number of replicas. This function blocks until the
// deployment is ready.
func RestartStargate(t *testing.T, releaseName, dcName, namespace string) {
	key := types.NamespacedName{Namespace: namespace, Name: releaseName + "-" + dcName + "-stargate"}
	retryInterval := 5 * time.Second
	scaleDownTimeout := 2 * time.Minute
	scaleUpTimeout := 5 * time.Minute

	g(t).Expect(WaitForDeploymentReady(t, key, retryInterval, scaleDownTimeout)).To(Succeed(), "failed waiting for Stargate to scale down")

	deployment := &appsv1.Deployment{}
	err := testClient.Get(context.Background(), key, deployment)
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed to get Stargate deployment: %s", err))

	originalCount := *deployment.Spec.Replicas

	patch := client.MergeFromWithOptions(deployment.DeepCopy(), client.MergeFromWithOptimisticLock{})
	count := int32(0)
	deployment.Spec.Replicas = &count

	err = testClient.Patch(context.Background(), deployment, patch)
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed to scale down Stargate: %s", err))

	patch = client.MergeFromWithOptions(deployment.DeepCopy(), client.MergeFromWithOptimisticLock{})
	deployment.Spec.Replicas = &originalCount

	err = testClient.Patch(context.Background(), deployment, patch)
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed to scale up Stargate: %s", err))

	g(t).Expect(WaitForDeploymentReady(t, key, retryInterval, scaleUpTimeout)).To(Succeed(), "failed waiting for Stargate to scale up")
}

// WaitForDeploymentReady Polls the deployment status until the Deployment is
// ready. Readiness is defined as .Status.Replicas == .Status.ReadyReplicas.
func WaitForDeploymentReady(t *testing.T, key types.NamespacedName, retryInterval, timeout time.Duration) bool {
	return g(t).Eventually(func() bool {
		deployment := &appsv1.Deployment{}
		if err := testClient.Get(context.Background(), key, deployment); err != nil {
			t.Logf("failed to get deployment %s: %s", key, err)
			return false
		}
		return deployment.Status.Replicas == deployment.Status.ReadyReplicas
	}, timeout, retryInterval).Should(BeTrue())
}

func WaitForCassDcToBeReady(t *testing.T, namespace string) {
	waitForCassDcToBe(t, namespace, cassdcapi.ProgressReady)
}

func cassDcResourceVersion(t *testing.T, namespace string) string {
	cassdcKey := types.NamespacedName{
		Name:      datacenterName,
		Namespace: namespace,
	}

	log.Printf("Checking cassandradatacenter %s resource version in namespace %s...", cassdcKey.Name, cassdcKey.Namespace)
	cassdc := &cassdcapi.CassandraDatacenter{}
	err := testClient.Get(context.Background(), cassdcKey, cassdc)
	if err != nil {
		t.Errorf("Failed getting cassdc: %s", err.Error())
		t.FailNow()
	}
	return cassdc.ResourceVersion

}

func waitForCassDcUpgrade(t *testing.T, namespace, initialResourceVersion string) {
	cassdcKey := types.NamespacedName{
		Name:      datacenterName,
		Namespace: namespace,
	}

	g(t).Eventually(func() bool {
		log.Printf("Checking cassandradatacenter %s resource version in namespace %s...", cassdcKey.Name, cassdcKey.Namespace)
		cassdc := &cassdcapi.CassandraDatacenter{}
		err := testClient.Get(context.Background(), cassdcKey, cassdc)
		if err != nil {
			t.Logf("Failed getting cassdc: %s", err.Error())
			return false
		}
		return cassdc.ResourceVersion != initialResourceVersion
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func waitForCassDcToBe(t *testing.T, namespace string, progress cassdcapi.ProgressState) {
	cassdcKey := types.NamespacedName{
		Name:      datacenterName,
		Namespace: namespace,
	}

	g(t).Eventually(func() bool {
		log.Printf("Checking cassandradatacenter %s state in namespace %s...", cassdcKey.Name, cassdcKey.Namespace)
		cassdc := &cassdcapi.CassandraDatacenter{}
		err := testClient.Get(context.Background(), cassdcKey, cassdc)
		if err != nil {
			t.Logf("Failed getting cassdc: %s", err.Error())
			return false
		}
		return cassdc.Status.CassandraOperatorProgress == progress &&
			(cassdc.GetConditionStatus(cassdcapi.DatacenterReady) == v1.ConditionTrue || progress == cassdcapi.ProgressUpdating)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func resourceWithLabelIsPresent(t *testing.T, namespace, resourceType string, labels map[string]string) bool {
	switch resourceType {
	case "pod":
		return CountPodsWithLabels(t, namespace, labels) == 1
	case "service":
		services := getServicesWithLabels(t, namespace, labels)
		if len(services.Items) == 1 {
			return true
		}
	default:
		log.Printf("Unsupported resource type for presence check: %s", resourceType)
		t.FailNow()
	}
	return false
}

func getPodsWithLabels(t *testing.T, namespace string, labels map[string]string) *v1.PodList {
	pods := &v1.PodList{}
	err := testClient.List(context.Background(), pods, client.InNamespace(namespace), client.MatchingLabels(labels))
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("Failed listing pods with labels %s", labels))
	return pods
}

func getServicesWithLabels(t *testing.T, namespace string, labels map[string]string) *v1.ServiceList {
	services := &v1.ServiceList{}
	err := testClient.List(context.Background(), services, client.InNamespace(namespace), client.MatchingLabels(labels))
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("Failed listing services with labels %s", labels))
	return services
}

func CountPodsWithLabels(t *testing.T, namespace string, labels map[string]string) int {
	pods := getPodsWithLabels(t, namespace, labels)
	return len(pods.Items)
}

func PodWithLabelsIsReady(t *testing.T, namespace string, label map[string]string) bool {
	g(t).Eventually(func() bool {
		return resourceWithLabelIsPresent(t, namespace, "pod", label)
	}, retryTimeout, retryInterval).Should(BeTrue())

	pods := getPodsWithLabels(t, namespace, label)
	if len(pods.Items) == 1 {
		return strings.ToLower(string(pods.Items[0].Status.Phase)) == "running"
	}
	return false
}

func WaitForPodWithLabelsToBeReady(t *testing.T, namespace string, labels map[string]string) {
	g(t).Eventually(func() bool {
		return PodWithLabelsIsReady(t, namespace, labels)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func DeployMinioAndCreateBucket(t *testing.T, bucketName string) {
	helmOptions := &helm.Options{
		KubectlOptions: getKubectlOptions("default"),
	}

	_, err := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "repo", "add", "minio", "https://helm.min.io/")
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed to add minio helm repo: %s", err))

	UninstallHelmRealeaseAndNamespace(t, "minio", "minio")

	values := fmt.Sprintf("accessKey=minio_key,secretKey=minio_secret,defaultBucket.enabled=true,defaultBucket.name=%s", bucketName)
	_, err = helm.RunHelmCommandAndGetOutputE(t, helmOptions, "install", "--set", values, "minio", "minio/minio", "-n", "minio", "--create-namespace")
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed to install the minio helm chart: %s", err))
}

func UninstallHelmRealeaseAndNamespace(t *testing.T, helmReleaseName, namespace string) {
	helmOptions := &helm.Options{
		KubectlOptions: getKubectlOptions("default"),
	}
	out, err := helm.RunHelmCommandAndGetOutputE(t, helmOptions, "list", "-n", namespace)
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed listing %s installs: %s", helmReleaseName, err))
	if strings.Contains(out, helmReleaseName) {
		_, err = helm.RunHelmCommandAndGetOutputE(t, helmOptions, "uninstall", helmReleaseName, "-n", namespace)
		g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed uninstalling %s: %s", helmReleaseName, err))
		DeleteNamespace(t, namespace)
		CheckNamespaceIsAbsent(t, namespace)
	}
}

func MinioServiceName(t *testing.T) string {
	minioService, err := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions("minio"), "get", "services", "-l", "app=minio", "-o", "jsonpath={.items[0].metadata.name}")
	g(t).Expect(err).To(BeNil())

	log.Printf("Minio service: %s", minioService)
	return minioService
}

func CheckResourceWithLabelsIsPresent(t *testing.T, namespace, resourceType string, labels map[string]string) {
	g(t).Eventually(func() bool {
		return resourceWithLabelIsPresent(t, namespace, resourceType, labels)
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func checkResourcePresence(t *testing.T, namespace, resourceType, name string) {
	switch resourceType {
	case "service":
		svc := &v1.Service{}
		key := types.NamespacedName{Namespace: namespace, Name: name}

		err := testClient.Get(context.Background(), key, svc)
		g(t).Expect(err).To(BeNil(), "failed to get service %s: %s", key, err)
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
	g(t).Expect(namespaceIsAbsent(namespace)).To(BeFalse())
}

func CheckSecretIsPresent(t *testing.T, namespace, secret string) {
	_, err := k8s.GetSecretE(t, getKubectlOptions(namespace), secret)
	g(t).Expect(err).To(BeNil())
}

func CheckNamespaceIsAbsent(t *testing.T, namespace string) {
	g(t).Eventually(func() bool {
		absent, err := namespaceIsAbsent(namespace)
		if err == nil {
			if absent {
				return true
			}
		} else {
			t.Logf("failed to check if namespace %s is absent: %s", namespace, err)
		}

		return false
	}, retryTimeout, retryInterval).Should(BeTrue())
}

func namespaceIsAbsent(namespace string) (bool, error) {
	if _, err := GetNamespace(namespace); err == nil || !apierrors.IsNotFound(err) {
		return false, err
	} else {
		return true, nil
	}
}

func namespaceIsTerminating(namespace string) (bool, error) {
	if ns, err := GetNamespace(namespace); err == nil {
		return ns.Status.Phase == v1.NamespaceTerminating, nil
	} else {
		return false, err
	}
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
	absent, err := namespaceIsAbsent(namespace)
	g(t).Expect(err).To(BeNil(), fmt.Sprintf("failed to check if namespace %s is absent: %s", namespace, err))

	if absent {
		log.Println(fmt.Sprintf("Creating namespace %s", namespace))
		k8s.CreateNamespace(t, getKubectlOptions("default"), namespace)
	} else {
		log.Println(Outline(fmt.Sprintf("Namespace %s already exists", namespace)))
	}
	return namespace
}

func getCassandraDatacenter(key types.NamespacedName) (*cassdcapi.CassandraDatacenter, error) {
	cassdc := &cassdcapi.CassandraDatacenter{}
	err := testClient.Get(context.Background(), key, cassdc)
	return cassdc, err
}

func WaitForCassandraDatacenterDeletion(t *testing.T, namespace string) {
	dcKey := types.NamespacedName{Namespace: namespace, Name: datacenterName}
	// Wait cassandradatacenter object to be actually deleted
	g(t).Eventually(func() bool {
		_, err := getCassandraDatacenter(dcKey)
		return apierrors.IsNotFound(err)
	}, retryTimeout, retryInterval).Should(BeTrue(), "cassandradatacenter object wasn't deleted within timeout")
}

func GetNamespace(name string) (*v1.Namespace, error) {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	if err := testClient.Get(context.Background(), types.NamespacedName{Namespace: name, Name: name}, ns); err == nil {
		return ns, nil
	} else {
		return nil, err
	}
}

func DeleteNamespace(t *testing.T, namespace string) {
	if ns, err := GetNamespace(namespace); err == nil {
		if err = testClient.Delete(context.Background(), ns); err != nil {
			t.Logf("failed to delete namespace %s: %s", namespace, err)
		}
	} else if apierrors.IsNotFound(err) {
		t.Logf("failed to delete namespace %s. namespace could not be retrieved: %s", namespace, err)
	}
}

func InstallTraefik(t *testing.T) {
	UninstallHelmRealeaseAndNamespace(t, "traefik", "traefik")
	kubectlOptions := getKubectlOptions("default")
	// Namespace doesn't exist yet, let's create it
	options := &helm.Options{KubectlOptions: kubectlOptions}

	// Add traefik repo and update repos
	helm.RunHelmCommandAndGetOutputE(t, options, "repo", "add", "traefik", "https://helm.traefik.io/traefik")
	helm.RunHelmCommandAndGetOutputE(t, options, "repo", "update")

	// Deploy traefik
	// helm install traefik traefik/traefik -n traefik --create-namespace -f docs/content/en/tasks/connect/ingress/kind-deployment/traefik.values.yaml
	valuesPath, _ := filepath.Abs("../../docs/content/en/tasks/connect/ingress/kind-deployment/traefik.values.yaml")
	_, err := helm.RunHelmCommandAndGetOutputE(t, options, "install", "traefik", "traefik/traefik", "-n", "traefik", "--create-namespace", "-f", valuesPath)
	g(t).Expect(err).To(BeNil())
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
	WaitForPodWithLabelsToBeReady(t, namespace, map[string]string{"app.kubernetes.io/managed-by": "reaper-operator"})
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
