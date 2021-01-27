{{/*
Expand the name of the chart.
*/}}
{{- define "k8ssandra-common.name" -}}
{{/*{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}*/}}
{{ include "common.names.name" . }}
{{- end }}

{{- define "k8ssandra-common.fullname" -}}
{{ include "common.names.fullname" . }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "k8ssandra-common.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*Common Labels Standard*/}}
{{- define "common.labels.standard" -}}
helm.sh/chart: {{ include "k8ssandra-common.chart" . }}
{{ include "k8ssandra-common.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}}}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

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

{{/*
Generate a password for use in a secret. The password is a random alphanumeric 20 character
string that is base 64 encoded.
*/}}
{{- define "k8ssandra-common.password" -}}
{{ randAlphaNum 20 | b64enc | quote }}
{{- end }}
