{{/*
Return the proper kratos image name
*/}}
{{- define "kratos.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.kratos.image "global" .Values.global) }}
{{- end -}}
