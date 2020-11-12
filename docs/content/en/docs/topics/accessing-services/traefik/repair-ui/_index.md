---
title: "Repair UI"
linkTitle: "Repair UI"
weight: 1
date: 2020-11-07
description: |
  Configuring Traefik to expose the Reaper Repair interface
---

## Prerequisites

1. Running Kubernetes cluster
   1. K8ssandra operators deployed
   1. Traefik deployed
   1. K8ssandra cluster deployed
1. DNS name where the repair service should be listening
1. Kubectl
1. Helm

## Helm Parameters

The `k8ssandra-cluster` Helm chart contains templates for Traefik `IngressRoute`
and `IngressRouteTCP` Custom Resources. These may be enabled at any time either
through a `values.yaml` file of command-line flags.

### `values.yaml`
```yaml
ingress:
  traefik:
    # Set to `true` to enable the templating of Traefik ingress custom resources
    enabled: false

    # Repair service
    repair: 
      # Note this will **only** work if `traefik.enabled` is also `true`
      enabled: true

      # Name of the Traefik entrypoints where we want to source traffic.
      entrypoints: 
        - web

      # Hostname Traefik should use for matching requests.
      host: repair.k8ssandra.cluster.local
```

// TODO - describe the importance of the DNS name, mention services like xip.io

## Enabling Traefik Ingress

### Command-line
```bash
# New Install
helm install cluster-name k8ssandra/k8ssandra-cluster \
  --set ingress.traefik.enabled=true \
  --set ingress.traefik.repair.host=repair.cluster-name.k8ssandra.cluster.local

# Existing Cluster
helm upgrade cluster-name k8ssandra/k8ssandra-cluster \
  --set ingress.traefik.enabled=true \
  --set ingress.traefik.repair.host=repair.cluster-name.k8ssandra.cluster.local
```

### `values.yaml`

// TODO - Discuss why using values.yaml is a good idea. Version control, CI/CD, etc

```bash
# New Install
helm install cluster-name k8ssandra/k8ssandra-cluster -f traefik.values.yaml

# Existing Cluster
helm upgrade cluster-name k8ssandra/k8ssandra-cluster -f traefik.values.yaml
```

## Validate Traefik Configuration

_Note this step is optional. The next step will also prove the configuration is working._

With the ingress routes configured and deployed to Kubernetes we can access the Traefik dashboard to validate the configuration has been picked up and is detecting the appropriate services.

1. Open your web browser and point it at the Traefik dashboard. This may require `kubectl port-forward` or the steps in our [Configuring Kind]({{< ref "configuring-kind" >}}) guide.
2. Navigate to the HTTP Routers page
    // TODO - Screenshot of routers page with repair rule
3. Navigate to the HTTP Services page
    // TODO - Screenshot of services page with repair service 

## Accessing Repair Interface

// TODO - insert reaper screenshot

With configuration complete and validated all that is left it to point your browser at the DNS name and access the GUI. Assuming the DNS name of `repair.cluster-name.k8ssandra.cluster.local` we would visit [http://repair.cluster-name.k8ssandra.cluster.local/webui](http://repair.cluster-name.k8ssandra.cluster.local/webui). Traefik receives the request, matches the `Host` header against the rule specified in our `IngressRoute` and proxies the request to the upstream service. Should this service go down and the pod get rescheduled everything will automatically update and continue functioning.
