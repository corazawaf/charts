# coraza-spoa

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.4.0](https://img.shields.io/badge/AppVersion-0.4.0-informational?style=flat-square)

A Helm chart for Kubernetes to deploy Coraza SPOA WAF for HAProxy

**Homepage:** <https://github.com/corazawaf/charts>

## Source Code

* <https://github.com/corazawaf/coraza-spoa>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` |  |
| autoscaling.enabled | bool | `false` |  |
| autoscaling.maxReplicas | int | `4` |  |
| autoscaling.minReplicas | int | `1` |  |
| autoscaling.targetCPUUtilizationPercentage | int | `80` |  |
| autoscaling.targetMemoryUtilizationPercentage | int | `80` |  |
| config | object | `{"applications":[{"directives":"Include @coraza.conf-recommended\nInclude @crs-setup.conf.example\nInclude /etc/coraza/before.conf\nInclude @owasp_crs/*.conf\nInclude /etc/coraza/after.conf\nSecRuleEngine On\n","log_file":"/dev/stdout","log_format":"json","log_level":"error","name":"default","response_check":false,"transaction_ttl_ms":60000}],"bind":"0.0.0.0:9000","default_application":"default","log_file":"/dev/stdout","log_format":"json","log_level":"info"}` | Coraza SPOA configuration: https://github.com/corazawaf/coraza-spoa/blob/main/example/coraza-spoa.yaml |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/corazawaf/coraza-spoa"` |  |
| image.tag | string | `""` | SemVer `X.X.X` or git `sha256:digest` |
| metrics.enabled | bool | `true` |  |
| metrics.port | int | `9100` |  |
| metrics.serviceMonitor.enabled | bool | `false` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| rules | object | `{"after":"# Include this file after default CRS rules\n","before":"# Include this file before default CRS rules\n","extra":"# Use this file for additional configuration flexibility\n"}` | Custom rules that will be mounted as separate files. |
| rules.after | string | `"# Include this file after default CRS rules\n"` | These rules will be mounted at: `/etc/coraza/after.conf` |
| rules.before | string | `"# Include this file before default CRS rules\n"` | These rules will be mounted at: `/etc/coraza/before.conf` |
| rules.extra | string | `"# Use this file for additional configuration flexibility\n"` | These rules will be mounted at: `/etc/coraza/extra.conf` |
| tolerations | list | `[]` |  |
| volumeMounts | list | `[]` |  |
| volumes | list | `[]` |  |
