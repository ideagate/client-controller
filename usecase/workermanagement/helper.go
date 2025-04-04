package workermanagement

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func constructWorkerDeploymentSpec(appName, version string) appsv1.Deployment {
	var (
		workerPort = int32(8080)
	)

	if version == "" {
		version = "latest"
	}

	// used for startup, readiness and liveness probe
	probeHandler := corev1.ProbeHandler{
		HTTPGet: &corev1.HTTPGetAction{Path: "/", Port: intstr.IntOrString{IntVal: workerPort}},
	}

	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: appName,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": appName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  appName,
							Image: "nginx:" + version,
							Ports: []corev1.ContainerPort{{ContainerPort: workerPort}},
							StartupProbe: &corev1.Probe{
								ProbeHandler:        probeHandler,
								InitialDelaySeconds: 10,
								PeriodSeconds:       5,
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler:        probeHandler,
								InitialDelaySeconds: 10,
								PeriodSeconds:       10,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler:        probeHandler,
								InitialDelaySeconds: 5,
								PeriodSeconds:       5,
							},
						},
					},
				},
			},
		},
	}
}
