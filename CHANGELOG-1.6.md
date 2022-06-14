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

* [CHANGE] [#1409](https://github.com/k8ssandra/k8ssandra/pull/1409) Upgrade cass-operator to v1.11.0
* [ENHANCEMENT] [#1415](https://github.com/k8ssandra/k8ssandra/pull/1415) Added s3_rgw support
* [BUGFIX] [#1380](https://github.com/k8ssandra/k8ssandra/issues/1380) Enable webhook functionality in cass-operator if cert-manager is installed.
* [BUGFIX] [#1404](https://github.com/k8ssandra/k8ssandra/issues/1404) Fix `garbage_collector` property to use G1GC
* [ENHANCEMENT] [#1135](https://github.com/k8ssandra/k8ssandra/issues/1135)  Add support for service account annotations in helm chart.