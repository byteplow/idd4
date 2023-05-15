{{/*
Return the proper hydra image name
*/}}
{{- define "hydra.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.hydra.image "global" .Values.global) }}
{{- end -}}
