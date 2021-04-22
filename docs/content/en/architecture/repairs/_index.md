---
title: "K8ssandra Architecture - Repairs"
linkTitle: "Repairs"
description: "Maintaining anti-entropy for your data, which is necessary for partition-tolerant distributed systems like your Kubernetes-managed Cassandra database."
---

Apache Cassandra® works tirelessly to provide consistent results for queries. There are a number of anti-entropy mechanisms running constantly to keep data in sync across replicas. Repair is one of these mechanisms. It is recommended that a complete repair cycle is run across the entire dataset once every ten days. In order to reduce the impact of analyzing the entire dataset at once many operators choose to spread out the repair process over the ten day period.

To that end, K8ssandra leverages the excellent [Reaper for Apache Cassandra](http://cassandra-reaper.io/)(Reaper) project from The Last Pickle to handle the scheduling, execution, and monitoring of repair tasks. Optionally ingress may be configured as part of the K8ssandra installation for external connectivity to the Reaper web interface.

![Reaper UI](reaper-ui.png)

## Next

Explore [Backups with Medusa for Apache Cassandra]({{< relref "backups" >}})
