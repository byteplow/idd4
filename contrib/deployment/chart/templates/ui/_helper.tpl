{{/*
Return the proper ui image name
*/}}
{{- define "ui.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.ui.image "global" .Values.global) }}
{{- end -}}
