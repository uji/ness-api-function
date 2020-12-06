package graph

import (
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/domain/usr"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	user   *usr.Usecase
	thread *thread.Usecase
}

func NewResolver(
	user *usr.Usecase,
	thread *thread.Usecase,
) *Resolver {
	return &Resolver{user, thread}
}
