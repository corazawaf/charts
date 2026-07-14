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
	caddyHelmChartPath   = "../charts/coraza-caddy"
	caddyHelmReleaseName = "coraza-caddy"
	caddyImageRepo       = "ghcr.io/corazawaf/coraza-caddy"
	caddyImageTagSemver  = "2.5.0"
	caddyImageTagSha256  = "sha256:91bc921dc03a2fc0fe69c5ab8fdd37369869970395fc65cbb981903d64359b04"
)

func TestCaddyDeployment(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"image.tag": caddyImageTagSemver,
		},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/deployment.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	Expect(deployment.Name).To(Equal(caddyHelmReleaseName))
	Expect(*deployment.Spec.Replicas).To(Equal(int32(1)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal(fmt.Sprintf("%s:%s", caddyImageRepo, caddyImageTagSemver)))
	Expect(deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(v1.PullPolicy("IfNotPresent")))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[0]).To(Equal("caddy"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[1]).To(Equal("run"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[2]).To(Equal("--config"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[3]).To(Equal("/etc/caddy/Caddyfile"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[4]).To(Equal("--adapter"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Command[5]).To(Equal("caddyfile"))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Command)).To(Equal(6))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Name).To(Equal("http"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(8080)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Ports)).To(Equal(1))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name).To(Equal("config"))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath).To(Equal("/etc/caddy"))
	Expect(deployment.Spec.Template.Spec.Volumes[0].Name).To(Equal("config"))
	Expect(deployment.Spec.Template.Spec.Volumes[0].ConfigMap.Name).To(Equal(caddyHelmReleaseName))
}

func TestCaddyDeploymentCustom(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"namespaceOverride":                           "ingress",
			"replicaCount":                                "3",
			"image.tag":                                   caddyImageTagSha256,
			"image.pullPolicy":                            "Always",
			"port":                                        "9080",
			"volumes[0].name":                             "tmp",
			"volumeMounts[0].name":                        "tmp",
			"volumeMounts[0].mountPath":                   "/tmp",
			"initContainers[0].name":                      "init",
			"initContainers[0].image":                     "busybox",
			"initContainers[0].command[0]":                "echo",
			"initContainers[0].volumeMounts[0].name":      "tmp",
			"initContainers[0].volumeMounts[0].mountPath": "/tmp",
			"sidecarContainers[0].name":                   "foo",
			"sidecarContainers[0].image":                  "busybox",
			"sidecarContainers[0].command[0]":             "echo",
		},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/deployment.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	Expect(deployment.ObjectMeta.Namespace).To(Equal("ingress"))
	Expect(*deployment.Spec.Replicas).To(Equal(int32(3)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Image).To(Equal(fmt.Sprintf("%s@%s", caddyImageRepo, caddyImageTagSha256)))
	Expect(deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(Equal(v1.PullPolicy("Always")))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Name).To(Equal("http"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(9080)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Ports)).To(Equal(1))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[1].Name).To(Equal("tmp"))
	Expect(deployment.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath).To(Equal("/tmp"))
	Expect(deployment.Spec.Template.Spec.Containers[1].Name).To(Equal("foo"))
	Expect(deployment.Spec.Template.Spec.Containers[1].Image).To(Equal("busybox"))
	Expect(deployment.Spec.Template.Spec.Containers[1].Command[0]).To(Equal("echo"))
	Expect(deployment.Spec.Template.Spec.Volumes[1].Name).To(Equal("tmp"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].Name).To(Equal("init"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].Image).To(Equal("busybox"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].Command[0]).To(Equal("echo"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].VolumeMounts[0].Name).To(Equal("tmp"))
	Expect(deployment.Spec.Template.Spec.InitContainers[0].VolumeMounts[0].MountPath).To(Equal("/tmp"))
}

func TestCaddyDeploymentMetrics(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"metrics.enabled": "true",
		},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/deployment.yaml"})

	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].Name).To(Equal("http"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort).To(Equal(int32(8080)))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].Name).To(Equal("metrics"))
	Expect(deployment.Spec.Template.Spec.Containers[0].Ports[1].ContainerPort).To(Equal(int32(2019)))
	Expect(len(deployment.Spec.Template.Spec.Containers[0].Ports)).To(Equal(2))
}

func TestCaddyService(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/service.yaml"})

	var service v1.Service
	helm.UnmarshalK8SYaml(t, output, &service)

	Expect(service.Name).To(Equal(caddyHelmReleaseName))
	Expect(service.Spec.Ports[0].Name).To(Equal("http"))
	Expect(service.Spec.Ports[0].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(service.Spec.Ports[0].Port).To(Equal(int32(8080)))
	Expect(len(service.Spec.Ports)).To(Equal(1))
}

func TestCaddyServiceCustomPort(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"port": "9080",
		},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/service.yaml"})

	var service v1.Service
	helm.UnmarshalK8SYaml(t, output, &service)

	Expect(service.Spec.Ports[0].Port).To(Equal(int32(9080)))
	Expect(len(service.Spec.Ports)).To(Equal(1))
}

func TestCaddyServiceMetrics(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"metrics.enabled": "true",
			"metrics.port":    "2019",
		},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/service.yaml"})

	var service v1.Service
	helm.UnmarshalK8SYaml(t, output, &service)

	Expect(service.Spec.Ports[0].Name).To(Equal("http"))
	Expect(service.Spec.Ports[0].Port).To(Equal(int32(8080)))
	Expect(service.Spec.Ports[1].Name).To(Equal("metrics"))
	Expect(service.Spec.Ports[1].Protocol).To(Equal(v1.Protocol("TCP")))
	Expect(service.Spec.Ports[1].Port).To(Equal(int32(2019)))
	Expect(len(service.Spec.Ports)).To(Equal(2))
}

func TestCaddyConfigMap(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/configmap.yaml"})

	var configmap v1.ConfigMap
	helm.UnmarshalK8SYaml(t, output, &configmap)

	Expect(configmap.Name).To(Equal(caddyHelmReleaseName))
	Expect(configmap.Data["Caddyfile"]).To(ContainSubstring("order coraza_waf first"))
	Expect(configmap.Data["Caddyfile"]).To(ContainSubstring(":8080"))
	Expect(configmap.Data["Caddyfile"]).To(ContainSubstring("load_owasp_crs"))
	Expect(configmap.Data["Caddyfile"]).To(ContainSubstring("SecRuleEngine On"))
}

func TestCaddyConfigMapCustom(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"caddyfile": ":9080 {\n  respond \"Hello!\"\n}\n",
		},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/configmap.yaml"})

	var configmap v1.ConfigMap
	helm.UnmarshalK8SYaml(t, output, &configmap)

	Expect(configmap.Name).To(Equal(caddyHelmReleaseName))
	Expect(configmap.Data["Caddyfile"]).To(ContainSubstring(":9080"))
}

func TestCaddyHorizontalPodAutoscaler(t *testing.T) {
	RegisterTestingT(t)

	options := &helm.Options{
		SetValues: map[string]string{
			"autoscaling.enabled": "true",
		},
	}

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/hpa.yaml"})

	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(t, output, &hpa)

	Expect(hpa.Name).To(Equal(caddyHelmReleaseName))
	Expect(*hpa.Spec.MinReplicas).To(Equal(int32(1)))
	Expect(hpa.Spec.MaxReplicas).To(Equal(int32(4)))
	Expect(*hpa.Spec.Metrics[0].Resource.Target.AverageUtilization).To(Equal(int32(80)))
	Expect(*hpa.Spec.Metrics[1].Resource.Target.AverageUtilization).To(Equal(int32(80)))
}

func TestCaddyHorizontalPodAutoscalerCustom(t *testing.T) {
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

	output := helm.RenderTemplate(t, options, caddyHelmChartPath, caddyHelmReleaseName, []string{"templates/hpa.yaml"})

	var hpa autoscalingv2.HorizontalPodAutoscaler
	helm.UnmarshalK8SYaml(t, output, &hpa)

	Expect(hpa.Name).To(Equal(caddyHelmReleaseName))
	Expect(*hpa.Spec.MinReplicas).To(Equal(int32(2)))
	Expect(hpa.Spec.MaxReplicas).To(Equal(int32(6)))
	Expect(*hpa.Spec.Metrics[0].Resource.Target.AverageUtilization).To(Equal(int32(75)))
	Expect(*hpa.Spec.Metrics[1].Resource.Target.AverageUtilization).To(Equal(int32(75)))
}
