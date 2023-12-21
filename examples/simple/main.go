// Package main implements a Server Composition Function that serves multiple
// ServerFunctions.
package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/alecthomas/kong"
	"github.com/crossplane/function-sdk-go"
	"github.com/crossplane/function-sdk-go/resource/composed"

	server "github.com/mistermx/crossplane-function-server"
	"github.com/mistermx/crossplane-function-server/apis/v1alpha1"
)

// CLI of this Function.
type CLI struct {
	Debug bool `short:"d" help:"Emit debug logs in addition to info logs."`

	Network     string `help:"Network on which to listen for gRPC connections." default:"tcp"`
	Address     string `help:"Address at which to listen for gRPC connections." default:":9443"`
	TLSCertsDir string `help:"Directory containing server certs (tls.key, tls.crt) and the CA used to verify client certificates (ca.crt)" env:"TLS_SERVER_CERTS_DIR"`
	Insecure    bool   `help:"Run without mTLS credentials. If you supply this flag --tls-server-certs-dir will be ignored."`
}

// Run this Function.
func (c *CLI) Run() error {
	log, err := function.NewLogger(c.Debug)
	if err != nil {
		return err
	}

	kingpin.FatalIfError(v1alpha1.AddToScheme(composed.Scheme), "Cannot add function server API to scheme")

	return function.Serve(
		server.NewServer(
			server.WithFunction("my-function", &MyFunction{log: log}),
			// more server functions can be registered here
		),
		function.Listen(c.Network, c.Address),
		function.MTLSCertificates(c.TLSCertsDir),
		function.Insecure(c.Insecure),
	)
}

func main() {
	ctx := kong.Parse(&CLI{}, kong.Description("A Crossplane Server Function."))
	ctx.FatalIfErrorf(ctx.Run())
}
