{{- define "cass-operator-certificate" }}
{{- if .Values.admissionWebhooks.customCertificate }}
{{- .Values.admissionWebhooks.customCertificate }}
{{- else }}
{{- printf "%s/%s-serving-cert" .Release.Namespace .Release.Name }}
{{- end }}
{{- end }}
