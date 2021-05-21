
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
{{- include "estafette-application.nginxIngressAnnotations" . }}
{{ include "estafette-application.dnsIngressAnnotations" . }}
{{- else if eq "iap" $visibility -}}
{{- include "estafette-application.gceIngressAnnotations" . }}
{{ include "estafette-application.dnsIngressAnnotations" . }}
{{- else if eq "apigee" $visibility -}}
kubernetes.io/ingress.class: "nginx-office"
{{- include "estafette-application.nginxIngressAnnotations" . }}
{{ include "estafette-application.dnsIngressAnnotations" . }}
{{- else if eq "public-whitelist" $visibility -}}
kubernetes.io/ingress.class: "nginx-office"
{{- include "estafette-application.nginxIngressAnnotations" . }}
{{ include "estafette-application.dnsIngressAnnotations" . }}
{{- end }}
{{- end }}

{{/*
Nginx annotations
*/}}
{{- define "estafette-application.nginxIngressAnnotations" -}}
{{- if or ($.Values.sidecars.openresty.enabled) (eq ($.Values.deployment.containerPort | int) 443 ) }}
nginx.ingress.kubernetes.io/backend-protocol: "HTTPS"
nginx.ingress.kubernetes.io/proxy-ssl-verify: "on"
{{- end -}}
{{/*TODO cannot be set for apigee*/}}
{{- if $.Values.ingress.allowHTTP -}}
nginx.ingress.kubernetes.io/ssl-redirect: "false"
{{- end}}
nginx.ingress.kubernetes.io/client-body-buffer-size: {{ .Values.ingress.nginxClientBodyBufferSize }}
nginx.ingress.kubernetes.io/proxy-body-size: {{ .Values.ingress.nginxProxyBodySize }}
nginx.ingress.kubernetes.io/proxy-buffers-number: {{ .Values.ingress.nginxProxyBuffersNumber }}
nginx.ingress.kubernetes.io/proxy-buffer-size: {{ .Values.ingress.nginxProxyBufferSize }}
nginx.ingress.kubernetes.io/proxy-connect-timeout: {{ .Values.ingress.nginxProxyConnectTimeout }}
nginx.ingress.kubernetes.io/proxy-send-timeout: {{ .Values.ingress.nginxProxySendTimeout }}
nginx.ingress.kubernetes.io/proxy-read-timeout: {{ .Values.ingress.nginxProxyReadTimeout }}
{{/*TODO check if mandatory for public-whitelist*/}}
{{- if $.Values.ingress.nginxWhitelistedIPS -}}
nginx.ingress.kubernetes.io/whitelist-source-range: {{ $.Values.ingress.nginxWhitelistedIPS | join "," }}
{{- end}}
{{- if  $.Values.ingress.nginxLoadBalanceAlgorithm -}}
nginx.ingress.kubernetes.io/load-balance: {{ $.Values.ingress.nginxLoadBalanceAlgorithm}}
{{- end }}
{{- end }}

{{/*
GCE annotations
*/}}
{{- define "estafette-application.gceIngressAnnotations" -}}
kubernetes.io/ingress.class: gce
kubernetes.io/ingress.allow-http: false
{{- end }}

{{/*
DNS annotations
*/}}
{{- define "estafette-application.dnsIngressAnnotations" -}}
{{- $visibility := include "estafette-application.visibilityLowerCase" . -}}
estafette.io/cloudflare-dns: true
estafette.io/cloudflare-proxy: {{ ne $visibility "iap" }}
{{- if not $.Values.ingress.cloudflareHostnames -}}
{{ fail "At least one cloudflare hostname has to be provided" }}
{{- end }}
estafette.io/cloudflare-hostnames: "{{ $.Values.ingress.cloudflareHostnames | join "," }}"
{{- end }}