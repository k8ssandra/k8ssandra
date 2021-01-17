{{/*
Expand the name of the chart.
*/}}
{{- define "k8ssandra-common.name" -}}
{{/*{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}*/}}
{{ include "common.names.name" . }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{/*{{- define "k8ssandra-common.fullname" -}}*/}}
{{/*{{- if .Values.fullnameOverride }}*/}}
{{/*{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}*/}}
{{/*{{- else }}*/}}
{{/*{{- $name := default .Chart.Name .Values.nameOverride }}*/}}
{{/*{{- if contains $name .Release.Name }}*/}}
{{/*{{- .Release.Name | trunc 63 | trimSuffix "-" }}*/}}
{{/*{{- else }}*/}}
{{/*{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}*/}}
{{/*{{- end }}*/}}
{{/*{{- end }}*/}}
{{/*{{- end }}*/}}

{{- define "k8ssandra-common.fullname" -}}
{{ include "common.names.fullname" . }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "k8ssandra-common.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{/*{{- define "k8ssandra-common.labels" -}}*/}}
{{/*helm.sh/chart: {{ include "k8ssandra-common.chart" . }}*/}}
{{/*{{ include "k8ssandra-common.selectorLabels" . }}*/}}
{{/*{{- if .Chart.AppVersion }}*/}}
{{/*app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}*/}}
{{/*{{- end }}*/}}
{{/*app.kubernetes.io/managed-by: {{ .Release.Service }}*/}}
{{/*{{- end }}*/}}

{{- define "k8ssandra-common.labels" }}
{{ include "common.labels.standard" . }}
app.kubernetes.io/part-of: k8ssandra-{{ .Release.Name }}-{{ .Release.Namespace }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "k8ssandra-common.selectorLabels" -}}
app.kubernetes.io/name: {{ include "k8ssandra-common.name" . | replace "\n" "" }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/part-of: k8ssandra-{{ .Release.Name }}-{{ .Release.Namespace }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "k8ssandra-common.serviceAccountName" -}}
{{- default (include "k8ssandra-common.fullname" .) .Values.serviceAccount.name }}
{{- end }}
