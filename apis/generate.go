//go:build generate
// +build generate

//go:generate go run -tags generate sigs.k8s.io/controller-tools/cmd/controller-gen paths=./... object

package apis

import (
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen" //nolint:typecheck
)
