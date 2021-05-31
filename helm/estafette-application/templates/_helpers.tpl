{{/*
Expand the name of the chart.
*/}}
{{- define "estafette-application.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "estafette-application.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "estafette-application.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "estafette-application.labels" -}}
helm.sh/chart: {{ include "estafette-application.chart" . }}
{{- include "estafette-application.appSelectorLabels" . -}}
{{- include "estafette-application.extraLabels" . -}}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
App selector labels
*/}}
{{- define "estafette-application.appSelectorLabels" }}
app: {{ include "estafette-application.name" . }}
{{- if $.Values.atomicId -}}
estafette.io/atomic-id: {{ $.Values.atomicId | quote }}
{{- end }}
{{- end }}

{{/*
Estafette custom labels
*/}}
{{- define "estafette-application.extraLabels" -}}
{{- range $key, $value := .Values.extraLabels }}
{{ $key }}: {{ $value }}
{{- end}}
{{- end }}

{{/*
App selector labels
*/}}
{{- define "estafette-application.trackSelectorLabel" }}
{{- if $.Values.track -}}
track: {{ $.Values.track }}
{{- end }}
{{- end }}

{{/*
Check if https is used on the container
*/}}
{{- define "estafette-application.usesHttps" -}}
{{- or ($.Values.sidecars.openresty.enabled) (eq ($.Values.deployment.containerPort | int) 443 ) -}}
{{- end }}

{{/*TODO revisit lables, it might need to take atomicId into account too?? */}}
{{/*
Generate default Hpa prometheus query
*/}}
{{- define "estafette-application.defaultHpaPromQuery" -}}
{{- printf "sum(rate(nginx_http_requests_total{app=%s}[5m])) by (app)" (include "estafette-application.name" . | quote ) -}}
{{- end }}

{{/*
Generate name with track for deployments and related manifests
*/}}
{{- define "estafette-application.nameWithTrack" -}}
{{- if $.Values.atomicId -}}
{{- $.Release.Name }}-{{ $.Values.atomicId -}}
{{- else if eq $.Values.track "stable" -}}
{{- $.Release.Name -}}-stable
{{- else if eq $.Values.track "canary" -}}
{{- $.Release.Name -}}-canary
{{- else -}}
{{- $.Release.Name -}}
{{- end }}
{{- end }}

{{/*
Generate name with track for deployments and related manifests
*/}}
{{- define "estafette-application.googleCloudCredentialsAppName" -}}
{{- default $.Release.Name $.Values.googleCloudCredentials.appName }}
{{- end }}

{{/*
Check if there are secrets to create
*/}}
{{- define "estafette-application.hasApplicationSecrets" -}}
{{- or $.Values.secrets.data $.Values.deployment.secretEnvironmentVariables -}}
{{- end }}

{{/*
Generate image pull secret .dockerconfigjson content
*/}}
{{- define "estafette-application.imagePullSecretCredentials" }}
{{- with $.Values.imagePullSecret }}
{{- printf "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"email\":\"%s\",\"auth\":\"%s\"}}}" .registry .username .password .email (printf "%s:%s" .username .password | b64enc) | b64enc }}
{{- end }}
{{- end }}