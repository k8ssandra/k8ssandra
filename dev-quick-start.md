# K8ssandra Developer - Quick Start

## About
This quick start guide is written for developers wanting to find out more about:
* Getting up and running with a basic IDE environment.
* Deploying to a local docker-based cluster environment (using kind).
* Understanding the K8ssandra project structure.
* Running unit tests.
* Troubleshooting.

## Environment setup
Let’s get started by setting up the foundations of a development environment.  Where possible, this guide attempts to maintain an operating system agnostic approach. 

Regarding K8ssandra supported operating systems, the following are supported: 
* MacOS
* Linux
* Windows 10
  

**Note:** For Windows, you may want to consider using [Windows Subsystem for Linux](https://docs.microsoft.com/en-us/windows/wsl/about) to utilize all the shell scripts provided.

Although K8ssandra is created to run in Kubernetes cloud environments, the scope of this article is focused on local machine development.  Checkout this great post [Requirements for Running K8ssandra](https://k8ssandra.io/blog/2021/03/10/requirements-for-running-k8ssandra-for-development/) for details on managing expectations if running on a resource constrained machine.

When installing an integrated development environment (IDE), it’s important to understand the file types you will be editing.  As such, let’s quickly cover the types of files involved with K8ssandra and the Kubernetes ecosystem.

### File types
K8ssandra is mainly composed of:
* YAML files used for declarative configurations.  In fact, many of the YAML files are a collection of Helm charts.  [Helm](https://helm.sh/) is an open source package management solution for Kubernetes and Helm charts are the packaging format used.

* [Go](https://golang.org/) (.go) files primarily used for construction of unit & integration tests.

* Documentation files.  From a document perspective, K8ssandra’s website [K8ssandra.io](https://k8ssandra.io), contains a wealth of information and is maintained with languages such as: HTML, Markdown, and CSS.

* Shell & Python scripts existing in the project for assisting with task automation and K8ssandra configurations.

With a  basic understanding of the file types involved, let's install an open source IDE that assists with constructing and editing K8ssandra files.

## Setup IDE

**Note:** If you already have an IDE setup supporting Kubernetes development, feel free to skip past this section.

One popular open source integrated development environment (IDE) that can be used to develop in the K8ssandra project is [Visual Studio Code](https://code.visualstudio.com/).  Another great alternative is [JetBrains' GoLand](https://www.jetbrains.com/go/download) IDE, which isn’t free for all users, but does provide an evaluation/trial period.

This guide will reference the free and open source **VS Code**.  

### VS Code installation
First, download the VS Code binaries and install the necessary extensions to assist with K8ssandra development.

[Download VS Code](https://code.visualstudio.com/) specific to your operating system. 
   
    
**Note**: at the time of this writing, version 1.54.2 was used.

Now, add the [Go extension](https://code.visualstudio.com/docs/languages/go) to your VS Code IDE.  This addition allows for management and compiling of Go source. This can be accomplished by navigating to **File**->**Preferences**->**Extensions** in VS Code and searching for “Go”. 

At the time of this writing, version **0.23.2** (not the nightly build) was used from the Go Team at Google.  
   
Follow the steps in the Go extension’s documentation, which includes a download of Go specific to your operating system.  

Select **1.14+** as a version.


### Install Git & GitHub
If you don’t already have Git installed and have a GitHub account, use the references below to get those in-place.

Reference the following to understand and install **Git**:
* [Setup Guide](https://github.com/git-guides/install-git)
* [Getting started with GitHub](https://github.com/join)


## Kubernetes environment
A quick session (10 minutes to complete) is provided below that will guide you through the various steps required to set up a running Kubernetes-based environment for K8ssandra. 

Once completed, return to this article and finish-up the rest of the activities.

> [K8ssandra Quickstarts](https://docs.k8ssandra.io/quickstarts/)


## Checkpoint
Before continuing on to next steps, you should have the following configured:

    ✔ VS Code - or something similar (IDE)
    ✔ Go (the Golang binaries)
    ✔ Git (a version control system)
    ✔ GitHub account (for forking and contributing)



## Repository
The K8ssandra GitHub repository resides [here](https://github.com/k8ssandra/k8ssandra).  

Once at the K8ssandra repository, click on the **Fork** button (top right corner of screen). This will fork the K8ssandra to your local repository.

Using Git, you can now pull the K8ssandra source code and documentation to your local machine.  Some developers choose to do this via the GitHub plugin-in in an IDE or from a command line.

Perform a clone.  Be sure to insert your GitHub repository name followed by k8ssandra.git.
      
>git clone https://github.com/your-github-repo-name/k8ssandra.git
   
Now you are ready to explore the K8ssandra project.

### Project composition
The K8ssandra project has a few important folders of interest.

* **Charts** - contains the Helm-based charts categorized by sub-chart.
* **Docs** - contains the layouts, & content for the [k8ssandra.io](https://k8ssandra.io) site.  
Checkout the docs README for more detail.
* **Scripts** - contains a collection of useful scripts that developers can build upon or just use directory when working on K8ssandra. Also included are documentation scripts for K8ssandra docs management.
* **Tests** - contains the set of unit, integration, and end-to-end testing for K8ssandra.


## Testing
Within the test directory are folders for managing and executing tests at the `unit`, `integration`, and `e2e` test levels.

As a working example, using a command line, navigate to the k8ssandra project root directory where the Makefile resides.

Issue the command

> make unit-test

Once complete, you should see something like the following:

> ok      github.com/k8ssandra/k8ssandra/tests/unit       47.156s


### Integration tests
Different integration test scenarios are available to check that all components work as expected.
When using kind, invoke the full stack scenario with the following command:

```
make kind-integ-test
```

This will delete any existing `k8ssandra_it` kind cluster and recreate it accordingly with the integration tests requirements.
In order to run the tests on a specific (and supported) Apache Cassandra version, set the `K8SSANDRA_CASSANDRA_VERSION` environment variable:

```
K8SSANDRA_CASSANDRA_VERSION=4.0.0 make kind-integ-test
```

Supported versions are referenced in [the Helm charts](https://github.com/k8ssandra/k8ssandra/blob/main/charts/k8ssandra/values.yaml#L2-L8).

If you want to run the integration tests on an running Kubernetes cluster (GKE, EKS, K3d, ...), run the following command instead:

```
make integ-test
```

Individual scenarios can be run for faster results if the change primarily impacts a specific component:

```
make kind-integ-test TESTS="TestReaperDeploymentScenario"
```

The Medusa tests are specific to a storage backend, which needs to be specified as follows:

```
make kind-integ-test TESTS="TestMedusaDeploymentScenario/\"Minio\""
```

By default, integration tests will cleanup the cluster from the deployed resources even in case of failure. To disable such cleanups and allow investigation, set the `CLUSTER_CLEANUP` argument to `success` which will only cleanup on success:

```
make kind-integ-test CLUSTER_CLEANUP="success" 
```

To fully disable clean up, use the following command:

```
make kind-integ-test CLUSTER_CLEANUP="never" 
```

## Troubleshooting advice
This section will be updated over time as K8ssandra grows.  Take a look to enhance your K8ssandra experience. 

### Common errors

You may experience a `missing in charts/ directory` error message.  

If so, you can utilize a K8ssandra script: `./scripts/update-helm-deps.sh`. This script assists with updating dependencies for each chart in an appropriate order.  

Be sure to run this script so the `./charts` folder is properly located. 

  
### Collecting useful information 

So you have an error after editing a K8ssandra configuration, or you want to inspect some things as you learn.  There are some useful commands that come in handy when needing to dig a bit deeper.  The examples assume you are using a k8ssandra namespace, but this can be adjusted as needed.

Issue the following `kubectl` command to view the `Management-api` logs.  Replace *cassandra-pod* with an actual pod instance name.

> kubectl logs *cassandra-pod* -c cassandra -n k8ssandra

Issue the following `kubectl` command to view the `Cassandra` logs.  Replace *cassandra-pod* with an actual pod instance name.

> kubectl logs *cassandra-pod* -c server-system-logger -n k8ssandra

Issue the following `kubectl` command to view `Medusa` logs.  Replace *cassandra-pod* with an actual pod instance name.

> kubectl logs *cassandra-pod* -c medusa -n k8ssandra

Issue the following `kubectl` command to describe the `CassandraDatacenter` resource.  This provides a wealth of information about the resource, which includes `aged events` that assist when trying to troubleshoot an issue.

> kubectl describe cassandradatacenter/dc1 -n k8ssandra

Gather container specific information for a pod.

 First, list out the pods scoped to the K8ssandra namespace and instance with a target release.

> kubectl get pods -l app.kubernetes.io/instance=*release-name* -n k8ssandra

Note: If you don't know the release name, look it up with:

> helm list -n k8ssandra

Next, targeting a specific pod, filter out `container` specific information. Replace the name of the pod with the pod of interest.

> kubectl describe pod/*pod-name* -n k8ssandra | grep container -C 3

A slight variation, list out pods having the label for a `cassandra` cluster.

> kubectl get pods -l cassandra.datastax.com/cluster=*release-name* -n k8ssandra

Now, using a pod-name returned, describe all the details.

> kubectl describe pod/*pod-name* -n k8ssandra


## Next steps
Now that you have a foundation for using K8ssandra, take a look at some other references to better understand what
 is available and where K8ssandra is headed.

* [Components](https://docs.k8ssandra.io/components/)
* [Tasks](https://docs.k8ssandra.io/tasks/)
* [Helm Charts](https://docs.k8ssandra.io/reference/helm-charts/)
* [K8ssandra Blog](https://k8ssandra.io/blog/)
* [Community](https://k8ssandra.io/community/)
* [Roadmap](https://docs.k8ssandra.io/roadmap/)
