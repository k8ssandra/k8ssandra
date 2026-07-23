# cass-operator

![Version: 0.66.0](https://img.shields.io/badge/Version-0.66.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.31.0](https://img.shields.io/badge/AppVersion-1.31.0-informational?style=flat-square)

Kubernetes operator which handles the provisioning and management of Apache Cassandra clusters.

**Homepage:** <https://github.com/k8ssandra/cass-operator>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Cass Operator Team | <cass-operator@datastax.com> | <https://github.com/k8ssandra/cass-operator> |
| K8ssandra Team | <k8ssandra-developers@googlegroups.com> | <https://github.com/k8ssandra> |

## Source Code

* <https://github.com/k8ssandra/cass-operator>
* <https://github.com/k8ssandra/k8ssandra/tree/main/charts/cass-operator>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../k8ssandra-common | k8ssandra-common | 0.29.2 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| global.clusterScoped | bool | `false` | Determines whether cass-operator only watch and manages CassandraDatacenters in the same namespace in which the operator is deployed or if watches and manages CassandraDatacenters across all namespaces. |
| global.watchNamespaces | list | `[]` | List of namespaces to watch for CassandraDatacenter resources. If empty, the operator will watch all namespaces in cluster-scope installation. |
| global.commonAnnotations | object | `{}` | Annotations to be added to all deployed resources. |
| global.commonLabels | object | `{}` | Labels to be added to all deployed resources. |
| global.clusterScopedResources | bool | `true` | Should we install cluster-scoped resources such as ClusterRole, ClusterRoleBinding |
| global.imageConfig | object | `{"defaults":{"pullPolicy":"IfNotPresent","registry":"docker.io"},"images":{"config-builder":{"name":"cass-config-builder","repository":"datastax","tag":"1.0-ubi"},"k8ssandra-client":{"name":"k8ssandra-client","repository":"k8ssandra","tag":"v0.8.13"},"system-logger":{"name":"system-logger","repository":"k8ssandra","tag":"v1.31.0"}},"types":{"cassandra":{"name":"cass-management-api","repository":"k8ssandra","suffix":"-ubi"},"dse":{"name":"dse-mgmtapi-6_8","repository":"datastax","suffix":"-ubi"}}}` | When set, this replaces the legacy imageConfig settings above. Set the default images used by all the k8ssandra operators. |
| global.imageConfig.images | object | `{"config-builder":{"name":"cass-config-builder","repository":"datastax","tag":"1.0-ubi"},"k8ssandra-client":{"name":"k8ssandra-client","repository":"k8ssandra","tag":"v0.8.13"},"system-logger":{"name":"system-logger","repository":"k8ssandra","tag":"v1.31.0"}}` | Default images used by cass-operator when creating Cassandra resources. |
| global.imageConfig.defaults | object | `{"pullPolicy":"IfNotPresent","registry":"docker.io"}` | defaults are used when other no override is present for the configured image or image type. All these values can be overridden under any image or image type. |
| global.imageConfig.types | object | `{"cassandra":{"name":"cass-management-api","repository":"k8ssandra","suffix":"-ubi"},"dse":{"name":"dse-mgmtapi-6_8","repository":"datastax","suffix":"-ubi"}}` | types are used as building blocks when serverVersion + serverType in the CassandraDatacenter is used |
| nameOverride | string | `""` | A name in place of the chart name which is used in the metadata.name of objects created by this chart. |
| fullnameOverride | string | `""` | A name in place of the value used for metadata.name in objects created by this chart. The default value has the form releaseName-chartName. |
| commonLabels | object | `{}` | Labels to be added to all deployed resources. |
| commonAnnotations | object | `{}` | Annotations to be added to all deployed resources. |
| replicaCount | int | `1` | Sets the number of cass-operator pods. |
| admissionWebhooks | object | `{"certificateSecret":"","enabled":true}` | Configures admission webhooks deployed with cass-operator. |
| admissionWebhooks.enabled | bool | `true` | Turns the admission webhooks on or off. We recommend keeping them on. |
| admissionWebhooks.certificateSecret | string | `""` | Externally managed Secret containing tls.crt, tls.key, and ca.crt for the webhook. When set, the chart does not create its self-signed webhook Certificate. The Secret must have the cert-manager.io/allow-direct-injection: "true" annotation. |
| image.registry | string | `"docker.io"` | Container registry containing the repository where the image resides |
| image.repository | string | `"k8ssandra/cass-operator"` | Docker repository for cass-operator |
| image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the operator container |
| image.tag | string | `"v1.31.0"` | Tag of the cass-operator image to pull from image.repository |
| imagePullSecrets | list | `[]` | References to secrets to use when pulling images. See: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/ |
| serviceAccount | object | `{"annotations":{},"name":""}` | Configures the ServiceAccount used by cass-operator. |
| serviceAccount.name | string | `""` | Name of the ServiceAccount. When empty, the chart-generated name is used. |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account. |
| podAnnotations | object | `{}` | Annotations for the cass-operator pod. |
| podSecurityContext | object | `{}` | PodSecurityContext for the cass-operator pod. See: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ |
| securityContext | object | `{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"runAsGroup":65534,"runAsNonRoot":true,"runAsUser":65534,"seccompProfile":{"type":"RuntimeDefault"}}` | SecurityContext for the cass-operator container. |
| securityContext.runAsNonRoot | bool | `true` | Run cass-operator container as non-root user |
| securityContext.runAsGroup | int | `65534` | Group for the user running the cass-operator container / process |
| securityContext.runAsUser | int | `65534` | User for running the cass-operator container / process |
| securityContext.readOnlyRootFilesystem | bool | `true` | Run cass-operator container having read-only root file system permissions. |
| securityContext.allowPrivilegeEscalation | bool | `false` | Do not allow privilege escalation for the cass-operator container. |
| securityContext.capabilities | object | `{"drop":["ALL"]}` | Linux capabilities for the cass-operator container. |
| securityContext.capabilities.drop | list | `["ALL"]` | Linux capabilities to drop. |
| securityContext.seccompProfile | object | `{"type":"RuntimeDefault"}` | Seccomp profile for the cass-operator container. |
| securityContext.seccompProfile.type | string | `"RuntimeDefault"` | Seccomp profile type. |
| resources | object | `{}` | Resources requests and limits for the cass-operator pod. We usually recommend not to specify default resources and to leave this as a conscious choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you want to specify resources, add `requests` and `limits` for `cpu` and `memory` while removing the existing `{}` |
| nodeSelector | object | `{}` | Node labels for operator pod assignment. Ref: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/  |
| tolerations | list | `[]` | Node tolerations for server scheduling to nodes with taints. Ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/  |
| metrics | object | `{"address":":8080","tls":{"enabled":false}}` | metrics allows to change the configuration of how metrics are exposed. Default is to expose /metrics endpoint at port :8080. If you wish to make this available only on localhost (such as when using kube-rbac-proxy to secure access to them), set the value to 127.0.0.1:8080 |
| metrics.address | string | `":8080"` | Address where the metrics endpoint binds. Use "0" to disable metrics. |
| metrics.tls | object | `{"enabled":false}` | TLS configuration for the metrics endpoint. |
| metrics.tls.enabled | bool | `false` | Enable TLS for metrics endpoint |
| disableCertManagerCheck | bool | `false` | Bypass the cert-manager API availability check. Prefer exposing cert-manager.io/v1 in Helm capabilities when rendering templates. |
| openshift | object | `{"mode":"auto"}` | OpenShift-specific behavior for the cass-operator pod. |
| openshift.mode | string | `"auto"` | OpenShift behavior mode. Supported values are: auto, enabled, disabled. In auto or enabled mode, cass-operator omits runAsUser and runAsGroup on OpenShift so restricted-v2 SCC can assign the namespace UID/GID range. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.14.2](https://github.com/norwoodj/helm-docs/releases/v1.14.2)
