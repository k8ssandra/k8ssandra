---
title: "Reaper for Apache Cassandra repair operations"
linkTitle: "Reaper"
weight: 3
description: K8ssandra deploys Reaper to enable anti-entropy repair operations on Apache Cassandra&reg; data.
---

Reaper for Apache Cassandra&reg; is deployed by a K8ssandra install, which includes the Reaper Operator [Helm chart]({{< relref "/reference/reaper-operator/" >}}).

If you haven't already installed K8ssandra, see the [quickstarts]({{< relref "/quickstarts/" >}}) and [install]({{< relref "/install" >}}) topics.

## Introduction

Reaper is an open source tool to schedule and orchestrate repairs of Apache Cassandra® clusters. Reaper was originally designed and open-sourced by [Spotify](https://www.slideshare.net/planetcassandra/spotify-automating-cassandra-repairs) in an attempt to automate repairs while applying best practices from their solid production experience.

Reaper enable you to maintain anti-entropy for your data, which is necessary for partition-tolerant distributed systems like your Kubernetes-managed Cassandra database.

Apache Cassandra works constantly to provide consistent results for queries. There are a number of anti-entropy mechanisms continuously running to keep data in sync across replicas. Repair is one of these mechanisms. We recommended running a complete repair cycle across the entire dataset **once every ten days**. In order to reduce the impact of analyzing the entire dataset at once, many operators choose to spread out the repair process over the ten day period.

To that end, K8ssandra leverages [Reaper for Apache Cassandra](http://cassandra-reaper.io/) from The Last Pickle to handle the scheduling, execution, and monitoring of repair tasks. Optionally, ingress may be configured as part of the K8ssandra installation for external connectivity to the Reaper web interface.

Here's an example from the Reaper Web UI. For more, see [Repair tasks]({{< relref "/tasks/repair/" >}}). 

![Reaper Web UI - Cluster dialog](reaper-ui.png)

## Repair challenges before Reaper

Anti-entropy repair has been performed traditionally using the `nodetool repair` command. It can be performed in two ways, full or incremental, and be configured to repair various ranges of tokens: 

* all
* primary range
* sub-range

Add to this task various anti-compaction triggers and the different validation compaction orchestration settings:

* sequential
* parallel
* datacenter aware

All of which often resulted in complexity for a mandatory repair operation that should be simple to run. 

In the 1.x/2.x days of Cassandra (and probably after that), some operators simply gave up on repairing their clusters due to the difficulties in completing the operation successfully, without impacting SLAs.

The main problems that were encountered then during repairs:

* High number of pending compactions and SSTables on disk
* Repairs taking longer than the tombstones GC grace period
* High cluster load due to repair pressure
* Blocked/never-ending repairs
* Inability to resume repair operations if failures occurred
* vnodes made the operation very long and challenging to perform

## Reaper performs safe repairs

Reaper was built to address those issues and make repairs as safe and reliable as possible. It splits the repair operations into evenly sized subranges, and schedules the operations, so that:

* All nodes are kept busy repairing small units of data if possible
* A single segment is running on a node at once
* Segments lasting too long are terminated and re-scheduled
* Failed segments get replayed in case of transient failure
* Pending compactions will be monitored to pause segment scheduling, preventing overload
* Repairs can be paused

Reaper also supports incremental repair - recommended for use starting with Cassandra 4.0. Since Cassandra 3.0, Reaper can create segments with several token ranges to reduce the overhead of vnodes on repairs. Such ranges will be repaired in a single job by Cassandra as segments will only contain ranges that are replicated on the same set of nodes.

## Reaper features

K8ssandra deploys the Reaper web UI. You can access it here, specifying `$REAPER_HOST` with the configured DNS name in your environment:

http://$REAPER_HOST:8080/webui/index.html

The web interface lets you:

* Add/remove clusters
* Manage repair schedules
* Run manual repairs and manage running repairs

Reaper collects and displays runtime Cassandra metrics, running compactions and ongoing streaming sessions.

Reaper comes with a scheduler for recurring repairs but can also perform on-demand one-off repairs.

Reaper also has the ability to listen and display live Cassandra’s emitted Diagnostic Events.

In Cassandra 4.0, internal system “diagnostic events” have become available, via the work done in CASSANDRA-12944. These allow operators to observe internal Cassandra events, for example in unit tests and with external tools. These diagnostic events provide operational monitoring and troubleshooting beyond logs and metrics.

Reaper can use Postgres and Cassandra itself as a storage backend for its data, and is capable of repairing all Cassandra versions since 1.2 up to the latest 4.0.

In order to make Reaper more efficient, segment orchestration was recently revamped and modernized. It opened for a long awaited feature: fully concurrent repairs for different keyspaces and tables.
These changes also introduced a long awaited feature by allowing fully concurrent repairs for different keyspaces/tables.

## Using the Reaper Web UI

For the steps to set up repair operations using the Reaper Web UI, see [Repair tasks]({{< relref "/tasks/repair/" >}}). 

For information about how secrets are created and used with Reaper authentication, see [K8ssandra security]({{< relref "/tasks/secure#reaper-security/" >}}).

For reference details, see the Reaper Operator [Helm chart]({{< relref "/reference/reaper-operator/" >}}).

## Next

See the other [components]({{< relref "/components/" >}}) deployed by K8ssandra. For information on using the deployed components, see the [Tasks]({{< relref "/tasks/" >}}) topics.
