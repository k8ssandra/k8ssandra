# K8ssandra.io - Website & Documentation Repository

This project is built with the [Hugo](https://gohugo.io/) static site generator and the [Docsy](https://github.com/google/docsy) theme.

## Version Management

With the introduction of the K8ssandra Operator, the documentation site has introduced a version switching capability, using the Hugo/Docsy capabilities described [here](https://www.docsy.dev/docs/adding-content/versioning/).

There three "versioned" docs sites:

"Common" content -- https://docs-temp.k8ssandra.io (at the official release of k8ssandra-operator this will be made available at https://docs.k8ssandra.io)

"v1" (Original `helm` focused implementation) content -- https://docs-v1.k8ssandra.io

"v2" (`k8ssandra-operator` implementation) content -- https://docs-v2.k8ssandra.io

## Dependencies

This project requires Node.js and NPM to be installed locally.  All other dependencies are provided via NPM.

The latest version of Node/NPM should work, the current automation uses Node 14.

You can install Node in a myriad of ways, check [here](https://nodejs.org/en/) for more information.

## Cloning the repo

Because of the way Docsy is included in projects as a git submodule it's necessary to clone the project in a recursive fashion, using a command like:

```
git clone --recurse-submodules https://github.com/k8ssandra/k8ssandra.git
```

## Branch Management

Docsy/Hugo versioning works through an association of site "version" to Git "branch".  The documentation content and tooling support are provided for each versioned side on a series of docs branches.

* `docs` branch -> `docs-temp.k8ssandra.io`/`docs-staging.k8ssandra.io`

* `docs-v1` branch -> `docs-v1.k8ssandra.io`/`docs-staging-v1.k8ssandra`

* `docs-v2` branch -> `docs-v2.k8ssandra.io`/`docs-staging-v2.k8ssandra`

## GitHub Actions

On each "docs" branch GitHub Action workflows are configured to provide automated PR checks, deployments to staging, and deployments to production.

On each PR, that content will be checked to confirm that a docs build succeeds with the proposed changes.

On each push, the new contents of the branch will be built and published to the `docs-staging*` site associated with that branch.

On the push of a semantically structured tag the tagged commit will be built and published to the `docs*` site associated with that branch.

* `docs-v*.*.*` tag -> `docs.k8ssandra.io`

* `docs-v1-v*.*.*` tag -> `docs-v1.k8ssandra.io`

* `docs-v2-v*.*.*` tag -> `docs-v2.k8ssandra.io`

## Development

### Install dependencies

Working against the particular docs branch of interesting: `docs`, `docs-v1`, or `docs-v2`...

From the `/docs` directory, run

```
npm install
```

This command will install the project dependencies, such as hugo-extended.

### Scripts

There are a number of utility scripts provided in package.json that can be executed for various purposes.  Different scripts are available on each docs branch to limit the tooling to act on the deployment target of interest for that branch and reduce the risk of accidental manual deployments to other/unintended version sites.

#### Start the Hugo server for local development

```
npm run start
```

This provides a live-refresh server that will automatically load any changes made in the source code.  

The local server will be available at http://localhost:1313/.

*Note: Only one version of the docs site can be run at any given time locally, because different versions are contained within different branches.*
