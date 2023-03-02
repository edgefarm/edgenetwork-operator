package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:resource:categories=foo
// +kubebuilder:resource:shortName=en
//
// +kubebuilder:printcolumn:name="NETWORK",type="string",JSONPath=".spec.network"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="FILE LIMIT",type="string",priority=1,JSONPath=".spec.limits.fileStorage"
// +kubebuilder:printcolumn:name="MEMORY LIMIT",type="string",priority=1,JSONPath=".spec.limits.inMemoryStorage"
type EdgeNetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              EdgeNetworkSpec   `json:"spec"`
	Status            EdgeNetworkStatus `json:"status,omitempty"`
}

type EdgeNetworkSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=3
	Network string `json:"network"`
	// +kubebuilder:validation:Required
	NodeSelecter map[string]string `json:"nodeSelector"`
	// +kubebuilder:validation:Required
	Limits Limits `json:"limits"`
}

type EdgeNetworkStatus struct {
	Replicas  int `json:"replicas,omitempty"`
	Succeeded int `json:"succeeded,omitempty"`
}

type Limits struct {
	// +kubebuilder:default="1G"
	// +kubebuilder:validation:Pattern=^\d+[GMKB]?$
	FileStorage string `json:"fileStorage"`
	// +kubebuilder:default="1G"
	// +kubebuilder:validation:Pattern=^\d+[GMKB]?$
	InMemoryStorage string `json:"inMemoryStorage"`
}
