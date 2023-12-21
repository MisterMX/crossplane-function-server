package main

import (
	"context"
	"log"

	"github.com/crossplane/function-sdk-go/logging"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	server "github.com/mistermx/crossplane-function-server"
)

type MyFunctionInput struct {
	APIGroups []interface{} `json:"apiGroups"`
	Resources []interface{} `json:"resources"`
}

type MyFunction struct {
	log logging.Logger
}

func (f *MyFunction) Run(ctx context.Context, req server.ServerFunctionRequest, res server.ServerFunctionResponse) error {
	// ServerFunction input can be anything that unmarshalls from a JSON blob.
	input := MyFunctionInput{}
	if err := req.GetInput(&input); err != nil {
		return errors.Wrap(err, "cannot parse input")
	}

	log.Default().Println("Calling function MyFunction")

	// For the sake of simplicity this example uses unstructured.Unstructured
	// to deploy a standard K8s ClusterRole.
	// However, it is possible to use any type that satisfies the runtime.Object
	// interface. This is especially interesting for managed resource types that
	// are shipped by Crossplane providers.
	err := res.SetComposed("clusterRole", &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRole",
			"metadata": map[string]interface{}{
				"name": "composed",
			},
			"rules": []interface{}{
				map[string]interface{}{
					// "apiGroups": []interface{}{""},
					// "resources": []interface{}{"pods"},
					"apiGroups": input.APIGroups,
					"resources": input.Resources,
					"verbs":     []interface{}{"get", "watch", "list"},
				},
			},
		},
	}, server.WithReadyIsReady)
	if err != nil {
		return err
	}
	return nil
}
