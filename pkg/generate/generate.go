package generate

import (
	"fmt"

	yurtv1alpha1 "github.com/edgefarm/edgenetwork-operator/pkg/api/yurtappdaemon"
	appsv1 "k8s.io/api/apps/v1"

	v1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	v1alpha1 "github.com/edgefarm/edgenetwork-operator/apis/edgenetwork/v1alpha1"
	json "github.com/edgefarm/edgenetwork-operator/pkg/json"
	"github.com/edgefarm/edgenetwork-operator/pkg/nats"
)

func Manifests(config *v1alpha1.EdgeNetwork) ([]runtime.Object, error) {
	response := []runtime.Object{}
	cm, err := getConfigMapForNats(config)
	if err != nil {
		return nil, err
	}
	response = append(response, cm)
	name := fmt.Sprintf("%s-%s", config.Spec.Network, config.Spec.SubNetwork)

	service := getService(config)
	response = append(response, service)

	yurtAppDaemon := &yurtv1alpha1.YurtAppDaemon{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps.openyurt.io/v1alpha1",
			Kind:       "YurtAppDaemon",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: yurtv1alpha1.YurtAppDaemonSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"network.edgefarm.io/type":                                               "leaf",
					fmt.Sprintf("name.network.edgefarm.io/%s", config.Spec.Network):          "",
					fmt.Sprintf("subnetwork.network.edgefarm.io/%s", config.Spec.SubNetwork): "",
				},
			},
			WorkloadTemplate: yurtv1alpha1.WorkloadTemplate{
				DeploymentTemplate: &yurtv1alpha1.DeploymentTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"network.edgefarm.io/type":                                               "leaf",
							fmt.Sprintf("name.network.edgefarm.io/%s", config.Spec.Network):          "",
							fmt.Sprintf("subnetwork.network.edgefarm.io/%s", config.Spec.SubNetwork): "",
						},
					},
					Spec: appsv1.DeploymentSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"network.edgefarm.io/type":                                               "leaf",
								fmt.Sprintf("name.network.edgefarm.io/%s", config.Spec.Network):          "",
								fmt.Sprintf("subnetwork.network.edgefarm.io/%s", config.Spec.SubNetwork): "",
							},
						},
						Template: v1.PodTemplateSpec{
							ObjectMeta: metav1.ObjectMeta{
								Labels: map[string]string{
									"network.edgefarm.io/type":                                               "leaf",
									fmt.Sprintf("name.network.edgefarm.io/%s", config.Spec.Network):          "",
									fmt.Sprintf("subnetwork.network.edgefarm.io/%s", config.Spec.SubNetwork): "",
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
										Name: "system-user-creds",
										VolumeSource: v1.VolumeSource{
											Secret: &v1.SecretVolumeSource{
												SecretName: config.Spec.ConnectionSecretRefs.SystemUserSecretRef.Name,
												Items: []v1.KeyToPath{{
													Key:  "creds",
													Path: "creds",
												}},
											},
										},
									},
									{
										Name: "system-account-user-creds",
										VolumeSource: v1.VolumeSource{
											Secret: &v1.SecretVolumeSource{
												SecretName: config.Spec.ConnectionSecretRefs.SysAccountUserSecretRef.Name,
												Items: []v1.KeyToPath{{
													Key:  "creds",
													Path: "creds",
												}},
											},
										},
									},
									{
										Name: "system-account-jwt",
										VolumeSource: v1.VolumeSource{
											Secret: &v1.SecretVolumeSource{
												SecretName: config.Spec.ConnectionSecretRefs.SysAccountUserSecretRef.Name,
												Items: []v1.KeyToPath{{
													Key:  "system-account-jwt",
													Path: "system-account-jwt",
												}},
											},
										},
									},
									{
										Name: "system-account-public-key",
										VolumeSource: v1.VolumeSource{
											Secret: &v1.SecretVolumeSource{
												SecretName: config.Spec.ConnectionSecretRefs.SysAccountUserSecretRef.Name,
												Items: []v1.KeyToPath{{
													Key:  "system-account-public-key",
													Path: "system-account-public-key",
												}},
											},
										},
									},
									{
										Name: "operator-jwt",
										VolumeSource: v1.VolumeSource{
											Secret: &v1.SecretVolumeSource{
												SecretName: config.Spec.ConnectionSecretRefs.SysAccountUserSecretRef.Name,
												Items: []v1.KeyToPath{{
													Key:  "operator-jwt",
													Path: "operator-jwt",
												}},
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
									{
										Name: "jwt",
										VolumeSource: v1.VolumeSource{
											EmptyDir: &v1.EmptyDirVolumeSource{},
										},
									},
								},
								Tolerations: getTolerations(config),
							},
						},
					},
				},
			},
			NodePoolSelector: config.Spec.NodePoolSelector,
		},
	}

	response = append(response, yurtAppDaemon)

	return response, nil
}

func getTolerations(config *v1alpha1.EdgeNetwork) []v1.Toleration {
	ret := []v1.Toleration{
		{
			Key:      "edgefarm.io",
			Operator: "Exists",
			Effect:   "NoSchedule",
		},
	}

	ret = append(ret, config.Spec.Tolerations...)
	return ret
}

func getService(config *v1alpha1.EdgeNetwork) *v1.Service {
	name := fmt.Sprintf("%s-%s", config.Spec.Network, config.Spec.SubNetwork)
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"network.edgefarm.io/type":                                               "leaf",
				fmt.Sprintf("name.network.edgefarm.io/%s", config.Spec.Network):          "",
				fmt.Sprintf("subnetwork.network.edgefarm.io/%s", config.Spec.SubNetwork): "",
			},
			Annotations: map[string]string{
				"openyurt.io/topologyKeys": "kubernetes.io/hostname",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "nats",
					Port:       4222,
					TargetPort: intstr.FromInt(4222),
					Protocol:   v1.ProtocolTCP,
				},
				{
					Name:       "nats-metrics",
					Port:       8222,
					TargetPort: intstr.FromInt(8222),
					Protocol:   v1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"network.edgefarm.io/type":                                               "leaf",
				fmt.Sprintf("name.network.edgefarm.io/%s", config.Spec.Network):          "",
				fmt.Sprintf("subnetwork.network.edgefarm.io/%s", config.Spec.SubNetwork): "",
			},
			Type: v1.ServiceTypeClusterIP,
		},
	}
}

func getNatsInitContainer(config *v1alpha1.EdgeNetwork) v1.Container {
	return v1.Container{
		Name:            "init",
		Image:           "nats:2.9.11-alpine",
		ImagePullPolicy: v1.PullIfNotPresent,
		Command: []string{
			"/bin/sh",
			"-c",
		},
		Args: []string{
			"cp /template/nats-server.conf /etc/nats/nats-server.conf && export ADDRESS_NO_PORT=$(echo $ADDRESS | awk -F \":\" '{print $1\":\"$2}') && sed -i 's/TEMPLATE_NODE_NAME/'\"$NODE_NAME\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_NETWORK/'\"$NETWORK\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_SUB_NETWORK/'\"$SUB_NETWORK\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_NETWORK/'\"$NETWORK\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_OPERATOR_JWT/'\"$OPERATOR_JWT\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_SYS_ACCOUNT_JWT/'\"$SYS_ACCOUNT_JWT\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_SYS_ACCOUNT_PUBLIC_KEY/'\"$SYS_ACCOUNT_PUBLIC_KEY\"'/g' /etc/nats/nats-server.conf && sed -i 's/TEMPLATE_ACCOUNT_PUBLIC_KEY/'\"$ACCOUNT_PUBLIC_KEY\"'/g' /etc/nats/nats-server.conf && sed -i 's#TEMPLATE_ADDRESS#'$ADDRESS_NO_PORT'#g' /etc/nats/nats-server.conf",
		},
		// Args: []string{
		// 	"sleep 1000000000",
		// },
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
				Name:  "NETWORK",
				Value: config.Spec.Network,
			},
			{
				Name:  "SUB_NETWORK",
				Value: config.Spec.SubNetwork,
			},
			{
				Name: "OPERATOR_JWT",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: config.Spec.ConnectionSecretRefs.SysAccountUserSecretRef.Name,
						},
						Key: "operator-jwt",
					},
				},
			},
			{
				Name: "SYS_ACCOUNT_JWT",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: config.Spec.ConnectionSecretRefs.SysAccountUserSecretRef.Name,
						},
						Key: "sys-account-jwt",
					},
				},
			},
			{
				Name: "SYS_ACCOUNT_PUBLIC_KEY",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: config.Spec.ConnectionSecretRefs.SysAccountUserSecretRef.Name,
						},
						Key: "sys-account-public-key",
					},
				},
			},
			{
				Name: "ACCOUNT_PUBLIC_KEY",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: config.Spec.ConnectionSecretRefs.SystemUserSecretRef.Name,
						},
						Key: "account-public-key",
					},
				},
			},
			{
				Name: "ADDRESS",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: config.Spec.ConnectionSecretRefs.SystemUserSecretRef.Name,
						},
						Key: "address",
					},
				},
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
		Name:            "nats",
		Image:           "nats:2.9.14-alpine",
		Ports:           []v1.ContainerPort{{ContainerPort: 4222}},
		ImagePullPolicy: v1.PullIfNotPresent,
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
				Name:      "system-user-creds",
				MountPath: "/system-user",
				ReadOnly:  true,
			},
			{
				Name:      "system-account-user-creds",
				MountPath: "/system-account-user",
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
		LeafNodes: nats.LeafNodesConfig{
			Remotes: []nats.LeafNodeRemoteConfig{
				{
					Url:         "TEMPLATE_ADDRESS:7422",
					Credentials: "/system-user/creds",
					Account:     "TEMPLATE_ACCOUNT_PUBLIC_KEY",
					DenyImports: []string{"local.>"},
					DenyExports: []string{"local.>"},
				},
				{
					Url:         "TEMPLATE_ADDRESS:7422",
					Credentials: "/system-account-user/creds",
					Account:     "TEMPLATE_SYS_ACCOUNT_PUBLIC_KEY",
				},
			},
		},
		Jetstream: nats.JetstreamConfig{
			MaxMemory: config.Spec.Limits.InMemoryStorage,
			MaxFile:   config.Spec.Limits.FileStorage,
			StoreDir:  "/data",
			Domain:    "TEMPLATE_NETWORK-TEMPLATE_SUB_NETWORK-TEMPLATE_NODE_NAME",
		},
		Operator:      "TEMPLATE_OPERATOR_JWT",
		SystemAccount: "TEMPLATE_SYS_ACCOUNT_PUBLIC_KEY",
		Resolver: &nats.ResolverConfig{
			Type:    "cache",
			Dir:     "/jwt",
			TTL:     "1h",
			Timeout: "2s",
		},
		ResolverPreload: map[string]string{
			"TEMPLATE_SYS_ACCOUNT_PUBLIC_KEY": "TEMPLATE_SYS_ACCOUNT_JWT",
		},
	}
	leafNatsConfigString, err := json.Marshal(leafNatsConfig, false)
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
