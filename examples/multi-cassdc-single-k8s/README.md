# Example: Multiple Cassandra datacenters in a single Kubernetes cluster

This example configuration deploys a Cassandra cluster consisting of two Cassandra datacenters with three nodes per datacenter. Each datacenter is deployed in a separate namespace within a shared Kubernetes cluster. 

## Background and setup
This example is designed to run on a Google Kubernetes Engine (GKE) cluster in the `us-west4` region. Please see the [Google GKE](https://docs.k8ssandra.io/install/gke/) page for instructions on setting up an appropriate cluster you can use with this example, tailoring the zone/region settings as needed. 

Alternatively, you can make this example more generic and run on another Kubernetes engine by removing the affinity lines in the yaml files.

These instructions assume you've installed Helm and that your Kubernetes context is set to the cluster you want to use. Consult the [Getting Started](https://k8ssandra.io/get-started) documentation for more information.

The instructions guide you through creating two Cassandra datacenters in a single Kubernetes cluster, one designated for transactional workloads (`txndc`), and one for analytics workloads (`analyticsdc`).

## Instructions

1. Create the namespaces and administrator credentials for each datacenter:

    ```
    kubectl create namespace txndc
    kubectl create secret generic cassandra-admin-secret --from-literal=username=cassandra-admin --from-literal=password=cassandra-admin-password -n txndc
    kubectl create namespace analyticsdc
    kubectl create secret generic cassandra-admin-secret --from-literal=username=cassandra-admin --from-literal=password=cassandra-admin-password -n analyticsdc
    ```

1. Install K8ssandra in the transactional namespace for the second datacenter using the [dc1.yaml](./dc1.yaml) config file:

    ```
    helm install txndc k8ssandra/k8ssandra -f dc1.yaml -n txndc  
    ```

1. After the first datacenter has initialized (you can watch status using `watch kubectl get pods -n txndc`), install the second datacenter in the analytics namespace using the [dc2.yaml](./dc2.yaml) config file. 

    ```
    helm install analyticsdc k8ssandra/k8ssandra -f dc2.yaml -n analyticsdc
    ```

1. Configure Cassandra keyspaces to replicate across both clusters. To do this, connect to a node and execute `cqlsh`:
    
    ```
    kubectl exec mixed-workload-dc1-rack1-sts-0 -n txndc -it -- cqlsh -u cassandra-admin -p cassandra-admin-password
    ```
   
    Then, at the `cqlsh` prompt, modify the keyspaces with the desired replication per datacenter:
    
    ```
    ALTER KEYSPACE system_auth WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': 3, 'dc2': 3}
    ALTER KEYSPACE system_distributed WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': 3, 'dc2': 3}
    ALTER KEYSPACE system_traces WITH replication = {'class': 'NetworkTopologyStrategy', 'dc1': 3, 'dc2': 3}
    ```

    Remember to use `NetworkTopologyStrategy` with appropriate replication per datacenter for any additional keyspaces you create for your application. Then you can exit `cqlsh`.
    
1. Make sure data is properly replicated to each node in the analytics datacenter using the `nodetool rebuild` command:

    ```
    kubectl exec mixed-workload-dc2-rack1-sts-0 -n analytics -- nodetool --username cassandra-admin --password cassandra-admin-password rebuild dc1
    kubectl exec mixed-workload-dc2-rack2-sts-0 -n analytics -- nodetool --username cassandra-admin --password cassandra-admin-password rebuild dc1
    kubectl exec mixed-workload-dc2-rack3-sts-0 -n analytics -- nodetool --username cassandra-admin --password cassandra-admin-password rebuild dc1
    ```

1. Use the `nodetool status` command to verify the second datacenter has joined the cluster:

    ```
    kubectl exec mixed-workload-dc1-rack1-sts-0 -n txndc -- nodetool --username cassandra-admin --password cassandra-admin-password status
   ```

## More information
For more detailed explanation of this configuration, see the [blog post](https://k8ssandra.io/blog/tutorials/deploy-a-multi-datacenter-apache-cassandra-cluster-in-kubernetes).

