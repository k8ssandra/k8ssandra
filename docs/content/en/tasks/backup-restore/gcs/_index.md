---
title: "Backup and restore with Google Cloud Storage"
linkTitle: "Google Cloud Storage"
no_list: true
weight: 3
description: Use Medusa to backup and restore Apache Cassandra® data in Kubernetes to Google Cloud Storage (GCS).
---

Medusa is a Cassandra backup and restore tool. It's packaged with K8ssandra and supports a variety of backends, including GCS for storage in Google Kubernetes Engine (GKE) environments.

## Introduction

Google Cloud Storage (GCS) is a RESTful online file storage web service for storing and accessing data on Google Cloud Platform (GCP) / GKE infrastructure. The service combines the performance and scalability of Google's cloud with advanced security and sharing capabilities.

For information about GCS, see the [Google Cloud Storage documentation](https://cloud.google.com/storage).

## Next steps

See the following reference topics:

* [Medusa Operator Helm Chart]({{< relref "/reference/helm-charts/medusa-operator" >}})
* [Backup Helm Chart]({{< relref "/reference/helm-charts/backup" >}})
* [Restore Helm Chart]({{< relref "/reference/helm-charts/restore" >}})