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

{{% alert title="Tip" color="success" %}}
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
kubectl get secret k8ssandra-superuser -o json | jq -r '.data.username' | base64 --decode
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

The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.

You can override the default username by setting the `stargate.cassandraUser.username` property.

Credentials are stored in a secret named `{cassandra.clusterName}-stargate`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-stargate -o json | jq -r '.data.username' | base64 --decode
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

The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.

You can override the default username by setting the `reaper.cassandraUser.username` property.

Credentials are stored in a secret named `{cassandra.clusterName}-reaper`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-reaper -o json | jq -r '.data.username' | base64 --decode
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

The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.

You can override the default username by setting the `medusa.cassandraUser.username` property.

Credentials are stored in a secret named `{cassandra.clusterName}-medusa`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-medusa -o json | jq -r '.data.username' | base64 --decode
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

K8ssandra uses secrets to store credentials. For every set of credentials, K8ssandra supplies a default username and a default password. K8ssandra generates a random, alphanumeric string 20 characters long for the password. These values are stored under the username and password keys in the secret.

Like other objects installed as part of the Helm release, secrets will be deleted when you run `helm uninstall`.

Each K8ssandra component (each deployment) has a username property to override the default username: 

* Cassandra -  `cassandra.auth.superuser.name`
* Stargate - `stargate.cassandraUser.username`
* Reaper - `reaper.cassandraUser.username`
* Medusa - `medusa.cassandraUser.username`

Alternatively, you can provide your own secret for each component:

* Cassandra -  `cassandra.auth.superuser.secret` 
* Stargate - `stargate.cassandraUser.secret`
* Reaper - `reaper.cassandraUser.secret`
* Medusa - `medusa.cassandraUser.secret`

The secret must have username and password keys. The secret also needs to exist in the same namespace in which you are installing K8ssandra.

Because the secret is created outside of the Helm release, it is not managed by Helm. Running `helm uninstall` will not delete the secret. 

{{% alert title="Recommendation" color="success" %}}
Be consistent in your handling of secrets. Either manage all of them outside of Helm, or let all of them be created by K8ssandra.
{{% /alert %}}

## Next

Learn how to [develop client]({{< relref "/tasks/develop" >}}) using the Stargate APIs. 

