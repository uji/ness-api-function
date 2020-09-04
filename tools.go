// +build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/rubenv/sql-migrate"
	_ "github.com/rubenv/sql-migrate/sql-migrate"
	_ "github.com/rubenv/sql-migrate/sqlparse"
)
