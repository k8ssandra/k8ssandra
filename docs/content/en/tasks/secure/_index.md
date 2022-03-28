---
title: "K8ssandra security"
linkTitle: "Secure"
no_list: true
weight: 3
description: K8ssandra security defaults, secrets, and options for Apache Cassandra&reg; authentication and role-based authorization.
---

**Note:** Please refer to [Enabling encryption]({{< relref "/tasks/secure/encryption/" >}}) in K8ssandra clusters; its description has been updated for K8ssandra Operator.

## Introduction

K8ssandra enables authentication and authorization by default. It uses the Cassandra default `PasswordAuthenticator` and `CassandraAuthorizer` functionality.  

{{% alert title="Tip" color="success" %}}
We recommend that you keep authentication enabled. Turning on auth for an existing cluster is a non-trivial exercise that may involve downtime.
{{% /alert %}}

Authentication can be disabled by setting `spec.auth` to `false` in the K8ssandraCluster object spec:

```
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: cluster1
spec:
  auth: false
  ...
```

## Cassandra security

With authentication enabled, K8ssandra configures a new, default superuser. The username defaults to `{metadata.name}-superuser`. 

{{% alert title="Note" color="success" %}}
K8ssandra disables and does not use the default superuser, `cassandra`.
{{% /alert %}}

The password is a random alphanumeric string 20 characters long.

You can override the default username and password by setting the `spec.cassandra.superuserSecretRef` property to an existing secret containing both entries.

Credentials are stored in a secret named `{metadata.name}-superuser`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-superuser -o json | jq -r '.data.username' | base64 --decode
```

```bash
kubectl get secret k8ssandra-superuser -o json | jq -r '.data.password' | base64 --decode
```

For more, see the [secrets]({{< relref "#secrets" >}}) section of this security topic.

## Stargate security

Stargate has no specific credentials. It uses the same superuser as defined for Cassandra.

## Reaper security

With authentication enabled, K8ssandra creates three distinct users for Reaper:

* `{metadata.name}-reaper` for Reaper's access to Cassandra using CQL
* `{metadata.name}-reaper-jmx` for Reaper's access to Cassandra using JMX
* `{metadata.name}-reaper-ui` for Reaper's UI credentials

The password is a random, alphanumeric string 20 characters long.

The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.

You can override the default username/password combinations by setting the following properties:

* `spec.reaper.cassandraUserSecretRef`
* `spec.reaper.jmxUserSecretRef`
* `spec.reaper.uiUserSecretRef`

Credentials are stored in a secret named `{metadata.name}-reaper`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

```bash
kubectl get secret k8ssandra-reaper -o json | jq -r '.data.username' | base64 --decode
```

```bash
kubectl get secret k8ssandra-reaper -o json | jq -r '.data.password' | base64 --decode
```

For more, see the [secrets]({{< relref "#secrets" >}}) section of this security topic.

## Medusa security

With authentication enabled, K8ssandra creates a default user for Medusa. The default username is `{metadata.name}-medusa`. 

The password is a random, alphanumeric string 20 characters long.

The user is created as a superuser because (at this time) K8ssandra does not support configuring authorization.

You can override the default credentials by setting the `spec.medusa.cassandraUserSecretRef` property, to point to an existing secret containing both the `username` and `password` entries.

Credentials are stored in a secret named `{metadata.name}-medusa`. If your cluster name is `k8ssandra`, for example, you can retrieve the username and password as follows:

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

For more, see the [secrets]({{< relref "#secrets" >}}) section of this security topic.

## Secrets

Kubernetes secrets let you store and manage sensitive information, such as passwords, OAuth tokens, and ssh keys. Storing confidential information in a Secret is safer and more flexible than putting it verbatim in a Pod definition or in a container image. You can use the secrets generated by K8ssandra components, or create your own secrets.

K8ssandra uses secrets to store credentials. For every set of credentials, K8ssandra supplies a default username and a default password. K8ssandra generates a random, alphanumeric string 20 characters long for the password. These values are stored under the username and password keys in the secret.

Like other objects installed as part of the Helm release, secrets will be deleted when you run `helm uninstall`.

Each K8ssandra component (each deployment) has a property to provide your own secret for each component:

* Cassandra
  * `spec.cassandra.superuserSecretRef` 
* Reaper
  * `spec.reaper.cassandraUserSecretRef`
  * `spec.reaper.jmxUserSecretRef`
  * `spec.reaper.uiUserSecretRef`
* Medusa
  * `spec.medusa.cassandraUserSecretRef`

The secret must have `username` and `password` keys. The secret also needs to exist in the same namespace in which you are installing K8ssandra.

{{% alert title="Recommendation" color="success" %}}
Be consistent in your handling of secrets. Either manage all by yourself, or let all of them be created by K8ssandra.
{{% /alert %}}

## JMX configuration and access 

By default, Cassandra restricts JMX access to localhost. The Reaper component that's deployed by K8ssandra requires JMX access for managing repairs. When Reaper is enabled and deployed, K8ssandra enables and configures remote JMX access.

### Remote Access

JMX access is configured in `/etc/cassandra/cassandra-env.sh`. The script checks the environment variable `LOCAL_JMX` to determine whether JMX access should be restricted to localhost. When Reaper is enabled, K8ssandra sets this `LOCAL_JMX` environment variable in the `cassandra` container to a value of `no`, in order to enable remote JMX access.

### JMX authentication

Cassandra turns on JMX authentication by default when remote access is enabled. JMX credentials are stored in `/etc/cassandra/jmxremote.password`. 

K8ssandra creates two sets of credentials - one for Reaper and one for the Cassandra default superuser.

### Reaper

The username of the Reaper JMX user defaults to `{metadata.name}-reaper-jmx`. The password is a random, 20 character alphanumeric string. Credentials are stored in a secret. See the [secrets]({{< relref "#secrets" >}}) section in this topic for an overview of how K8ssandra manages secrets. The username and secret can be overridden with the following property:

* `spec.reaper.jmxUserSecretRef`

### Cassandra default superuser

K8ssandra creates JMX credentials for the default superuser. The username and password are the same as those for Cassandra.

If you change the Cassandra superuser credentials through `cqlsh` for example, the JMX credentials are not updated to the new values. 

### nodetool

When JMX authentication is enabled, you need to specify the username and password options with `nodetool`, as follows:

```bash
nodetool -u <username> -pw <password> status
```

### JMX authorization - not supported at this time

K8ssandra currently does not support JMX authorization.

## Next steps

* Explore other K8ssandra Operator [tasks]({{< relref "/tasks" >}}).
* See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Operator Custom Resource Definitions (CRDs) and the single K8ssandra Operator Helm chart. 
