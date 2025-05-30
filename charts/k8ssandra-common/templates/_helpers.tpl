{{/*
Expand the name of the chart.
*/}}
{{- define "k8ssandra-common.name" -}}
{{ include "common.names.name" . }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "k8ssandra-common.fullname" -}}
{{ include "common.names.fullname" . }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "k8ssandra-common.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "k8ssandra-common.labels" }}
{{ include "common.labels.standard" . }}
app.kubernetes.io/part-of: k8ssandra-{{ .Release.Name }}-{{ .Release.Namespace }}
{{- $commonLabels := dict -}}
{{- if .Values.global.commonLabels -}}
{{- $commonLabels = merge $commonLabels .Values.global.commonLabels -}}
{{- end -}}
{{- if .Values.commonLabels -}}
{{- $commonLabels = merge $commonLabels .Values.commonLabels -}}
{{- end -}}
{{- if $commonLabels }}
{{ toYaml $commonLabels | trim }}
{{- end }}
{{- end }}

{{- define "k8ssandra-common.annotations" }}
{{- $context := . -}}
{{- if .context -}}
{{- $context = .context -}}
{{- end -}}
{{- $commonAnnotations := dict -}}
{{- if $context.Values.global.commonAnnotations -}}
{{- $commonAnnotations = merge $commonAnnotations $context.Values.global.commonAnnotations -}}
{{- end -}}
{{- if $context.Values.commonAnnotations -}}
{{- $commonAnnotations = merge $commonAnnotations $context.Values.commonAnnotations -}}
{{- end -}}
{{- if .annotations -}}
{{- $commonAnnotations = merge $commonAnnotations .annotations -}}
{{- end -}}
{{- if $commonAnnotations -}}
{{- toYaml $commonAnnotations | trim }}
{{- end -}}
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
Create the service account.
*/}}
{{- define "k8ssandra-common.serviceAccount" -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "k8ssandra-common.serviceAccountName" . }}
  labels: {{ include "k8ssandra-common.labels" . | indent 4 }}
  {{- $annotations := include "k8ssandra-common.annotations" (dict "context" . "annotations" .Values.serviceAccount.annotations) }}
  {{- if $annotations }}
  annotations:
    {{- $annotations | nindent 4 }}
  {{- end }}
{{- if semverCompare ">=1.24-0" .Capabilities.KubeVersion.GitVersion }}
secrets:
  - name: {{ include "k8ssandra-common.serviceAccountName" . }}-token
{{- end }}
{{- if .Values.imagePullSecrets }}
imagePullSecrets:
{{ tpl (toYaml .Values.imagePullSecrets) . }}
{{- end }}
{{- end }}

{{- define "k8ssandra-common.flattenedImage" -}}
{{- if (not (empty .)) }}
{{- if (not .repository) }}
{{- fail (print "The repository property must be defined and in scope for the flattenedImage template.") }}
{{- end }}

{{- $registry := default "docker.io" .registry }}
{{- $repository := .repository }}
{{- $tag := default "latest" .tag }}

{{- printf "%s/%s:%s" $registry $repository $tag }}
{{- end }}
{{- end }}

{{/*
Generate a password for use in a secret. The password is a random alphanumeric 20 character
string that is base 64 encoded.
*/}}
{{- define "k8ssandra-common.password" -}}
{{ randAlphaNum 20 | b64enc | quote }}
{{- end }}

