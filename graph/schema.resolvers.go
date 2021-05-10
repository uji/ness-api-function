package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/uji/ness-api-function/domain/thread"
	"github.com/uji/ness-api-function/graph/generated"
	"github.com/uji/ness-api-function/graph/model"
)

func (r *mutationResolver) CreateThread(ctx context.Context, input model.NewThread) (*model.Thread, error) {
	res, err := r.thread.Create(ctx, thread.CreateRequest{
		Title: input.Title,
	})
	if err != nil {
		return nil, err
	}
	return cnvThread(res), nil
}

func (r *mutationResolver) OpenThread(ctx context.Context, input model.OpenThread) (*model.Thread, error) {
	res, err := r.thread.Open(ctx, thread.OpenRequest{
		ThreadID: input.ThreadID,
	})
	return cnvThread(res), err
}

func (r *mutationResolver) CloseThread(ctx context.Context, input model.CloseThread) (*model.Thread, error) {
	res, err := r.thread.Close(ctx, thread.CloseRequest{
		ThreadID: input.ThreadID,
	})
	return cnvThread(res), err
}

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Threads(ctx context.Context, input model.GetThreadsInput) ([]*model.Thread, error) {
	clsd := null.BoolFromPtr(input.Closed)
	size := 30
	if input.Size != nil {
		size = *input.Size
	}
	from := 0
	if input.From != nil {
		from = *input.From
	}
	word := ""
	if input.Word != nil {
		word = *input.Word
	}
	thrds, err := r.thread.Get(ctx, thread.GetRequest{
		Closed: clsd,
		Size:   size,
		From:   from,
		Word:   word,
	})
	if err != nil {
		return nil, err
	}
	res := make([]*model.Thread, len(thrds))
	for i, thrd := range thrds {
		res[i] = cnvThread(thrd)
	}
	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
