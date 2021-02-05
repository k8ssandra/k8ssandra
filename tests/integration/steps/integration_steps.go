package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var namespace = "k8ssandra"
var repairId string

func RunShellCommand(command *exec.Cmd) error {
	err := command.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = command.Wait()
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

func GetServicesWithLabel(label string) (*v1.ServiceList, error) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	clientset, err := k8s.GetKubernetesClientFromOptionsE(GinkgoT(), kubectlOptions)
	Expect(err).To(BeNil())
	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: label})
	return services, err
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

func deployCluster(customValues string, helmValues map[string]string) {
	clusterChartPath, err := filepath.Abs("../../charts/k8ssandra")
	if err != nil {
		Fail("Couldn't find the absolute path for K8ssandra charts")
	}

	customChartPath, err := filepath.Abs("charts/" + customValues)
	if err != nil {
		Fail(fmt.Sprintf("Couldn't find the absolute path for custom values: %s", customValues))
	}

	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	helmOptions := &helm.Options{}

	helmOptions = &helm.Options{
		// Enable traefik to allow redirections for testing
		SetValues:      helmValues,
		KubectlOptions: k8s.NewKubectlOptions("", "", namespace),
		ValuesFiles:    []string{customChartPath},
	}

	releaseName := "k8ssandra"
	defer timeTrack(time.Now(), fmt.Sprintf("Installing and starting k8ssandra"))
	helm.Install(GinkgoT(), helmOptions, clusterChartPath, releaseName)

	// Wait for cass-operator pod to be ready
	attempts := 0
	maxAttempts := 10
	for {
		attempts++
		clientset, _ := k8s.GetKubernetesClientFromOptionsE(GinkgoT(), kubectlOptions)
		pods, _ := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=cass-operator"})
		if len(pods.Items) == 1 {
			k8s.RunKubectl(GinkgoT(), kubectlOptions, "wait", "--for=condition=Ready", "pod", "-l", "app.kubernetes.io/name=cass-operator", "--timeout=1800s")
			break
		} else if attempts > maxAttempts {
			Fail("Couldn't find cass-operator pod")
		}
		time.Sleep(20 * time.Second)
	}

	// Wait for CassandraDatacenter to be ready..
	k8s.RunKubectl(GinkgoT(), kubectlOptions, "wait", "--for=condition=Ready", "cassandradatacenter/dc1", "--timeout=1800s")
}

func GetUsernamePassword(secretName, ns string) credentials {
	kubectlOptions := k8s.NewKubectlOptions("", "", ns)
	secret := k8s.GetSecret(GinkgoT(), kubectlOptions, secretName)
	username := secret.Data["username"]
	password := secret.Data["password"]
	creds := credentials{string(username), string(password)}
	return creds
}

func runCassandraQueryAndGetOutput(query string) string {
	cqlCredentials := GetUsernamePassword("k8ssandra-superuser", namespace)
	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	// Get reaper service
	output, _ := k8s.RunKubectlAndGetOutputE(GinkgoT(), kubectlOptions, "exec", "-it", "k8ssandra-dc1-default-sts-0", "--", "/opt/cassandra/bin/cqlsh", "--username", cqlCredentials.username, "--password", cqlCredentials.password, "-e", query)
	return output
}

func WaitForPodWithLabelToBeReady(label string, waitTime time.Duration, maxAttempts int) {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	attempts := 0
	for {
		attempts++
		getCassandraPodOutput, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), kubectlOptions, "get", "pods", "-l", label)
		if err == nil && !strings.HasPrefix(getCassandraPodOutput, "No resources found") {
			break
		}
		if attempts > maxAttempts {
			Fail(fmt.Sprintf("Pod with label '%s' didn't start within timeout", label))
		}
		time.Sleep(waitTime)
	}
	k8s.RunKubectl(GinkgoT(), kubectlOptions, "wait", "--for=condition=Ready", "pod", "-l", label, "--timeout=1800s")
}

func GetStargateService() v1.Service {
	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	clientset, err := k8s.GetKubernetesClientFromOptionsE(GinkgoT(), kubectlOptions)
	if err != nil {
		Fail("Couldn't get k8s client")
	}

	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	for _, service := range services.Items {
		if strings.HasSuffix(service.ObjectMeta.Name, "-stargate-service") {
			return service
		}
	}
	Fail(fmt.Sprintf("Couldn't find the Stargate service"))
	return v1.Service{}
}

var (
	Info    = Yellow
	Outline = Purple
	Success = Green
	Running = Teal
	Failed  = Red
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
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

func deleteKindCluster() {
	RunShellCommand(exec.Command("kind", "delete", "cluster"))
}

// MinIO related functions
func GetMinioServiceName() string {
	kubectlOptions := k8s.NewKubectlOptions("", "", "minio")
	minioService, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), kubectlOptions, "get", "services", "-l", "app=minio", "-o", "jsonpath='{.items[0].metadata.name}'")
	if err != nil {
		Fail("Failed identifying the minio service")
	}
	log.Println(fmt.Sprintf("Minio service: %s", minioService))

	return minioService
}

///////////////////////////////////
// Test steps
///////////////////////////////////

// To create a new step, use the following template:
//
// func iDoSomethingImportantStep(arg1, arg2 string) {
//	By(fmt.Sprintf("Describe what the step does")
//
//
//
//	step implementation
// }

func ICanCheckThatResourceOfTypeWithLabelIsPresentInNamespaceStep(resourceType, label string) {
	By(fmt.Sprintf("I can check that a resource of type %s with label %s is present", resourceType, label))

	attempts := 0
	maxAttempts := 2
	for {
		attempts++
		switch resourceType {
		case "service":
			services, _ := GetServicesWithLabel(label)
			if len(services.Items) > 0 {
				return
			}
		default:
			Fail(fmt.Sprintf("Unsupported resource type for presence check: %s", resourceType))
		}
		// Not found yet
		time.Sleep(10 * time.Second)
		if attempts > maxAttempts {
			Fail(fmt.Sprintf("Resource of type %s with label %s was not found in namespace %s", resourceType, label, namespace))
		}
	}
}

func ICanCheckThatResourceOfTypeWithNameIsPresentInNamespaceStep(resourceType, name string) {
	By(fmt.Sprintf("I can check that resource %s of type %s is present", name, resourceType))

	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	switch resourceType {
	case "service":
		k8s.GetService(GinkgoT(), kubectlOptions, name)
	default:
		Fail(fmt.Sprintf("Unsupported resource type: %s", resourceType))
	}
}

func AKindClusterIsRunningAndReachableStep(clusterType string) {
	By("Starting a kind cluster with " + clusterType)
	deleteKindCluster()
	var kindClusterShell string
	switch clusterType {
	case "one worker":
		kindClusterShell, _ = filepath.Abs("scripts/cluster_one_worker.sh")
	case "three workers":
		kindClusterShell, _ = filepath.Abs("scripts/cluster_three_workers.sh")
	default:
		Fail(fmt.Sprintf("Kind cluster creation shell script not found for %s", clusterType))
	}
	err := RunShellCommand(exec.Command(kindClusterShell))

	Expect(err).To(BeNil())
}

func ICanDeleteTheKindClusterStep() {
	By("I can delete the kind cluster")

	deleteKindCluster()
}

func IDeployAClusterWithOptionsInTheNamespaceUsingTheValuesStep(options, customValues string) {
	By(fmt.Sprintf("I can deploy a cluster with %s options using the %s values", options, customValues))

	helmValues := map[string]string{}
	if options == "default" {
		helmValues = map[string]string{
			"reaper.ingress.host": "repair.localhost",
		}
	}
	if options == "minio" {
		serviceName := GetMinioServiceName()
		helmValues = map[string]string{
			"medusa.storage_properties.host": fmt.Sprintf("%s.minio.svc.cluster.local", strings.ReplaceAll(serviceName, "'", "")),
		}
	}
	deployCluster(customValues, helmValues)
}

func IDeployAClusterWithCassandraHeapAndMBStargateHeapUsingTheValuesStep(options, cassandraHeap, stargateHeap, customValues string) {
	By(fmt.Sprintf("I can deploy a cluster with %s options, %s Cassandra Heap and %s Stargate heap using the %s values", options, cassandraHeap, stargateHeap, customValues))

	splitOptions := strings.Split(options, "-")
	medusaEnabled := "true"
	reaperEnabled := "true"
	monitoringEnabled := "true"

	if Find(splitOptions, "nomedusa") {
		medusaEnabled = "false"
	}

	if Find(splitOptions, "noreaper") {
		reaperEnabled = "false"
	}

	if Find(splitOptions, "nomonitoring") {
		monitoringEnabled = "false"
	}

	newGenSize, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(cassandraHeap, "M", ""), "G", ""))
	helmValues := map[string]string{
		"cassandra.heap.size":           cassandraHeap,
		"cassandra.heap.newGenSize":     strconv.Itoa(newGenSize/2) + "M",
		"stargate.heapMB":               strings.ReplaceAll(stargateHeap, "M", ""),
		"medusa.enabled":                medusaEnabled,
		"reaper.enabled":                reaperEnabled,
		"reaper-operator.enabled":       reaperEnabled,
		"kube-prometheus-stack.enabled": monitoringEnabled,
	}
	deployCluster(customValues, helmValues)
}

func ICanSeeTheNamespaceInTheListOfNamespacesStep() {
	By(fmt.Sprintf("I can see the %s namespace in the list of namespaces", namespace))

	kubectlOptions := k8s.NewKubectlOptions("", "", "default")
	_, err := k8s.GetNamespaceE(GinkgoT(), kubectlOptions, namespace)
	if err != nil {
		Fail(fmt.Sprintf("Couldn't find namespace %s", namespace))
	}
}

func ICanSeeTheSecretInTheListOfSecretsInTheNamespaceStep(secret string) {
	By(fmt.Sprintf("I can see the %s secret in the namespaces", secret))

	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	_, err := k8s.GetSecretE(GinkgoT(), kubectlOptions, secret)
	if err != nil {
		Fail(fmt.Sprintf("Couldn't find secret %s", secret))
	}
}

func ICannotSeeTheNamespaceInTheListOfNamespacesStep() {
	By(fmt.Sprintf("I cannot see the %s namespace", namespace))

	kubectlOptions := k8s.NewKubectlOptions("", "", "default")
	attempts := 0
	maxAttempts := 10
	for {
		attempts++
		namespaceObject, _ := k8s.GetNamespaceE(GinkgoT(), kubectlOptions, namespace)
		if namespaceObject == nil {
			// namespace was deleted
			break
		}

		if namespaceObject.Status.Phase == v1.NamespaceTerminating {
			// namespace is terminating, which is good enough
			break
		}

		time.Sleep(10 * time.Second)
		if attempts > maxAttempts {
			Fail(fmt.Sprintf("namespace %s was supposed to be deleted but was found in the k8s cluster", namespace))
		}
	}
}

func ICreateTheNamespaceStep() {
	By(fmt.Sprintf("I create the %s namespaces", namespace))

	namespace = fmt.Sprintf("k8ssandra%s", time.Now().Format("2006010215040507"))
	log.Println(fmt.Sprintf("Creating namespace %s", namespace))
	kubectlOptions := k8s.NewKubectlOptions("", "", "default")
	k8s.CreateNamespace(GinkgoT(), kubectlOptions, namespace)
}

func IDeleteTheNamespaceStep() {
	By(fmt.Sprintf("I delete the namespace"))

	kubectlOptions := k8s.NewKubectlOptions("", "", "default")
	k8s.DeleteNamespace(GinkgoT(), kubectlOptions, namespace)
}

func IInstallTraefikStep() {
	By(fmt.Sprintf("I install Traefik"))

	kubectlOptions := k8s.NewKubectlOptions("", "", "default")
	options := &helm.Options{KubectlOptions: kubectlOptions}

	// Add traefik repo and update repos
	helm.RunHelmCommandAndGetOutputE(GinkgoT(), options, "repo", "add", "traefik", "https://helm.traefik.io/traefik")
	helm.RunHelmCommandAndGetOutputE(GinkgoT(), options, "repo", "update")

	// Deploy traefik
	// helm install traefik traefik/traefik -n traefik --create-namespace -f docs/content/en/docs/topics/ingress/traefik/kind-deployment/traefik.values.yaml
	valuesPath, _ := filepath.Abs("../../docs/content/en/docs/topics/ingress/traefik/kind-deployment/traefik.values.yaml")
	_, err := helm.RunHelmCommandAndGetOutputE(GinkgoT(), options, "install", "traefik", "traefik/traefik", "-n", "traefik", "--create-namespace", "-f", valuesPath)
	Expect(err).To(BeNil())
}

type credentials struct {
	username string
	password string
}

func ICanSeeThatTheKeyspaceExistsInCassandraInNamespaceStep(keyspace string) {
	By(fmt.Sprintf("I can see that the %s keyspace exists in Cassandra", keyspace))

	keyspaces := runCassandraQueryAndGetOutput("describe keyspaces")

	// Check that the keyspace exists in the list of keyspaces
	Expect(keyspaces).Should(ContainSubstring(keyspace))
}

func IWaitForTheReaperPodToBeReadyInNamespaceStep() {
	By(fmt.Sprintf("I wait for the Reaper pod to be ready"))

	WaitForPodWithLabelToBeReady("app.kubernetes.io/managed-by=reaper-operator", 30*time.Second, 10)
}

func ICanReadRowsInTheTableInTheKeyspaceStep(nbRows int, tableName, keyspaceName string) {
	By(fmt.Sprintf("I can read %d rows in table %s.%s", nbRows, keyspaceName, tableName))

	output := runCassandraQueryAndGetOutput(fmt.Sprintf("SELECT id FROM %s.%s", keyspaceName, tableName))

	// Check that we have the right number of rows
	Expect(output).Should(ContainSubstring(fmt.Sprintf("(%d rows)", nbRows)))
}

func ICreateTheTableInTheKeyspaceStep(tableName, keyspaceName string) {
	By(fmt.Sprintf("I can create table %s.%s", keyspaceName, tableName))

	runCassandraQueryAndGetOutput(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s with replication = {'class':'SimpleStrategy', 'replication_factor':1};", keyspaceName))
	runCassandraQueryAndGetOutput(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s(id timeuuid PRIMARY KEY, val text);", keyspaceName, tableName))
}

func ILoadRowsInTheTableInTheKeyspaceStep(nbRows int, tableName, keyspaceName string) {
	By(fmt.Sprintf("I load %d rows in table %s.%s", nbRows, keyspaceName, tableName))

	for i := 0; i < nbRows; i++ {
		runCassandraQueryAndGetOutput(fmt.Sprintf("INSERT INTO %s.%s(id,val) values(now(), '%d');", keyspaceName, tableName, i))
	}
}

// Medusa related functions
func ICreateTheMedusaSecretInTheNamespaceApplyingTheFileStep(secretFile string) {
	By(fmt.Sprintf("I create the Medusa secret applying the %s secret file", secretFile))

	home, _ := os.UserHomeDir()
	medusaSecretPath, _ := filepath.Abs(strings.Replace(secretFile, "~", home, 1))
	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	k8s.KubectlApply(GinkgoT(), kubectlOptions, medusaSecretPath)
}

func IPerformABackupWithMedusaNamedStep(backupName string) {
	By(fmt.Sprintf("I perform a backup with Medusa named %s", backupName))

	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	backupChartPath, err := filepath.Abs("../../charts/backup")
	if err != nil {
		Fail(fmt.Sprintf("Couldn't find the absolute path for backup charts"))
	}

	helmOptions := &helm.Options{
		SetValues: map[string]string{
			"name":                     backupName,
			"cassandraDatacenter.name": "dc1",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespace),
	}
	helm.Install(GinkgoT(), helmOptions, backupChartPath, "test")

	// Wait for the backup to be finished
	// kubectl get cassandrabackup backup1 -o jsonpath='{.status.finished}' -n k8ssandra2021021306435807
	attempts := 0
	for {
		attempts++
		output, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), kubectlOptions, "get", "cassandrabackup", backupName, "-o", "jsonpath='{.status.finished}'")
		if err == nil && len(output) > 0 {
			var nodes []string
			json.Unmarshal([]byte(strings.ReplaceAll(output, "'", "")), &nodes)
			if len(nodes) == 1 {
				return
			}
		}
		if attempts > 12 {
			Fail(fmt.Sprintf("Backup didn't succeed within timeout: %s", err))
		}
		time.Sleep(10 * time.Second)
	}
}

func IRestoreTheBackupNamedUsingMedusaStep(backupName string) {
	By(fmt.Sprintf("I restore the backup with Medusa named %s", backupName))

	restoreChartPath, err := filepath.Abs("../../charts/restore")
	if err != nil {
		Fail(fmt.Sprintf("Couldn't find the absolute path for restore charts"))
	}

	helmOptions := &helm.Options{
		SetValues: map[string]string{
			"backup.name":              backupName,
			"cassandraDatacenter.name": "dc1",
			"name":                     "restore-test",
		},
		KubectlOptions: k8s.NewKubectlOptions("", "", namespace),
	}
	helm.Install(GinkgoT(), helmOptions, restoreChartPath, "restore-test")
	// Give a little time for the cassandraDatacenter resource to be recreated
	time.Sleep(60 * time.Second)
	WaitForPodWithLabelToBeReady("app.kubernetes.io/managed-by=cass-operator", 30*time.Second, 10)
}

// Reaper related steps
func ICanCheckThatAClusterNamedWasRegisteredInReaperInNamespaceStep(clusterName string) {
	By(fmt.Sprintf("I can check that a cluster named %s was registered in Reaper", clusterName))

	restClient := resty.New()
	attempts := 0
	for {
		attempts++
		response, err := restClient.R().Get("http://repair.localhost:8080/cluster")
		if err != nil {
			log.Println(fmt.Sprintf("The HTTP request failed with error %s", err))
		} else {
			data := response.Body()
			log.Println(fmt.Sprintf("Reaper response: %s", data))
			var clusters []string
			json.Unmarshal([]byte(data), &clusters)
			if len(clusters) > 0 {
				Expect(clusterName).Should(Equal(clusters[0]))
				return
			}
		}
		time.Sleep(30 * time.Second)
		if attempts >= 10 {
			break
		}
	}
	Fail("Cluster wasn't properly registered in Reaper")
}

func ICanCancelTheRunningRepairStep() {
	By(fmt.Sprintf("I can cancel the running repair"))

	restClient := resty.New()
	// Start the previously created repair run
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		Put(fmt.Sprintf("http://repair.localhost:8080/repair_run/%s/state/ABORTED", repairId))

	log.Println(fmt.Sprintf("Reaper response: %s", response.Body()))
	log.Println(fmt.Sprintf("Reaper status code: %d", response.StatusCode()))

	if err != nil || response.StatusCode() != 200 {
		Fail(fmt.Sprintf("Failed aborting repair %s: %s / %s", repairId, err, response.Body()))
	}
}

func ITriggerARepairOnTheKeyspaceStep(keyspace string) {
	By(fmt.Sprintf("I trigger a repair on keyspace %s", keyspace))

	restClient := resty.New()

	// Create the repair run
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParams(map[string]string{
			"clusterName":  "k8ssandra",
			"keyspace":     keyspace,
			"owner":        "k8ssandra",
			"segmentCount": "5",
		}).
		Post("http://repair.localhost:8080/repair_run")

	data := response.Body()
	log.Println(fmt.Sprintf("Reaper response: %s", data))
	var reaperResponse interface{}
	err2 := json.Unmarshal(data, &reaperResponse)

	if err != nil || err2 != nil {
		Fail(fmt.Sprintf("The REST request or response parsing failed with error %s %s: %s", err, err2, data))
	}

	reaperResponseMap := reaperResponse.(map[string]interface{})
	repairId = fmt.Sprintf("%s", reaperResponseMap["id"])
	// Start the previously created repair run
	response, err = restClient.R().
		SetHeader("Content-Type", "application/json").
		Put(fmt.Sprintf("http://repair.localhost:8080/repair_run/%s/state/RUNNING", repairId))

	log.Println(fmt.Sprintf("Reaper response: %s", response.Body()))
	log.Println(fmt.Sprintf("Reaper status code: %d", response.StatusCode()))

	if err != nil || response.StatusCode() != 200 {
		Fail(fmt.Sprintf("Failed starting repair %s: %s / %s", repairId, err, response.Body()))
	}
}

func IWaitForAtLeastOneSegmentToBeProcessedStep() {
	By(fmt.Sprintf("I wait for at least one segment to be processed"))

	restClient := resty.New()
	attempts := 0
	for {
		attempts++
		response, err := restClient.R().Get(fmt.Sprintf("http://repair.localhost:8080/repair_run/%s/segments", repairId))
		if err != nil {
			log.Println(fmt.Sprintf("The HTTP request failed with error %s", err))
		}

		if strings.Contains(fmt.Sprintf("%s", response.Body()), "\"state\":\"DONE\"") {
			// We have at least one completed segment
			return
		}

		time.Sleep(30 * time.Second)
		if attempts >= 10 {
			// Too many attempts, failed test.
			log.Println(fmt.Sprintf("Latest segment list from Reaper: %s", response.Body()))
			break
		}
	}
	Fail(fmt.Sprintf("No repair segment was fully processed within timeout"))
}

func ICanRunACyclesStressTestWithReadsAndAOpssRateWithinTimeoutStep(stressCycles, percentRead string, rate, timeout int) {
	By(fmt.Sprintf("I run a %s cycles stess test with %s reads and %d ops/s within %d seconds", stressCycles, percentRead, rate, timeout))

	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)
	cqlCredentials := GetUsernamePassword("k8ssandra-superuser", namespace)
	parsedReadRatio, _ := strconv.Atoi(strings.ReplaceAll(percentRead, "%", ""))
	readRatio := parsedReadRatio / 10
	writeRatio := 10 - readRatio

	jobName := fmt.Sprintf("nosqlbench-%s", strings.ToLower(random.UniqueId()))
	k8s.RunKubectl(GinkgoT(), kubectlOptions, "create", "job", "--image=nosqlbench/nosqlbench", jobName,
		"--", "java", "-jar", "nb.jar", "cql-iot", "rampup-cycles=1k", fmt.Sprintf("cyclerate=%d", rate),
		fmt.Sprintf("username=%s", cqlCredentials.username), fmt.Sprintf("password=%s", cqlCredentials.password),
		fmt.Sprintf("main-cycles=%s", stressCycles), "hosts=k8ssandra-dc1-stargate-service", "--progress", "console:1s", "-v",
		fmt.Sprintf("write_ratio=%d", writeRatio), fmt.Sprintf("read_ratio=%d", readRatio), "async=100")

	defer timeTrack(time.Now(), fmt.Sprintf("nosqlbench stress test with %d ops/s", rate))
	k8s.RunKubectl(GinkgoT(), kubectlOptions, "wait", "--for=condition=complete", fmt.Sprintf("--timeout=%ds", timeout), fmt.Sprintf("job/%s", jobName))

	output := RunShellCommandAndGetOutput(
		exec.Command("bash", "-c", fmt.Sprintf("kubectl logs job/%s -n %s | grep -e cqliot_default_main.cycles.servicetime -e cqliot_default_main.cycles.responsetime", jobName, namespace)))
	log.Println(Outline(output))
}

func IWaitForTheStargatePodsToBeReadyStep() {
	By(fmt.Sprintf("I wait for the Stargate pods to be ready"))

	kubectlOptions := k8s.NewKubectlOptions("", "", namespace)

	attempts := 0
	for {
		attempts++
		output, err := k8s.RunKubectlAndGetOutputE(GinkgoT(), kubectlOptions, "rollout", "status", "deployment", "k8ssandra-dc1-stargate")
		if err == nil {
			if strings.HasSuffix(output, "successfully rolled out") {
				return
			}
		}
		time.Sleep(30 * time.Second)
		if attempts >= 10 {
			// Too many attempts, failed test.
			break
		}
	}
	Fail("Stargate deployment didn't roll out within timeout")
}

func IDeployMinIOUsingHelmAndCreateTheBucketStep(bucketName string) {
	By(fmt.Sprintf("I deploy MinIO using Helm and create the %s bucket", bucketName))

	helmOptions := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "default"),
	}
	helm.RunHelmCommandAndGetOutputE(GinkgoT(), helmOptions, "repo", "add", "minio", "https://helm.min.io/")

	helm.RunHelmCommandAndGetOutputE(GinkgoT(), helmOptions, "install",
		"--set", fmt.Sprintf("accessKey=minio_key,secretKey=minio_secret,defaultBucket.enabled=true,defaultBucket.name=%s", bucketName),
		"--generate-name", "minio/minio", "-n", "minio", "--create-namespace")
}
