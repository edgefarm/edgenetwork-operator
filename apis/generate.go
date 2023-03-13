//go:build generate
// +build generate

// Remove existing CRDs
//go:generate rm -rf ../manifests/crds

// Generate deepcopy methodsets and CRD manifests
//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen +crd:generateEmbeddedObjectMeta=true paths=./... crd:crdVersions=v1 output:artifacts:config=../manifests/crds

package apis

import (
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen" //nolint:typecheck
)
