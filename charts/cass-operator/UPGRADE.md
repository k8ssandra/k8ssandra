# cass-operator

## Upgrading from older release of cass-operator chart

If you run into a following problem when updating from Chart version older than 0.47.0 to a newer one:

```console
Error: UPGRADE FAILED: Unable to continue with update: CustomResourceDefinition "cassandradatacenters.cassandra.datastax.com" in namespace "" exists and cannot be imported into the current release: invalid ownership metadata; label validation error: missing key "app.kubernetes.io/managed-by": must be set to "Helm"; annotation validation error: missing key "meta.helm.sh/release-name": must be set to "cass-operator"; annotation validation error: missing key "meta.helm.sh/release-namespace": must be set to "cass-operator"
```

Run the following commands:

```
kubectl label --overwrite crd cassandradatacenters.cassandra.datastax.com app.kubernetes.io/managed-by=Helm
kubectl label --overwrite crd cassandratasks.control.k8ssandra.io app.kubernetes.io/managed-by=Helm
kubectl annotate --overwrite crd cassandradatacenters.cassandra.datastax.com meta.helm.sh/release-namespace=cass-operator
kubectl annotate --overwrite crd cassandratasks.control.k8ssandra.io meta.helm.sh/release-namespace=cass-operator
kubectl annotate --overwrite crd cassandratasks.control.k8ssandra.io meta.helm.sh/release-name="cass-operator"
kubectl annotate --overwrite crd cassandradatacenters.cassandra.datastax.com meta.helm.sh/release-name="cass-operator"
```

Replace cass-operator in the ``release-name="cass-operator"`` part with your current Helm release name. After that, run the upgrade command again.

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