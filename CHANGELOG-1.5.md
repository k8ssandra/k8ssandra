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
* [CHANGE] Deprecate all K8ssandra v1.x charts, refer users to k8ssandra-operator.

* [CHANGE] Update cass-operator to v1.10.5 as well as system-logger to v1.13.1 and cass-config-builder to 1.0.5

* [CHANGE] Upgrade Grafana to v7.5.17
* [CHANGE] Upgrade Reaper to v3.2.1
* [CHANGE] Upgrade cass-operator to v1.10.5
* [CHANGE] Upgrade Stargate to v1.0.68

## v1.5.1 - 2022-07-08

* [CHANGE] Upgrade cass-operator to v1.10.4

## v1.5.0 - 2022-04-09

* [FEATURE] Enable the use of ZGC (Z Garbage Collector) for Cassandra 4.0
* [CHANGE] Upgrade cass-operator to v1.10.3
* [CHANGE] Upgrade Stargate to v1.0.52
* [CHANGE] Upgrade Medusa to v0.12.2
* [CHANGE] Upgrade Reaper to v3.1.1
* [CHANGE] Upgrade Management API to v0.1.37
* [BUGFIX] Fix GC subsettings mappings for Cassandra 4.0
