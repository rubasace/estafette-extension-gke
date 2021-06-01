
{{/*
Get the visibility in lowercase to match ignoring case or fail if unknown
*/}}
{{- define "estafette-application.visibilityLowerCase" -}}
{{- $.Values.ingress.visibility | lower -}}
{{- end }}

{{/*
Convenience method to validate the visibility value
*/}}
{{- define "estafette-application.validateVisibility" -}}
{{- $validVisibilities := list "private" "iap" "apigee" "public-whitelist" "public" "esp" "espv2" . -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- if not (has $visibility $validVisibilities ) -}}
{{- printf "value %s for $.Values.ingress.visibility is not valid" $.Values.ingress.visibility | fail }}
{{- end }}
{{- end }}

{{/*
Ingress properties depending on visibility
*/}}
{{- define "estafette-application.ingressAnnotations" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- if eq "private" $visibility -}}
kubernetes.io/ingress.class: "nginx-office"
{{- include "estafette-application.nginxCommonAnnotations" . }}
{{- else if eq "iap" $visibility -}}
{{- include "estafette-application.gceIngressAnnotations" . }}
{{- else if eq "apigee" $visibility -}}
kubernetes.io/ingress.class: "nginx-office"
{{- include "estafette-application.nginxCommonAnnotations" . }}
{{- else if eq "public-whitelist" $visibility -}}
kubernetes.io/ingress.class: "nginx-office"
{{- include "estafette-application.nginxCommonAnnotations" . }}
{{- end }}
{{- end }}

{{/*
Nginx common
*/}}
{{- define "estafette-application.nginxCommonAnnotations" -}}
nginx.ingress.kubernetes.io/client-body-buffer-size: {{ .Values.ingress.nginxClientBodyBufferSize | quote }}
nginx.ingress.kubernetes.io/proxy-body-size: {{ .Values.ingress.nginxProxyBodySize | quote }}
nginx.ingress.kubernetes.io/proxy-buffers-number: {{ .Values.ingress.nginxProxyBuffersNumber | quote }}
nginx.ingress.kubernetes.io/proxy-buffer-size: {{ .Values.ingress.nginxProxyBufferSize | quote }}
nginx.ingress.kubernetes.io/proxy-connect-timeout: {{ .Values.ingress.nginxProxyConnectTimeout | quote }}
nginx.ingress.kubernetes.io/proxy-send-timeout: {{ .Values.ingress.nginxProxySendTimeout | quote }}
nginx.ingress.kubernetes.io/proxy-read-timeout: {{ .Values.ingress.nginxProxyReadTimeout | quote }}
{{- if  $.Values.ingress.nginxLoadBalanceAlgorithm -}}
nginx.ingress.kubernetes.io/load-balance: {{ $.Values.ingress.nginxLoadBalanceAlgorithm | quote}}
{{- end }}
{{- end }}

{{/*
Nginx annotations
*/}}
{{- define "estafette-application.nginxSslAnnotations" -}}
{{- if eq (include "estafette-application.usesHttps" .) "true" -}}
nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
nginx.ingress.kubernetes.io/proxy-ssl-verify: "on"
{{- end -}}
{{- if $.Values.ingress.allowHTTP -}}
nginx.ingress.kubernetes.io/ssl-redirect: "false"
{{- end}}
{{- end }}

{{/*
GCE annotations
*/}}
{{- define "estafette-application.gceIngressAnnotations" -}}
kubernetes.io/ingress.class: "gce"
kubernetes.io/ingress.allow-http: "false"
{{- end }}

{{/*
Check if main ingress has to be rendered
*/}}
{{- define "estafette-application.usesMainIngress" -}}
{{/* TODO move these validations away from charts (values.schema.json???) */}}
{{- include "estafette-application.validateVisibility" . -}}
{{/*  TODO end of validations that should be moved away*/}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- has $visibility (list "private" "iap" "apigee" "public-whitelist") -}}
{{- end }}

{{/*
Check if backend-config annotation has to be enabled on the service
*/}}
{{- define "estafette-application.usesBackendConfig" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- has $visibility (list "iap") -}}
{{- end }}

{{/*
Check if access to service has to be restricted to certain IP ranges
*/}}
{{- define "estafette-application.limitsTrustedIPRanges" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- has $visibility (list "esp" "espv2" "public") -}}
{{- end }}

{{/*
Check if prometheus probe has to be enabled by default
*/}}
{{- define "estafette-application.defaultUsePrometheusProbe" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- not (has $visibility (list "esp" "espv2")) -}}
{{- end }}

{{/*
Check if prometheus probe has to be enabled by default
*/}}
{{- define "estafette-application.usesESP" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- has $visibility (list "esp" "espv2") -}}
{{- end }}

{{/*
Decide the service type based on the selected visibility
*/}}
{{- define "estafette-application.serviceType" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
{{- if eq "iap" $visibility -}}
NodePort
{{- else if has $visibility (list "esp" "espv2" "public") -}}
LoadBalancer
{{- else -}}
ClusterIP
{{- end }}
{{- end }}