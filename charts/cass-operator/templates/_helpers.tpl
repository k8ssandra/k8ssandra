{{- define "cass-operator-certificate" }}
{{- if .Values.admissionWebhooks.customCertificate }}
{{- .Values.admissionWebhooks.customCertificate }}
{{- else }}
{{- printf "%s/%s-serving-cert" .Release.Namespace (include "k8ssandra-common.fullname" .) }}
{{- end }}
{{- end }}
