{{/*
Common labels
*/}}
{{- define "estafette-application.labels" -}}
helm.sh/chart: {{ include "estafette-application.chart" . }}
{{ include "estafette-application.appSelectorLabels" . }}
app.kubernetes.io/instance: {{ include "estafette-application.nameWithTrack" . }}
{{- if $.Values.releaseData.version }}
app.kubernetes.io/version: {{ $.Values.releaseData.version | quote }}
app.kubernetes.io/managed-by: {{ $.Values.releaseData.releaseService | default .Release.Service }}
{{- end }}
{{ include "estafette-application.releaseLabels" . }}
{{ include "estafette-application.legacyLabels" . }}
{{ include "estafette-application.extraLabels" . }}
{{- end }}

{{/*
App selector labels
*/}}
{{- define "estafette-application.appSelectorLabels" -}}
app.kubernetes.io/name: {{ include "estafette-application.name" . }}
{{- if $.Values.releaseData.atomicId }}
estafette.io/atomic-id: {{ $.Values.releaseData.atomicId | quote }}
{{- end }}
{{- end }}

{{/*
Estafette custom labels
*/}}
{{- define "estafette-application.extraLabels" -}}
{{- range $key, $value := .Values.extraLabels -}}
{{ $key }}: {{ $value }}
{{- end}}
{{- end }}

{{/*
App selector labels
*/}}
{{- define "estafette-application.trackSelectorLabel" }}
{{- if $.Values.releaseData.track -}}
track: {{ $.Values.releaseData.track }}
{{- end }}
{{- end }}

{{/*
Pod specific labels
*/}}
{{- define "estafette-application.releaseLabels" }}
{{- with $.Values.releaseData -}}
{{- if .releaseId }}
estafette.io/release-id: {{ .releaseId | quote }}
{{- end }}
{{- if .triggeredBy }}
estafette.io/triggered-by: {{ .triggeredBy | quote }}
{{- end }}
{{- if .gitRepository }}
estafette.io/git-repository: {{ .gitRepository | quote }}
{{- end }}
{{- if .gitBranch }}
estafette.io/git-branch: {{ .gitBranch | quote }}
{{- end }}
{{- if .gitRevision }}
estafette.io/git-revision: {{ .gitRevision | quote }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Legacy labels that will be removed in favor of the suggested k8s ones
*/}}
{{- define "estafette-application.legacyLabels" }}
app: {{ include "estafette-application.name" . }}
{{- if $.Values.releaseData.version -}}
version: {{ $.Values.releaseData.version | quote }}
{{- end }}
{{- end }}