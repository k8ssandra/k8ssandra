package tests

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	dc "../../cassandra/dc/typed"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestK8ssandraClusterTemplate(t *testing.T) {

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

	var dcStruct dc.CassDC
	helm.UnmarshalK8SYaml(t, renderedOutput, &dcStruct)

	require.Equal(t, "CassandraDatacenter", dcStruct.Kind)
	require.Equal(t, clusterName, dcStruct.Spec.ClusterName)
	require.Equal(t, metadataName, dcStruct.Metadata.Name)
}
