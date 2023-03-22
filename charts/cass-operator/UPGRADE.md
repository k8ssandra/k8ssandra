# cass-operator

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