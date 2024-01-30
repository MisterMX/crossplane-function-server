package main

import (
	"context"
	_ "embed"
	"testing"

	"github.com/crossplane/function-sdk-go/logging"
	fntesting "github.com/mistermx/crossplane-function-server/testing"
)

var (
	//go:embed testdata/composite.yaml
	composite []byte

	//go:embed testdata/input.yaml
	input []byte

	//go:embed testdata/expected-resources.yaml
	expectedResources []byte
)

func TestFunction(t *testing.T) {
	fntesting.TestFunction(
		t, &MyFunction{log: logging.NewNopLogger()},
		fntesting.WithObservedCompositeYAML(composite),
		fntesting.WithInputYAML(input),
		fntesting.ExpectDesiredResourcesYAML(expectedResources),
	).Run(context.Background())
}
