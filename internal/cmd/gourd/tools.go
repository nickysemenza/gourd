//go:build tools
// +build tools

package main

import (
	_ "github.com/air-verse/air"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/ory/go-acc"
	_ "golang.org/x/tools/cmd/stringer"
)
