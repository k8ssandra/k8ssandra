# Changelog

Changelog for K8ssandra, new PRs should update the `main / unreleased` section with entries in the order:

```markdown
* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX]
```

## v1.1.0 - 2021-04-09
* [CHANGE] #657 Make latest Stargate version available
* [BUGFIX] #616 Upgrade to Medusa 0.10.0 to fix scale up issues after a backup was restored
* [CHANGE] #533 Remove Jolokia integration
* [CHANGE] #630 Upgrade to medusa-operator 0.2.0
* [CHANGE] #613 Mount Cassandra pod labels in volume
* [CHANGE] #611 Shut down cluster by default with in-place restores
* [CHANGE] #637 Update Management API image locations
* [ENHANCEMENT] #576 Add option to disable Cassandra logging sidecar
* [ENHANCEMENT] #530 Upgrade Reaper to 2.2.2 and Medusa to 0.9.1
* [ENHANCEMENT] #510 Add docs and examples in values.yaml
* [ENHANCEMENT] #504 split dashboards into separate configmaps
* [ENHANCEMENT] #436 Upgrade Stargate to 1.0.11 and add a `preStop` lifecycle hook to improve behavior when reducing the number of Stargate replicas in the presence of live traffic
* [ENHANCEMENT] #419 Add automation for stable and next release streams
* [ENHANCEMENT] #239 Developer documentation
* [ENHANCEMENT] #547 Add support for additionalSeeds in the CassandraDatacenter
* [ENHANCEMENT] #606 Support installation of operators only, disabling the Cassandra cluster creation
* [BUGFIX] #475 Fix Cassandra config clobbering when enabling Medusa
* [BUGFIX] #396 cqlsh commands show warnings
* [BUGFIX] #516 Fix issue with scripts not being checked out before attempting to run them.
* [BUGFIX] #517 Removed GitHub Actions for pre-releasing off of main
* [BUGFIX] #475 Fix Cassandra config clobbering when enabling Medusa
* [BUGIFX] #590 Create cass-operator webhook secret
* [BUGFIX] #602 Fix indentation error in example backup-restore-values.yaml
* [BUGFIX] #623 `helm uninstall` can leave CassandraDatacenter behind