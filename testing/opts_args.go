package testing

import (
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/meta"
	fnapi "github.com/crossplane/function-sdk-go/proto/v1beta1"
	evtv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// WithContextValue sets the expected context field to value.
func WithContextValue(key string, value any) TestFunctionOpt {
	return func(tc *FunctionTest) {
		val := mustStructValue(value)
		tc.args.context.Fields[key] = val
	}
}

// WithContextValueYAML reads a value from a single YAML document and sets it
// as value of the given context field.
func WithContextValueYAML(key string, rawYAML []byte) TestFunctionOpt {
	var val any
	if err := yaml.Unmarshal(rawYAML, &val); err != nil {
		panic(err.Error())
	}
	return WithContextValue(key, val)
}

// WithContextValueYAML reads a value from a JSON document and sets it
// as value of the given context field.
func WithContextValueJSON(key string, rawJSON []byte) TestFunctionOpt {
	var val any
	if err := json.Unmarshal(rawJSON, &val); err != nil {
		panic(err.Error())
	}
	return WithContextValue(key, val)
}

// WithInput sets the input that is passed to the function run.
// It accepts any value that can be marshaled to JSON.
func WithInput(input any) TestFunctionOpt {
	raw, err := json.Marshal(input)
	if err != nil {
		panic(err.Error())
	}
	return WithInputJSON(raw)
}

// WithInputYAML is the same as [WithInput] but accepts raw YAML.
func WithInputYAML(inputYaml []byte) TestFunctionOpt {
	rawJson, err := yaml.ToJSON(inputYaml)
	if err != nil {
		panic(err.Error())
	}
	return WithInputJSON(rawJson)
}

// WithInputJSON is the same as [WithInput] but accepts raw JSON.
func WithInputJSON(inputJson []byte) TestFunctionOpt {
	return func(tc *FunctionTest) {
		tc.args.input = evtv1.JSON{Raw: inputJson}
	}
}

// WithObservedResourceObject adds o to the observed state passed to the
// function.
func WithObservedResourceObject(name string, o runtime.Object) TestFunctionOpt {
	return func(tc *FunctionTest) {
		str := mustObjectAsStruct(o)
		tc.args.observedResources[name] = &fnapi.Resource{
			Resource: str,
		}
	}
}

// WithObservedResourceYAML reads an object from a single YAML document and adds
// it to the observed state passed to the function.
func WithObservedResourceYAML(name string, rawYaml []byte) TestFunctionOpt {
	u := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(rawYaml, u); err != nil {
		panic(err.Error())
	}
	return WithObservedResourceObject(name, u)
}

// WithObservedResourceJSON reads an object from a single JSON document and adds
// it to the observed state passed to the function.
func WithObservedResourceJSON(name string, rawJSON []byte) TestFunctionOpt {
	u := &unstructured.Unstructured{}
	if err := json.Unmarshal(rawJSON, u); err != nil {
		panic(err.Error())
	}
	return WithObservedResourceObject(name, u)
}

// AnnotationKeyResourceName is the key of the annotation that defines the
// resource name.
const AnnotationKeyResourceName = "fn-server.test/resource-name"

// WithObservedResourcesYAML reads all objects from a multi-document YAML and
// passes them with the observed state to the function.
//
// It uses the annotation [AnnotationKeyResourceName] to determine
// the name of the resource.
func WithObservedResourcesYAML(rawYAML []byte) TestFunctionOpt {
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

			str := mustObjectAsStruct(u)
			tc.args.observedResources[key] = &fnapi.Resource{
				Resource: str,
			}
		}
	}
}

// WithObservedCompositeObject sets the observed composite to the given object.
func WithObservedCompositeObject(o runtime.Object) TestFunctionOpt {
	return func(tc *FunctionTest) {
		str := mustObjectAsStruct(o)
		tc.args.observedComposite = &fnapi.Resource{
			Resource: str,
		}
	}
}

// WithObservedCompositeYAML reads an object from a single YAML document and
// passes it as observed composite to the function.
func WithObservedCompositeYAML(rawYaml []byte) TestFunctionOpt {
	u := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(rawYaml, u); err != nil {
		panic(err.Error())
	}
	return WithObservedCompositeObject(u)
}

// WithObservedCompositeJSON reads an object from a JSON document and
// passes it as observed composite to the function.
func WithObservedCompositeJSON(rawJSON []byte) TestFunctionOpt {
	u := &unstructured.Unstructured{}
	if err := json.Unmarshal(rawJSON, u); err != nil {
		panic(err.Error())
	}
	return WithObservedCompositeObject(u)
}
