# Vagrant Deployed Cluster

## VM Resources

The VM is deployed with the following system resources:

* 4 vCPUs
* 8GB Memory

## Required Components

The VM is bootstrapped with the following basic components, all ultimately required to deploy and command the k8ssandra cluster.

* Docker
* Kind
* Helm
* Kubectl

The following ports are forwarded from the host machine to the VM:

* 8080
* 8443
* 9000
* 9042
* 9142

**The VM will fail to start if conflicts on these ports are detected with the host machine.**

## k8ssandra Deployment

The following charts will be installed via helm:

* traefik (namespace: traefik)
* k8ssandra (namespace: default)
* k8ssandra-cluster (namespace: default)

The k8ssanda-cluster is configuration to support the following:

* Ingress via traefik
* Cassandra 3.11.7
* Size = 1
* Repair
* Backup
* Stargate (replicas = 1)

## Cassandra User Authentication

The Cassandra username and password can be accessed by executing the following commands.

### Username

```
kubectl get secret k8ssandra-cluster-superuser -o jsonpath="{.data.username}" | base64 --decode
```

### Password

```
kubectl get secret k8ssandra-cluster-superuser -o jsonpath="{.data.password}" | base64 --decode
```

## Access URLs

The k8ssandra services deployed can be accessed at the following URLs:

[Traefik](http://127.0.0.1.nip.io:9000/dashboard/#/)

[Grafana](http://grafana.127.0.0.1.nip.io:8080/)

[Prometheus](http://prometheus.127.0.0.1.nip.io:8080/)

[Reaper](http://repair.127.0.0.1.nip.io:8080/webui/)

[Stargate Authentication APIs](http://auth.127.0.0.1.nip.io:8080)

[Stargate REST APIs](http://auth.127.0.0.1.nip.io:8080)

[Stargate GraphQL APIs](http://graphql.127.0.0.1.nip.io:8080/)

[Stargate GraphQL Playground](http://graphql.127.0.0.1.nip.io:8080/playground)

