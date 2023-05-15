{{/*
Return the proper keto image name
*/}}
{{- define "keto.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.keto.image "global" .Values.global) }}
{{- end -}}
