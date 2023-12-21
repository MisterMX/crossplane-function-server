package v1alpha1

import (
	evtv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServerInputSpec defines request data for a ServerFunction call.
type ServerInputSpec struct {
	// FunctionName is the name of the ServerFunction to be invoked.
	FunctionName string `json:"functionName"`

	// Input is the request payload that should be passed to the function.
	// It can contain any kind of valid JSON data.
	Input evtv1.JSON `json:"input"`
}

// +kubebuilder:object:root=true

// ServerInput defines request data for a ServerFunction call.
type ServerInput struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServerInputSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// // ServerInputList contains a list of ServerInputs
// type CompositeServerInputList struct {
// 	metav1.TypeMeta `json:",inline"`
// 	metav1.ListMeta `json:"metadata,omitempty"`
// 	Items           []CompositeServerInput `json:"items"`
// }

// // Repository type metadata.
// var (
// 	CompositeServerInputKind             = "ServerInput"
// 	CompositeServerInputGroupKind        = schema.GroupKind{Group: CRDGroup, Kind: CompositeServerInputKind}.String()
// 	CompositeServerInputKindAPIVersion   = CompositeServerInputKind + "." + GroupVersion.String()
// 	CompositeServerInputGroupVersionKind = GroupVersion.WithKind(CompositeServerInputKind)
// )

// func init() {
// 	SchemeBuilder.Register(&CompositeServerInput{}, &CompositeServerInputList{})
// }
