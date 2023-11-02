package pypik8s

import (
	apiv1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/jmeidam/pypi-server-local-k8s/gok8s/pkg/utils"
)

func podspec(name string, secretsObjName string, pvcName string, image string) apiv1.PodTemplateSpec {

	ports := []apiv1.ContainerPort{
		{
			Name: name,
			ContainerPort: 80,
		},
	}

	// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#envvar-v1-core
	env_variables := []apiv1.EnvVar {
		{
			Name: "PYPI_USER",
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#envvarsource-v1-core
			ValueFrom: &apiv1.EnvVarSource{
				// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#secretkeyselector-v1-core
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: secretsObjName,
					},
					Key: "username",
				},
			},
		},
		{
			Name: "PYPI_PASS",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{
						Name: secretsObjName,
					},
					Key: "password",
				},
			},
		},
	}

	podspec := apiv1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": name,
			},
		},
		// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#podspec-v1-core
		Spec: apiv1.PodSpec{
			// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#container-v1-core
			Containers: []apiv1.Container{
				{
					Name:  name,
					Image: image,
					Ports: ports,
					Env: env_variables,
					ImagePullPolicy: apiv1.PullAlways,
					SecurityContext: &apiv1.SecurityContext{
						Privileged: utils.BoolPtr(true),
					},
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "pypi-packages",
							MountPath: "/pypi-server/packages",
						},
					},
				},
			},
			Volumes: []apiv1.Volume{
				{
					Name: "pypi-packages",
					VolumeSource: apiv1.VolumeSource{
						PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
							ClaimName: pvcName,
						},
					},
				},
			},
		},
	}

	return podspec
}

func strategy() appsv1.DeploymentStrategy{
	return appsv1.DeploymentStrategy{
		Type: appsv1.DeploymentStrategyType("RollingUpdate"),
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &intstr.IntOrString{
				Type:   intstr.Type(1),
				StrVal: "25%",
			},
			MaxSurge: &intstr.IntOrString{
				Type:   intstr.Type(1),
				StrVal: "25%",
			},
		},
	}
}

func Deployment(appname string, secretsObjName string, pvcName string, image string) *appsv1.Deployment{

	dplt := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: appname,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appname,
				},
			},
			Template: podspec(appname, secretsObjName, pvcName, image),
			Strategy: strategy(),
		},
	}
	return dplt
}
