---
title: "K8ssandra Architecture - Backups"
linkTitle: "Backups"
description: "Backing up your data helps you prepare for when the unthinkable happens."
---

Even with the heightened availability of Apache Cassandra® a proper backup schedule and testing of restore procedures is good practice in case catastrophe strikes. With distributed systems backups can be tricky, there's the timing of the snapshot process on all nodes, correlation of data files to remote storage, and eventual restore.

K8ssandra provides Helm charts for taking backups or triggering the restoration of data. This is accomplished via the [Medusa for Apache Cassandra](https://github.com/thelastpickle/cassandra-medusa) project from The Last Pickle and Spotify.

## Next

Dig into accessing your information through [Data APIs via Stargate]({{< relref "/architecture/stargate" >}}).
