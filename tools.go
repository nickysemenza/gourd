// +build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/cosmtrek/air"
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/mgechev/revive"
	_ "golang.org/x/tools/cmd/stringer"
)
