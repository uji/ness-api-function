package graph

import "example.com/ness-api-function/domain/thread"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	thread *thread.Usecase
}

func NewResolver(thread *thread.Usecase) *Resolver {
	return &Resolver{thread}
}
