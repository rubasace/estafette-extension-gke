{{/*
Calculate the service port name
*/}}
{{- define "estafette-application.defaultTlsCertificateSecretName" -}}
{{- if $.Values.ingress.useWildcardCertificate -}}
{{ $.Values.ingress.wildCardTlsSecretName }}
{{- else -}}
{{ $.Release.Name }}-letsencrypt-certificate
{{- end }}
{{- end }}

{{/*
Return hostnames
*/}}
{{- define "estafette-application.hostnames" -}}
{{- if eq (include "estafette-application.usesMainIngress" .) "true" -}}
{{- ternary (list (printf "%s.%s.travix.com" $.Release.Name $.Values.environment)) $.Values.ingress.hosts ($.Values.ingress.hosts | empty) | join "," -}}
{{- end }}
{{- end }}
{{/*

{{/*
Return internal hostnames
*/}}
{{- define "estafette-application.internalHostnames" -}}
{{- ternary (list (printf "%s.%s.internal.travix.io" $.Release.Name $.Values.environment)) $.Values.ingress.internalHosts ($.Values.ingress.internalHosts | empty) | join "," -}}
{{- end }}
{{/*

{{/*
Return apigee hostnames
*/}}
{{- define "estafette-application.apigeeHostnames" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- if eq $visibility "apigee" -}}
{{- ternary (list (printf "%s-apigee.%s.travix.com" $.Release.Name $.Values.environment)) $.Values.ingress.apigeeHosts ($.Values.ingress.apigeeHosts | empty) | join "," -}}
{{- end }}
{{- end }}
{{/*

Respect provided secretName for tls or, otherwise, provide default name
*/}}
{{- define "estafette-application.tlsCertificateSecretName" -}}
{{- $.Values.ingress.tlsSecretName | default (include "estafette-application.defaultTlsCertificateSecretName" .) -}}
{{- end }}

{{/*
Calculate the service port name
*/}}
{{- define "estafette-application.servicePortName" -}}
{{- if eq (include "estafette-application.usesHttps" .) "true" -}}
https
{{- else -}}
web
{{- end }}
{{- end }}
