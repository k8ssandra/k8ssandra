---
title: "Access Monitoring"
linkTitle: "Access Monitoring"
weight: 1
date: 2020-11-11
description: 
---

K8ssandra provides preconfigured GitHub-hosted templates and build scripts for metrics reporter dashboards using [Prometheus](https://operatorhub.io/operator/prometheus) and [Grafana](https://operatorhub.io/operator/grafana-operator). The dashboards allow you to check the health of open-source Apache Cassandra® resources in your Kubernetes cluster.

## Tools

[Metrics dashboard for Cassandra in Kubernetes](https://github.com/datastax/metric-collector-for-apache-cassandra/tree/master/dashboards/k8s-build)

## Prerequisites

Use `git clone` to clone the [repo](https://github.com/datastax/metric-collector-for-apache-cassandra) for your environment and follow the steps in this topic.

## Steps

### Python scripts for dashboards and configurations

The dashboards plus Prometheus and Grafana configuration files are transformed via Python scripts under [bin](https://github.com/datastax/metric-collector-for-apache-cassandra/tree/master/dashboards/k8s-build/bin).

Run:

`bin/clean.py && bin/build.py`

The generated files will integrate with the Custom Resources defined by the Prometheus and Grafana operators that are available on Operator Hub.

*Note:* The Python-generated files are written to the `generated` directory.

### Prometheus Operator setup

The Prometheus Operator handles the orchestration, configuration, and deployment of Kubernetes resources required for a High Availability (HA) Prometheus installation. Rather than specifying a list of Cassandra nodes in a JSON file, this setup directs Prometheus to monitor a Kubernetes Service that exposes all nodes via DNS. This mapping of hosts is handled by a ServiceMonitor Custom Resource defined by the operator.

The following steps illustrate how to install the Prometheus Operator, deploy a service monitor pointed at a Cassandra or DSE cluster (with metric relabeling), and deploy a HA Prometheus cluster connected to the service monitor.

# Install OperatorHub Lifecycle Manager (OLM)
`curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.15.1/install.sh | bash -s 0.15.1`

# Install Prometheus Operator
`kubectl create -f dashboards/k8s-build/generated/prometheus/operator.yaml`

# Configure and install the Service Monitor
`kubectl apply -f dashboards/k8s-build/generated/prometheus/service_monitor.yaml`

# Configure and install the Prometheus deployment
The Prometheus Custom Resource maps the deployment to all service monitors with the label `cassandra.datastax.com/cluster: cluster-name`
```
serviceMonitorSelector:
    matchLabels:
      cassandra.datastax.com/cluster: cluster-name 
```

Edit the `cluster-name` for your environment, and apply the instance.yaml file to the cluster.
`kubectl apply -f dashboards/k8s-build/generated/prometheus/instance.yaml`

### Grafana Operator setup

# Install OperatorHub Lifecycle Manager (OLM)
`curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.15.1/install.sh | bash -s 0.15.1`

# Install Grafana Operator
`kubectl create -f dashboards/k8s-build/generated/grafana/operator.yaml`

# Configure and install the GrafanaDataSource
`kubectl apply -f dashboards/k8s-build/generated/grafana/datasource.yaml`

# Configure and install the GrafanaDashboard
Before installing the GrafanaDashboard , edit the YAML with appropriate labels. In this example, a label of `app=grafana` is used. If needed, modify for your environment.

With the configuration file updated, apply the resource to the cluster.
`kubectl apply -f dashboards/k8s-build/generated/grafana/`

# Configure and install the Grafana deployment
The following section in the Grafana Custom Resource maps the deployment to all dashboards with the label `app=grafana`.

```
  dashboardLabelSelector:
    - matchExpressions:
        - key: app
          operator: In
          values:
            - grafana
```


With the configuration file updated, apply the resource to the cluster.
``kubectl apply -f dashboards/k8s-build/generated/grafana/instance.yaml``

# Check the Grafana instance
Port forward to the grafana instance and check it out at http://127.0.0.1:3000/ (username: admin, password: secret)

**Note:** Never use documented credentials in production environments!

`kubectl port-forward svc/grafana-service 3000`

## Next

Access the [Repair Web interface](docs/topics/accessing-repair-interface/).
