# k8ssandra

![Version: 1.4.0-SNAPSHOT](https://img.shields.io/badge/Version-1.4.0--SNAPSHOT-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

Provisions and configures an instance of the entire K8ssandra stack. This includes Apache Cassandra, Stargate, Reaper, Medusa, Prometheus, and Grafana.

**Homepage:** <https://k8ssandra.io/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| K8ssandra Team | k8ssandra-developers@googlegroups.com | https://github.com/k8ssandra |

## Source Code

* <https://github.com/k8ssandra/k8ssandra>
* <https://github.com/k8ssandra/k8ssandra/tree/main/charts/k8ssandra>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../cass-operator | cass-operator | 0.30.0 |
| file://../k8ssandra-common | k8ssandra-common | 0.28.4 |
| file://../k8ssandra-operator | k8ssandra-operator | 0.30.1 |
| file://../medusa-operator | medusa-operator | 0.30.1 |
| file://../reaper-operator | reaper-operator | 0.32.2|
| https://prometheus-community.github.io/helm-charts | kube-prometheus-stack | 12.11.3 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| cassandra.enabled | bool | `true` | Enables installation of Cassandra cluster. Set to false if you only wish to install operators. |
| cassandra.version | string | `"4.0.1"` | The Cassandra version to use. The supported versions include the following:    - 3.11.7    - 3.11.8    - 3.11.9    - 3.11.10    - 3.11.11    - 4.0.0    - 4.0.1 |
| cassandra.versionImageMap | object | `{"3.11.10":"k8ssandra/cass-management-api:3.11.10-v0.1.27","3.11.11":"k8ssandra/cass-management-api:3.11.11-v0.1.33","3.11.7":"k8ssandra/cass-management-api:3.11.7-v0.1.33","3.11.8":"k8ssandra/cass-management-api:3.11.8-v0.1.33","3.11.9":"k8ssandra/cass-management-api:3.11.9-v0.1.27","4.0.0":"k8ssandra/cass-management-api:4.0.0-v0.1.33","4.0.1":"k8ssandra/cass-management-api:4.0.1-v0.1.33"}` | Specifies the image to use for a particular Cassandra version. Exercise care and caution with changing these values! cass-operator is not designed to work with arbitrary Cassandra images. It expects the cassandra container to be running management-api images. If you do want to change one of these mappings, the new value should be a management-api image. |
| cassandra.image | object | `{}` | Overrides the default image mappings. This is intended for advanced use cases like development or testing. By default the Cassandra version has to be one that is in versionImageMap. Template rendering will fail if the version is not in the map. When you set the image directly, the version mapping check is skipped. Note that you are still constrained to the versions supported by cass-operator. |
| cassandra.securityContext | object | `{}` | Security context override for Cassandra container.
| cassandra.podSecurityContext | object | `{}` | Security context override for Cassandra pod |
| cassandra.baseConfig | object | `{}` | Cassandra base init container |
| cassandra.baseConfig.securityContext | object | `{}` | Security context override for base init container. |
| cassandra.configBuilder | object | `{"image":{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"datastax/cass-config-builder","tag":"1.0.4"},"securityContext":{}}` | The server-config-init init container |
| cassandra.configBuilder.securityContext | object | `{}` | Security context override for server-config-init container. |
| cassandra.configBuilder.image.registry | string | `"docker.io"` | Container registry for the config builder |
| cassandra.configBuilder.image.repository | string | `"datastax/cass-config-builder"` | Repository for cass-config-builder image |
| cassandra.configBuilder.image.tag | string | `"1.0.4"` | Tag of the config builder image to pull from image repository |
| cassandra.configBuilder.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the config builder image |
| cassandra.jmxCredentialsConfig | object | `{"image":{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"busybox","tag":"1.33.1"},"securityContext":{}}` | The jmx-credentials init container that configures JMX credentials. |
| cassandra.jmxCredentialsConfig.securityContext | object | `{}` | Security context override for jmx init container. |
| cassandra.jmxCredentialsConfig.image.registry | string | `"docker.io"` | Container registry for the jmx-credentials container |
| cassandra.jmxCredentialsConfig.image.repository | string | `"busybox"` | Repository for jmx-credentials container image |
| cassandra.jmxCredentialsConfig.image.tag | string | `"1.33.1"` | Tag of the jmx-credentials image to pull from image repository |
| cassandra.jmxCredentialsConfig.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the jmx-credentials image |
| cassandra.clusterName | string | `""` | The ServiceAccount to use for Cassandra pods. If not defined, defaults to the default account for the namespace. serviceAccount: -- Cluster name defaults to release name when not specified. |
| cassandra.auth | object | `{"cacheUpdateIntervalMillis":3600000,"cacheValidityPeriodMillis":3600000,"enabled":true,"superuser":{"secret":"","username":""}}` | Authentication and authorization related settings. |
| cassandra.auth.enabled | bool | `true` | Enables or disables authentication and authorization. This also enables/disables JMX authentication. Note that if Reaper is enabled JMX authentication will still be enabled even if auth is disabled here. This is because Reaper requires remote JMX access. |
| cassandra.auth.superuser | object | `{"secret":"","username":""}` | Configures the default Cassandra superuser when authentication is enabled. If neither `superuser.secret` nor `superuser.username` are set, then a user and a secret with the user's credentials will be created. The username and secret name will be of the form {clusterName}-superuser. The password will be a random 20 character password. If `superuser.secret` is set, then the Cassandra user will be created from the contents of the secret. If `superuser.secret` is not set and if `superuser.username` is set, a secret will be generated using the specified username. The password will be generated as previously described. JMX credentials will also be created for the superuser. The same username/password that is used here will be used for JMX. If you change the Cassandra superuser credentials through cqlsh for example, the JMX credentials will not be updated. You need to update the credentials via helm upgrade in order for the change to propagate to JMX. This will be fixed in https://github.com/k8ssandra/k8ssandra/issues/323. |
| cassandra.auth.cacheValidityPeriodMillis | int | `3600000` | Cache entries validity period in milliseconds. cassandra.yaml has settings for roles, permissions, and credentials caches. This property will configure the validity period for all three. |
| cassandra.auth.cacheUpdateIntervalMillis | int | `3600000` | Cache entries update period in milliseconds. cassandra.yaml has settings for roles, permissions, and credentials caches. This property will configure the update interval for all three. |
| cassandra.cassandraLibDirVolume.storageClass | string | `"standard"` | Storage class for persistent volume claims (PVCs) used by the underlying cassandra pods. Depending on your Kubernetes distribution this may be named "standard", "hostpath", or "localpath". Run `kubectl get storageclass` to identify what is available in your environment. |
| cassandra.cassandraLibDirVolume.size | string | `"5Gi"` | Size of the provisioned persistent volume per node. It is recommended to keep the total amount of data per node to approximately 1 TB. With room for compactions this value should max out at ~2 TB. This recommendation is highly dependent on data model and compaction strategies in use. Consider testing with your data model to find an optimal value for your usecase. |
| cassandra.allowMultipleNodesPerWorker | bool | `false` | Permits running multiple Cassandra pods per Kubernetes worker. If enabled resources.limits and resources.requests **must** be defined. |
| cassandra.additionalSeeds | list | `[]` | Optional additional contact points for the Cassandra cluster to connect to. |
| cassandra.additionalServiceConfig | object | `{}` | Optional AdditionalServiceConfig allows to define additional parameters that are included in the created Services. Note, user can override values set by cass-operator and doing so could break cass-operator functionality. Avoid label "cass-operator" and anything that starts with "cassandra.datastax.com/" |
| cassandra.loggingSidecar.enabled | bool | `true` | Set to false if you do not want to deploy the server-system-logger container. |
| cassandra.loggingSidecar.image | object | `{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"k8ssandra/system-logger","tag":"6c64f9c4"}` | The server-system-logger container image |
| cassandra.loggingSidecar.image.registry | string | `"docker.io"` | Container registry fo the system logger |
| cassandra.loggingSidecar.image.repository | string | `"k8ssandra/system-logger"` | Repository for the system-logger image |
| cassandra.loggingSidecar.image.tag | string | `"6c64f9c4"` | Tag of the system-logger image to pull from image repository |
| cassandra.loggingSidecar.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the system-logger image |
| cassandra.heap | object | `{}` | Optional cluster-level heap configuration, can be overridden at `datacenters` level. Options are commented out for reference. Note that k8ssandra does not automatically apply default values for heap size. It instead defers to Cassandra's out of box defaults. |
| cassandra.gc | object | `{"cms":{},"g1":{}}` | Optional cluster-level garbage collection configuration. It can be overridden at the datacenter level. |
| cassandra.gc.cms | object | `{}` | GC configuration for the CMS collector. |
| cassandra.gc.g1 | object | `{}` | Controls the size of the two survivor spaces in the heap's young generation. survivorRatio: 8 -- The number of times an object survives a minor collection before being promoted to the old generation. maxTenuringThreshold: 1 -- A major collection starts if the occupancy of the old generation exceeds this percentage. initiatingOccupancyFraction: 75 -- The time in milliseconds that CMS threads wait for young GC. waitDuration: 10000 -- GC configuration for the G1 collector. |
| cassandra.resources | object | `{}` | Sets a target value for desired maximum pause time. maxGcPauseMillis: 500 -- Sets the heap occupancy threshold that triggers a marking cycle. initiatingHeapOccupancyPercent: 70 -- Set the number of stop the world (STW) worker threads. parallelGcThreads: 16 -- Set the number of stop the world (STW) worker threads. concurrentGcThreads: 16 -- Resource requests for each Cassandra pod. See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/ for background on managing resources. |
| cassandra.tolerations | list | `[]` |  |
| cassandra.datacenters[0].name | string | `"dc1"` | Name of the datacenter |
| cassandra.datacenters[0].size | int | `1` | Number of nodes within the datacenter. This value should, at a minimum, match the number of racks and be no less than 3 for non-development environments. |
| cassandra.datacenters[0].racks | list | `[{"affinityLabels":{},"name":"default"}]` | The replication factor for keyspaces in the datacenter. Triggers the even token distribution algorithm for num_tokens and the replication factor. Note that this property is for Cassandra 4.0 and later. Setting this property with Cassandra 3.11.x will result in a chart validation error. When the Cassandra version is 4.0, this property is enabled by default with a value of 3. allocateTokensForLocalRF: 3 -- Specifies the racks for the data center, if unset the datacenter will be composed of a single rack named `default`. The number of racks should equal the replication factor of your application keyspaces. Cassandra will ensure that replicas are spread across racks versus having multiple replicas within the same rack. For example, let's say we are using RF = 3 with a 9 node cluster and 3 racks (and 3 nodes per rack). There will be one replica of the dataset spread across each rack. |
| cassandra.datacenters[0].racks[0].name | string | `"default"` | Identifier for the rack, this may align with the labels used to control where resources are deployed for this rack. For example, if a rack is limited to a single availability zone the identifier may be the name of that AZ (eg us-east-1a). |
| cassandra.datacenters[0].racks[0].affinityLabels | object | `{}` | an optional set of labels that are used to pin Cassandra pods to specific k8s worker nodes via affinity rules. See https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes-using-node-affinity/ for background on using affinity rules. topology.kubernetes.io/zone is a well-known k8s label used by cloud providers to indicate the failure zone in which a k8s worker node is running. The following example illustrates how you can pin racks to specific failure zones. racks: - name: r1   affinityLabels:     topology.kubernetes.io/zone: us-east1-b - name: r2   affinityLabels:     topology.kubernetes.io/zone: us-east1-a - name: r3   affinityLabels:     topology.kubernetes.io/zone: us-east1-c |
| cassandra.datacenters[0].heap | object | `{}` | Optional datacenter-level heap setting, overrides cluster-level setting `cassandra.heap`. Options are commented out for reference. Note that k8ssandra does not automatically apply default values for heap size. It instead defers to Cassandra's out of box defaults. |
| cassandra.datacenters[0].gc | object | `{"cms":{},"g1":{}}` | Optional datacenter-level garbage collection configuration. |
| cassandra.datacenters[0].gc.cms | object | `{}` | Optional GC configuration for the CMS collector |
| cassandra.datacenters[0].gc.g1 | object | `{}` | Controls the size of the two survivor spaces in the heap's young generation. survivorRatio: 8 -- The number of times an object survives a minor collection before being promoted to the old generation. maxTenuringThreshold: 1 -- A major collection starts if the occupancy of the old generation exceeds this percentage. initiatingOccupancyFraction: 75 -- The time in milliseconds that CMS threads wait for young GC. waitDuration: 10000 -- Optional GC configuration for the G1 collector |
| cassandra.ingress | object | `{"enabled":false,"host":null,"method":"traefik","traefik":{"entrypoint":"cassandra"}}` | Sets a target value for desired maximum pause time. maxGcPauseMillis: 500 -- Sets the heap occupancy threshold that triggers a marking cycle. initiatingHeapOccupancyPercent: 70 -- Set the number of stop the world (STW) worker threads. parallelGcThreads: 16 -- Sets the number of parallel marking threads. concurrentGcThreads: 16 Cassandra native transport ingress support |
| cassandra.ingress.enabled | bool | `false` | Enables Cassandra Traefik ingress definitions. Note that this is mutually exclusive with stargate.ingress.cassandra.enabled |
| cassandra.ingress.method | string | `"traefik"` | Determines which TCP-based ingress custom resources to template out. Currently only `traefik` is supported |
| cassandra.ingress.host | string | `nil` | Optional hostname used to match requests. Warning: many native Cassandra clients, notably including cqlsh, initialize their connection by querying for the cluster's contactPoints, and thereafter communicate to the cluster using those names/IPs rather than whatever host was specified to the client. In order for clients to work correctly through ingress with a host filter, this means that the host filter must match the hostnames specified in the contactPoints. This value must be a DNS-resolvable hostname and not an IP address. To avoid this issue, leave this setting blank. |
| cassandra.ingress.traefik.entrypoint | string | `"cassandra"` | Traefik entrypoint where traffic is sourced. See https://doc.traefik.io/traefik/routing/entrypoints/ |
| stargate.enabled | bool | `true` | Enable Stargate resources as part of this release |
| stargate.version | string | `"1.0.29"` | version of Stargate to deploy. This is used in conjunction with cassandra.version to select the Stargate container image. If stargate.image is set, this value has no effect. |
| stargate.image | object | `{}` | Sets the Stargate container image. This value must be compatible with the value provided for stargate.clusterVersion. If left blank (recommended), k8ssandra will derive an appropriate image based on cassandra.clusterVersion. |
| stargate.replicas | int | `1` | Number of Stargate instances to deploy. This value may be scaled independently of Cassandra cluster nodes. Each instance handles API and coordination tasks for inbound queries. |
| stargate.waitForCassandra | object | `{"image":{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"alpine","tag":"3.12.2"}}` | The wait-for-cassandra init container in the Stargate Deployment |
| stargate.waitForCassandra.image.registry | string | `"docker.io"` | Image registry for the wait-for-cassandra container |
| stargate.waitForCassandra.image.repository | string | `"alpine"` | Image repository for the wait-for-cassandra container |
| stargate.waitForCassandra.image.tag | string | `"3.12.2"` | Tag of the image to pull from image.repository |
| stargate.waitForCassandra.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the wait-for-cassandra container |
| stargate.heapMB | int | `256` | The service account to use for Stargate pods. Defaults to the default account for the namespace. serviceAccount: -- Sets the heap size Stargate will use in megabytes. Memory request and limit for the pod will be set to this value x2 and x4, respectively. |
| stargate.cpuReqMillicores | int | `200` | Sets the CPU request for the Stargate pod in millicores. |
| stargate.cpuLimMillicores | int | `1000` | Sets the CPU limit for the Stargate pod in millicores. |
| stargate.livenessInitialDelaySeconds | int | `30` | Sets the initial delay in seconds for the Stargate liveness probe. |
| stargate.readinessInitialDelaySeconds | int | `30` | Sets the initial delay in seconds for the Stargate readiness probe. |
| stargate.cassandraUser | object | `{"secret":"","username":""}` | Configures the Cassandra user used by Stargate when authentication is enabled. If neither `cassandraUser.secret` nor `cassandraUser.username` are set, then a Cassandra user and a secret will be created. The username will be `stargate`. The secret name will be of the form `{clusterName}-stargate`. The password will be a random 20 character password. If `cassandraUser.secret` is set, then the Cassandra user will be created from the contents of the secret. If `cassandraUser.secret` is not set and if `cassandraUser.username` is set, a secret will be generated using the specified username. The password will be generated as previously described. |
| stargate.ingress.host | string | `nil` | Optional hostname used to match requests. Warning: many native Cassandra clients, notably including cqlsh, initialize their connection by querying for the cluster's contactPoints, and thereafter communicate to the cluster using those names/IPs rather than whatever host was specified to the client. In order for clients to work correctly through ingress with a host filter, this means that the host filter must match the hostnames specified in the contactPoints. This value must be a DNS-resolvable hostname and not an IP address. To avoid this issue, leave this setting blank, or override it to "" (empty string) for stargate.ingress.cassandra.host. This note does not apply to clients of Stargate's auth, REST, or GraphQL APIs. |
| stargate.ingress.enabled | bool | `false` | Enables all Stargate ingresses. Note: This must be true for any Stargate ingress to function. |
| stargate.ingress.auth.enabled | bool | `true` | Enables Stargate authentication ingress. Note: stargate.ingress.enabled must also be true. |
| stargate.ingress.auth.host | string | `nil` | Optional hostname used to match requests, overriding stargate.ingress.host if set |
| stargate.ingress.rest.enabled | bool | `true` | Enables Stargate REST ingress. Note: stargate.ingress.enabled must also be true. |
| stargate.ingress.rest.host | string | `nil` | Optional hostname used to match requests, overriding stargate.ingress.host if set |
| stargate.ingress.graphql.enabled | bool | `true` | Enables Stargate GraphQL API ingress. Note: stargate.ingress.enabled must also be true. |
| stargate.ingress.graphql.host | string | `nil` | Optional hostname used to match requests, overriding stargate.ingress.host if set |
| stargate.ingress.graphql.playground.enabled | bool | `true` | Enables GraphQL playground ingress.  Note: stargate.ingress.enabled and stargate.ingress.graphql.enabled must also be true. |
| stargate.ingress.cassandra.enabled | bool | `true` | Enables C* native protocol ingress with Traefik. Note that this is mutually exclusive with cassandra.ingress.enabled, and stargate.ingress.enabled must also be true. |
| stargate.ingress.cassandra.method | string | `"traefik"` | Determines which TCP-based ingress custom resources to template out. Currently only `traefik` is supported |
| stargate.ingress.cassandra.host | string | `nil` | Optional hostname used to match requests. Warning: many native Cassandra clients, notably including cqlsh, initialize their connection by querying for the cluster's contactPoints, and thereafter communicate to the cluster using those names/IPs rather than whatever host was specified to the client. In order for clients to work correctly through ingress with a host filter, this means that the host filter must match the hostnames specified in the contactPoints. This value must be a DNS-resolvable hostname and not an IP address. To avoid this issue, leave this setting blank, or if stargate.ingress.host is set, override it here to "" (empty string). This note does not apply to clients of Stargate's auth, REST, or GraphQL APIs. |
| stargate.ingress.cassandra.traefik | object | `{"entrypoint":"cassandra"}` | Parameters used by the Traefik IngressRoute custom resource |
| stargate.ingress.cassandra.traefik.entrypoint | string | `"cassandra"` | Traefik entrypoint where traffic is sourced. See https://doc.traefik.io/traefik/routing/entrypoints/ |
| stargate.affinity | object | `{}` | Affinity to apply to the Stargate pods. See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity for background |
| stargate.tolerations | list | `[]` | Tolerations to apply to the Stargate pods. See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/ for background. |
| reaper.securityContext | object | `{}` | Security context override for reaper container.
| reaper.schemaInitContainerConfig | object | `{}` | Security context override for reaper schema init container.
| reaper.configInitContainerConfig | object | `{}` | Security context override for reaper config init container.
| reaper.podSecurityContext | object | `{}` | Security context override for reaper pod |
| reaper.autoschedule | bool | `false` | When enabled, Reaper automatically sets up repair schedules for all non-system keypsaces. Repear monitors the cluster so that as keyspaces are added or removed repair schedules will be added or removed respectively. |
| reaper.autoschedule_properties | object | `{}` | Additional autoscheduling properties. Allows you to customize the schedule rules for autoscheduling. Properties are the same as accepted by the Reaper. |
| reaper.enabled | bool | `true` | Enable Reaper resources as part of this release. Note that Reaper uses Cassandra's JMX APIs to perform repairs. When Reaper is enabled, Cassandra will also be configured to allow remote JMX access. JMX authentication will be configured in Cassandra with credentials only created for Reaper in order to limit access. |
| reaper.image | object | `{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"thelastpickle/cassandra-reaper","tag":"2.3.1"}` | The name of the service account to use for Reaper pods. Defaults to the the default account. serviceAccount: Configures the Reaper container image to use. Exercise care when changing the Reaper image. Reaper is deployed and managed by reaper-operator. You will need to make sure that the image is compatible with reaper-operator. |
| reaper.image.registry | string | `"docker.io"` | Image registry for reaper |
| reaper.image.repository | string | `"thelastpickle/cassandra-reaper"` | Image repository for reaper |
| reaper.image.tag | string | `"2.3.1"` | Tag of the reaper image to pull from |
| reaper.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the reaper container |
| reaper.cassandraUser | object | `{"secret":"","username":""}` | Configures the Cassandra user used by Reaper when authentication is enabled. If neither cassandraUser.secret nor cassandraUser.username are set, then a Cassandra user and a secret with the user's credentials will be created. The username will be reaper. The secret name will be of the form {clusterName}-reaper. The password will be a random 20 character password. If cassandraUser.secret is set, then the Cassandra user will be created from the contents of the secret. If cassandraUser.secret is not set and if cassandraUser.username is set, a secret will be generated using the specified username. The password will be generated as previously described. |
| reaper.jmx | object | `{"secret":"","username":""}` | Configures JMX access to the Cassandra cluster. Reaper requires remote JMX access to perform repairs. The Cassandra cluster will be configured with remote JMX access enabled when Reaper is deployed. The JMX access will be configured to use authentication. If neither `jmx.secret` nor `jmx.username` are set, then a default user and secret with the user's credentials will be created. |
| reaper.jmx.username | string | `""` | Username that Reaper will use for JMX access. If left blank a random, alphanumeric string will be generated. |
| reaper.ingress.enabled | bool | `false` | Enables Reaper ingress definitions. When enabled, you must specify a value for reaper.ingress.host. |
| reaper.ingress.host | string | `nil` | Hostname to use for routing requests to the repair UI. If using a local deployment consider leveraging dynamic DNS services like xip.io. Example: `repair.127.0.0.1.xip.io` will return `127.0.0.1` for DNS requests routing requests to your local machine. This is required when reaper.ingress.enabled is true. |
| reaper.ingress.method | string | `"traefik"` |  |
| reaper.ingress.traefik.entrypoint | string | `"web"` | Traefik entrypoint where traffic is sourced. See https://doc.traefik.io/traefik/routing/entrypoints/ |
| reaper.affinity | object | `{}` | Affinity to apply to the Reaper pods. See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity for background |
| reaper.tolerations | list | `[]` | Tolerations to apply to the Reaper pods. See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/ for background. |
| medusa.enabled | bool | `false` | Enable Medusa resources as part of this release. If enabled, `bucketName` and `storageSecret` **must** be defined. |
| medusa.securityContext | object | `{}` | Security context override for Medusa container.
| medusa.restoreInitContainerConfig | object | `{}` | Security context override for Medusa restore init container.
| medusa.image.registry | string | `"docker.io"` | Image registry for medusa |
| medusa.image.repository | string | `"k8ssandra/medusa"` | Image repository for medusa |
| medusa.image.tag | string | `"0.11.1"` | Tag of the medusa image to pull from |
| medusa.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the medusa container |
| medusa.cassandraUser | object | `{"secret":"","username":""}` | Configures the Cassandra user used by Medusa when authentication is enabled. If neither `cassandraUser.secret` nor `cassandraUser.username` are set, then a Cassandra user and a secret will be created. The username will be medusa. The secret name will be of the form {clusterName}-medusa. The password will be a random 20 character password. If `cassandraUser.secret` is set, then the Cassandra user will be created from the contents of the secret. If `cassandraUser.secret` is not set and if `cassandraUser.username` is set, a secret will be generated using the specified username. The password will be generated as previously described. |
| medusa.multiTenant | bool | `false` | Enables usage of a bucket across multiple clusters. |
| medusa.storage | string | `"s3"` | API interface used by the object store. Supported values include `s3`, 's3_compatible', `google_storage` and 'azure_blobs'. For file system storage, i.e., a pod volume mount, use 'local' and set the podStorage properties. Note that 'local' does not necessarily imply a local volume. It could also be network attached storage. It is simply accessed through the file system. |
| medusa.storage_properties | object | `{}` | Optional properties for storage. Supported values depend on the type of the storage. |
| medusa.bucketName | string | `"awstest"` |  |
| medusa.storageSecret | string | `"medusa-bucket-key"` | Name of the Kubernetes `Secret` that stores the key file for the storage provider's API. If using 'local' storage, this value is ignored. |
| medusa.podStorage | object | `{}` | To use a locally mounted volumes for backups, the Cassandra pods must have a PVC where to write the backups to. |
| monitoring.grafana.provision_dashboards | bool | `true` | Enables the creation of configmaps containing Grafana dashboards. If leveraging the kube-prometheus-stack subchart this value should be `true`. See https://helm.sh/docs/chart_template_guide/subcharts_and_globals/ for background on subcharts. |
| monitoring.prometheus.provision_service_monitors | bool | `true` | Enables the creation of Prometheus Operator ServiceMonitor custom resources. If you are not using the kube-prometheus-stack subchart or do not have the ServiceMonitor CRD installed on your cluster, set this value to `false`. |
| monitoring.serviceMonitors.namespace | string | `nil` |  |
| cleaner | object | `{"image":{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"k8ssandra/k8ssandra-tools","tag":"latest"}}` | The cleaner is a pre-delete hook that that ensures objects with finalizers get deleted. For example, cass-operator sets a finalizer on the CassandraDatacenter. Kubernetes blocks deletion of an object until all of its finalizers are cleared. In the case of the CassandraDatacenter object, cass-operator removes the finalizer. The problem is that there are no ordering guarantees with helm uninstall which means that the cass-operator deployment could be deleted before the CassandraDatacenter. The cleaner ensures that the CassandraDatacenter is deleted before cass-operator. |
| cleaner.image | object | `{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"k8ssandra/k8ssandra-tools","tag":"latest"}` | Uncomment to specify the name of the service account to use for the cleaner. Defaults to <release-name>-cleaner-k8ssandra serviceAccount: |
| cleaner.image.registry | string | `"docker.io"` | Image registry for the cleaner |
| cleaner.image.repository | string | `"k8ssandra/k8ssandra-tools"` | Image repository for the cleaner |
| cleaner.image.tag | string | `"latest"` | Tag of the cleaner image to pull from |
| cleaner.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the cleaner container |
| client | object | `{"image":{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"k8ssandra/k8ssandra-tools","tag":"latest"}}` | k8ssandra-client provides CLI utilities, but also certain functions such as upgradecrds that allow modifying the running instances |
| client.image | object | `{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"k8ssandra/k8ssandra-tools","tag":"latest"}` | Uncomment to specify the name of the service account to use for the client tools image. Defaults to <release-name>-crd-upgrader-k8ssandra. serviceAccount: |
| client.image.registry | string | `"docker.io"` | Image registry for the client |
| client.image.repository | string | `"k8ssandra/k8ssandra-tools"` | Image repository for the client |
| client.image.tag | string | `"latest"` | Tag of the client image to pull from |
| client.image.pullPolicy | string | `"IfNotPresent"` | Pull policy for the client container |
| k8ssandra-operator.enabled | bool | `false` | Enables the k8ssandra-operator as part of this release. This is experimental deployment option, do not use in production. |
| cass-operator.enabled | bool | `true` | Enables the cass-operator as part of this release. If this setting is disabled no Cassandra resources will be deployed. |
| reaper-operator.enabled | bool | `true` | Enables the reaper-operator as part of this release. If this setting is disabled no repair resources will be deployed. |
| kube-prometheus-stack.enabled | bool | `true` | Controls whether the kube-prometheus-stack chart is used at all. Disabling this parameter prevents all monitoring components from being installed. |
| kube-prometheus-stack.coreDns.enabled | bool | `false` |  |
| kube-prometheus-stack.kubeApiServer.enabled | bool | `false` |  |
| kube-prometheus-stack.kubeControllerManager.enabled | bool | `false` |  |
| kube-prometheus-stack.kubeDns.enabled | bool | `false` |  |
| kube-prometheus-stack.kubeEtcd.enabled | bool | `false` |  |
| kube-prometheus-stack.kubeProxy.enabled | bool | `false` |  |
| kube-prometheus-stack.kubeScheduler.enabled | bool | `false` |  |
| kube-prometheus-stack.kubeStateMetrics.enabled | bool | `false` |  |
| kube-prometheus-stack.kubelet.enabled | bool | `false` |  |
| kube-prometheus-stack.nodeExporter.enabled | bool | `false` |  |
| kube-prometheus-stack.alertmanager.enabled | bool | `false` |  |
| kube-prometheus-stack.alertmanager.serviceMonitor.selfMonitor | bool | `false` |  |
| kube-prometheus-stack.prometheusOperator.enabled | bool | `true` |  |
| kube-prometheus-stack.prometheusOperator.namespaces | object | `{"additional":[],"releaseNamespace":true}` | Locks Prometheus operator to this namespace. Changing this setting may result in a non-namespace scoped deployment. |
| kube-prometheus-stack.prometheusOperator.serviceMonitor | object | `{"selfMonitor":false}` | Monitoring of prometheus operator |
| kube-prometheus-stack.prometheus.enabled | bool | `true` | Provisions an instance of Prometheus as part of this release |
| kube-prometheus-stack.prometheus.prometheusSpec | object | `{"externalUrl":"","routePrefix":"/"}` | Allows for tweaking of the Prometheus installation's configuration. Common parameters include `externalUrl: http://localhost:9090/prometheus` and `routePrefix: /prometheus` for running Prometheus resources under a specific path (`/prometheus` in this example). |
| kube-prometheus-stack.prometheus.prometheusSpec.routePrefix | string | `"/"` | Prefixes all Prometheus routes with the specified value. It is useful for ingresses which do not rewrite URLs. |
| kube-prometheus-stack.prometheus.prometheusSpec.externalUrl | string | `""` | An external URL at which Prometheus will be reachable. |
| kube-prometheus-stack.prometheus.ingress.enabled | bool | `false` | Enable templating of ingress resources for external prometheus traffic |
| kube-prometheus-stack.prometheus.ingress.paths | list | `[]` | Path-based routing rules, `/prometheus` is possible if the appropriate changes are made to `prometheusSpec` |
| kube-prometheus-stack.prometheus.serviceMonitor.selfMonitor | bool | `false` |  |
| kube-prometheus-stack.grafana.enabled | bool | `true` | Provisions an instance of Grafana and wires it up with a DataSource referencing this Prometheus installation |
| kube-prometheus-stack.grafana.ingress.enabled | bool | `false` | Generates ingress resources for the Grafana instance |
| kube-prometheus-stack.grafana.ingress.path | string | `nil` | Path-based routing rules, '/grafana' is possible if appropriate changes are made to `grafana.ini` |
| kube-prometheus-stack.grafana.adminUser | string | `"admin"` | Username for accessing the provisioned Grafana instance |
| kube-prometheus-stack.grafana.adminPassword | string | `"secret"` | Password for accessing the provisioned Grafana instance |
| kube-prometheus-stack.grafana.serviceMonitor.selfMonitor | bool | `false` | Whether the Grafana instance should be monitored |
| kube-prometheus-stack.grafana.defaultDashboardsEnabled | bool | `false` | Default dashboard installation |
| kube-prometheus-stack.grafana.plugins | list | `["grafana-polystat-panel"]` | Additional plugins to be installed during Grafana startup, `grafana-polystat-panel` is used by the default Cassandra dashboards. |
| kube-prometheus-stack.grafana."grafana.ini" | object | `{}` | Customization of the Grafana instance. To listen for Grafana traffic under a different url set `server.root_url: http://localhost:3000/grafana` and `serve_from_sub_path: true`. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.5.0](https://github.com/norwoodj/helm-docs/releases/v1.5.0)
