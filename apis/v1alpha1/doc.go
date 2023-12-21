// +kubebuilder:object:generate=true
// Package v1alpha1 is the v1alpha1 version of the server.fn.crossplane.io API.
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Package type metadata.
const (
	CRDGroup   = "server.fn.crossplane.io"
	CRDVersion = "v1alpha1"
)

var (
	// GroupVersion is the API Group Version used to register the objects
	GroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
