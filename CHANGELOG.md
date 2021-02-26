# Changelog

Changelog for K8ssandra, new PRs should update the `main / unreleased` section with entries in the order:

```markdown
* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX]
```

When cutting a new release of the parent `k8ssandra` chart update the `main / unreleased` heading to the tag being generated and date `## vX.Y.Z - YYYY-MM-DD` and create a new placeholder section for  `main / unreleased` entries.

## v1.0.0 - 2021-02-26

* [ENHANCEMENT] #444 Upgrade cass-operator to 1.6.0
* [ENHANCEMENT] #450 Add CHANGELOG.md

## v0.60.3 - 2021-02-25

* [ENHANCEMENT] #429 Update reaper-operator version and make resource names more consistent
* [ENHANCEMENT] #435 Make secret template name consistent
* [BUGFIX] #432 Update name of reaper ingress service reference
* [BUGFIX] #414 Do not generate new passwords if secret already exists

## v0.59.0 - 2021-02-24

* [ENHANCEMENT] #423 Upgrade to Reaper 2.2.1 and Medusa 0.9.0
* [BUGFIX] #426 Fix issue with tarball dependencies

## v0.58.0 - 2021-02-23

* [FEATURE] #409 Add support for `configOverride`
* [FEATURE] #398 Add support for running a subset of tests
* [FEATURE] #419 Add support for Stable and Next release streams
* [ENHANCEMENT] #403 Update Stargate version

## v0.55.0 - 2021-02-22

* [FEATURE] #397 Add support for Stargate and C* 4.0
* [FEATURE] #326 Generate JMX credentials for C* superuser
* [ENHANCEMENT] #401 Update Reaper version

## v0.54.0 - 2021-02-22

* [FEATURE] #355 Make garbage collection configurable

## v0.53.0 - 2021-02-19

* [FEATURE] #336 Support running non-root images

## v0.52.0 - 2021-02-19

* [ENHANCEMENT] #361 Refactor Stargate ingress

## v0.51.0 - 2021-02-18

* [FEATURE] #382 Provide default `num_tokens` based on C* version
* [ENHANCEMENT] #317 Support for `s3_compatible` settings
