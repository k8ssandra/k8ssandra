# K8ssandra.io - Website & Documentation Repository

This project is built with the [Hugo](https://gohugo.io/) static site generator and the [Docsy](https://github.com/google/docsy) theme.

The documentation produced from this repository can be viewed at:

Development Version -- https://docs.k8ssandra-dev.io

Production Version -- https://docs.k8ssandra.io

## Dependencies

This project requires Node.js and NPM to be installed locally.  All other dependencies are provided via NPM.

The latest version of Node/NPM should work, the current automation uses Node 14.

You can install Node in a myriad of ways, check [here](https://nodejs.org/en/) for more information.

## Cloning the repo

Because of the way Docsy is included in projects as a git submodule it's necessary to clone the project in a recursive fashion, using a command like:

```
git clone --recurse-submodules https://github.com/k8ssandra/k8ssandra.git
```

## Development

### Install dependencies

From the `/docs` directory, run

```
npm install
```

This command will install the project dependencies, such as hugo-extended.

### Scripts

There are a number of utility scripts provided in package.json that can be executed for various purposes:

#### Start the Hugo server for local development

```
npm run start
```

This provides a live-refresh server that will automatically load any changes made in the source code.  

The local server will be available at http://localhost:1313/.

#### Use Hugo to build the site

```
npm run build:dev
```

or

```
npm run build:prod
```

#### Cleanup the build artifacts previously produced

```
npm run clean
```

#### Deploy to docs.k8ssandra-dev.io

**This requires local gcloud authentication and permissions**

```
npm run deploy:dev
```

This will replace all content currently hosted at docs.k8ssandra-dev.io -- use with caution.
