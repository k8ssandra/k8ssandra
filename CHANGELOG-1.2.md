# Changelog

Changelog for K8ssandra, new PRs should update the ` unreleased` section with entries in the order:

```markdown
* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX]
```

When cutting a new release of the parent `k8ssandra` chart update the `unreleased` heading to the tag being generated and date `## vX.Y.Z - YYYY-MM-DD` and create a new placeholder section for  `unreleased` entries.

## unreleased

* [BUGFIX] #678 Upgrade to Medusa 0.10.1 fixing failed backups after a restore
* [CHANGE] #726 Upgrade to Cassandra 4.0-rc1
* [ENHANCEMENT] #560 Add the ability to attach additional PVs for medusa backups
* [ENHANCEMENT] #693 Update cass-operator to v1.7.0
