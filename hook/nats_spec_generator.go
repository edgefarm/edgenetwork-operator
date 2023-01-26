package main

import (
	"leaf-nats-controller/nats"

	appsv1 "k8s.io/api/apps/v1"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

type NatsSpecGenerator struct{}

func (c *NatsSpecGenerator) GenerateManifestResponse(request *SyncRequest) (*SyncResponse, error) {
	response := &SyncResponse{}

	response.Children = append(response.Children, c.getConfigMapForNats(request))

	daemonSet := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: request.Parent.Name, // TODO: which name should it be?
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": request.Parent.Name, // TODO: which name should it be?
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":                      request.Parent.Name, // TODO: necessary? which name should it be?
						"network.edgefarm.io/type": "leaf",
						"network.edgefarm.io/name": request.Parent.Name, // TODO: add domain here?
					},
				},
				
				Spec: v1.PodSpec{
					InitContainers: []v1.Container{c.getNatsInitContainer()},

					Containers: []v1.Container{
						{
							Name:  "leaf-nats",
							Image: "nats:2.9.11-alpine",
							Ports: []v1.ContainerPort{{ContainerPort: 4111}}, // TODO: Which Ports?
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/nats",
									ReadOnly:  true,
								},
								// {
								// 	Name:      "creds",
								// 	MountPath: "/creds",
								// 	ReadOnly:  true,
								// },
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "config-origin",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "leaf-nats-config",
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
						// {
						// 	Name: "creds",
						// 	VolumeSource: v1.VolumeSource{
						// 		Secret: &v1.SecretVolumeSource{
						// 			SecretName: request.Parent.Spec.Network, // TODO: change name later! Must be generated once pushed from vault to cluster.
						// 		},
						// 	},
						// },
					},
				},
			},
		},
	}

	response.Children = append(response.Children, daemonSet)

	return response, nil
}

func (*NatsSpecGenerator) getNatsInitContainer() v1.Container {
	return v1.Container{
		Name:  "leaf-nats-init",
		Image: "bash:5.0.17",
		Command: []string{
			"/bin/sh",
			"-c",
		},
		Args: []string{
			"cp /template/nats-server.conf /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_NATS_DOMAIN/'\"$NODE_NAME\"'/g' /etc/nats/nats-server.conf",
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
		},
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "config-origin",
				MountPath: "/template",
				ReadOnly:  true,
			},
			{
				Name:      "config",
				MountPath: "/etc/nats",
			},
		},
	}
}

func (*NatsSpecGenerator) getConfigMapForNats(request *SyncRequest) *v1.ConfigMap {
	leafNatsConfig := &nats.ServerConfig{
		Listen: "localhost:4111",
		LeafNodes: nats.LeafNodesConfig{
			Remotes: []nats.LeafNodeRemoteConfig{
				{
					Url:         "nats://secret@host.minikube.internal:7422",
					Credentials: "/creds/user.creds",
				},
			},
		},
		Jetstream: nats.JetstreamConfig{
			MaxMemory: request.Parent.Spec.Limits.InMemoryStorage,
			MaxFile:   request.Parent.Spec.Limits.FileStorage,
			StoreDir:  "/data", // TODO: Need to be a path which is not in use of another JS. May be a volumeClaim.
			Domain:    "TEMPLATE_NATS_DOMAIN",
		},
	}
	leafNatsConfigString, _ := json.Marshal(leafNatsConfig)

	configMap := &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "leaf-nats-config",
		},
		Data: map[string]string{
			"nats-server.conf": string(leafNatsConfigString),
		},
	}

	return configMap
}
