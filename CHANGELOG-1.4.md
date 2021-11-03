# Changelog

Changelog for K8ssandra, new PRs should update the `unreleased` section below with entries sorted by type, in the 
following order:

```markdown
* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX]
```

If two entries have the same type, they should be sorted from newest to oldest (the newest comes first, the oldest comes 
last).

When cutting a new release of the parent `k8ssandra` chart update the `unreleased` heading to the tag being generated 
and date `## vX.Y.Z - YYYY-MM-DD` and create a new placeholder section for  `unreleased` entries.

## unreleased
* [CHANGE] Update to Reaper 3.0.0
* [CHANGE] #1118 Update to Stargate 1.0.40
* [CHANGE] #950 Update to cass-operator 1.8.0-rc.1
* [CHANGE] Update to cass-operator v1.8.0
* [ENHANCEMENT] #1179 Make `JAVA_OPTS` configurable for Stargate
* [ENHANCEMENT] Apply customizable filters on table level metrics in MCAC
* [ENHANCEMENT] #1140 securityContext defaults for operators and security foundations 
* [ENHANCEMENT] #1150 Bring reaper resources and CRDs up to date with main reaper-operator repo; operator-sdk 1.6.1/controller-runtime 0.9.2.
* [ENHANCEMENT] #1150 Update CRD versions to v1 from v1beta1 allowing compatibility with k8s 1.22.
* [ENHANCEMENT] #1083 Add support for full query logging (Cassandra 4.0.0 feature)
* [ENHANCEMENT] #1083 Add support for audit logging (Cassandra 4.0.0 feature)
* [ENHANCEMENT] #1083 Add support for client backpressure (Cassandra 4.0.0 feature)
* [ENHANCEMENT] #1083 Add support for deployment of Cassandra 4.0.1
* [ENHANCEMENT] #959 Root file system in Cassandra pod read only; security context for containers.
* [ENHANCEMENT] #874 expose cass-operator AdditionalServiceConfig in k8ssandra helm chart values.yaml
* [BUGFIX] #1181 Update k8ssandra-operator chart
* [BUGFIX] #1129 CassOperator kills C* pods with due to incorrect memory
* [BUGFIX] #1066 Azure backups are broken due to missing azure-cli deps
* [BUGFIX] #1012 reaper-operator's role.yaml has more data than it should, causing role name conflicts
* [BUGFIX] #1018 reaper image registry typo and jvm typo fixed
* [BUGFIX] #1029 Do not change num_tokens when upgrading
* [BUGFIX] Ensure Cassandra 4x is compatible with Stargate deployments by including `allow_alter_rf_during_range_movement` in config
