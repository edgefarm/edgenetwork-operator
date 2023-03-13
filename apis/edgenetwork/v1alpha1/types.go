// +groupName=network.edgefarm.io
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=edgefarm,path=edgenetworks,singular=edgenetwork,shortName=en
// +kubebuilder:printcolumn:name="NETWORK",type="string",JSONPath=".spec.network"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="DESIRED",type="string",JSONPath=".status.desired"
// +kubebuilder:printcolumn:name="CURRENT",type="string",JSONPath=".status.current"
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="FILE LIMIT",type="string",priority=1,JSONPath=".spec.limits.fileStorage"
// +kubebuilder:printcolumn:name="MEMORY LIMIT",type="string",priority=1,JSONPath=".spec.limits.inMemoryStorage"
type EdgeNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              EdgeNetworkSpec   `json:"spec"`
	Status            EdgeNetworkStatus `json:"status,omitempty"`
}

// The spec to define an edge network
type EdgeNetworkSpec struct {
	// The name of the network
	// +kubebuilder:validation:Required
	Network string `json:"network"`

	//The address of the server.
	// Example: "example.com"
	// +kubebuilder:validation:Required
	Address string `json:"address"`

	// Indicates the node selector to form the node pool.
	// A pool's nodeSelectorTerm is not allowed to be updated.
	// +kubebuilder:validation:Required
	NodeSelectorTerm corev1.NodeSelectorTerm `json:"nodeSelectorTerm,omitempty"`

	// Indicates the tolerations the pods under this pool have.
	// A pool's tolerations is not allowed to be updated.
	// +kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Hardware limits for the edge network
	// +kubebuilder:validation:Required
	Limits Limits `json:"limits"`

	// Reference to the secret containing the credentials to connect to the network
	// +kubebuilder:validation:Required
	ConnectionSecretRef *corev1.LocalObjectReference `json:"connectionSecretRef"`
}

// Defines memory/storage limits to use
type Limits struct {
	// +kubebuilder:default="1G"
	// +kubebuilder:validation:Pattern=^\d+[GMKB]?$
	// +kubebuilder:validation:Required
	// How much disk space is available for data that is stored on disk
	FileStorage string `json:"fileStorage"`
	// +kubebuilder:default="1G"
	// +kubebuilder:validation:Pattern=^\d+[GMKB]?$
	// How much memory is available for data that is stored in memory
	// +kubebuilder:validation:Required
	InMemoryStorage string `json:"inMemoryStorage"`
}

type EdgeNetworkStatus struct {
	// The amount of desired participants in the edge network
	Desired string `json:"desired,omitempty"`
	// The amount of current participants in the edge network
	Current string `json:"current,omitempty"`
	// The amount of ready participants in the edge network
	Ready string `json:"ready,omitempty"`
}
