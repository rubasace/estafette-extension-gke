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

{{/*
Calculate if the deployment has volume mounts
*/}}
{{- define "estafette-application.hasVolumeMounts" -}}
{{ or (eq (include "estafette-application.hasSecrets" . ) "true") .MountConfigmap .MountServiceAccountSecret .MountPayloadLogging .MountAdditionalVolumes }}
{{- end }}


{{/* TODO take sidecars into account (params.go:707) */}}
{{/*
Calculate if the release has secrets
*/}}
{{- define "estafette-application.hasSecrets" -}}
{{ or $.Values.deployment.secretEnvironmentVariables $.Values.secrets.data }}
{{- end }}