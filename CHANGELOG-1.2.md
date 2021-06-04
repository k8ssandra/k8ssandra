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

## v1.2.0 - 2021-06-02

* [CHANGE] #726 Upgrade to Cassandra 4.0-rc1
* [FEATURE] #673, #698 Make tolerations configurable
* [ENHANCEMENT] #654 Reduce initial delay of Stargate readiness probe
* [ENHANCEMENT] #560 Add the ability to attach additional PVs for medusa backups
* [ENHANCEMENT] #693 Update cass-operator to v1.7.0
* [ENHANCEMENT] #732 Make allocate_tokens_for_replication_factor configurable
* [BUGFIX] #678 Upgrade to Medusa 0.10.1 fixing failed backups after a restore
* [FEATURE] #224 Allow defining additional settings for autoscheduling in Reaper, update reaper-operator to v0.3.1
