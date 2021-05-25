
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
