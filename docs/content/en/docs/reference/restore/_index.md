---
title: "Restore Helm Chart"
linkTitle: "restore"
weight: 4
description: >
  Handles the scheduling an execution of an ad-hoc restore.
---

```yaml
# The name of the CassandraRestore
name: restore

# Name of the backup to restore
backup:
  name: backup

# Name of the target CassandraDatacenter where the data should be restored
cassandraDatacenter:
  name: dc1

# An in-place restore will restore the backup to the source cluster. Note that this will
# trigger a rolling restart of the cluster.
inPlace: true
```
