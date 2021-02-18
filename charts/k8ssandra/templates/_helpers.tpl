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