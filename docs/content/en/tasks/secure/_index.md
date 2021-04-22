---
title: "Secure K8ssandra managed assets"
linkTitle: "Secure K8ssandra"
no_list: true
weight: 2
description: K8ssandra security defaults, secrets, and options for Apache Cassandra&reg; authentication and role-based authorization.
---

Intro sentences here. Content is TBS.

## Introduction

Content TBS. 

## Authentication in Cassandra deployed by K8ssandra cass-operator

Content TBS. 

By default, cass-operator renmoves the username &amp; password credentials of cassandra cassandra.

## Authorization roles used by K8ssandra deployments

Content TBS. 

K8ssandra provides a `superuser` role that ...

Reaper - repairs over JMX authorization (remote enabled). Must provide credentials to run `nodetool`. Note caveat of changed creds not being propagated for JMX superuser.

Stargate uses...

Medusa uses ... 

## Secrets

How we create them...

Where they’re used...

Per user role, can specify the username and K8ssandra generates the secret value. 

How to extract passwords...

Alternatively, you can provide your own secret...

What happens on uninstall...

## Next

Learn how to [develop client apps]({{< relref "/tasks/develop" >}}) in a Kubernetes cluster that is managed by K8ssandra. 
