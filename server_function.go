package server

import (
	"context"

	fnapi "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// A ServerFunction is a high-level subroutine of native Crossplane Go function.
type ServerFunction interface {
	// Run executes the ServerFunction for the given request.
	Run(ctx context.Context, req ServerFunctionRequest, res ServerFunctionResponse) error
}

// ServerFunctionRequest provides ways to easily the request payload of a
// ServerFunction call.
type ServerFunctionRequest interface {
	// GetNativeRequest returns the native function request of the
	// underlying Crossplane function-sdk-go.
	GetNativeRequest() *fnapi.RunFunctionRequest

	// GetInput for this server function.
	//
	// Note that this is not the input of the Crossplane function-sdk-go.
	// To receive that use GetNativeRequest().GetInput().
	GetInput(target any) error

	// GetComposite copies the current state of the composite resource
	// into the given target object.
	GetComposite(target runtime.Object) error

	// GetComposed copies the current state of the composed resource identified
	// by the given name.
	//
	// If a no composed resource with the given name exists, target remains
	// unchanged.
	GetComposed(name string, target runtime.Object) error
}

// ServerFunctionResponse provides ways to easily define the response payload
// of a ServerFunction call.
type ServerFunctionResponse interface {
	// SetComposite save the given object as desired state of the composite
	// resource for the given name.
	SetComposite(o runtime.Object, mods ...ResourceModifier) error

	// GetComposite gets the current state of the composite resource of this
	// response and writes its contents into the given target object.
	GetComposite(target runtime.Object) error

	// SetCompositeRaw sets the desired response state directly using the
	// native SDK types.
	//
	// This is useful if the desired resource state already exists as a native
	// type so there is no need to work with runtime.Objects.
	SetCompositeRaw(res *fnapi.Resource)

	// SetComposed saves the given composed object as desired composed object
	// identified by the given name for this function's response.
	SetComposed(name string, o runtime.Object, mods ...ResourceModifier) error

	// GetComposed looks up the composed resource in the current response object
	// and writes its contents into the given target object.
	GetComposed(name string, target runtime.Object) error

	// SetComposedRaw sets the desired composed resource state directly using
	// the native SDK types.
	//
	// This is useful if the desired resource state already exists as a native
	// type so there is no need to work with runtime.Objects.
	SetComposedRaw(name string, res *fnapi.Resource)

	// SetContextField sets the value of the context field key to the given
	// value. The passed value must be convertable to protobuf.
	SetContextField(key string, value any) error

	// SetNativeResults of the underlying SDK requests.
	SetNativeResults(results []*fnapi.Result)
}

// ResourceModifier applies modifications to a bare-metal Crossplane function
// resource.
type ResourceModifier func(r *fnapi.Resource)

// WithConnectionDetails sets the given connection details to a resource.
func WithConnectionDetails(connectionDetails map[string][]byte) ResourceModifier {
	return func(r *fnapi.Resource) {
		r.ConnectionDetails = connectionDetails
	}
}

// WithReady applies the desired ready state to a resource.
func WithReady(ready fnapi.Ready) ResourceModifier {
	return func(r *fnapi.Resource) {
		r.Ready = ready
	}
}

// WithReady shorthand
var (
	WithReadyIsReady     = WithReady(fnapi.Ready_READY_TRUE)
	WithReadyIsNotReady  = WithReady(fnapi.Ready_READY_FALSE)
	WithReadyUnspecified = WithReady(fnapi.Ready_READY_UNSPECIFIED)
)
