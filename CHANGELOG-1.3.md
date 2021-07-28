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

## v1.3.0 - 2021-07-27

* [FEATURE] #949 Upgrade from Management API v0.1.26 to v0.1.27 to provide support for Cassandra 4.0.0 GA, make 4.0.0 the default
* [CHANGE] Upgrade to reaper-operator 0.3.3 and Reaper 2.3.0
* [CHANGE] Upgrade from Stargate 1.0.18 to 1.0.29
* [CHANGE] Upgrade from Medusa 0.10.1 to 0.11.0
* [CHANGE] Upgrade from Reaper 2.2.2 to 2.2.5
* [CHANGE] #812 Integrate Fossa component/license scanning
* [CHANGE] #905 Upgrade medusa-operator to v0.3.3
* [FEATURE] #890 Upgrade from Management API v0.1.25 to v0.1.26 to provide support for Cassandra 4.0.0-RC2
* [FEATURE] #617 Make affinity configurable for Stargate
* [FEATURE] #847 Make affinity configurable for Reaper
* [FEATURE] #952 Experimental support for custom init containers
* [ENHANCEMENT] #844 Allow configuring the namespace of service monitors
* [ENHANCEMENT] #29 Detect IEC formatted c* heap.size and heap.newGenSize; return error identifying issue  
* [ENHANCEMENT] #420 Add support for private registries
* [ENHANCEMENT] #685 Add support for Medusa backups on Azure
* [BUGFIX] #853 Fix property name in scaling docs
* [BUGFIX] #870 Hot replace disallowed characters in generated secret names
* [BUGFIX] #412 Stargate metrics don't show up in the dashboards
