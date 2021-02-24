# Changelog

Changelog for K8ssandra, new PRs should update the `main / unreleased` section with entries in the order:

```markdown
* [CHANGE]
* [FEATURE]
* [ENHANCEMENT]
* [BUGFIX]
```

When cutting a new release of the parent `k8ssandra` chart update the `main / unreleased` heading to the tag being generated and create a new placeholder section.

## main / unreleased

* [FEATURE] #409 Add support for `configOverride`
* [FEATURE] #398 Add support for running a subset of tests
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
