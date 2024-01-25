package server

import (
	"context"
	"encoding/json"

	fnapi "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mistermx/crossplane-function-server/apis/v1alpha1"
)

// Server is a special Crossplane function which acts as a router that
// distributes request among several subroutines (aka server functions).
//
// A server receives a dedicated input payload by the composition that tells
// the server which subroutine to call.
type Server struct {
	fnapi.UnimplementedFunctionRunnerServiceServer

	functions map[string]ServerFunction
}

func (s *Server) RunFunction(ctx context.Context, req *fnapi.RunFunctionRequest) (*fnapi.RunFunctionResponse, error) {
	serverInput := &v1alpha1.ServerInput{}
	if err := resource.AsObject(req.Input, serverInput); err != nil {
		return nil, errors.Wrap(err, "cannot parse input")
	}

	fn, exists := s.functions[serverInput.Spec.FunctionName]
	if !exists {
		return nil, errors.Errorf("no function with name %q", serverInput.Spec.FunctionName)
	}

	fnReq := RunServerFunctionRequest{
		Req:         req,
		ServerInput: serverInput,
	}
	fnRes := RunServerFunctionResponse{}
	if err := fn.Run(ctx, &fnReq, &fnRes); err != nil {
		return nil, errors.Wrapf(err, "error while running subroutine function %q", serverInput.Spec.FunctionName)
	}

	res := &fnapi.RunFunctionResponse{
		Desired: req.GetDesired(),
		Context: req.GetContext(),
	}
	if fnRes.DesiredComposite != nil {
		res.Desired.Composite = fnRes.DesiredComposite
	}
	if fnRes.DesiredComposed != nil {
		res.Desired.Resources = fnRes.DesiredComposed
	}
	if fnRes.DesiredContext != nil {
		res.Context = fnRes.DesiredContext
	}
	if fnRes.Results != nil {
		res.Results = fnRes.Results
	}
	return res, nil
}

type RunServerFunctionRequest struct {
	Req         *fnapi.RunFunctionRequest
	ServerInput *v1alpha1.ServerInput
}

func (r *RunServerFunctionRequest) GetNativeRequest() *fnapi.RunFunctionRequest {
	return r.Req
}

func (r *RunServerFunctionRequest) GetComposite(target runtime.Object) error {
	return resource.AsObject(r.Req.GetObserved().GetComposite().GetResource(), target)
}

func (r *RunServerFunctionRequest) GetComposed(name string, target runtime.Object) error {
	resources := r.Req.GetObserved().GetResources()
	res, exists := resources[name]
	if !exists {
		return nil
	}
	return resource.AsObject(res.GetResource(), target)
}

func (r *RunServerFunctionRequest) GetInput(target any) error {
	return json.Unmarshal(r.ServerInput.Spec.Input.Raw, target)
}

type RunServerFunctionResponse struct {
	DesiredComposite *fnapi.Resource
	DesiredComposed  map[string]*fnapi.Resource
	DesiredContext   *structpb.Struct
	Results          []*fnapi.Result
}

func (r *RunServerFunctionResponse) SetCompositeRaw(res *fnapi.Resource) {
	r.DesiredComposite = res
}

func (r *RunServerFunctionResponse) SetComposite(o runtime.Object, mods ...ResourceModifier) error {
	raw, err := resource.AsStruct(o)
	if err != nil {
		return err
	}
	r.DesiredComposite = &fnapi.Resource{
		Resource: raw,
	}
	for _, m := range mods {
		m(r.DesiredComposite)
	}
	return nil
}

func (r *RunServerFunctionResponse) GetComposite(target runtime.Object) error {
	if r.DesiredComposite == nil || r.DesiredComposite.Resource == nil {
		return nil // Return an error here?
	}
	return resource.AsObject(r.DesiredComposite.Resource, target)
}

func (r *RunServerFunctionResponse) SetComposedRaw(name string, res *fnapi.Resource) {
	if r.DesiredComposed == nil {
		r.DesiredComposed = map[string]*fnapi.Resource{}
	}
	r.DesiredComposed[name] = res
}

func (r *RunServerFunctionResponse) SetComposed(name string, o runtime.Object, mods ...ResourceModifier) error {
	raw, err := resource.AsStruct(o)
	if err != nil {
		return err
	}
	res := &fnapi.Resource{
		Resource: raw,
	}
	for _, m := range mods {
		m(res)
	}
	r.SetComposedRaw(name, res)
	return nil
}

func (r *RunServerFunctionResponse) GetComposed(name string, target runtime.Object) error {
	state, exists := r.DesiredComposed[name]
	if !exists {
		return errNotFound(name)
	}
	return resource.AsObject(state.Resource, target)
}

func (r *RunServerFunctionResponse) SetContextField(key string, value any) error {
	if r.DesiredContext == nil {
		r.DesiredContext = &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		}
	}
	raw, err := structpb.NewValue(value)
	if err != nil {
		return errors.Wrap(err, "cannot convert context value to protobuf")
	}
	r.DesiredContext.Fields[key] = raw
	return nil
}

func (r *RunServerFunctionResponse) SetNativeResults(results []*fnapi.Result) {
	r.Results = results
}
