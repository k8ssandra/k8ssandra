---
title: "Upgrade K8ssandra"
linkTitle: "Upgrade"
no_list: true
weight: 2
description: How to upgrade K8ssandra to use the latest release, new settings, or both.
---

You can easily upgrade your K8ssandra cluster. In this topic, we'll describe how to:

* Update an existing K8ssandra repo to the latest release
* Take action based on an upgrade consideration for K8ssandra 1.1.0
* Upgrade from the single-node Cassandra instance to a 3-node Cassandra instance, as an example

## Introduction

Upgrading a K8ssandra instance is a multi-step process:

1. Use `helm repo update` to ensure your Kubernetes cluster is using the latest K8ssandra software.
1. Follow the steps for any upgrade considerations, as is the case when going from K8ssandra 1.0.0 to 1.1.0. 
1. Update or create a configuration YAML file with the changes you want to apply the cluster.
1. Apply the changes using the `helm upgrade` command.

An assumption here is that you previously installed K8ssandra, with a command such as:

```
helm repo add k8ssandra https://helm.k8ssandra.io/stable/
```

## Update the K8ssandra repo

To update your installed version of K8ssandra, enter:

```
helm repo update
```

**Output:**

```bash
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "k8ssandra" chart repository
Update Complete. ⎈Happy Helming!⎈
```

For example, because K8ssandra released 1.1.0 on 09-Apr-2021, the `helm repo update` command automatically gets the latest software. 

## Upgrade notice for K8ssandra 1.1.0

As cited in the K8ssandra [release notes]({{< relref "/release-notes/#upgrade-notice" >}}), upgrading from K8ssandra 1.0.0 to 1.1.0 causes a StatefulSet update, which has the effect of a rolling restart. This situation could require you to perform a manual restart of all Stargate nodes after the Cassandra cluster is back online. 

To manually restart Stargate nodes:

1. Get the Deployment object in your Kubernetes environment:
   ```bash
   kubectl get deployment | grep stargate
   ```
2. Scale down with this command:
   ```bash
   kubectl scale deployment <stargate-deployment> --replicas 0
   ```
3. Run this next command and wait until it reports 0/0 ready replicas. This should happen within a couple seconds.
   ```bash
   kubectl get deployment <stargate-deployment>
   ```
4. Now scale up with:
   ```bash
    kubectl scale deployment <stargate-deployment> --replicas 1
    ```

## Upgrade example with node scaling

To upgrade your single-node instance to a 3-node instance:

1. Create a new `k8ssandra-upgrade.yaml` file with the following configuration fragment:

    ```yaml
    cassandra:
      datacenters:
      - name: dc1
        size: 3
    ```

    The cluster size has increased from `1` to `3`

    {{% alert title="Tip" color="success" %}}
You only need the YAML statements pertinent to the upgrade. You don't need to duplicate the entire original configuration file.
    {{% /alert %}}

1. Upgrade the cluster using the `helm upgrade` command:

    ```bash
    helm upgrade -f k8ssandra-upgrade.yaml k8ssandra k8ssandra/k8ssandra
    ```

    **Output**:

    ```bash
    Release "k8ssandra" has been upgraded. Happy Helming!
    NAME: k8ssandra
    LAST DEPLOYED: Mon Apr 12 11:17:37 2021
    NAMESPACE: default
    STATUS: deployed
    REVISION: 2
    ```

    Notice that the REVISION is now at `2`. It will increment each time you run a `helm upgrade` command. 

{{% alert title="Tip" color="success" %}}
For insights into the underlying operations that occur with scaling, see [Scale your Cassandra cluster]({{< relref "scale" >}}).
{{% /alert %}}


1. Monitor `kubectl get pods` until the new Cassandra nodes are up and running:

    ```bash
    kubectl get pods
    ```

    **Output**:

    ```bash
    NAME                                                  READY   STATUS      RESTARTS   AGE
    k8ssandra-cass-operator-6666588dc5-qpvzg              1/1     Running     4          2d2h
    k8ssandra-dc1-default-sts-0                           2/2     Running     0          76m
    k8ssandra-dc1-default-sts-1                           2/2     Running     0          3m29s
    k8ssandra-dc1-default-sts-2                           2/2     Running     0          3m28s
    k8ssandra-dc1-stargate-6f7f5d6fd6-sblt8               1/1     Running     13         2d2h
    k8ssandra-grafana-6c4f6577d8-hsbf7                    2/2     Running     0          3m32s
    k8ssandra-kube-prometheus-operator-5556885bd6-st4fp   1/1     Running     4          2d2h
    k8ssandra-reaper-k8ssandra-5b6cc959b7-zzlzr           1/1     Running     22         2d2h
    k8ssandra-reaper-k8ssandra-schema-47qzk               0/1     Completed   0          2d2h
    k8ssandra-reaper-operator-cc46fd5f4-85mk5             1/1     Running     5          2d2h
    prometheus-k8ssandra-kube-prometheus-prometheus-0     2/2     Running     9          2d2h
    ```

   Eventually you should see two additional K8ssandra pods with the extensions `-sts-1` and `-sts-2` in `RUNNING` status.

## Next steps

Explore other K8ssandra [tasks]({{< relref "/tasks" >}}).

See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra Helm charts, and a glossary.  
