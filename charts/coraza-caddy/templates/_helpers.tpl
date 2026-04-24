{{/*
Expand the name of the chart.
*/}}
{{- define "coraza-caddy.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "coraza-caddy.fullname" -}}
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
{{- define "coraza-caddy.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "coraza-caddy.labels" -}}
helm.sh/chart: {{ include "coraza-caddy.chart" . }}
{{ include "coraza-caddy.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "coraza-caddy.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coraza-caddy.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Image tag
*/}}
{{- define "coraza-caddy.imageTag" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag }}
{{- $prefix := ternary "@" ":" (hasPrefix "sha256" $tag) }}
{{- printf "%s%s" $prefix $tag }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "coraza-caddy.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "coraza-caddy.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
