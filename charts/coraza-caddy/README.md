# coraza-caddy

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 20260412](https://img.shields.io/badge/AppVersion-20260412-informational?style=flat-square)

A Helm chart for Kubernetes to deploy Coraza Caddy WAF

**Homepage:** <https://github.com/corazawaf/charts>

## Source Code

* <https://github.com/corazawaf/coraza-caddy>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity rules for pod scheduling |
| autoscaling | object | `{"enabled":false,"maxReplicas":4,"minReplicas":1,"targetCPUUtilizationPercentage":80,"targetMemoryUtilizationPercentage":80}` | Autoscaling configuration |
| autoscaling.enabled | bool | `false` | Enable autoscaling |
| autoscaling.maxReplicas | int | `4` | Maximum number of replicas |
| autoscaling.minReplicas | int | `1` | Minimum number of replicas |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | Target CPU utilization percentage |
| autoscaling.targetMemoryUtilizationPercentage | int | `80` | Target memory utilization percentage |
| caddyfile | string | See values.yaml | Caddyfile configuration |
| fullnameOverride | string | `""` | Override the full name of the chart |
| image | object | `{"pullPolicy":"IfNotPresent","repository":"ghcr.io/coreruleset/coraza-crs","tag":"caddy-alpine"}` | Image configuration |
| image.pullPolicy | string | `"IfNotPresent"` | Image pull policy |
| image.repository | string | `"ghcr.io/coreruleset/coraza-crs"` | Image repository |
| image.tag | string | `"caddy-alpine"` | Image tag. Rolling tag `caddy-alpine` always points to the latest build; use a pinned stable tag (e.g. `4-caddy-alpine-20260412`) in production. |
| imagePullSecrets | list | `[]` | Reference to one or more secrets to use for pulling images |
| initContainers | list | `[]` | Init containers to add to the pod |
| livenessProbe | object | `{"failureThreshold":3,"initialDelaySeconds":5,"periodSeconds":10,"successThreshold":1,"tcpSocket":{"port":"http"},"timeoutSeconds":5}` | Liveness probe configuration |
| metrics | object | `{"enabled":false,"port":2019,"serviceMonitor":{"enabled":false}}` | Metrics configuration |
| metrics.enabled | bool | `false` | Enable metrics endpoint (Caddy admin API) |
| metrics.port | int | `2019` | Metrics port (Caddy admin API) |
| metrics.serviceMonitor | object | `{"enabled":false}` | ServiceMonitor configuration |
| metrics.serviceMonitor.enabled | bool | `false` | Enable ServiceMonitor for Prometheus Operator |
| nameOverride | string | `""` | Override the name of the chart |
| namespaceOverride | string | `""` | Override the namespace |
| nodeSelector | object | `{}` | Node selector for pod scheduling |
| podAnnotations | object | `{}` | Annotations to add to the pod |
| podDisruptionBudget | object | `{"enabled":false}` | Pod Disruption Budget configuration |
| podDisruptionBudget.enabled | bool | `false` | Enable PodDisruptionBudget |
| podLabels | object | `{}` | Labels to add to the pod |
| podSecurityContext | object | `{}` | Pod security context |
| port | int | `8080` | HTTP port that Caddy listens on (must match the port in the Caddyfile) |
| priorityClassName | string | `""` | Priority class name for the pod |
| readinessProbe | object | `{"failureThreshold":3,"initialDelaySeconds":5,"periodSeconds":10,"successThreshold":1,"tcpSocket":{"port":"http"},"timeoutSeconds":5}` | Readiness probe configuration |
| replicaCount | int | `1` | Number of replicas |
| resources | object | `{}` | Resource requests and limits |
| securityContext | object | `{}` | Container security context |
| serviceAccount | object | `{"annotations":{},"automountServiceAccountToken":false,"create":true,"name":""}` | ServiceAccount configuration |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.automountServiceAccountToken | bool | `false` | Specifies whether to automount the service account token |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| sidecarContainers | list | `[]` | Sidecar containers to add to the pod |
| terminationGracePeriodSeconds | int | `30` | Termination grace period in seconds |
| tolerations | list | `[]` | Tolerations for pod scheduling |
| topologySpreadConstraints | list | `[]` | Topology spread constraints for pod scheduling |
| volumeMounts | list | `[]` | Additional volume mounts for the main container |
| volumes | list | `[]` | Additional volumes to add to the pod |
