package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	v1 "k8s.io/api/core/v1"
)

var (
	spoaHelmChartPath   = "../charts/coraza-spoa"
	spoaHelmReleaseName = "coraza-spoa"
	spoaImageRepo       = "ghcr.io/corazawaf/coraza-spoa"
	spoaImageTagSemver  = "0.4.0"
	spoaImageTagSha256  = "sha256:91bc921dc03a2fc0fe69c5ab8fdd37369869970395fc65cbb981903d64359b04"
)

func TestSpoaDeployment(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"image.tag": spoaImageTagSemver,
		},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/deployment.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	Expect(deployment.Name).To(Equal(spoaHelmReleaseName))
	Expect(*deployment.Spec.Replicas).To(Equal(int32(1)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal(fmt.Sprintf("%s:%s", spoaImageRepo, spoaImageTagSemver)))
	Expect(deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(v1.PullPolicy("IfNotPresent")))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[0]).To(Equal("/coraza-spoa"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[1]).To(Equal("--config=/etc/coraza/config.yaml"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[2]).To(Equal("--autoreload"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[3]).To(Equal("--metrics-addr=:9100"))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Command)).To(Equal(4))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Name).To(Equal("spoa"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(9000)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].Name).To(Equal("metrics"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort).To(Equal(int32(9100)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Ports)).To(Equal(2))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("config"))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/etc/coraza"))
	Expect(deployment.Spec.Template.Spec.Volumes[0].Name).To(Equal("config"))
	Expect(deployment.Spec.Template.Spec.Volumes[0].ConfigMap.Name).To(Equal(spoaHelmReleaseName))
}

func TestSpoaDeploymentCustom(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"namespaceOverride":                           "ingress",
			"replicaCount":                                "3",
			"image.tag":                                   spoaImageTagSha256,
			"image.pullPolicy":                            "Always",
			"metrics.port":                                "9200",
			"volumes[0].name":                             "crs",
			"volumeMounts[0].name":                        "crs",
			"volumeMounts[0].mountPath":                   "/etc/crs",
			"initContainers[0].name":                      "crs",
			"initContainers[0].image":                     "busybox",
			"initContainers[0].command[0]":                "echo",
			"initContainers[0].volumeMounts[0].name":      "crs",
			"initContainers[0].volumeMounts[0].mountPath": "/etc/crs",
			"sidecarContainers[0].name":                   "foo",
			"sidecarContainers[0].image":                  "busybox",
			"sidecarContainers[0].command[0]":             "echo",
		},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/deployment.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	Expect(deployment.ObjectMeta.Namespace).To(Equal("ingress"))
	Expect(*deployment.Spec.Replicas).To(Equal(int32(3)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal(fmt.Sprintf("%s@%s", spoaImageRepo, spoaImageTagSha256)))
	Expect(deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(v1.PullPolicy("Always")))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[0]).To(Equal("/coraza-spoa"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[1]).To(Equal("--config=/etc/coraza/config.yaml"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[2]).To(Equal("--autoreload"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[3]).To(Equal("--metrics-addr=:9200"))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Command)).To(Equal(4))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Name).To(Equal("spoa"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(9000)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].Name).To(Equal("metrics"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort).To(Equal(int32(9200)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Ports)).To(Equal(2))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[1].Name).To(Equal("crs"))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath).To(Equal("/etc/crs"))
	Expect(deployment.Spec.Template.Spec.Containers[1].Name).To(Equal("foo"))
	Expect(deployment.Spec.Template.Spec.Containers[1].Image).To(Equal("busybox"))
	Expect(deployment.Spec.Template.Spec.Containers[1].Command[0]).To(Equal("echo"))
	Expect(deployment.Spec.Template.Spec.Volumes[1].Name).To(Equal("crs"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].Name).To(Equal("crs"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].Image).To(Equal("busybox"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].Command[0]).To(Equal("echo"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].VolumeMounts[0].Name).To(Equal("crs"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].VolumeMounts[0].MountPath).To(Equal("/etc/crs"))
}

func TestSpoaService(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/service.yaml"})

	var service v1.Service
	helm.UnmarshalK8SYaml(t, output, &service)

	Expect(service.Name).To(Equal(spoaHelmReleaseName))
	Expect(service.Spec.Ports[0].Name).To(Equal("spoa"))
	Expect(service.Spec.Ports[0].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(service.Spec.Ports[0].Port).To(Equal(int32(9000)))
	Expect(service.Spec.Ports[1].Name).To(Equal("metrics"))
	Expect(service.Spec.Ports[1].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(service.Spec.Ports[1].Port).To(Equal(int32(9100)))
	Expect(len(service.Spec.Ports)).To(Equal(2))
}

func TestSpoaServiceCustomPorts(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"config.bind":  "0.0.0.0:9300",
			"metrics.port": "9200",
		},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/service.yaml"})

	var service v1.Service
	helm.UnmarshalK8SYaml(t, output, &service)

	Expect(service.Spec.Ports[0].Port).To(Equal(int32(9300)))
	Expect(service.Spec.Ports[1].Port).To(Equal(int32(9200)))
	Expect(len(service.Spec.Ports)).To(Equal(2))
}

func TestSpoaServiceNoMetrics(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"metrics.enabled": "false",
		},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/service.yaml"})

	var service v1.Service
	helm.UnmarshalK8SYaml(t, output, &service)

	Expect(len(service.Spec.Ports)).To(Equal(1))
}

func TestSpoaConfigMap(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/configmap.yaml"})

	var configmap v1.ConfigMap
	helm.UnmarshalK8SYaml(t, output, &configmap)

	Expect(configmap.Name).To(Equal(spoaHelmReleaseName))
	Expect(configmap.Data["config.yaml"]).To(HavePrefix("applications:\n"))
	Expect(configmap.Data["config.yaml"]).To(HaveSuffix("log_level: info\n"))
	Expect(configmap.Data["before.conf"]).To(Equal("# Include this file before default CRS rules\n"))
	Expect(configmap.Data["after.conf"]).To(Equal("# Include this file after default CRS rules\n"))
	Expect(configmap.Data["extra.conf"]).To(Equal("# Use this file for additional configuration flexibility"))
}

func TestSpoaConfigMapCustom(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"config.log_level": "error",
			"rules.before":     "before",
			"rules.after":      "after",
			"rules.extra":      "extra",
		},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/configmap.yaml"})

	var configmap v1.ConfigMap
	helm.UnmarshalK8SYaml(t, output, &configmap)

	Expect(configmap.Name).To(Equal(spoaHelmReleaseName))
	Expect(configmap.Data["config.yaml"]).To(ContainSubstring("log_level: error\n"))
	Expect(configmap.Data["before.conf"]).To(Equal("before\n"))
	Expect(configmap.Data["after.conf"]).To(Equal("after\n"))
	Expect(configmap.Data["extra.conf"]).To(Equal("extra"))
}

func TestSpoaHorizontalPodAutoscaler(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled": "true",
		},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/hpa.yaml"})

	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(t, output, &hpa)

	Expect(hpa.Name).To(Equal(spoaHelmReleaseName))
	Expect(*hpa.Spec.MinReplicas).To(Equal(int32(1)))
	Expect(hpa.Spec.MaxReplicas).To(Equal(int32(4)))
	Expect(*hpa.Spec.Metrics[0].Resource.Target.AverageUtilization).To(Equal(int32(80)))
	Expect(*hpa.Spec.Metrics[1].Resource.Target.AverageUtilization).To(Equal(int32(80)))
}

func TestSpoaHorizontalPodAutoscalerCustom(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled":                           "true",
			"autoscaling.minReplicas":                       "2",
			"autoscaling.maxReplicas":                       "6",
			"autoscaling.targetCPUUtilizationPercentage":    "75",
			"autoscaling.targetMemoryUtilizationPercentage": "75",
		},
	}

	output := helm.RenderTemplate(t, options, spoaHelmChartPath, spoaHelmReleaseName, []string{"templates/hpa.yaml"})

	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(t, output, &hpa)

	Expect(hpa.Name).To(Equal(spoaHelmReleaseName))
	Expect(*hpa.Spec.MinReplicas).To(Equal(int32(2)))
	Expect(hpa.Spec.MaxReplicas).To(Equal(int32(6)))
	Expect(*hpa.Spec.Metrics[0].Resource.Target.AverageUtilization).To(Equal(int32(75)))
	Expect(*hpa.Spec.Metrics[1].Resource.Target.AverageUtilization).To(Equal(int32(75)))
}
