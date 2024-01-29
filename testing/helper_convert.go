package testing

import (
	fnapi "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func convertResultsToMap(r []*fnapi.Result) []map[string]interface{} {
	if r == nil {
		return nil
	}
	res := make([]map[string]interface{}, len(r))
	for i, rr := range r {
		if rr == nil {
			continue
		}
		res[i] = map[string]interface{}{
			"Message":  rr.Message,
			"Severity": rr.Severity,
		}
	}
	return res
}

func convertResourcesMapToUnstructured(r map[string]*fnapi.Resource) map[string]*unstructured.Unstructured {
	if r == nil {
		return nil
	}
	res := map[string]*unstructured.Unstructured{}
	for k, v := range r {
		res[k] = convertResourceToUnstructured(v)
	}
	return res
}

func convertResourceToUnstructured(r *fnapi.Resource) *unstructured.Unstructured {
	if r == nil {
		return nil
	}
	u := &unstructured.Unstructured{}
	if err := resource.AsObject(r.Resource, u); err != nil {
		panic(err)
	}
	return u
}
