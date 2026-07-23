# cass-operator

## Upgrading to the next chart release (0.66)

This release removes the unused `image.repositoryOverride`, `image.registryOverride`, and `image.namespaceOverride` values. Configure Cassandra image defaults and overrides through `global.imageConfig` instead since the older ImageConfig v1 is removed from the current version of cass-operator.

If you wish to disable the cert-manager managed self signed certificate, the process has changed a little. Set `admissionWebhooks.certificateSecret` to use an externally managed webhook TLS certificates. Setting that secret will override our default method of creating a self-signed certificate for the webhooks. Note that the Secret must have certain structure at this point:

* Must contain `tls.crt`, `tls.key`, and `ca.crt`
* Have the `cert-manager.io/allow-direct-injection: "true"` annotation.

This does not remove the requirement to have cert-manager installed in the cluster.

### Steps to renegerate using kustomize in cass-operator

You need to manually modify and verify the labels and metadata names, these steps do not automate them.

#### CRDs from cass-operator:

From cass-operator directory, assuming k8ssandra is checked out at ../k8ssandra:

```
scripts/release-helm-chart.sh version
```

Replace version with the intended tag, without the "v" prefix.

#### To build RBACs:

Add to kustomization.yaml the following to simplify the name verifications and namespace removals (one or two are left behind which you need to manually cleanup as well as remove some extra ' characters)

```yaml
namePrefix: '{{ include "k8ssandra-common.fullname" . }}-'
namespace: 
```

```
kustomize build config/rbac  --output ../k8ssandra/charts/cass-operator/templates/
```

#### The deployment and config:

Add configMap patch to kustomization.yaml:

```yaml
patchesStrategicMerge:
- ../default/manager_config_patch.yaml
```

Then build with output:

```
kustomize build config/manager  --output ../k8ssandra/charts/cass-operator/templates/
```
