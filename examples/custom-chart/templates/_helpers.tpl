{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "k8ssandra.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Cluster name definition.
*/}}
{{- define "k8ssandra.clusterName" -}}
{{- $clusterName := lower .Values.cassandra.clusterName | replace " " "-" | replace "_" "-" }}
{{- $matchAll := mustRegexFindAll "[a-z]([-a-z0-9]*[a-z0-9])?" $clusterName -1 }}
{{- $final := join "" $matchAll }}
{{- default .Release.Name $final }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "k8ssandra.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "k8ssandra.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "k8ssandra.labels" -}}
{{- include "k8ssandra-common.labels" . -}}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "k8ssandra.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "k8ssandra.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "k8ssandra.datacenterName" -}}
{{ (index .Values.cassandra.datacenters 0).name }}
{{- end }}

{{- define "k8ssandra.datacenterSize" -}}
{{ (index .Values.cassandra.datacenters 0).size }}
{{- end }}

{{/*
Given a dict with "overrideHost" and "defaultHost", return overrideHost if it is set, and otherwise return defaultHost.
Interpret "*" and "" as meaning "match any host" and return empty string in that case.
Intended for intermediate use from other helper functions and not directly from templates; see below.
*/}}
{{- define "k8ssandra.overridableHost" -}}
  {{- if not (kindIs "invalid" .overrideHost) }}
    {{- if ne .overrideHost "*" }}
      {{- .overrideHost }}
    {{- else }}
      {{- "" }}
    {{- end }}
  {{- else if not (kindIs "invalid" .defaultHost) }}
    {{- if ne .defaultHost "*" }}
      {{- .defaultHost }}
    {{- else }}
      {{- "" }}
    {{- end }}
  {{- else }}
    {{- "" }}
  {{- end }}
{{- end }}

{{/*
Return the ingress host that should be used for Stargate's auth API.
*/}}
{{- define "k8ssandra.stargateIngressAuthHost" -}}
{{include "k8ssandra.overridableHost" (dict "defaultHost" .Values.stargate.ingress.host "overrideHost" .Values.stargate.ingress.auth.host)}}
{{- end }}

{{/*
Return the ingress host that should be used for Stargate's REST API.
*/}}
{{- define "k8ssandra.stargateIngressRestHost" -}}
{{include "k8ssandra.overridableHost" (dict "defaultHost" .Values.stargate.ingress.host "overrideHost" .Values.stargate.ingress.rest.host)}}
{{- end }}

{{/*
Return the ingress host that should be used for Stargate's GraphQL API.
*/}}
{{- define "k8ssandra.stargateIngressGraphqlHost" -}}
{{include "k8ssandra.overridableHost" (dict "defaultHost" .Values.stargate.ingress.host "overrideHost" .Values.stargate.ingress.graphql.host)}}
{{- end }}

{{/*
Return the ingress host that should be used for Stargate's Cassandra/CQL interface.
*/}}
{{- define "k8ssandra.stargateIngressCassandraHost" -}}
{{include "k8ssandra.overridableHost" (dict "defaultHost" .Values.stargate.ingress.host "overrideHost" .Values.stargate.ingress.cassandra.host)}}
{{- end }}

{{/*
Create the jvm options based on heap properties specified.
Expecting that c* heap.size and heap.newGenSize are NOT in IEC format.
*/}}
{{- define "k8ssandra.configureJvmHeap" -}}
{{- $datacenter := (index .Values.cassandra.datacenters 0) -}}
{{- if $datacenter.heap }}
  {{- if $datacenter.heap.size }}
      {{- if (regexMatch "^(([0]\\.\\d*)+|(^[1-9]\\d*)+(\\.\\d+)?)(?i)(k|m|g|e|p|t){1}$" (print $datacenter.heap.size) ) }}
        {{- nindent 6 (print "initial_heap_size: " $datacenter.heap.size) }}
        {{- nindent 6 (print "max_heap_size: " $datacenter.heap.size) }}
      {{- else }}
        {{- fail "Specify datacenter.heap.size using one of these suffixes: E, P, T, G, M, K. Format: <NUMBER>[.<NUMBER>]<SUFFIX>" }}
      {{- end }}
  {{- end }}
  {{- if $datacenter.heap.newGenSize }}
      {{- if (regexMatch "^(([0]\\.\\d*)+|(^[1-9]\\d*)+(\\.\\d+)?)(?i)(k|m|g|e|p|t){1}$" (print $datacenter.heap.newGenSize) ) }}
        {{- nindent 6 (print "heap_size_young_generation: " $datacenter.heap.newGenSize) }}
      {{- else }}
        {{- fail "Specify datacenter.heap.newGenSize using one of these suffixes: E, P, T, G, M, K. Format: <NUMBER>[.<NUMBER>]<SUFFIX>" }}
      {{- end }}
  {{- end }}
{{- else if .Values.cassandra.heap }}
  {{- if .Values.cassandra.heap.size }}
    {{- if (regexMatch "^(([0]\\.\\d*)+|(^[1-9]\\d*)+(\\.\\d+)?)(?i)(k|m|g|e|p|t){1}$" (print .Values.cassandra.heap.size)) }}
      {{- nindent 6 (print "initial_heap_size: " .Values.cassandra.heap.size) }}
      {{- nindent 6 (print "max_heap_size: " .Values.cassandra.heap.size) }}
    {{- else }}
      {{- fail "Specify cassandra.heap.size using one of these suffixes: E, P, T, G, M, K. Format: <NUMBER>[.<NUMBER>]<SUFFIX>" }}
    {{- end }}
  {{- end }}
  {{- if .Values.cassandra.heap.newGenSize }}
     {{- if (regexMatch "^(([0]\\.\\d*)+|(^[1-9]\\d*)+(\\.\\d+)?)(?i)(k|m|g|e|p|t){1}$" (print .Values.cassandra.heap.newGenSize)) }}
        {{- nindent 6 (print "heap_size_young_generation: " .Values.cassandra.heap.newGenSize) }}
     {{- else }}
        {{- fail "Specify cassandra.heap.newGenSize using one of these suffixes: E, P, T, G, M, K. Format: <NUMBER>[.<NUMBER>]<SUFFIX>" }}
     {{- end }}
  {{- end }}
{{- end }}
{{- end }}

{{/*
Set default num_tokens based on the server version
*/}}
{{- define "k8ssandra.default_num_tokens" -}}
{{- $datacenter := (index .Values.cassandra.datacenters 0) -}}
{{- if .Release.IsInstall }}
  {{- /*
  If num_tokens is set we simply use that value.
  */}}
  {{- if $datacenter.num_tokens }}
    {{- nindent 6 (print "num_tokens: " $datacenter.num_tokens) }}
  {{- else }}
  {{- /*
  If num_tokens is not set we calculate the default value based on the C* version.
  */}}
  {{- if hasPrefix "3.11" .Values.cassandra.version }}
    {{- nindent 6 (print "num_tokens: 256") }}
  {{- else }}
    {{- nindent 6 (print "num_tokens: 16") }}
  {{- end }}
{{- end }}
{{- else }}
  {{ $numTokens := "" }}
  {{ $datacenterObj := (lookup "cassandra.datastax.com/v1beta1" "CassandraDatacenter" .Release.Namespace $datacenter.name) }}
  {{- if $datacenterObj }}
    {{- if $datacenterObj.spec }}
      {{- $config := $datacenterObj.spec.config }}
      {{- $cassandraYaml := (get $config "cassandra-yaml") }}
      {{- if $cassandraYaml }}
        {{- if $cassandraYaml.num_tokens }}
          {{- $numTokens = $cassandraYaml.num_tokens }}
        {{- end }}
      {{- end }}
    {{- end }}
  {{- end }}

  {{- /*
    If upgrading and num_tokens is set simply use that value.
  */}}
  {{- if $numTokens }}
    {{- if and $datacenter.num_tokens (ne (int $datacenter.num_tokens) (int $numTokens)) }}
      {{- fail (printf "num_tokens cannot be changed once the CassandraDatacenter is created. The actual value is %d, but the specified value is %d" (int $datacenter.num_tokens) $numTokens) }}
    {{- end }}
    {{- nindent 6 (print "num_tokens: " $numTokens) }}
  {{- else }}
    {{- if $datacenter.num_tokens }}
      {{- nindent 6 (print "num_tokens: " $numTokens) }}
    {{- else }}
      {{- if hasPrefix "3.11" .Values.cassandra.version }}
        {{- nindent 6 (print "num_tokens: 256") }}
      {{- else }}
        {{- nindent 6 (print "num_tokens: 16") }}
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
{{- end }}

{{- define "medusa.configMapName" -}}
{{ .Release.Name }}-medusa
{{- end }}

{{/*
Generatea a random, alphanumeric password that is 20 characters long.
**/}}
{{- define "k8ssandra.password" -}}
{{ randAlphaNum 20 }}
{{- end }}

{{/*
Gets the superuser secret name.
*/}}
{{- define "k8ssandra.superuserSecretName" -}}
{{- if .Values.cassandra.auth.superuser.secret -}}
{{ .Values.cassandra.auth.superuser.secret }}
{{- else }}
{{- include "k8ssandra.clusterName" . | replace " " "-" | replace "_" "-" }}-superuser
{{- end }}
{{- end }}

{{/*
Gets the reaper user secret name.
*/}}
{{- define "k8ssandra.reaperUserSecretName" }}
{{- if .Values.reaper.cassandraUser.secret -}}
{{ .Values.reaper.cassandraUser.secret }}
{{- else }}
{{- include "k8ssandra.clusterName" . | replace " " "-" | replace "_" "-" }}-reaper
{{- end }}
{{- end }}

{{/*
Gets the reaper jmx user secret name.
*/}}
{{- define "k8ssandra.reaperJmxUserSecretName" }}
{{- if .Values.reaper.jmx.secret -}}
{{ .Values.reaper.jmx.secret }}
{{- else }}
{{- include "k8ssandra.clusterName" . | replace " " "-" | replace "_" "-" }}-reaper-jmx
{{- end }}
{{- end }}

{{/*
Gets the medus user secret name.
*/}}
{{- define "k8ssandra.medusaUserSecretName" }}
{{- if .Values.medusa.cassandraUser.secret -}}
{{ .Values.medusa.cassandraUser.secret }}
{{- else }}
{{- include "k8ssandra.clusterName" . | replace " " "-" | replace "_" "-" }}-medusa
{{- end }}
{{- end }}

{{/*
Gets the stargate user secret name.
*/}}
{{- define "k8ssandra.stargateUserSecretName" }}
{{- if .Values.stargate.cassandraUser.secret -}}
{{ .Values.stargate.cassandraUser.secret }}
{{- else }}
{{- include "k8ssandra.clusterName" . | replace " " "-" | replace "_" "-" }}-stargate
{{- end }}
{{- end }}

{{/*
Creates Cassandra auth environment variables if authentication is enabled.
*/}}
{{- define "medusa.cassandraAuthEnvVars" -}}
{{- if .Values.cassandra.auth.enabled }}
  {{- if .Values.medusa.cassandraUser.secret }}
    {{- nindent 10 "- name: CQL_USERNAME" }}
    {{- nindent 12 "valueFrom:" }}
    {{- nindent 14 "secretKeyRef:" }}
    {{- nindent 16 (print "name: " .Values.medusa.cassandraUser.secret) }}
    {{- nindent 16 "key: username" }}
    {{- nindent 10 "- name: CQL_PASSWORD" }}
    {{- nindent 12 "valueFrom:" }}
    {{- nindent 14 "secretKeyRef:" }}
    {{- nindent 16 (print "name: " .Values.medusa.cassandraUser.secret) }}
    {{- nindent 16 "key: password" }}
  {{- else }}
    {{- nindent 10 "- name: CQL_USERNAME" -}}
    {{- nindent 12 "valueFrom:" }}
    {{- nindent 14 "secretKeyRef:" }}
    {{- nindent 16 (print "name: " (include "k8ssandra.medusaUserSecretName" . )) }}
    {{- nindent 16 "key: username" }}
    {{- nindent 10 "- name: CQL_PASSWORD" }}
    {{- nindent 12 "valueFrom:" }}
    {{- nindent 14 "secretKeyRef:" }}
    {{- nindent 16 (print "name: " (include "k8ssandra.medusaUserSecretName" . )) }}
    {{- nindent 16 "key: password" }}
  {{- end -}}
{{- end }}
{{- end }}

{{/*
Add garbage collection settings based on the following rules in the order listed:

  1) If no GC is enabled do nothing.
  2) DC-level GC overrides cluster-level GC.
  3) If multiple garbage collectors are enabled at the DC-level, fail.
  4) If a DC-level garbage collector is enabled, use it.
  5) If multple garbage collectors are enabled at the cluster-level, fail.
  6) If a cluster-level garbage collector is enabled, use it.
*/}}
{{- define "k8ssandra.configureGc" -}}
{{- $datacenter := (index .Values.cassandra.datacenters 0) }}
{{- $clusterCmsEnabled := .Values.cassandra.gc.cms.enabled }}
{{- $clusterG1Enabled := .Values.cassandra.gc.g1.enabled }}
{{- $dcCmsEnabled := false }}
{{- $dcG1Enabled := false }}
{{- $gc := "" -}}

{{- if $datacenter.gc }}
  {{- if $datacenter.gc.cms }}
    {{ $dcCmsEnabled = $datacenter.gc.cms.enabled }}
  {{- end }}
  {{- if $datacenter.gc.g1 }}
    {{ $dcG1Enabled = $datacenter.gc.g1.enabled }}
  {{- end }}
{{- end }}

{{- $dcGcEnabled := or $dcCmsEnabled $dcG1Enabled }}
{{- $clusterGcEnabled := or $clusterCmsEnabled $clusterG1Enabled }}

{{- if and $dcCmsEnabled $dcG1Enabled }}
  {{- fail "Only one of the CMS and G1 garbage collectors can be enabled" }}
{{- end }}

{{- if not $dcGcEnabled }}
  {{- if and $clusterCmsEnabled $clusterG1Enabled }}
    {{- fail "Only one of the CMS and G1 garbage collectors can be enabled" }}
  {{- end }}
{{- end }}

{{- if $dcGcEnabled }}
  {{- if $dcCmsEnabled }}
    {{ include "k8ssandra.gcCms" $datacenter.gc.cms }}
  {{- else if $dcG1Enabled }}
    {{ include "k8ssandra.gcG1" $datacenter.gc.g1 }}
  {{- end }}
{{- else if $clusterGcEnabled }}
  {{- if $clusterCmsEnabled }}
    {{ include "k8ssandra.gcCms" .Values.cassandra.gc.cms }}
  {{- else if $clusterG1Enabled }}
    {{ include "k8ssandra.gcG1" .Values.cassandra.gc.g1 }}
  {{- end }}
{{- end }}
{{- end -}}

{{- define "k8ssandra.gcCms" -}}
{{- indent 2  "garbage_collector: CMS" -}}
{{- if .survivorRatio }}
  {{ indent 4 (cat "survivor_ratio:" .survivorRatio) }}
{{- end }}
{{- if .maxTenuringThreshold }}
  {{ indent 4 (cat "max_tenuring_threshold:" .maxTenuringThreshold) }}
{{- end }}
{{- if .initiatingOccupancyFraction }}
  {{ indent 4 (cat "cms_initiating_occupancy_fraction:" .initiatingOccupancyFraction) }}
{{- end }}
{{- if .waitDuration }}
  {{ indent 4 (cat "cms_wait_duration:" .waitDuration) }}
{{- end }}
{{- end -}}

{{- define "k8ssandra.gcG1" -}}
{{- indent 2 "garbage_collector: G1" -}}
{{- if .setUpdatingPauseTimePercent }}
  {{ indent 4 (cat "g1r_set_updating_pause_time_percent:" .setUpdatingPauseTimePercent) }}
{{- end }}
{{- if .maxGcPauseMillis }}
  {{ indent 4 (cat "max_gc_pause_millis:" .maxGcPauseMillis) }}
{{- end }}
{{- if .initiatingHeapOccupancyPercent }}
  {{ indent 4 (cat "initiating_heap_occupancy_percent:" .initiatingHeapOccupancyPercent) }}
{{- end }}
{{- if .parallelGcThreads }}
  {{ indent 4 (cat "parallel_gc_threads:" .parallelGcThreads) }}
{{- end }}
{{- if .concurrentGcThreads }}
  {{ indent 4 (cat "conc_gc_threads:" .concurrentGcThreads) }}
{{- end }}
{{- end -}}

{{/*
Cassandra image (Management API image and version tag).
*/}}
{{- define "k8ssandra.cassandraImage" -}}
{{- if .Values.cassandra.image }}
{{- include "k8ssandra-common.flattenedImage" .Values.cassandra.image -}}
{{ else }}
{{- (printf (get .Values.cassandra.versionImageMap .Values.cassandra.version)) }}
{{- end }}
{{- end }}