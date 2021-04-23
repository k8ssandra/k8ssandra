---
title: "K8ssandra security"
linkTitle: "Secure"
no_list: true
weight: 3
description: K8ssandra security defaults, secrets, and options for Apache Cassandra&reg; authentication and role-based authorization.
---

This topic describes how K8ssandra supports Cassandra's authentication and authorization features. Also, how Kubernetes secrets are created and used by the K8ssandra deployed components: Cassandra (via cass-operator), Stargate, Reaper, and Medusa. 

## Introduction

K8ssandra enables authentication and authorization by default. It uses the Cassandra default `PasswordAuthenticator` and `CassandraAuthorizer` functionality.  

{{% alert title="Important" color="success" %}}
We recommend that you keep authentication enabled. Turning on auth for an existing cluster is a non-trivial exercise that may involve downtime.
{{% /alert %}}

## Cassandra security

With authentication enabled, K8ssandra configures a new, default superuser. The username defaults to `{cassandra.clusterName}-superuser`. 

{{% alert title="Note" color="success" %}}
K8ssandra disables and does not use the default superuser, `cassandra`.
{{% /alert %}}

The password is a random alphanumeric string 20 characters long.

You can override the default username by setting the `cassandra.auth.superuser.username` property.

Credentials are stored in a secret named `{cassandra.clusterName}-superuser`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-superuser -o json | jq -r '.data.username | base64 --decode
```

```bash
kubectl get secret k8ssandra-superuser -o json | jq -r '.data.password' | base64 --decode
```

You can override both the username and password by providing your own secret and setting `cassandra.auth.superuser.secret` to its name.

If you provide your own secret it will not be managed by Helm. Helm will not do anything with it when you run `helm upgrade` or `helm uninstall`, for example.

For more, see the [secrets]({{< relref "#secrets" >}}) section of this security topic.

If both `cassandra.auth.superuser.username` and `cassandra.auth.superuser.secret` are set, `cassandra.auth.superuser.secret` takes precedence.

## Stargate security

With authentication enabled, K8ssandra creates a default user for Stargate. The default username is `stargate`. 

The password is a random, alphanumeric string 20 characters long.

{{% alert title="Note" color="success" %}}
The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.
{{% /alert %}}

You can override the default username by setting the `stargate.cassandraUser.username` property.

Credentials are stored in a secret named `{cassandra.clusterName}-stargate`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-stargate -o json | jq -r '.data.username | base64 --decode
```

```bash
$ kubectl get secret k8ssandra-stargate -o json | jq -r '.data.password' | base64 --decode
```

You can override both the username and password by providing your own secret and setting `stargate.cassandraUser.secret` to its name. 

If you provide your own secret it will not be managed by Helm. Helm will not do anything with it when you run `helm upgrade` or `helm uninstall`, for example.

For more, see the [secrets]({{< relref "#secrets" >}}) section of this security topic.

If both `stargate.cassandraUser.username` and `stargate.cassandraUser.superuser.secret` are set, `stargate.cassandraUser.secret` takes precedence.

## Reaper security

With authentication enabled, K8ssandra creates a default user for Reaper. The default username is `reaper`. 

The password is a random, alphanumeric string 20 characters long.

{{% alert title="Note" color="success" %}}
The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.
{{% /alert %}}

You can override the default username by setting the `reaper.cassandraUser.username` property.

Credentials are stored in a secret named `{cassandra.clusterName}-reaper`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-reaper -o json | jq -r '.data.username | base64 --decode
```

```bash
kubectl get secret k8ssandra-reaper -o json | jq -r '.data.password' | base64 --decode
```

You can override both the username and password by providing your own secret and setting `reaper.cassandraUser.secret` to its name.

If you provide your own secret it will not be managed by Helm. Helm will not do anything with it when you run `helm upgrade` or `helm uninstall`, for example.

For more, see the [secrets]({{< relref "#secrets" >}}) section of this security topic.

If both `reaper.cassandraUser.username` and `reaper.cassandraUser.superuser.secret` are set, `reaper.cassandraUser.secret` takes precedence.

## Medusa security

With authentication enabled, K8ssandra creates a default user for Medusa. The default username is `medusa`. 

The password is a random, alphanumeric string 20 characters long.

{{% alert title="Note" color="success" %}}
The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.
{{% /alert %}}

You can override the default username by setting the `medusa.cassandraUser.username` property.

Credentials are stored in a secret named `{cassandra.clusterName}-medusa`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-medusa -o json | jq -r '.data.username | base64 --decode
```

```bash
kubectl get secret k8ssandra-medusa -o json | jq -r '.data.password' | base64 --decode
```

You can override both the username and password by providing your own secret and setting `medusa.cassandraUser.secret` to its name.

If you provide your own secret it will not be managed by Helm. Helm will not do anything with it when you run `helm upgrade` or `helm uninstall`, for example.

For more, see the [secrets]({{< relref "#secrets" >}}) section of this security topic.

If both `medusa.cassandraUser.username` and `medusa.cassandraUser.superuser.secret` are set, `medusa.cassandraUser.secret` takes precedence.

## Secrets

Kubernetes secrets let you store and manage sensitive information, such as passwords, OAuth tokens, and ssh keys. Storing confidential information in a Secret is safer and more flexible than putting it verbatim in a Pod definition or in a container image. You can use the secrets generated by K8ssandra components, or create your own secrets.

TODO: More info about secrets in K8ssandra.

## Next

Learn how to [develop client]({{< relref "/tasks/develop" >}}) using the Stargate APIs. 
