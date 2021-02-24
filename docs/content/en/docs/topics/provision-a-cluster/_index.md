---
title: "Scale your Cassandra Cluster"
linkTitle: "Scale Cassandra"
weight: 4
description: Steps to provision a cluster in Kubernetes
---

## Tools

[helm](https://helm.sh/docs/intro/install/)

## Prerequisites

* A Kubernetes environment
* k8ssandra installed and running in Kubernetes - see [Getting Started]({{< ref "getting-started" >}})

## Steps

### Use helm to get the running configuration

For many basic configuration options, you may change values in the deployed YAML files. For example, you can scale up or scale down, as needed, by updated the YAML.

Let's check the currently running values. First let's get the list of installed charts that we installed in [Getting Started]({{< ref "getting-started" >}}):

`helm list`
```
NAME               	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART                  	APP VERSION
k8ssandra         	default  	1       	2020-11-11 17:05:20.010071 -0700 MST	deployed	k8ssandra-0.2.0        	3.11.7     
```

Now specify the name of the installed cluster to get the full manifest. Notice how helm displays the properties defined in each deployed YAML file. Example:

`helm get manifest k8ssandra`

```
---
# Source: k8ssandra/templates/reaper-operator/service_account.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: k8ssandra-a-reaper-operator-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
---
# Source: k8ssandra/templates/reaper-jmx-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: k8ssandra-a-reaper-secret-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
type: Opaque
data:
  username: "WlZralBpaTY5dQ=="
  password: "Y0FFM3FWZmt1UQ=="
---
# Source: k8ssandra/templates/reaper-operator/leader_election_role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8ssandra-a-reaper-operator-leader-election-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - configmaps/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
---
# Source: k8ssandra/templates/reaper-operator/role.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8ssandra-a-reaper-operator-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - services
  verbs:
  - create
  - get
  - list
  - watch
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - batch
  resources:
  - jobs
  verbs:
  - create
  - get
  - list
  - watch
- apiGroups:
  - cassandra.datastax.com
  resources:
  - cassandradatacenters
  verbs:
  - create
  - get
  - list
  - watch
- apiGroups:
  - reaper.cassandra-reaper.io
  resources:
  - reapers
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - reaper.cassandra-reaper.io
  resources:
  - reapers/status
  verbs:
  - get
  - patch
  - update
---
# Source: k8ssandra/templates/reaper-operator/leader_election_role_binding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8ssandra-a-reaper-operator-lead-election-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8ssandra-a-reaper-operator-leader-election-k8ssandra
subjects:
- kind: ServiceAccount
  name: k8ssandra-a-reaper-operator-k8ssandra
---
# Source: k8ssandra/templates/reaper-operator/role_binding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8ssandra-a-reaper-operator-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8ssandra-a-reaper-operator-k8ssandra
subjects:
  - kind: ServiceAccount
    name: k8ssandra-a-reaper-operator-k8ssandra
---
# Source: k8ssandra/templates/reaper-operator/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8ssandra-a-reaper-operator-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
spec:
  replicas: 1
  selector:
    matchLabels:
      name: k8ssandra-a-reaper-operator-k8ssandra
  template:
    metadata:
      labels:
        name: k8ssandra-a-reaper-operator-k8ssandra
    spec:
      serviceAccountName: k8ssandra-a-reaper-operator-k8ssandra
      containers:
        - args:
          - --enable-leader-election
          command:
            - /manager
          env:
          - name: WATCH_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          image: docker.io/thelastpickle/reaper-operator
          name: reaper-operator
          resources:
            limits:
              cpu: 100m
              memory: 30Mi
            requests:
              cpu: 100m
              memory: 20Mi
      terminationGracePeriodSeconds: 10
---
# Source: k8ssandra/templates/cassdc.yaml
# Sized to work on 3 k8s workers nodes with 1 core / 4 GB RAM
# See neighboring example-cassdc-full.yaml for docs for each parameter
apiVersion: cassandra.datastax.com/v1beta1
kind: CassandraDatacenter
metadata:
  name: dc1
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
  annotations:
    reaper.cassandra-reaper.io/instance: k8ssandra-a-reaper-k8ssandra
spec:
  clusterName: k8ssandra
  serverType: cassandra
  serverVersion: "3.11.7"
  managementApiAuth:
    insecure: {}
  size: 1
  storageConfig:
    cassandraDataVolumeClaimSpec:
      storageClassName: standard
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 5Gi
  config:    
    jvm-options:
      initial_heap_size: "800M"
      max_heap_size: "800M"
  podTemplateSpec:
    spec:
      initContainers:
        - name: jmx-credentials
          image: busybox
          imagePullPolicy: IfNotPresent
          env:
            - name: JMX_USERNAME
              valueFrom:
                secretKeyRef:
                  name: k8ssandra-a-reaper-secret-k8ssandra
                  key: username
            - name: JMX_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: k8ssandra-a-reaper-secret-k8ssandra
                  key: password
          args:
            - /bin/sh
            - -c
            - echo -n "$JMX_USERNAME $JMX_PASSWORD" > /config/jmxremote.password
          volumeMounts:
            - mountPath: /config
              name: server-config
      containers:
        - name: cassandra
          env:
            - name: LOCAL_JMX
              value: "no"
---
# Source: k8ssandra/templates/reaper.yaml
apiVersion: reaper.cassandra-reaper.io/v1alpha1
kind: Reaper
metadata:
  name: k8ssandra-a-reaper-k8ssandra
  labels:
    helm.sh/chart: k8ssandra-0.2.0
    app.kubernetes.io/name: k8ssandra
    app.kubernetes.io/instance: k8ssandra-a
    app.kubernetes.io/version: "3.11.7"
    app.kubernetes.io/managed-by: Helm
spec:
  image: thelastpickle/cassandra-reaper:2.0.5
  serverConfig:
    storageType: cassandra
    jmxUserSecretName: k8ssandra-a-reaper-secret-k8ssandra
    cassandraBackend:
      clusterName: k8ssandra
      replication:
        networkTopologyStrategy:
          dc1: 1
      # This is a bit of a hack. We really should not be specifying the service name here as it is
      # implementation detail of cass-operator. reaper-operator needs to be updated to simply take
      # the name of the CassandraDatacenter here.
      cassandraService: k8ssandra-dc1-service
```


### Scale up the cluster

Use the following command to find the `size` property:

`helm get manifest k8ssandra-a | grep size`

In this example, it returns:

```
  size: 1
      initial_heap_size: "800M"
      max_heap_size: "800M"
```

Notice the value of `size: 1` in cassdc.yaml. This is the Cassandra DataCenter definition. 

To scale up, you could change the `size` to 3. Example with helm:

`helm upgrade k8ssandra-a k8ssandra/k8ssandra --set size=3 --reuse-values`

Note: using `--reuse-values` to ensure keeping settings from previous `helm upgrade`.

```
Release "k8ssandra-a" has been upgraded. Happy Helming!
NAME: k8ssandra-a
LAST DEPLOYED: Thu Nov 12 07:13:33 2020
NAMESPACE: default
STATUS: deployed
REVISION: 2
TEST SUITE: None
```

Verify the upgrade:

`helm get manifest k8ssandra-a | grep size`           

```
size: 3
      initial_heap_size: "800M"
      max_heap_size: "800M"
```

### Scale down the cluster

Similarly, to scale down, lower the current `size` to conserve cloud resource costs, if the new value can support your computing requirements in Kubernetes.  Example:

`helm upgrade k8ssandra-a k8ssandra/k8ssandra --set size=1 --reuse-values`
```
Release "k8ssandra-a" has been upgraded. Happy Helming!
NAME: k8ssandra-a
LAST DEPLOYED: Thu Nov 12 07:18:15 2020
NAMESPACE: default
STATUS: deployed
REVISION: 4
TEST SUITE: None
```

Again, verify the upgrade:

`helm get manifest k8ssandra-a | grep size`
```
  size: 1
      initial_heap_size: "800M"
      max_heap_size: "800M"
```

## Next

Use Medusa to [backup and restore]({{< ref "/docs/topics/restore-a-backup/" >}}) data from/to a Cassandra database. 
