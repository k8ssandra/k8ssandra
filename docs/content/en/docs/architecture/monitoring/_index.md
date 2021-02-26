---
title: "Monitoring"
linkTitle: "Monitoring"
weight: 2
description: > 
  How do you know if your cluster is healthy?
---

When running applications in Kubernetes observability is key. With K8ssandra and cass-operator, each Apache Cassandra® pod is configured with the DataStax [Metrics Collector for Apache Cassandra](https://github.com/datastax/metric-collector-for-apache-cassandra). This exposes Cassandra node-level metrics in the Prometheus format covering everything from operations per second and latency to compaction throughput and heap usage.

![Grafana Overview](grafana-overview.png)

K8ssandra goes a step further providing a deployment of Prometheus and Grafana for the storage and visualization of these metrics including:

* Prometheus `ServiceMonitor` resource with Metric Relabeling directed at the
  K8ssandra cluster's Service
* Grafana `DataSource` resource configured to reference the deployed Prometheus
  instance
* Grafana `Dashboard` resources

![Grafana Cluster](grafana-cluster.png)

## Metric Data Retention

Metric data collected in the cluster is retained within Prometheus for 24 hours.

## Next

Check out [Repairs with Reaper for Apache Cassandra]({{< ref "repairs" >}})
