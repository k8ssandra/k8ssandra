---
title: "Contribution guidelines"
linkTitle: "Contribute"
weight: 9
description: How to contribute to the K8ssandra open-source code and documentation.
---

We welcome contributions from the K8ssandra community! 

## Code contributions

The overall procedure:

1. Start on https://github.com/k8ssandra/k8ssandra.
2. Fork the repo by clicking the **Fork** button in the GitHub UI.
3. Make your changes locally on your fork. Git commit and push only to your fork.
4. Wait for CI to run successfully in GitHub Actions before submitting a PR.
5. Submit a Pull Request (PR) with your forked updates.
6. If you're not yet ready for a review, add "WIP" to the PR name to indicate it's a work in progress.  
7. Wait for the automated PR workflow to do some checks. Members of the K8ssandra community will review your PR and decide whether to approve and merge it.

Also, we encourage you to submit Issues, starting from https://github.com/k8ssandra/k8ssandra/issues. Add a label to help categorize the issue, such as the complexity level, component name, and other labels you'll find in the repo's Issues display. 

## Setting up CI in GitHub Actions

CI will run on push to any branch of your forked repository. In order to run CI jobs involving S3 or
Azure buckets (for Medusa), the following GitHub secrets need to be set up in the fork's
organization:

- `K8SSANDRA_MEDUSA_BUCKET_NAME`: Name of the S3 bucket or Azure container used to store Medusa 
  backups.
- `K8SSANDRA_MEDUSA_BUCKET_REGION`: Region of the S3 bucket (for Amazon S3 tests only).
- `K8SSANDRA_MEDUSA_SECRET_S3`: A valid Kubernetes secret definition to access the S3 bucket (for
  Amazon S3 tests only). See below for an example.
- `K8SSANDRA_MEDUSA_SECRET_AZURE`: A valid Kubernetes secret definition to access the Azure storage 
  account (for Azure tests only). See below for an example.

### Kubernetes secret definitions for Medusa tests 

When testing Medusa on Amazon S3, the GitHub `K8SSANDRA_MEDUSA_SECRET_S3` secret _must_  contain a
valid Kubernetes secret definition similar to the one below:

```yaml
apiVersion: v1
kind: Secret
metadata:
  # The name must be medusa-bucket-key!
  name: medusa-bucket-key
type: Opaque
stringData:
  # The entry must be named medusa_s3_credentials!
  medusa_s3_credentials: |-
    [default]
    aws_access_key_id = <your key id>
    aws_secret_access_key = <your key secret>
```

See [Backup and restore with Amazon S3]({{< relref "/tasks/backup-restore/amazon-s3/" >}}) for
more information.

When testing Medusa on Azure, the GitHub `K8SSANDRA_MEDUSA_SECRET_AZURE` secret _must_ contain a
valid Kubernetes secret definition similar to the one below:

```yaml
apiVersion: v1
kind: Secret
metadata:
  # The name must be medusa-bucket-key!
  name: medusa-bucket-key
type: Opaque
stringData:
  # The entry must be named medusa_azure_credentials.json!
  medusa_azure_credentials.json: |-
    {
        "storage_account": "<your storage account name>",
        "key": "<your storage account key>"
    }
```

See [Backup and restore with Azure Storage]({{< relref "/tasks/backup-restore/azure/" >}}) for
more information.

## Prevent CI from running for a specific commit

CI runs can be disabled for specific commits by adding `[skip ci]` in the commit message. This can be useful when pushing commits on WIP branches for backup purposes.

## Documentation contributions and build environment

We use [Hugo](https://gohugo.io/) to format and generate this website, the [Docsy](https://github.com/google/docsy) theme for styling and site structure, and [Google Cloud Storage](https://console.cloud.google.com/) to manage the deployment of the site.

Hugo is an open-source static site generator that provides us with templates, content organization in a standard directory structure, and a website generation engine. You write the pages in Markdown (or HTML if you want), and Hugo wraps them up into a website.

All submissions, including submissions by project members, require review. We use GitHub pull requests for this purpose. Consult
[GitHub Help](https://help.github.com/articles/about-pull-requests/) for more information on using pull requests.

Most of the documentation source files are authored in [markdown](https://daringfireball.net/projects/markdown/) plus some special properties unique to Hugo.

### Deploying to K8ssandra.io

Here's a quick guide to updating the docs. It assumes you're familiar with the GitHub workflow and you're happy to use the automated preview of your doc updates:

1. **Fork** the [K8ssandra docs repo](https://github.com/k8ssandra/k8ssandra.git) on GitHub. 
1. Make your changes and send a pull request (PR).
1. If you're not yet ready for a review, add "WIP" to the PR name to indicate 
  it's a work in progress. (**Don't** add the Hugo property 
  "draft = true" to the page front matter.)
1. Wait for the automated PR workflow to do some checks.
1. Continue updating your doc and pushing your changes until you're happy with 
  the content.
1. When you're ready for a review, add a comment to the PR, and remove any
  "WIP" markers.
1. After the Pull Request is reviewed and merged it will be deployed automatically to the development version of the docs site (https://docs.k8ssandra-dev.io). There is usually a delay of 10 or more minutes between deployment and when the updates are online.

Deployments to the production site (https://docs.k8ssandra.io) are made by pushing a semantic version structured tag in the form of `v{major}.{minor}.{patch}-docs`.  This process is done one when docs content has stabilized and is ready to be made generally available. 

### Updating a single page

If you've just spotted something you'd like to change while using the docs, K8ssandra.io has a shortcut for you:

1. Click **Edit this page** in the top right hand corner of the page.
1. If you don't already have an up to date fork of the project repo, you are prompted to get one - click **Fork this repository and propose changes** or **Update your Fork** to get an up to date version of the project to edit. The appropriate page in your fork is displayed in edit mode.
1. Follow the rest of the [Deploying to K8ssandra.io](#deploying-to-k8ssandraio) process above to make and propose your changes.

### Previewing your changes locally

If you want to run your own local Hugo server to preview your changes as you work:

1. Install Node/NPM using the instructions from the [NodeJS download](https://nodejs.org/en/download/) site (the latest LTS release version should be sufficient).
1. Fork the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra) repo into your own project, then create a local copy using `git clone`. Don’t forget to use `--recurse-submodules` or you won’t pull down some of the code you need to generate a working site.

    ```bash
    git clone --recurse-submodules https://github.com/k8ssandra/k8ssandra.git
    ```
1. Navigate to the `docs` directory within the k8ssandra repository
1. Run `npm install` to install all of the dependencies required for building/testing with Hugo
1. Run `npm run start` to start a live hugo server locally.  As changes are made and saved the site will automatically be rebuilt and be made available through a browser refresh.

    ```bash
    % npm run start

    > k8ssandra@0.2.0 start
    > hugo server --bind 0.0.0.0

    port 1313 already in use, attempting to use an available port
    Start building sites … 
    hugo v0.89.4-AB01BA6E+extended darwin/amd64 BuildDate=2021-11-17T08:24:09Z VendorInfo=gohugoio
    WARN 2021/12/09 08:55:07 Page.Dir is deprecated and will be removed in a future release. Use .File.Dir

                      | EN   
    -------------------+------
      Pages            | 121  
      Paginator pages  |   0  
      Non-page files   | 109  
      Static files     |  51  
      Processed images |   0  
      Aliases          |   0  
      Sitemaps         |   1  
      Cleaned          |   0  

    Built in 690 ms
    Watching for changes in /k8ssandra/k8ssandra/docs/{assets,content,i18n,layouts,package.json,static,themes}
    Watching for config changes in /k8ssandra/k8ssandra/docs/config.toml, /Users/jeffdinoto/github/jdonenine/forks/k8ssandra/k8ssandra/docs/themes/docsy/config.toml
    Environment: "development"
    Serving pages from memory
    Running in Fast Render Mode. For full rebuilds on change: hugo server --disableFastRender
    Web Server is available at http://localhost:56711/ (bind address 0.0.0.0)
    Press Ctrl+C to stop
    ```

1. Continue with the usual GitHub workflow to edit files, commit them, push the
  changes up to your fork, and create a pull request.

### Creating an issue

If you've found a problem in the docs, but you're not sure how to fix it yourself, please create an issue in the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra/issues). You can also create an issue about a specific page by clicking the **Create Issue** button in the top right hand corner of the page.

## Next steps

Refer to these useful resources:

* [Docsy user guide](https://www.docsy.dev/docs/): All about Docsy, including how it manages navigation, look and feel, and multi-language support.
* [Hugo documentation](https://gohugo.io/documentation/): Comprehensive reference for Hugo.
* [Github Hello World!](https://guides.github.com/activities/hello-world/): A basic introduction to GitHub concepts and workflow.
