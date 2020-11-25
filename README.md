# K8ssandra
A distribution of Cassandra made for Kubernetes

# Overview
K8ssandra provides a full, open source stack for running and managing Cassandra in Kubernetes.

## Cassandra
K8ssandra packages and deploys [Apache Cassandra](https://cassandra.apache.org/).

## Monitoring
Monitoring includes the collection, the storage, and the visualization of metrics. With that in mind, K8ssandra integrates the following components.

### Metric Collector for Apache Cassandra (MCAC)
[MCAC](https://github.com/datastax/metric-collector-for-apache-cassandra) collects and aggregate Cassandra and OS-level metrics that can easily be stored in Prometheus.

### Prometheus
[Prometheus](https://prometheus.io/) is a very popular time-series, metrics database that is used extensively both inside of as well as outside of Kubernetes deployments.

### Grafana
[Grafana](https://grafana.com/) is the de facto standard for dashboards.

## Repairs
[Reaper](http://cassandra-reaper.io/) is used to schedule and manage repairs in Cassandra.

## Backup & Restore


# Implementation
K8ssandra is essentially an aggregation of several components that  together comprise the stack described above.

* [cass-operator](https://github.com/datastax/cass-operator)
* [prometheus-operator](https://github.com/prometheus-operator/prometheus-operator)
* [grafana-operator](https://github.com/integr8ly/grafana-operator)
* [reaper-operator](https://github.com/thelastpickle/reaper-operator)
* [helm](https://helm.sh)

Testing