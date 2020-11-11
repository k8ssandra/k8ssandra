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

## Next
