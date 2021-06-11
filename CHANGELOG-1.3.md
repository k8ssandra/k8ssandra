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

* [CHANGE] #812 Integrate Fossa component/license scanning
* [BUGFIX] #853 Fix property name in scaling docs
* [BUGFIX] #412 Stargate metrics don't show up in the dashboards
