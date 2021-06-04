# K8ssandra upfront validations
Author: @jeffbanks

## Enhancement
Request for validator tools to be used to notify of invalid K8ssandra configurations.  Rules will be customizable and support detection of issues earlier in the deployment process.

For example, in the K8ssandra jvm-option heap settings there is a requirement to support the following property structure when customized.

```yaml
cassandra:
  # cluster-level settings
  heap:
    # min & max equal
    size: 700M
    newGenSize: 1.4Gi
  # dc specific settings
  dataCenters:
    - name: dc1
      size: 1
      heap:
        size: 600M
        newGenSize: 1.2Gi
    - name: dc2
      size: 1
```

In the example, when required values are not specified, a heap property exists without any child properties.

These types of issues could be detected, perhaps as a collection of issues, and reported to be addressed by a K8ssandra monitor.


## Proposal

Use of a configuration linting tool to assist with catching customized configuration problems in advance of issues discovered when applying Kubernetes objects.

Use of complementary static-analysis config tools for Kubernetes standards verification.

### Tools for consideration

* kube-score - Kubernetes object analysis with recommendations for improved reliability and security.
* config-lint  - configuration lint, available as a CLI for CI/CD purposes and can be used as a Golang module in unit and integration tests.
* kubeval - allows for YAML manifest validation against specific versions of API Schemas.


## kube-score

Kubernetes object analysis with recommendations for improved reliability and security.

### Features:

* Open-source and available under the MIT-license.

* Static code analysis of your Kubernetes object definitions.

* Output contains a list of recommendations to make applications more secure and resilient.

* [Checks included](https://github.com/zegl/kube-score/blob/master/README_CHECKS.md)

* [Repository](https://github.com/zegl/kube-score)

### Usage

```
kube-score [action] --flags

    Actions:
            **score **  Checks all files in the input, and gives them a score and recommendations
            **list **   Prints a CSV list of all available score checks
            **version **Print the version of kube-score
            **help **   Print this message


    Flags for score:

         --disable-ignore-checks-annotations   Set to true to disable the effect of the 'kube-score/ignore' annotations

         --enable-optional-test strings        Enable an optional test, can be set multiple times

         --exit-one-on-warning                 Exit with code 1 in case of warnings

         --help                                Print help

         --ignore-container-cpu-limit          Disables the requirement of setting a container CPU limit

         --ignore-container-memory-limit       Disables the requirement of setting a container memory limit

         --ignore-test strings                 Disable a test, can be set multiple times

         --kubernetes-version string           Setting the kubernetes-version will affect the checks ran against the manifests. Set this to the version of Kubernetes that you're using in production for the best results. (default "v1.18")

     -o, --output-format string                Set to 'human', 'json' or 'ci'. If set to ci, kube-score will output the program in a format that is easier to parse by other programs. (default "human")

         --output-version string               Changes the version of the --output-format. The 'json' format has version 'v2' (default) and 'v1' (deprecated, will be removed in v1.7.0). The 'human' and 'ci' formats has only version 'v1' (default). If not explicitly set, the default version for that particular output format will be used.

     -v, --verbose count                       Enable verbose output, can be set multiple times for increased verbosity.

```

### Example kube-score w/ K8ssandra

View the list of score criteria.

```kube-score list```

This prints a CSV list of all available score checks.

Example:
```
stable-version, all, Checks if the object is using a deprecated apiVersion, default
label-values, all, Validates label values, default
```

Generate the target configuration.

```helm template k8ssandra/k8ssandra > kt.yaml```

### Example portion of .yaml generated

```yaml
# Source: k8ssandra/templates/cassandra/cassdc.yaml
apiVersion: cassandra.datastax.com/v1beta1
kind: CassandraDatacenter
metadata:
  name: dc1
  labels:     
    app.kubernetes.io/name: k8ssandra
    helm.sh/chart: k8ssandra-1.1.0
    app.kubernetes.io/instance: RELEASE-NAME
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/part-of: k8ssandra-RELEASE-NAME-default
  annotations:
    reaper.cassandra-reaper.io/instance: RELEASE-NAME-reaper
spec:
  clusterName: RELEASE-NAME
  serverType: cassandra
  serverVersion: "3.11.10"
  dockerImageRunsAsCassandra: true
  serverImage: k8ssandra/cass-management-api:3.11.10-v0.1.24
  managementApiAuth:
    insecure: {}
  size: 1
  racks:
  - name: default
  storageConfig:
    cassandraDataVolumeClaimSpec:
      storageClassName: standard
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 5Gi
  allowMultipleNodesPerWorker: false
  users:
    - secretName: RELEASE-NAME-reaper
      superuser: true
    - secretName: RELEASE-NAME-stargate
      superuser: true
  config:
    cassandra-yaml:
      num_tokens: 256
      authenticator: PasswordAuthenticator
      authorizer: CassandraAuthorizer
      role_manager: CassandraRoleManager
      roles_validity_in_ms: 3.6e+06
      roles_update_interval_in_ms: 3.6e+06
      permissions_validity_in_ms: 3.6e+06
      permissions_update_interval_in_ms: 3.6e+06
      credentials_validity_in_ms: 3.6e+06
      credentials_update_interval_in_ms: 3.6e+06
    jvm-options:
      additional-jvm-opts:
        - "-Dcassandra.system_distributed_replication_dc_names=dc1"
        - "-Dcassandra.system_distributed_replication_per_dc=1"
  podTemplateSpec:
    spec:
      initContainers:
      - name: base-config-init
        image: k8ssandra/cass-management-api:3.11.10-v0.1.24
        imagePullPolicy: IfNotPresent
        command:
          - /bin/sh
        args:
          - -c
          - cp -r /etc/cassandra/* /cassandra-base-config/
        volumeMounts:
          - name: cassandra-config
            mountPath: /cassandra-base-config/
      - name: server-config-init
      - name: jmx-credentials
        image: busybox
        imagePullPolicy: IfNotPresent
        env:
          - name: REAPER_JMX_USERNAME
            valueFrom:
              secretKeyRef:
                name: RELEASE-NAME-reaper-jmx
                key: username
          - name: REAPER_JMX_PASSWORD
            valueFrom:
              secretKeyRef:
                name: RELEASE-NAME-reaper-jmx
                key: password
          - name: SUPERUSER_JMX_USERNAME
            valueFrom:
              secretKeyRef:
                name: RELEASE-NAME-superuser
                key: username
          - name: SUPERUSER_JMX_PASSWORD
            valueFrom:
              secretKeyRef:
                name: RELEASE-NAME-superuser
                key: password
        args:
          - /bin/sh
          - -c
          - echo "$REAPER_JMX_USERNAME $REAPER_JMX_PASSWORD" > /config/jmxremote.password && echo "$SUPERUSER_JMX_USERNAME $SUPERUSER_JMX_PASSWORD" >> /config/jmxremote.password
        volumeMounts:
          - mountPath: /config
            name: server-config
      containers:
      - name: cassandra
        env:
          - name: LOCAL_JMX
            value: "no"
      volumes:
      - name: cassandra-config
        emptyDir: {}
```

### Score the target configuration.

```kube-score score kt.yaml > kt.score```

#### Example kube-score w/ k8ssandra cass-operator deployment

deployment.yaml

```yaml
# Source: k8ssandra/charts/cass-operator/templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: RELEASE-NAME-cass-operator
  labels:     
    app.kubernetes.io/name: cass-operator
    helm.sh/chart: cass-operator-0.29.1
    app.kubernetes.io/instance: RELEASE-NAME
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/part-of: k8ssandra-RELEASE-NAME-default
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: cass-operator
      app.kubernetes.io/instance: RELEASE-NAME
      app.kubernetes.io/part-of: k8ssandra-RELEASE-NAME-default
  template:
    metadata:
      labels:        
        app.kubernetes.io/name: cass-operator
        helm.sh/chart: cass-operator-0.29.1
        app.kubernetes.io/instance: RELEASE-NAME
        app.kubernetes.io/managed-by: Helm
        app.kubernetes.io/part-of: k8ssandra-RELEASE-NAME-default
    spec:
      serviceAccountName: RELEASE-NAME-cass-operator
      securityContext:
        {}
      containers:
        - name: cass-operator
          securityContext:
            readOnlyRootFilesystem: true
            runAsGroup: 65534
            runAsNonRoot: true
            runAsUser: 65534
          image: "docker.io/datastax/cass-operator:1.6.0"
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - mountPath: /tmp/k8s-webhook-server/serving-certs
              name: cass-operator-certs-volume
              readOnly: false
            - mountPath: /tmp/
              name: tmpconfig-volume
              readOnly: false
          livenessProbe:
            exec:
              command:
                - pgrep
                - ".*operator"
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            exec:
              command:
                - stat
                - "/tmp/operator-sdk-ready"
            initialDelaySeconds: 5
            periodSeconds: 5
            timeoutSeconds: 5
            failureThreshold: 1
          resources:
            {}
          env:
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "cass-operator"
            - name: SKIP_VALIDATING_WEBHOOK
              value: "TRUE"
      volumes:
        - name: tmpconfig-volume
          emptyDir:
            medium: "Memory"
        - name: cass-operator-certs-volume
          secret:
            secretName: RELEASE-NAME-webhook
```

#### kube-score output (human)

```
apps/v1/Deployment RELEASE-NAME-cass-operator                                 💥
    [OK] Stable version
    [OK] Label values
    [CRITICAL] Container Resources
        · cass-operator -> CPU limit is not set
            Resource limits are recommended to avoid resource DDOS. Set resources.limits.cpu
        · cass-operator -> Memory limit is not set
            Resource limits are recommended to avoid resource DDOS. Set resources.limits.memory
        · cass-operator -> CPU request is not set
            Resource requests are recommended to make sure that the application can start and run without crashing. Set
            resources.requests.cpu
        · cass-operator -> Memory request is not set
            Resource requests are recommended to make sure that the application can start and run without crashing. Set
            resources.requests.memory
    [OK] Container Image Tag
    [CRITICAL] Container Image Pull Policy
        · cass-operator -> ImagePullPolicy is not set to Always
            It's recommended to always set the ImagePullPolicy to Always, to make sure that the imagePullSecrets are always correct, and to
            always get the image you want.
    [CRITICAL] Pod NetworkPolicy
        · The pod does not have a matching NetworkPolicy
            Create a NetworkPolicy that targets this pod to control who/what can communicate with this pod. Note, this feature needs to be
            supported by the CNI implementation used in the Kubernetes cluster to have an effect.
    [OK] Pod Probes
        · The pod is not targeted by a service, skipping probe checks.
    [OK] Container Security Context

```

#### kube-score output (ci)

```
[OK] RELEASE-NAME-cass-operator apps/v1/Deployment
[OK] RELEASE-NAME-cass-operator apps/v1/Deployment
[CRITICAL] RELEASE-NAME-cass-operator apps/v1/Deployment: The pod does not have a matching NetworkPolicy
[OK] RELEASE-NAME-cass-operator apps/v1/Deployment: The pod is not targeted by a service, skipping probe checks.
[OK] RELEASE-NAME-cass-operator apps/v1/Deployment
[CRITICAL] RELEASE-NAME-cass-operator apps/v1/Deployment: (cass-operator) CPU limit is not set
[CRITICAL] RELEASE-NAME-cass-operator apps/v1/Deployment: (cass-operator) Memory limit is not set
[CRITICAL] RELEASE-NAME-cass-operator apps/v1/Deployment: (cass-operator) CPU request is not set
[CRITICAL] RELEASE-NAME-cass-operator apps/v1/Deployment: (cass-operator) Memory request is not set
[OK] RELEASE-NAME-cass-operator apps/v1/Deployment
[CRITICAL] RELEASE-NAME-cass-operator apps/v1/Deployment: (cass-operator) ImagePullPolicy is not set to Always
[SKIPPED] RELEASE-NAME-cass-operator apps/v1/Deployment: Skipped because the deployment has less than 2 replicas
[SKIPPED] RELEASE-NAME-cass-operator apps/v1/Deployment: Skipped because the deployment has less than 2 replicas
[SKIPPED] RELEASE-NAME-cass-operator apps/v1/Deployment: Skipped because the deployment is not targeted by a HorizontalPodAutoscaler

```


## config-lint

Stelligent config-lint is an open source command line tool to lint configuration files in a variety of formats: JSON, YAML, Terraform, and Kubernetes.

### Features:

* Enables developers with abiltiy to validate configuration files.
* Provides custom validations and built-in rules to ensure configurations meet best practices.
* Custom files can be used for other formats.
* Built-in rules provided for Terraform


#### Example Kubernetes

```yaml
version: 1
description: Rules for Kubernetes spec files
type: Kubernetes
files:
  - "*.yml"
rules:

  - id: ALLOW_KIND
    severity: FAILURE
    message: Allowed kinds
    resource: "*"
    assertions:
      - key: kind
        op: in
        value: Pod,Policy,ServiceAccount,NetworkPolicy
    tags:
      - kind

  - id: POD_CONTAINERS
    severity: FAILURE
    message: Pod must include containers
    resource: Pod
    assertions:
      - key: spec.containers
        op: present
    tags:
      - pod

  - id: POD_SECURITY_CONTEXT
    severity: FAILURE
    message: Pod should set securityContent
    resource: Pod
    assertions:
      - key: spec.securityContext.runAsNonRoot
        op: eq
        value: true
      - key: spec.securityContext.readOnlyRootFilesystem
        op: eq
        value: true
    tags:
      - pod
      - security

  - id: DEFAULT_NAMESPACE
    severity: FAILURE
    message: Policy should not use default namespace
    resource: Policy
    assertions:
      - key: spec.namespace
        op: ne
        value: default
    tags:
      - policy

  - id: NETWORK
    severity: FAILURE
    message: Network policy should include from pods
    resource: NetworkPolicy
    assertions:
      - key: spec.allowIncoming.from[].pods
        op: present
    tags:
      - network
  - id: DOCKER_REGISTRY
    severity: FAILURE
    message: Pods should pull from one of these docker registries
    resource: Pod
    assertions:
     - every:
         key: spec.containers
         expressions:
           - or:
             - key: image
               op: starts-with
               value: <private docker registry url 1>
             - key: image
               op: starts-with
               value: <private docker registry url 2>
    tags:
      - pod
```
Example for generic value assertions

```yaml

rules:
…

- id: REQ_COLOR
    message: Missing or invalid color
    severity: FAILURE
    resource: ui-form
    assertions:
      - key: color
        op: in
        value: red,blue,green

```


## kubeval

Used to validate a Kubernetes YAML file against the relevant schema.

```
Kubeval --exit-on-error 

Usage
  kubeval <file> [file...] [flags]

Flags:
  -d, --directories strings         A comma-separated list of directories to recursively search for YAML documents
      --exit-on-error               Immediately stop execution when the first error is encountered
  -f, --filename string             filename to be displayed when testing manifests read from stdin (default "stdin")
      --force-color                 Force colored output even if stdout is not a TTY
  -h, --help                        help for kubeval
      --ignore-missing-schemas      Skip validation for resource definitions without a schema
  -v, --kubernetes-version string   Version of Kubernetes to validate against (default "master")
      --openshift                   Use OpenShift schemas instead of upstream Kubernetes
  -o, --output string               The format of the output of this script. Options are: [stdout json]
      --schema-location string      Base URL used to download schemas. Can also be specified with the environment variable KUBEVAL_SCHEMA_LOCATION
      --skip-kinds strings          Comma-separated list of case-sensitive kinds to skip when validating against schemas
      --strict                      Disallow additional properties not in schema
      --version                     version for kubeval
```

### Features:

* Relies on schemas generated from the Kubernetes API, as CRD support is not available.
* Supports validation of 1..* K8s configuration files by using target directory.  
* Allows for specific versions of k8s to be specified.
* Used as part of a development workflow locally or in CI pipelines.

When using Helm, kubeval can utilize source template comments to report the relevant pass or fail output.

```
PASS - charts/templates/k8ssandra.yaml contains a valid Service.
```

## Pros / Cons

Identifying the advantages and disadvantages of each tool.

### config-lint

Pros:
* Provides customization for non-K8s configurations.

Cons: 
* TODO

### kube-score

Pros: 
* Good for scoring/covering best practices.

Cons: 
* TODO


### Kubeval

Pros: 
* Versions of Kubernetes can be specified **Major.Minor.Patch**.
* Schemas can be targeted offline for efficiency.

Cons: 
* Doesn't support CRDs with current version, but a flag is available to ignore them.

## Summary ##

* TODO

## References

* https://pkg.go.dev/github.com/stelligent/config-lint
* https://stelligent.github.io/config-lint/#/
* https://www.patricia-anong.com/blog/2019/2/28/linting-kubernetes-deployments-using-config-lint
* https://github.com/zegl/kube-score

