# k8ssandra-operator

### Steps to renegerate using kustomize in k8ssandra-operator

You need to manually modify and verify the labels and metadata names, these steps do not automate them.

#### CRDs from k8ssandra-operator:

```
kustomize build config/crd  --output ../k8ssandra/charts/k8ssandra-operator/crds/
```

#### To build RBACs:

Add to kustomization.yaml the following to simplify the name verifications and namespace removals (one or two are left behind which you need to manually cleanup as well as remove some extra ' characters)

```yaml
namePrefix: '{{ include "k8ssandra-common.fullname" . }}-'
namespace: 
```

```
kustomize build config/rbac  --output ../k8ssandra/charts/k8ssandra-operator/templates/
```

#### The deployment and config:

Add configMap patch to kustomization.yaml:

```yaml
patchesStrategicMerge:
- ../default/manager_config_patch.yaml
```

Then build with output:

```
kustomize build config/manager  --output ../k8ssandra/charts/k8ssandra-operator/templates/
```
