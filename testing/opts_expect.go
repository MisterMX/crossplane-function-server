package testing

import (
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	fnapi "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceModifier modifies a [fnapi.Resource].
type ResourceModifier func(res *fnapi.Resource)

// WithReady sets the ready state of an [fnapi.Resource].
func WithReady(ready fnapi.Ready) ResourceModifier {
	return func(res *fnapi.Resource) { res.Ready = ready }
}

// WithConnectionDetails sets the connection details of an [fnapi.Resource].
func WithConnectionDetails(cd map[string][]byte) ResourceModifier {
	return func(res *fnapi.Resource) { res.ConnectionDetails = cd }
}

// ExpectDesiredCompositeObject expects the given [runtime.Object] as desired
// composite as result of the function.
func ExpectDesiredCompositeObject(o runtime.Object, mods ...ResourceModifier) TestFunctionOpt {
	return func(tc *FunctionTest) {
		res := &fnapi.Resource{
			Resource: mustObjectAsStruct(o),
		}
		for _, m := range mods {
			m(res)
		}
		tc.want.desiredComposite = res
	}
}

// ExpectDesiredCompositeYAML is the same as [ExpectDesiredCompositeObject] but
// reads the object from a single YAML document.
func ExpectDesiredCompositeYAML(rawYAML []byte, mods ...ResourceModifier) TestFunctionOpt {
	return ExpectDesiredCompositeObject(mustUnstructuredFromYAML(rawYAML), mods...)
}

// ExpectDesiredCompositeJSON is the same as [ExpectDesiredCompositeObject] but
// reads the object from a JSON document.
func ExpectDesiredCompositeJSON(rawJSON []byte, mods ...ResourceModifier) TestFunctionOpt {
	return ExpectDesiredCompositeObject(mustUnstructuredFromJSON(rawJSON), mods...)
}

// ExpectDesiredResourceObject adds an object to the expected outcome of a
// function.
func ExpectDesiredResourceObject(name string, o runtime.Object, mods ...ResourceModifier) TestFunctionOpt {
	return func(tc *FunctionTest) {
		res := &fnapi.Resource{
			Resource: mustObjectAsStruct(o),
		}
		for _, m := range mods {
			m(res)
		}
		tc.want.desiredResources[name] = res
	}
}

// ExpectDesiredResourceYAML is the same as [ExpectDesiredResourceObject] but
// reads the object from a single YAML document.
func ExpectDesiredResourceYAML(name string, rawYAML []byte, mods ...ResourceModifier) TestFunctionOpt {
	return ExpectDesiredResourceObject(name, mustUnstructuredFromYAML(rawYAML), mods...)
}

// ExpectDesiredResourceJSON is the same as [ExpectDesiredResourceObject] but
// reads the object from a JSON document.
func ExpectDesiredResourceJSON(name string, rawJSON []byte, mods ...ResourceModifier) TestFunctionOpt {
	return ExpectDesiredResourceObject(name, mustUnstructuredFromJSON(rawJSON), mods...)
}

// ExpectedDesiredResourcesYAML reads all objects from a multi-document YAML and
// expected them as desired resources from the function.
//
// It uses the annotation [AnnotationKeyResourceName] to determine
// the name of the resource.
func ExpectDesiredResourcesYAML(rawYAML []byte) TestFunctionOpt {
	return func(tc *FunctionTest) {
		uList, err := unmarshalObjectsYAML(rawYAML)
		if err != nil {
			panic(err.Error())
		}
		for _, u := range uList {
			key, exists := u.GetAnnotations()[AnnotationKeyResourceName]
			if !exists || key == "" {
				panic("resource has no name annotation")
			}
			meta.RemoveAnnotations(u, AnnotationKeyResourceName)

			// If the name annotation was the only annotation in the resource,
			// delete the entire field to avoid creating unnecessary diffs.
			if len(u.GetAnnotations()) == 0 {
				u.SetAnnotations(nil)
			}

			str := mustObjectAsStruct(u)
			tc.want.desiredResources[key] = &fnapi.Resource{
				Resource: str,
				// TODO: Set connection details and ready state
			}
		}
	}
}

// ExpectResults expects a list of [fnapi.Result] from a function.
func ExpectResults(results []*fnapi.Result) TestFunctionOpt {
	return func(tc *FunctionTest) { tc.want.results = results }
}

// ExpectError expects an error from a TestFunctionOpt.
func ExpectError(err error) TestFunctionOpt {
	return func(tc *FunctionTest) { tc.want.err = err }
}
