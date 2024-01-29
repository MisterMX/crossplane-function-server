package testing

import (
	"context"
	"testing"

	fnapi "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/google/go-cmp/cmp"
	"github.com/mistermx/crossplane-function-server/apis/v1alpha1"
	"google.golang.org/protobuf/types/known/structpb"
	evtv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	server "github.com/mistermx/crossplane-function-server"
)

type FunctionTest struct {
	t  *testing.T
	fn server.ServerFunction

	// want function request
	args struct {
		input             evtv1.JSON
		observedResources map[string]*fnapi.Resource
		observedComposite *fnapi.Resource
		context           *structpb.Struct
	}
	// want function response
	want struct {
		desiredResources map[string]*fnapi.Resource
		desiredComposite *fnapi.Resource
		err              error
		results          []*fnapi.Result
	}
}

type TestFunctionOpt func(tc *FunctionTest)

func TestFunction(t *testing.T, fn server.ServerFunction, opts ...TestFunctionOpt) *FunctionTest {
	tc := &FunctionTest{
		t:  t,
		fn: fn,
	}
	tc.args.observedResources = map[string]*fnapi.Resource{}
	tc.args.context = &structpb.Struct{
		Fields: map[string]*structpb.Value{},
	}
	tc.want.desiredResources = map[string]*fnapi.Resource{}
	tc.want.results = []*fnapi.Result{}
	for _, o := range opts {
		o(tc)
	}
	return tc
}

func (t *FunctionTest) Run(ctx context.Context) {
	req := server.RunServerFunctionRequest{
		Req: &fnapi.RunFunctionRequest{
			Observed: &fnapi.State{
				Composite: t.args.observedComposite,
				Resources: t.args.observedResources,
			},
			Context: t.args.context,
		},
		ServerInput: &v1alpha1.ServerInput{
			Spec: v1alpha1.ServerInputSpec{
				Input: t.args.input,
			},
		},
	}
	res := server.RunServerFunctionResponse{
		DesiredComposed: map[string]*fnapi.Resource{},
		DesiredContext:  &structpb.Struct{},
	}

	err := t.fn.Run(ctx, &req, &res)

	// Compare test results with expected outcome and print any assert failure
	// as a plain diff in the test results:

	if diff := cmp.Diff(convertResourceToUnstructured(t.want.desiredComposite), convertResourceToUnstructured(res.DesiredComposite)); diff != "" {
		t.t.Errorf("Composite: -want +got\n%s\n", diff)
	}
	if diff := cmp.Diff(convertResourcesMapToUnstructured(t.want.desiredResources), convertResourcesMapToUnstructured(res.DesiredComposed)); diff != "" {
		t.t.Errorf("Resources: -want +got\n%s\n", diff)
	}
	if diff := cmp.Diff(convertResultsToMap(t.want.results), convertResultsToMap(res.Results)); diff != "" {
		t.t.Errorf("Results: -want +got\n%s\n", diff)
	}
	if diff := cmp.Diff(t.want.err, err); diff != "" {
		t.t.Errorf("Error: -want +got\n%s\n", diff)
	}
}
