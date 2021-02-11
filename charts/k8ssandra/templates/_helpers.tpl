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
{{- default .Release.Name .Values.cassandra.clusterName }}
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
*/}}
{{- define "k8ssandra.configureJvmHeap" -}}
{{- $datacenter := (index .Values.cassandra.datacenters 0) -}}
{{- if $datacenter.heap }}
  {{- if $datacenter.heap.size }}
      initial_heap_size: {{ $datacenter.heap.size }}
      max_heap_size: {{ $datacenter.heap.size }}
  {{- end }}
  {{- if $datacenter.heap.newGenSize }}
      heap_size_young_generation: {{ $datacenter.heap.newGenSize }}
  {{- end }}
{{- else if .Values.cassandra.heap }}
  {{- if .Values.cassandra.heap.size  }}
      initial_heap_size: {{ .Values.cassandra.heap.size }}
      max_heap_size: {{ .Values.cassandra.heap.size }}
  {{- end }}
  {{- if  .Values.cassandra.heap.newGenSize }}
      heap_size_young_generation: {{ .Values.cassandra.heap.newGenSize }}
  {{- end }}
{{- end }}
{{- end }}

{{/*
Set default num_tokens based on the server version
*/}}
{{- define "k8ssandra.default_num_tokens" -}}
{{- $datacenter := (index .Values.cassandra.datacenters 0) -}}
    {{- if $datacenter.num_tokens }}
      num_tokens: {{ $datacenter.num_tokens }}
    {{- else }}
    {{- if hasPrefix "3.11" .Values.cassandra.version }}
      num_tokens: 256
    {{- else }}
      num_tokens: 16
    {{- end }}
    {{- end }}
{{- end }}

{{- define "medusa.configMapName" -}}
{{ .Release.Name }}-medusa
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
    {{- nindent 16 (print "name: " (include "k8ssandra.clusterName" . ) "-medusa") }}
    {{- nindent 16 "key: username" }}
    {{- nindent 10 "- name: CQL_PASSWORD" }}
    {{- nindent 12 "valueFrom:" }}
    {{- nindent 14 "secretKeyRef:" }}
    {{- nindent 16 (print "name: " (include "k8ssandra.clusterName" . ) "-medusa") }}
    {{- nindent 16 "key: password" }}
  {{- end -}}
{{- end }}
{{- end }}

{- define "k8sandra.checkGcEnabled" -}}
{{- if .Values.cassandra.gc }}
  {{- if and .Values.cassandra.gc.cms.enabled .Values.cassandra.gc.g1.enabled }}
    {{- fail "Only one of the CMS and G1 garbage collectors can be enabled" }}
  {{- end }}
{{- end }}
{{- end }}

{{- define "k8ssandra.configureGc" -}}
{{- $datacenter := (index .Values.cassandra.datacenters 0) }}
{{- $gc := "" -}}
{{- if $datacenter.gc }}
  {{ $gc = $datacenter.gc }}
{{- else -}}
  {{ $gc = .Values.cassandra.gc }}
{{- end -}}

{{- if $gc.cms.enabled -}}
  {{- nindent 6 "garbage_collector: CMS" -}}
  {{ with $gc.cms }}
    {{- if .survivorRatio }}
      survivor_ratio: {{ .survivorRatio }}
    {{- end }}
    {{- if .maxTenuringThreshold }}
      max_tenuring_threshold: {{ .maxTenuringThreshold }}
    {{- end }}
    {{- if .initiatingOccupancyFraction }}
      cms_initiating_occupancy_fraction: {{ .initiatingOccupancyFraction }}
    {{- end }}
    {{- if .waitDuration }}
      cms_wait_duration: {{ .waitDuration }}
    {{- end }}
  {{- end }}
{{- else if $gc.g1.enabled -}}
  {{- nindent 6 "garbage_collector: G1" -}}
  {{ with $gc.g1 }}
    {{- if .setUpdatingPauseTimePercent }}
      g1r_set_updating_pause_time_percent: {{ .setUpdatingPauseTimePercent }}
    {{- end }}
    {{- if .maxGcPauseMillis }}
      max_gc_pause_millis: {{ .maxGcPauseMillis }}
    {{- end }}
    {{- if .initiatingHeapOccupancyPercent }}
      initiating_heap_occupancy_percent: {{ .initiatingHeapOccupancyPercent }}
    {{- end }}
    {{- if .parallelGcThreads }}
      parallel_gc_threads: {{ .parallelGcThreads }}
    {{- end }}
    {{- if .concurrentGcThreads }}
      conc_gc_threads: {{ .concurrentGcThreads }}
    {{- end }}
  {{- end }}
{{- end }}
{{- end }}
