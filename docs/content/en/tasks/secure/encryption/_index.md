---
title: "Enabling encryption"
linkTitle: "Enabling encryption"
toc_hide: true
no_list: true
weight: 1
description: Setting up encryption in K8ssandra clusters.
---

Apache Cassandra&reg; offers the ability to encrypt internode communications and client-to-node communications separately. This topic explains how to set up and configure encryption in K8ssandra clusters.

## Prerequisites

* A supported Kubernetes 1.19+ environment, either local (kind, K3D, minikube) or via a cloud provider:
  * Amazon Elastic Kubernetes Service (EKS)
  * DigitalOcean Kubernetes (DOKS)
  * Google Kubernetes Engine (GKE) in a Google Cloud project
  * Microsoft Azure Kubernetes Service (AKS)
* **K8ssandra Operator** has been installed - see the [install]({{< relref "install" >}}) topics
* An SSL encryption store, as covered in the next section

## Generating SSL encryption stores

If you do not have a set of encryption stores available, follow the instructions in [this TLP blog post](https://thelastpickle.com/blog/2021/06/15/cassandra-certificate-management-part_1-how-to-rotate-keys.html). More specifically, use [this script](https://github.com/thelastpickle/cassandra-toolbox/tree/main/generate_cluster_ssl_stores) to generate the SSL stores.

You could clone the [cassandra-toolbox](https://github.com/thelastpickle/cassandra-toolbox) GitHub repository, and create a `cert.conf` file with the following format:

```conf
[ req ]
distinguished_name     = req_distinguished_name
prompt                 = no
output_password        = MyPassWord123!
default_bits           = 2048

[ req_distinguished_name ]
C                      = FR
ST                     = IDF
L                      = Paris
O                      = YourCompany
OU                     = SSLTestCluster
CN                     = SSLTestClusterRootCA
emailAddress           = youraddress@whatever.com
```

Next, run:

```bash
./generate_cluster_ssl_stores.sh -v 10000 -g cert.conf
```

The `-v` value above sets the validity of the generated certificates in days. 

{{% alert title="Tip" color="success" %}}
Don't set this `-v` days value too low. Doing so would require you to rotate the certificates too often; it's not a trivial operation.
{{% /alert %}}

The command output should be a folder containing a keystore, a truststore, and a file containing their respective passwords.

Rename the keystore file to `keystore`, and rename the truststore file to `truststore`. Then create a Kubernetes secret with the following command:

```bash
kubectl create secret generic server-encryption-stores --from-file=keystore --from-literal=keystore-password=<keystore password> --from-file=truststore --from-literal=truststore-password=<truststore password> -o yaml > server-encryption-stores.yaml
```

Replace the `<keystore password>` and `<truststore password>` above with each store's actual password.

{{% alert title="Tip" color="success" %}}
You can repeat the above procedure to generate encryption stores for client-to-node encryption, changing the secret name appropriately.
{{% /alert %}}

## Creating a cluster with internode encryption

In order to create a K8ssandra cluster with encryption, first create a namespace and the encryption stores secrets previously generated in it.

In the `K8ssandraCluster` manifest, you will need to configure encryption settings in the `config/cassandraYaml` section.

Also, you'll need to reference the encryption stores' secrets under:

* `cassandra/serverEncryptionStores` 
* *Or*`cassandra/clientEncryptionStores`

Server encryption and client encryption are different entities. They both have their own keystore/truststore pair.
The "or" here shows that you can turn on either independently, or both. Server is for internode communications encryption, and client is for client-to-node communications encryption.

Example:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  cassandra:
    serverVersion: "4.0.1"
    storageConfig:
      cassandraDataVolumeClaimSpec:
        storageClassName: standard
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 5Gi
    config:
      cassandraYaml:
        server_encryption_options:
            internode_encryption: all
            require_client_auth: true
            ...
            ...
        client_encryption_options:
            enabled: true
            require_client_auth: true
            ...
            ...
    datacenters:
      - metadata:
          name: dc1
        size: 3
    serverEncryptionStores:
      keystoreSecretRef:
        name: server-encryption-stores
      truststoreSecretRef:
        name: server-encryption-stores
    clientEncryptionStores:
      keystoreSecretRef:
        name: client-encryption-stores
      truststoreSecretRef:
        name: client-encryption-stores
```

Enabling client-to-node encryption will also encrypt JMX communications. Running Cassandra `nodetool` commands will then require additional arguments to pass the encryption stores and their passwords.

{{% alert title="Note" color="success" %}}
Again, server (internode) and client (client-to-node) encryption are totally independent and can be enabled/disabled individually, as well as use different encryption stores.
{{% /alert %}}

## Stargate and Reaper encryption

Stargate and Reaper will both inherit from Cassandra's encryption settings without any additional change to the manifest.

An encrypted cluster with both Stargate and Reaper would be deployed with the following manifest:

```yaml
apiVersion: k8ssandra.io/v1alpha1
kind: K8ssandraCluster
metadata:
  name: test
spec:
  cassandra:
    serverVersion: "4.0.1"
    storageConfig:
      cassandraDataVolumeClaimSpec:
        storageClassName: standard
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 5Gi
    config:
      cassandraYaml:
        server_encryption_options:
            internode_encryption: all
            require_client_auth: true
            ...
            ...
        client_encryption_options:
            enabled: true
            require_client_auth: true
            ...
            ...
    datacenters:
      - metadata:
          name: dc1
        size: 3
    serverEncryptionStores:
      keystoreSecretRef:
        name: server-encryption-stores
      truststoreSecretRef:
        name: server-encryption-stores
    clientEncryptionStores:
      keystoreSecretRef:
        name: client-encryption-stores
      truststoreSecretRef:
        name: client-encryption-stores
  stargate:
    size: 1
  reaper:
    deploymentMode: SINGLE
```

## Next steps

Explore other K8ssandra [tasks]({{< relref "/tasks" >}}).

See the [Reference]({{< relref "/reference" >}}) topics for information about K8ssandra charts, and a glossary. 
