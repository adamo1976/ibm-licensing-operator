//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package service

import (
	"strconv"

	operatorv1alpha1 "github.com/ibm/ibm-licensing-operator/pkg/apis/operator/v1alpha1"
	res "github.com/ibm/ibm-licensing-operator/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getLicensingEnvironmentVariables(spec operatorv1alpha1.IBMLicensingSpec) []corev1.EnvVar {
	var httpsEnableString = strconv.FormatBool(spec.HTTPSEnable)
	var environmentVariables = []corev1.EnvVar{
		{
			Name:  "NAMESPACE",
			Value: spec.InstanceNamespace,
		},
		{
			Name:  "DATASOURCE",
			Value: spec.Datasource,
		},
		{
			Name:  "HTTPS_ENABLE",
			Value: httpsEnableString,
		},
	}
	if spec.IsDebug() {
		environmentVariables = append(environmentVariables, corev1.EnvVar{
			Name:  "logging.level.com.ibm",
			Value: "DEBUG",
		})
	}
	if spec.HTTPSEnable {
		environmentVariables = append(environmentVariables, corev1.EnvVar{
			Name:  "HTTPS_CERTS_SOURCE",
			Value: string(spec.HTTPSCertsSource),
		})
	}
	if spec.IsMetering() {
		environmentVariables = append(environmentVariables, corev1.EnvVar{
			Name:  "METERING_URL",
			Value: "https://metering-server:4002/api/v1/metricData",
		})
	}
	if spec.Sender != nil {
		environmentVariables = append(environmentVariables, []corev1.EnvVar{
			{
				Name:  "CLUSTER_ID",
				Value: spec.Sender.ClusterID,
			},
			{
				Name: "HUB_TOKEN",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: spec.Sender.ReporterSecretToken,
						},
						Key: ReporterSecretTokenKeyName,
					},
				},
			},
			{
				Name:  "HUB_URL",
				Value: spec.Sender.ReporterURL,
			},
			{
				Name:  "CLUSTER_NAME",
				Value: spec.Sender.ClusterName,
			},
		}...)
	}
	return environmentVariables
}

func getProbeScheme(spec operatorv1alpha1.IBMLicensingSpec) corev1.URIScheme {
	if spec.HTTPSEnable {
		return "HTTPS"
	}
	return "HTTP"
}

func getProbeHandler(spec operatorv1alpha1.IBMLicensingSpec) corev1.Handler {
	var probeScheme = getProbeScheme(spec)
	return corev1.Handler{
		HTTPGet: &corev1.HTTPGetAction{
			Path: "/",
			Port: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: licensingServicePort.IntVal,
			},
			Scheme: probeScheme,
		},
	}
}

func getMeteringSecretCheckScript() string {
	script := `while true; do
  echo "$(date): Checking for metering secret"
  ls /opt/metering/certs/* && break
  echo "$(date): Required metering secret not found ... try again in 30s"
  sleep 30
done
echo "$(date): All required secrets exist"
`
	return script
}

func GetLicensingInitContainers(spec operatorv1alpha1.IBMLicensingSpec, isOpenShift bool) []corev1.Container {
	containers := []corev1.Container{}
	if spec.IsMetering() {
		baseContainer := getLicensingContainerBase(spec, isOpenShift)
		meteringSecretCheckContainer := corev1.Container{}
		baseContainer.DeepCopyInto(&meteringSecretCheckContainer)
		meteringSecretCheckContainer.Name = "metering-check-secret"
		meteringSecretCheckContainer.Command = []string{
			"sh",
			"-c",
			getMeteringSecretCheckScript(),
		}
		containers = append(containers, meteringSecretCheckContainer)
	}
	if isOpenShift && spec.HTTPSEnable && spec.HTTPSCertsSource == operatorv1alpha1.OcpCertsSource {
		baseContainer := getLicensingContainerBase(spec, isOpenShift)
		ocpSecretCheckContainer := corev1.Container{}

		baseContainer.DeepCopyInto(&ocpSecretCheckContainer)
		ocpSecretCheckContainer.Name = res.OcpCheckString
		ocpSecretCheckContainer.Command = []string{
			"sh",
			"-c",
			res.GetOCPSecretCheckScript(),
		}
		containers = append(containers, ocpSecretCheckContainer)
	}
	return containers
}

func getLicensingContainerBase(spec operatorv1alpha1.IBMLicensingSpec, isOpenShift bool) corev1.Container {
	container := res.GetContainerBase(spec.Container)
	if spec.SecurityContext != nil {
		container.SecurityContext.RunAsUser = &spec.SecurityContext.RunAsUser
	}
	container.VolumeMounts = getLicensingVolumeMounts(spec, isOpenShift)
	container.Env = getLicensingEnvironmentVariables(spec)
	container.Ports = []corev1.ContainerPort{
		{
			ContainerPort: licensingServicePort.IntVal,
			Protocol:      corev1.ProtocolTCP,
		},
	}
	return container
}

func GetLicensingContainer(spec operatorv1alpha1.IBMLicensingSpec, isOpenShift bool) corev1.Container {
	container := getLicensingContainerBase(spec, isOpenShift)
	probeHandler := getProbeHandler(spec)
	container.Name = "license-service"
	container.LivenessProbe = res.GetLivenessProbe(probeHandler)
	container.ReadinessProbe = res.GetReadinessProbe(probeHandler)
	return container
}
