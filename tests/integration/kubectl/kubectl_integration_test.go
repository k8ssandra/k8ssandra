package kubectl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
	authv1 "k8s.io/api/authorization/v1"
)

// An example of verifying the kubectl is operational and that
// authorization for basic operations.
func TestKubectlAuthorization(t *testing.T) {

	t.Parallel()

	ns := fmt.Sprintf("test-ns-%s", strings.ToLower(random.UniqueId()))
	kubectlOptions := k8s.NewKubectlOptions("", "", ns)

	adminGetServiceAction := authv1.ResourceAttributes{
		Namespace: ns,
		Verb:      "get",
		Resource:  "service",
	}

	adminGetPodAction := authv1.ResourceAttributes{
		Namespace: ns,
		Verb:      "get",
		Resource:  "pod",
	}

	require.True(t, k8s.CanIDo(t, kubectlOptions, adminGetServiceAction))
	require.True(t, k8s.CanIDo(t, kubectlOptions, adminGetPodAction))
}
