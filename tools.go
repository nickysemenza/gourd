//go:build tools
// +build tools

package main

import (
	_ "github.com/cosmtrek/air"
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/ory/go-acc"
	_ "golang.org/x/tools/cmd/stringer"
)
