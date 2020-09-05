// +build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/golang/mock"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/rubenv/sql-migrate"
	_ "github.com/rubenv/sql-migrate/sql-migrate"
	_ "github.com/rubenv/sql-migrate/sqlparse"
	_ "github.com/xo/xo"
	_ "golang.org/x/tools/cmd/goimports"
)
