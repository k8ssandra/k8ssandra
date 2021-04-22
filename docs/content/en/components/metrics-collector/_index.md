---
title: "Metrics Collector"
linkTitle: "Metrics Collector"
weight: 5
description: Metrics Collector for Apache Cassandra&reg;
---

## Introduction

Metrics Collector for Apache Cassandra (MCAC) is deployed to your Kubernetes environment via the K8ssandra Helm chart installation. If you haven't already, see the [quickstart]({{< relref "/quickstarts/" >}}) and [install]({{< relref "/install" >}}) topics.

MCAC aggregates OS and Cassandra metrics along with diagnostic events to facilitate problem resolution and remediation. MCAC supports existing Apache Cassandra clusters and is a self contained drop in agent.  K8ssandra provides preconfigured Grafana dashboards to visualize the collected metrics. 

* Built on [collectd](https://collectd.org), a popular, well-supported, open source metric collection agent.
   With over 90 plugins, you can tailor the solution to collect metrics most important to you and ship them to
   wherever you need.

* Easily added to Cassandra nodes as a Java agent (via the K8ssandra deployment), Apache Cassandra sends metrics and other structured events to collectd over a local Unix socket.  

* Fast and efficient.  It can track over 100k unique metric series per node (that is, hundreds of Cassandra tables).

* Comes with extensive dashboards out of the box, built on [Prometheus](http://prometheus.io) and [Rrafana](http://grafana.com). The Cassandra dashboards let you aggregate latency accurately across all nodes, dc or rack, down to an individual table.   

Sample Grafana displays of collected metrics: 

### Sample of overview metrics in Grafana

![Overview metrics displayed in Grafana](overview.png) 

### Sample of OS metrics in Grafana

![OS metrics displayed in Grafana](os.png)

### Sample of cluster metrics in Grafana

![Cluster metrics ](cluster.png)

## Design principles

* Little or no performance impact to Cassandra
* Simple to install via K8ssandra, and self managed
* Collect all OS and Cassandra metrics by default
* Keep historical metrics on node for analysis
* Provide useful integration with Prometheus and Grafana

## Cassandra version supported:

The supported versions of Apache Cassandra: 2.2+ (2.2.X, 3.0.X, 3.11.X, 4.0)

## FAQs

  1. Where is the list of all Cassandra metrics?

     The full list is located on [Apache Cassandra docs](https://cassandra.apache.org/doc/latest/operating/metrics.html) site.
     The names are automatically changed from CamelCase to snake_case.

     In the case of Prometheus the metrics are further renamed based on [relabel config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config) which live in the
     [prometheus.yaml](https://github.com/datastax/metric-collector-for-apache-cassandra/blob/master/dashboards/prometheus/prometheus.yaml) file in the MCAC repo.

  2. How can I filter out metrics I don't care about?

     Please read the [metric-collector.yaml](https://github.com/datastax/metric-collector-for-apache-cassandra/blob/master/config/metric-collector.yaml) section in the MCAC GitHub repo on how to add filtering rules.

  3. What is the datalog? And what is it for?

     The datalog is a space limited JSON based structured log of metrics and events which are optionally kept on each node.  
     It can be useful to diagnose issues that come up with your cluster.  If you wish to use the logs yourself, 
     there's a [script](https://github.com/datastax/metric-collector-for-apache-cassandra/blob/master/scripts/datalog-parser.py) included on the MCAC GitHub repo to parse these logs which can be analyzed or piped into [jq](https://stedolan.github.io/jq/).

     Alternatively, we offer free support for issues, and these logs can help our support engineers help diagnose your problem.


## Next

See the other [components]({{< relref "/components/" >}}) deployed by K8ssandra. For information on using the deployed components, see the [Tasks]({{< relref "/tasks/" >}}) topics.
