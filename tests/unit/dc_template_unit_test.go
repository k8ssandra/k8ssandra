package tests

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestK8ssandraClusterTemplate(t *testing.T) {

	var renderedMap map[string]interface{}
	helmChartPath, err := filepath.Abs("../../charts/k8ssandra-cluster")
	metadataName := fmt.Sprintf("test-meta-name-%s", strings.ToLower(random.UniqueId()))
	clusterName := fmt.Sprintf("test-cluster-name-%s", strings.ToLower(random.UniqueId()))
	require.NoError(t, err)

	options := &helm.Options{
		SetStrValues:   map[string]string{"name": metadataName, "clusterName": clusterName},
		KubectlOptions: k8s.NewKubectlOptions("", "", "k8ssandra"),
	}

	renderedOutput := helm.RenderTemplate(
		t, options, helmChartPath, "k8ssandra-test",
		[]string{"templates/cassdc.yaml"},
	)
	jsonRendered, err := yaml.YAMLToJSON([]byte(renderedOutput))
	require.NoError(t, json.Unmarshal(jsonRendered, &renderedMap))

	metadata := renderedMap["metadata"].(map[string]interface{})
	spec := renderedMap["spec"].(map[string]interface{})

	require.NotNil(t, metadata)
	require.NotNil(t, spec)

	require.Equal(t, "cassandra", spec["serverType"])
	require.Equal(t, "CassandraDatacenter", renderedMap["kind"])
	require.Equal(t, "cassandra.datastax.com/v1beta1", renderedMap["apiVersion"])
	require.Equal(t, clusterName, spec["clusterName"])
	require.Equal(t, metadataName, metadata["name"])

}
