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

* [ENHANCEMENT] Added s3_rgw support
* [CHANGE] Upgrade cass-operator to v1.11.0
* [CHANGE] cass-operator HELM-Chart: Upgrade requires manual action if registryOverride was used before. registryOverride is now repositoryOverride.
* [ENHANCEMENT] Enable webhook functionality in cass-operator if cert-manager is installed.
