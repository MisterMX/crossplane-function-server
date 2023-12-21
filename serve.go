package server

// ServerOption that configures a function Server.
type ServerOption func(server *Server)

// NewServer create a new Server instance that implements the Crossplane
// Function interface and is able to serve multiple subfunctions
// (aka server functions) at the same time.
//
// ServerFunctions are registered via the WithFunction option:
//
//	server.NewServer(
//		server.WithFunction(&MyFunction{})
//		server.WithFunction(&MyOtherFunction{})
//		// ...
//	)
func NewServer(opts ...ServerOption) *Server {
	server := &Server{
		functions: map[string]ServerFunction{},
	}
	for _, o := range opts {
		o(server)
	}
	return server
}

// WithFunction registeres a ServerFunction at a Server with a given name.
func WithFunction(name string, fn ServerFunction) ServerOption {
	return func(server *Server) {
		server.functions[name] = fn
	}
}
