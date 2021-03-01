---
title: "Contribution guidelines"
linkTitle: "Contribute"
weight: 6
description: >
  How to contribute to K8ssandra documentation.
---

We use [Hugo](https://gohugo.io/) to format and generate this website, the [Docsy](https://github.com/google/docsy) theme for styling and site structure, and [Google Cloud Storage](https://console.cloud.google.com/) to manage the deployment of the site.

Hugo is an open-source static site generator that provides us with templates, content organization in a standard directory structure, and a website generation engine. You write the pages in Markdown (or HTML if you want), and Hugo wraps them up into a website.

All submissions, including submissions by project members, require review. We use GitHub pull requests for this purpose. Consult
[GitHub Help](https://help.github.com/articles/about-pull-requests/) for more information on using pull requests.

## Deploying to K8ssandra.io

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
1. After the Pull Request is reviewed and merged it will be deployed automatically. There is usually a delay of 10 or more minutes between deployment and when the updates are online. 

## Updating a single page

If you've just spotted something you'd like to change while using the docs, K8ssandra.io has a shortcut for you:

1. Click **Edit this page** in the top right hand corner of the page.
1. If you don't already have an up to date fork of the project repo, you are prompted to get one - click **Fork this repository and propose changes** or **Update your Fork** to get an up to date version of the project to edit. The appropriate page in your fork is displayed in edit mode.
1. Follow the rest of the [Deploying to K8ssandra.io](#deploying-to-k8ssandraio) process above to make and propose your changes.

## Previewing your changes locally

If you want to run your own local Hugo server to preview your changes as you work:

1. Follow the instructions to install [Hugo](https://gohugo.io/) and any other tools you need. You'll need at least **Hugo version 0.45** (we recommend using the most recent available version), and it must be the **extended** version, which supports SCSS.
1. Fork the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra) repo into your own project, then create a local copy using `git clone`. Don’t forget to use `--recurse-submodules` or you won’t pull down some of the code you need to generate a working site.

    ```bash
    git clone --recurse-submodules --depth 1 https://github.com/k8ssandra/k8ssandra.git
    ```

1. Run `hugo server` in the docs site root directory, such as your `~/github/k8ssandra/docs` directory. By default your site will be available at http://localhost:1313/. Now that you're serving your site locally, Hugo will watch for changes to the content and automatically refresh your site.
1. Continue with the usual GitHub workflow to edit files, commit them, push the
  changes up to your fork, and create a pull request.

## Creating an issue

If you've found a problem in the docs, but you're not sure how to fix it yourself, please create an issue in the [K8ssandra repo](https://github.com/k8ssandra/k8ssandra/issues). You can also create an issue about a specific page by clicking the **Create Issue** button in the top right hand corner of the page.

## Useful resources

* [Docsy user guide](https://www.docsy.dev/docs/): All about Docsy, including how it manages navigation, look and feel, and multi-language support.
* [Hugo documentation](https://gohugo.io/documentation/): Comprehensive reference for Hugo.
* [Github Hello World!](https://guides.github.com/activities/hello-world/): A basic introduction to GitHub concepts and workflow.
