

![Version: 0.26.0](https://img.shields.io/badge/Version-0.26.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| K8ssandra Team | k8ssandra-developers@googlegroups.com | https://github.com/k8ssandra |

## Source Code

* <https://github.com/k8ssandra/medusa-operator>
* <https://github.com/k8ssandra/k8ssandra/tree/main/charts/backup>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| name | string | `"backup"` | Name of the CassandraBackup custom resource |
| cassandraDatacenter.name | string | `"dc1"` | Name of the CassandraDatacenter custom resource where the backup is sourced. |
