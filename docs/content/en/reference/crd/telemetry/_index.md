---
title: "Telemetry CRD"
linkTitle: "Telemetry CRD"
simple_list: false
weight: 6
description: >
  Telemetry Custom Resource Definition (CRD) reference for use with K8ssandra Operator.
---

### Custom Resources



* [PrometheusTelemetrySpec](#prometheustelemetryspec)
* [TelemetrySpec](#telemetryspec)

#### PrometheusTelemetrySpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| enabled | Enable the creation of Prometheus serviceMonitors for this resource (Cassandra or Stargate). | bool | false |
| commonLabels | CommonLabels are applied to all serviceMonitors created. | map[string]string | false |

[Back to Custom Resources](#custom-resources)

#### TelemetrySpec



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| prometheus |  | *[PrometheusTelemetrySpec](#prometheustelemetryspec) | false |

[Back to Custom Resources](#custom-resources)
