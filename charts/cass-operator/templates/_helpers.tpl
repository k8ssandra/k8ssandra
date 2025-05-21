{{- define "cass-operator.certificateName" }}
{{- if .Values.admissionWebhooks.customCertificate }}
{{- .Values.admissionWebhooks.customCertificate }}
{{- else }}
{{- printf "%s/%s-serving-cert" .Release.Namespace (include "k8ssandra-common.fullname" .) }}
{{- end }}
{{- end }}

{{- define "cass-operator.watchNamespaces" -}}
{{- if .Values.global.watchNamespaces -}}
{{ join "," .Values.global.watchNamespaces }}
{{- else -}}
{{ printf "\"\"" }}
{{- end -}}
{{- end -}}
