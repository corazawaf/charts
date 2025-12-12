{{/*
Expand the name of the chart.
*/}}
{{- define "coraza-spoa.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "coraza-spoa.fullname" -}}
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
{{- define "coraza-spoa.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "coraza-spoa.labels" -}}
helm.sh/chart: {{ include "coraza-spoa.chart" . }}
{{ include "coraza-spoa.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "coraza-spoa.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coraza-spoa.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Image tag
*/}}
{{- define "coraza-spoa.imageTag" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag }}
{{- $prefix := ternary "@" ":" (hasPrefix "sha256" $tag) }}
{{- printf "%s%s" $prefix $tag }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "coraza-spoa.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "coraza-spoa.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
