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

	fnReq := serverFunctionReq{
		req:         req,
		serverInput: serverInput,
	}
	fnRes := serverFunctionRes{}
	if err := fn.Run(ctx, &fnReq, &fnRes); err != nil {
		return nil, errors.Wrapf(err, "error while running subroutine function %q", serverInput.Spec.FunctionName)
	}

	res := &fnapi.RunFunctionResponse{
		Desired: req.GetDesired(),
		Context: req.GetContext(),
	}
	if fnRes.desiredComposite != nil {
		res.Desired.Composite = fnRes.desiredComposite
	}
	if fnRes.desiredComposed != nil {
		res.Desired.Resources = fnRes.desiredComposed
	}
	if fnRes.desiredContext != nil {
		res.Context = fnRes.desiredContext
	}
	if fnRes.results != nil {
		res.Results = fnRes.results
	}
	return res, nil
}

type serverFunctionReq struct {
	req         *fnapi.RunFunctionRequest
	serverInput *v1alpha1.ServerInput
}

func (r *serverFunctionReq) GetNativeRequest() *fnapi.RunFunctionRequest {
	return r.req
}

func (r *serverFunctionReq) GetComposite(target runtime.Object) error {
	return resource.AsObject(r.req.GetObserved().GetComposite().GetResource(), target)
}

func (r *serverFunctionReq) GetComposed(name string, target runtime.Object) error {
	resources := r.req.GetObserved().GetResources()
	res, exists := resources[name]
	if !exists {
		return nil
	}
	return resource.AsObject(res.GetResource(), target)
}

func (r *serverFunctionReq) GetInput(target any) error {
	return json.Unmarshal(r.serverInput.Spec.Input.Raw, target)
}

type serverFunctionRes struct {
	desiredComposite *fnapi.Resource
	desiredComposed  map[string]*fnapi.Resource
	desiredContext   *structpb.Struct
	results          []*fnapi.Result
}

func (r *serverFunctionRes) SetCompositeRaw(res *fnapi.Resource) {
	r.desiredComposite = res
}

func (r *serverFunctionRes) SetComposite(o runtime.Object, mods ...ResourceModifier) error {
	raw, err := resource.AsStruct(o)
	if err != nil {
		return err
	}
	r.desiredComposite = &fnapi.Resource{
		Resource: raw,
	}
	for _, m := range mods {
		m(r.desiredComposite)
	}
	return nil
}

func (r *serverFunctionRes) SetComposedRaw(name string, res *fnapi.Resource) {
	if r.desiredComposed == nil {
		r.desiredComposed = map[string]*fnapi.Resource{}
	}
	r.desiredComposed[name] = res
}

func (r *serverFunctionRes) SetComposed(name string, o runtime.Object, mods ...ResourceModifier) error {
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

func (r *serverFunctionRes) SetContextField(key string, value any) error {
	if r.desiredContext == nil {
		r.desiredContext = &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		}
	}
	raw, err := structpb.NewValue(value)
	if err != nil {
		return errors.Wrap(err, "cannot convert context value to protobuf")
	}
	r.desiredContext.Fields[key] = raw
	return nil
}

func (r *serverFunctionRes) SetNativeResults(results []*fnapi.Result) {
	r.results = results
}
