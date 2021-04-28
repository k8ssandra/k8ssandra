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

* [ENHANCEMENT] #560 Add the ability to attach additional PVs for medusa backups
* [BUGFIX] #678 Backup fails after a restore
