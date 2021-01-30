# K8ssandra
[K8ssandra](https://k8ssandra.io/) is a simple to manage, production-ready,
distribution of [Apache Cassandra](https://cassandra.apache.org/) and
[Stargate](https://stargate.io/) that is ready for 
[Kubernetes](https://kubernetes.io/). It is built on a foundation of rock-solid 
open-source projects covering both the transactional and operational aspects of
Cassandra deployments. This project is distributed as a collection of
[Helm](https://helm.sh/) charts. Feel free to fork the repo and contribute. If
you're looking to install K8ssandra head over to the [Getting
Started](https://k8ssandra.io/docs/getting-started/) guide.

## Components
K8ssandra is composed of a number of sub-charts each representing a component in
the K8ssandra stack. The default installation is focused on developer
deployments with all of the features enabled and configured for running with a
minimal set of resources. Many of these components may be deployed
independently in a centralized fashion. Below is a list of the components in the
K8ssandra stack with links to the appropriate projects.

### Apache Cassandra
K8ssandra packages and deploys [Apache Cassandra](https://cassandra.apache.org/)
via the [cass-operator](https://github.com/datastax/cass-operator) project. Each
Cassandra container has the [Management API for Apache Cassandra
(MAAC)](https://github.com/datastax/management-api-for-apache-cassandra) and
[Metrics Collector for Apache
Cassandra(MCAC)](https://github.com/datastax/metric-collector-for-apache-cassandra)
pre-installed and configured to come up automatically.

### Stargate
[Stargate](https://stargate.io/) provides a collection of horizontally scalable
API endpoints for interacting with Cassandra databases. Developers may leverage
REST and GraphQL alongside the traditional CQL interfaces. With Stargate
operations teams gain the ability to independently scale coordination (Stargate)
and data (Cassandra) layers. In some use-cases, this has resulted in a lower TCO and
smaller infrastructure footprint.

### Monitoring
Monitoring includes the collection, storage, and visualization of
metrics. Along with the previously mentioned MCAC, K8ssandra utilizes
[Prometheus](https://prometheus.io/) and [Grafana](https://grafana.com/) for the
storage and visualization of metrics. Installation and management of these
pieces is handled by the [Kube Prometheus
Stack](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack)
Helm chart.

### Repairs
The Last Pickle [Reaper](http://cassandra-reaper.io/) is used to schedule and
manage repairs in Cassandra. It provides a web interface to visualize repair
progress and manage activity.

### Backup & Restore

Another project from The Last Pickle,
[Medusa](https://github.com/thelastpickle/cassandra-medusa), manages the backup
and restore of K8ssandra clusters. 

## Next Steps

If you are looking to run K8ssandra in your [Kubernetes](https://kubernetes.io/) 
environment check out the
[getting started](https://k8ssandra.io/docs/getting-started/) guide. We are
always looking for contributions to the docs, helm charts, and underlying
components.
