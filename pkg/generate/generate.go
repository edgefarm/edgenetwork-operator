package generate

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"

	v1alpha1 "github.com/edgefarm/edgenetwork-operator/apis/edgenetwork/v1alpha1"
	"github.com/edgefarm/edgenetwork-operator/pkg/nats"
)

// TODO: Generate YurtAppSet -> Deployment
func Manifests(config *v1alpha1.EdgeNetwork) ([]runtime.Object, error) {
	response := []runtime.Object{}
	cm, err := getConfigMapForNats(config)
	if err != nil {
		return nil, err
	}
	response = append(response, cm)

	daemonSet := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Spec.Network,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": config.Spec.Network,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":                      config.Spec.Network,
						"network.edgefarm.io/type": "leaf",
						"network.edgefarm.io/name": config.Spec.Network,
					},
				},

				Spec: v1.PodSpec{
					InitContainers: []v1.Container{getNatsInitContainer(config)},
					Containers:     []v1.Container{getNatsContainer()},
					Volumes: []v1.Volume{
						{
							Name: "config-template",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: fmt.Sprintf("%s-leaf-nats-config", config.Spec.Network),
									},
								},
							},
						},
						{
							Name: "config",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "creds",
							VolumeSource: v1.VolumeSource{
								Secret: &v1.SecretVolumeSource{
									SecretName: config.Spec.ConnectionSecretRef.Name,
								},
							},
						},
						{
							Name: "data",
							VolumeSource: v1.VolumeSource{
								EmptyDir: &v1.EmptyDirVolumeSource{
									SizeLimit: func() *resource.Quantity {
										q := resource.MustParse(config.Spec.Limits.FileStorage)
										return &q
									}(),
								},
							},
						},
					},
					Affinity: &v1.Affinity{
						NodeAffinity: &v1.NodeAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
								NodeSelectorTerms: []v1.NodeSelectorTerm{
									config.Spec.NodeSelectorTerm,
								},
							},
						},
					},
					Tolerations: config.Spec.Tolerations,
				},
			},
		},
	}

	response = append(response, daemonSet)

	return response, nil
}

func getNatsInitContainer(config *v1alpha1.EdgeNetwork) v1.Container {
	return v1.Container{
		Name:  "init",
		Image: "nats:2.9.11-alpine",
		Command: []string{
			"/bin/sh",
			"-c",
		},
		Args: []string{
			"cp /template/nats-server.conf /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_NATS_DOMAIN/'\"$NODE_NAME\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_NETWORK_NAME/'\"$NETWORK_NAME\"'/g' /etc/nats/nats-server.conf",
		},
		Env: []v1.EnvVar{
			{
				Name: "NODE_NAME",
				ValueFrom: &v1.EnvVarSource{
					FieldRef: &v1.ObjectFieldSelector{
						FieldPath: "spec.nodeName",
					},
				},
			},
			{
				Name:  "NETWORK_NAME",
				Value: config.Spec.Network,
			},
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "config-template",
				MountPath: "/template",
				ReadOnly:  true,
			},
			{
				Name:      "config",
				MountPath: "/etc/nats",
				ReadOnly:  false,
			},
		},
	}
}

func getNatsContainer() v1.Container {
	return v1.Container{
		Name:  "nats",
		Image: "nats:2.9.11-alpine",
		Ports: []v1.ContainerPort{{ContainerPort: 4222}},
		// Command: []string{"/bin/sh", "-c", "--"},
		// Args:    []string{"while true; do sleep 30; done;"},
		Args: []string{"-c", "/etc/nats/nats-server.conf"},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "config",
				MountPath: "/etc/nats",
				ReadOnly:  true,
			},
			{
				Name:      "creds",
				MountPath: "/creds",
				ReadOnly:  true,
			},
			{
				Name:      "data",
				MountPath: "/data",
				ReadOnly:  false,
			},
		},
	}
}

func getConfigMapForNats(config *v1alpha1.EdgeNetwork) (*v1.ConfigMap, error) {
	leafNatsConfig := &nats.ServerConfig{
		Listen: "localhost:4222",
		LeafNodes: nats.LeafNodesConfig{
			Remotes: []nats.LeafNodeRemoteConfig{
				{
					Url:         fmt.Sprintf("nats://%s:7422", config.Spec.Address),
					Credentials: "/creds/creds",
				},
			},
		},
		Jetstream: nats.JetstreamConfig{
			MaxMemory: config.Spec.Limits.InMemoryStorage,
			MaxFile:   config.Spec.Limits.FileStorage,
			StoreDir:  "/data",
			Domain:    "TEMPLATE_NETWORK_NAME-TEMPLATE_NATS_DOMAIN",
		},
	}
	leafNatsConfigString, err := json.Marshal(leafNatsConfig)
	if err != nil {
		return nil, err
	}

	configMap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-leaf-nats-config", config.Spec.Network),
		},
		Data: map[string]string{
			"nats-server.conf": string(leafNatsConfigString),
		},
	}

	return configMap, nil
}
