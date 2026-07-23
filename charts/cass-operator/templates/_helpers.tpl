{{- define "cass-operator.certificateName" }}
{{- printf "%s/%s-serving-cert" .Release.Namespace (include "k8ssandra-common.fullname" .) }}
{{- end }}

{{- define "cass-operator.webhookCertificateSecret" -}}
{{- default (printf "%s-webhook-server-cert" (include "k8ssandra-common.fullname" .)) .Values.admissionWebhooks.certificateSecret -}}
{{- end -}}

{{- define "cass-operator.watchNamespaces" -}}
{{- if .Values.global.watchNamespaces -}}
{{ join "," .Values.global.watchNamespaces }}
{{- else -}}
{{ printf "\"\"" }}
{{- end -}}
{{- end -}}
