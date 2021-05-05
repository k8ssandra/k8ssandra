package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/gruntwork-io/terratest/modules/k8s"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Stargate related functions
func StargateService(t *testing.T, namespace string) (v1.Service, error) {
	clientset, err := k8s.GetKubernetesClientFromOptionsE(t, getKubectlOptions(namespace))
	g(t).Expect(err).To(BeNil())

	services, _ := clientset.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	for _, service := range services.Items {
		if label, ok := service.ObjectMeta.Labels["app"]; ok {
			if label == fmt.Sprintf("%s-%s-stargate", releaseName, datacenterName) {
				return service, nil
			}
		}
	}
	return v1.Service{}, fmt.Errorf("failed finding Stargate service")
}

func WaitForStargatePodReady(t *testing.T, namespace string) {
	g(t).Eventually(func() bool {
		return stargatePodIsReady(t, namespace)
	}, retryTimeout, retryInterval).Should(BeTrue(), "Stargate deployment didn't roll out within timeout")
}

func stargatePodIsReady(t *testing.T, namespace string) bool {
	output, err := k8s.RunKubectlAndGetOutputE(t, getKubectlOptions(namespace), "rollout", "status", "deployment", fmt.Sprintf("%s-%s-stargate", releaseName, datacenterName))
	if err == nil {
		if strings.HasSuffix(output, "successfully rolled out") {
			return true
		}
	}
	return false
}

func WaitForAuthEndpoint(t *testing.T) {
	g(t).Eventually(func() bool {
		return authEndpointIsReachable(t)
	}, 2*time.Minute, retryInterval).Should(BeTrue())
}

func authEndpointIsReachable(t *testing.T) bool {
	restClient := resty.New()
	response, err := restClient.R().Get("http://stargate.127.0.0.1.nip.io:8081/v1/auth")
	if err != nil {
		log.Printf("Failed connecting to Stargate auth endpoint: %s", err.Error())
		return false
	}
	return response.StatusCode() == 405 // This endpoint should be invoked with a Post, we expect a 405 with a Get
}

func GenerateStargateAuthToken(t *testing.T, namespace string) string {
	credentials := ExtractUsernamePassword(t, "k8ssandra-superuser", namespace)
	creds := fmt.Sprintf("{\"username\":\"%s\", \"password\":\"%s\"}", credentials.username, credentials.password)
	restClient := resty.New()
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(creds).
		Post("http://stargate.127.0.0.1.nip.io:8081/v1/auth")
	g(t).Expect(err).To(BeNil(), "Failed connecting to Stargate")
	stargateResponse := response.Body()
	var genericJson map[string]interface{}
	json.Unmarshal(stargateResponse, &genericJson)
	g(t).Expect(genericJson["authToken"]).ToNot(BeNil(), fmt.Sprintf("Expected response to have authToken property: %#v", genericJson))
	return genericJson["authToken"].(string)
}

func CreateStargateDocumentNamespace(t *testing.T, token string) string {
	docNamespace := fmt.Sprintf("stargate%s", time.Now().Format("2006010215040507"))
	docNamespaceJson := fmt.Sprintf("{\"name\":\"%s\"}", docNamespace)
	restClient := resty.New()
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Cassandra-Token", token).
		SetBody(docNamespaceJson).
		Post("http://stargate.127.0.0.1.nip.io:8082/v2/schemas/namespaces")
	g(t).Expect(err).To(BeNil(), "Failed creating Stargate document namespace")
	stargateResponse := string(response.Body())
	g(t).Expect(stargateResponse).To(Equal(docNamespaceJson), fmt.Sprintf("Unexpected response from Stargate: %s", stargateResponse))
	return docNamespace
}

const (
	awesomeMovieDirector = "Zack Snyder"
	awesomeMovieName     = "Watchmen"
)

func WriteStargateDocument(t *testing.T, token, docNamespace string) string {
	awesomeMovieDocument := fmt.Sprintf("{\"Director\":\"%s\",\"Name\":\"%s\"}", awesomeMovieDirector, awesomeMovieName)
	documentId := fmt.Sprintf("watchmen%s", time.Now().Format("2006010215040507"))
	restClient := resty.New()
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Cassandra-Token", token).
		SetBody(awesomeMovieDocument).
		Put(fmt.Sprintf("http://stargate.127.0.0.1.nip.io:8082/v2/namespaces/%s/collections/movies/%s", docNamespace, documentId))
	g(t).Expect(err).To(BeNil(), "Failed creating Stargate document")
	stargateResponse := string(response.Body())
	expectedResponse := fmt.Sprintf("{\"documentId\":\"%s\"}", documentId)
	g(t).Expect(stargateResponse).To(Equal(expectedResponse))
	return documentId
}

func CheckStargateDocumentExists(t *testing.T, token, docNamespace, documentId string) {
	document := readStargateDocument(t, token, docNamespace, documentId)
	var genericJson map[string]interface{}
	json.Unmarshal(document, &genericJson)
	g(t).Expect(genericJson["documentId"].(string)).To(Equal(documentId))
	g(t).Expect(genericJson["data"].(map[string]interface{})["Director"].(string)).To(Equal(awesomeMovieDirector))
	g(t).Expect(genericJson["data"].(map[string]interface{})["Name"].(string)).To(Equal(awesomeMovieName))
}

func readStargateDocument(t *testing.T, token, docNamespace, documentId string) []byte {
	restClient := resty.New()
	response, err := restClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Cassandra-Token", token).
		Get(fmt.Sprintf("http://stargate.127.0.0.1.nip.io:8082/v2/namespaces/%s/collections/movies/%s", docNamespace, documentId))
	g(t).Expect(err).To(BeNil(), "Failed reading Stargate document")
	return response.Body()
}
