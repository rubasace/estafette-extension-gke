{{/*
Process tolerations
*/}}
{{- define "estafette-application.tolerations" -}}
{{- if $.Values.deployment.tolerations.preferPreemptibles }}
- key:      "cloud.google.com/gke-preemptible"
  operator: "Equal"
  value:    "true"
  effect:   "NoSchedule"
{{- end }}
{{- if $.Values.deployment.tolerations.useWindowsNodes }}
- key:      "node.kubernetes.io/os"
  operator: "Equal"
  value:    "windows"
  effect:   "NoSchedule"
{{- end }}
{{- if $.Values.deployment.tolerations.extraTolerations }}
{{ toYaml $.Values.deployment.tolerations.extraTolerations }}
{{- end }}
{{- end }}

{{/*
Calculate if the deployment has volume mounts
*/}}
{{- define "estafette-application.prestopSleepEnabledByDefault" -}}
{{ eq "linux" ($.Values.os | required "os is mandatory") }}
{{- end }}

{{/* TODO take sidecars into account (params.go:707) */}}
{{/*
Calculate if the release has secrets
*/}}
{{- define "estafette-application.hasSecrets" -}}
{{ or $.Values.deployment.secretEnvironmentVariables $.Values.secrets.data }}
{{- end }}

{{/*TODO only return true for deployment*/}}
{{/*
Calculate if the release has secrets
*/}}
{{- define "estafette-application.hasSslCertificate" -}}
true
{{- end }}

{{/*TODO investigate why this doesn't work*/}}
{{/*
Calculate if the release has config maps
*/}}
{{- define "estafette-application.hasConfigMaps" -}}
{{- or $.Values.configmaps.data ($.Files.Glob "externalFiles/configmaps/*") -}}
{{- end }}