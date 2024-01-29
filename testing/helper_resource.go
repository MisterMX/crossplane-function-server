package testing

import (
	"encoding/json"

	"github.com/crossplane/function-sdk-go/resource"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func mustObjectAsStruct(o runtime.Object) *structpb.Struct {
	s, err := resource.AsStruct(o)
	if err != nil {
		panic(err.Error())
	}
	return s
}

func mustStructValue(in any) *structpb.Value {
	val, err := structpb.NewValue(in)
	if err != nil {
		panic(err.Error())
	}
	return val
}

func mustObjectAsUnstructured(o runtime.Object) *unstructured.Unstructured {
	raw, err := json.Marshal(o)
	if err != nil {
		panic(err.Error())
	}
	u := &unstructured.Unstructured{}
	if err := json.Unmarshal(raw, u); err != nil {
		panic(err.Error())
	}
	return u
}
