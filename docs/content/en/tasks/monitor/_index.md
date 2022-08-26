---
title: "Monitor Cassandra"
linkTitle: "Monitor"
weight: 6
description: "Access tools to monitor your Apache Cassandra® cluster running in Kubernetes."
---


# Monitoring using Prometheus

While K8ssandra v1 managed the deployment of the kube-prometheus stack, that ability was removed in k8ssandra-operator. Both Prometheus and Grafana installations are to be handled separately.
The following guide will show you how to install Prometheus and Grafana on your Kubernetes cluster using the prometheus-operator and a Grafana custom deployment, but this can be achieved using the kube-prometheus stack as well.

## Installing and configuring Prometheus for monitoring

`k8ssandra-operator` has integrations with Prometheus which allow for the simple rollout of Prometheus ServiceMonitors for both Stargate, Cassandra Datacenters and Reaper.
ServiceMonitors are custom resources of [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) which describe the set of targets to be monitored by Prometheus.

### Prerequisites

The following guide assumes k8ssandra-operator is already installed, and a K8ssandraCluster object was created with the following manifest, in the `k8ssandra-operator` namespace:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
  namespace: k8ssandra-operator
spec:
  cassandra:
    serverVersion: "4.0.3"
    serverImage: k8ssandra/cass-management-api:4.0.3
    storageConfig:
      cassandraDataVolumeClaimSpec:
        storageClassName: standard
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 1Gi
    config:
      jvmOptions:
        heapSize: 512M
    datacenters:
      - metadata:
          name: dc1
        size: 3
    mgmtAPIHeap: 64Mi 
  stargate:
    size: 1
  reaper:
    keyspace: reaper_db
```
*Download this manifest [here](k8ssandra.yaml).*

Wait for the pods to come up in the `k8ssandra-operator` namespace and fully start.

To use Prometheus for monitoring, you need to have the prometheus-operator installed on your Kubernetes (k8s) cluster.
The prometheus-operator installs the ServiceMonitor CRD, which is the integration point we use to tell Prometheus how to find the Stargate and Cassandra pods and what endpoints on those pods to scrape.

### Install prometheus-operator

*Skip to the next section if you already have the prometheus-operator installed*  
  
Install the Prometheus operator by running the following command:

```
kubectl create -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/master/bundle.yaml
```

This will install all the CRDs along with a `prometheus` service account (SA), a cluster role and a cluster role binding to the `prometheus` SA.
It will also create a `prometheus-operator` deployment in the `default` namespace.

### Deploy the ServiceMonitor resources for the K8ssandraCluster

Now that the ServiceMonitor CRD exists in the cluster, we can deploy them for Cassandra, Stargate and Reaper by updating the `.spec.cassandra`, `.spec.stargate` and `.spec.reaper` fields of the K8ssandraCluster CR with:

```yaml
    telemetry:
      prometheus:
        enabled: true
```

After apply the patch, running `kubectl get servicemonitor -n k8ssandra-operator` should return three ServiceMonitor resources.  
You can selectively enable service monitor creation for each component without any requirement to enable them all.  

*Note: Reaper's telemetry block was added in K8ssandra v1.2.0 and Reaper v3.2.0.*

### Create/Update a Prometheus deployment

Create the `prometheus` namespace: `kubectl create namespace prometheus`
  
If the Prometheus deployment is to be created in another namespace than the one containing the K8ssandra pods, make sure it has the ability to monitor the pods in the K8ssandra namespace. You may also have multiple CassandraDatacenter resources existing in multiple namespaces, with a single Prometheus instance monitoring all of them.

This can be done by adding a label to the K8ssandra namespaces such as: `kubectl label namespace/k8ssandra-operator "k8ssandra.io/monitor=true"`.

We also need a service account for our prometheus instance with a cluster role and cluster role binding that allow it to access the necessary resources in all namespaces. You'll find all the necessary resources in this [manifest](prometheus-rbac.yaml).

Then create a Prometheus deployment in the namespace of your choice, including a service monitor namespace selector matching our label, and a service to access our Prometheus instances. We will use the `prometheus` namespace here:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
  namespace: prometheus
spec:
  evaluationInterval: 30s
  image: quay.io/prometheus/prometheus:v2.22.1
  nodeSelector:
    kubernetes.io/os: linux
  replicas: 1
  resources:
    requests:
      memory: 400Mi
  scrapeInterval: 30s
  securityContext:
    fsGroup: 2000
    runAsNonRoot: true
    runAsUser: 1000
  serviceAccountName: prometheus
  serviceMonitorNamespaceSelector:
    matchLabels:
      k8ssandra.io/monitor: "true"
  serviceMonitorSelector: {}
  version: v2.22.1
---
apiVersion: v1
kind: Service
metadata:
  name: prometheus
  labels:
    app: prometheus
  namespace: prometheus
spec:
  ports:
  - name: web
    port: 9090
    targetPort: web
  selector:
    app.kubernetes.io/name: prometheus
  sessionAffinity: ClientIP
```
*Download this manifest [here](prometheus.yaml).*

  
If you already have a Prometheus deployment in your cluster, make sure the `serviceMonitorNamespaceSelector` field is set to match the label we added to the K8ssandra namespaces, or set to `{}` to match all namespaces.
  
After Prometheus starts, you can port-forward it with `kubectl port-forward svc/prometheus 9090`. Go to [http://localhost:9090/service-discovery](http://localhost:9090/service-discovery) and check if the service monitors are detected, and that there are active targets. In this case, we should have three active targets for Cassandra and one active target for both Stargate and Reaper.

### Filtering metrics

Cassandra provides a lot of metrics which can create some overload, especially when there are many tables in a cluster. [Filtering rules for MCAC](https://github.com/datastax/metric-collector-for-apache-cassandra/blob/master/config/metric-collector.yaml#L9-L72) can be defined in the telemetry spec:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  cassandra:
    telemetry: 
      prometheus:
        enabled: true
        mcacMetricFilters:
          - "deny:org.apache.cassandra.metrics.Table"
          - "allow:org.apache.cassandra.metrics.Table.LiveSSTableCount"
```

When no filter is explicitly defined in the spec, default K8ssandra v1.x filters will be applied:

```
 - "deny:org.apache.cassandra.metrics.Table"
 - "deny:org.apache.cassandra.metrics.table"
 - "allow:org.apache.cassandra.metrics.table.live_ss_table_count"
 - "allow:org.apache.cassandra.metrics.Table.LiveSSTableCount"
 - "allow:org.apache.cassandra.metrics.table.live_disk_space_used"
 - "allow:org.apache.cassandra.metrics.table.LiveDiskSpaceUsed"
 - "allow:org.apache.cassandra.metrics.Table.Pending"
 - "allow:org.apache.cassandra.metrics.Table.Memtable"
 - "allow:org.apache.cassandra.metrics.Table.Compaction"
 - "allow:org.apache.cassandra.metrics.table.read"
 - "allow:org.apache.cassandra.metrics.table.write"
 - "allow:org.apache.cassandra.metrics.table.range"
 - "allow:org.apache.cassandra.metrics.table.coordinator"
 - "allow:org.apache.cassandra.metrics.table.dropped_mutations"
```

## Installing and configuring Grafana for monitoring

The K8ssandra project provides Grafana dashboards for monitoring Cassandra and Stargate.

### Create the Prometheus Grafana datasource

Create the datasource using a configMap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-datasource
  namespace: prometheus
data:
  prometheus.yaml: |
    apiVersion: 1
    datasources:
    - name: Prometheus
      type: prometheus
      url: http://prometheus.prometheus.svc.cluster.local:9090
      access: proxy
      isDefault: false
      jsonData:
        timeInterval: 30s
```
*Download this manifest [here](prometheus-datasource.yaml).*

This datasource will point to our `prometheus` service in the `prometheus` namespace.

### Create the Grafana dashboards

The Grafana dashboards and dashboard providers will be created as configMaps as well. The following manifests can be applied in the `prometheus` namespace to create the required objects: [grafana-config-maps.yaml](grafana-config-maps.yaml).

### Create the Grafana deployment

All the required objects are created in the `prometheus` namespace, and we can now proceed with creating the Grafana deployment and the corresponding service:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: grafana
  name: grafana
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      containers:
      - env:
        - name: GF_INSTALL_PLUGINS
          value: grafana-polystat-panel
        - name: GF_SECURITY_ADMIN_PASSWORD
          value: admin123
        image: grafana/grafana:7.5.15
        name: grafana
        ports:
        - containerPort: 3000
        volumeMounts:
        - mountPath: /var/lib/grafana
          name: grafana-storage
        - mountPath: /etc/grafana/provisioning/datasources
          name: grafana-datasources
        - name: grafana-dashboard-providers
          mountPath: /etc/grafana/provisioning/dashboards
        - name: grafana-dashboard-stargate
          mountPath: /var/lib/grafana/dashboards-stargate
        - name: grafana-overview-dashboard
          mountPath: /var/lib/grafana/dashboards-overview
        - name: grafana-condensed-dashboard
          mountPath: /var/lib/grafana/dashboards-condensed
      volumes:
        - name: grafana-storage
          emptyDir: {}
        - name: grafana-datasources
          configMap:
            name: prometheus-datasource
        - name: grafana-dashboard-providers
          configMap:
            name: grafana-dashboard-providers
        - name: grafana-dashboard-stargate
          configMap:
            name: k8ssandra-stargate-dashboard
        - name: grafana-overview-dashboard
          configMap:
            name: k8ssandra-cassandra-overview-dashboard
        - name: grafana-condensed-dashboard
          configMap:
            name: k8ssandra-cassandra-condensed-dashboard
      dnsPolicy: ClusterFirst
      restartPolicy: Always
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 100%
      maxSurge: 25%
---
apiVersion: v1
kind: Service
metadata:
  name: grafana-service
spec:
  ports:
    - name: native
      protocol: TCP
      port: 80
      targetPort: 3000
  selector:
    app: grafana
  type: ClusterIP
  sessionAffinity: None
```
*Download this manifest [here](grafana.yaml).*

The above manifest will set the Grafana credentials to `admin/admin123`. Please adjust these settings for production deployments.
It will also install the `grafana-polystat-panel` plugin which is used in the overview dashboard, and mount all the required configMaps.

You can port-forward the Grafana service to access the dashboard at [http://localhost:3000](http://localhost:3000): `kubectl port-forward svc/grafana-service 3000:3000`

You should then see the following list of available dashboards:
![Dashboard list](grafana-dashboard-list.png)

Clicking on the Overview Dashboard should get you to the following screen:
![Overview Dashboard](grafana-overview-dashboard.png)

## Next steps

* Explore other K8ssandra Operator [tasks]({{< relref "/tasks" >}}).
* See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Operator Custom Resource Definitions (CRDs) and the single K8ssandra Operator Helm chart. 

