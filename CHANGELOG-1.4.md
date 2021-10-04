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
* [ENHANCEMENT] #1083 Add support for deployment of Cassandra 4.0.1
* [ENHANCEMENT] #959 Root file system in Cassandra pod read only; security context for containers.
* [BUGFIX] #969 Prometheus "out-of-order timestamp" error due to metrics relabeling conflict 
* [BUGFIX] #1066 Azure backups are broken due to missing azure-cli deps
* [BUGFIX] #1012 reaper-operator's role.yaml has more data than it should, causing role name conflicts
* [BUGFIX] #1018 reaper image registry typo and jvm typo fixed
* [BUGFIX] #1029 Do not change num_tokens when upgrading
* [ENHANCEMENT] #874 expose cass-operator AdditionalServiceConfig in k8ssandra helm chart values.yaml
* [CHANGE] #1118 Update to Stargate 1.0.35
* [CHANGE] #950 Update to cass-operator 1.8.0-rc.1
* [CHANGE] Update to cass-operator v1.8.0