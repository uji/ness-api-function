package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"example.com/ness-api-function/domain/thread"
	"example.com/ness-api-function/graph/generated"
	"example.com/ness-api-function/graph/model"
)

func (r *mutationResolver) CreateThread(ctx context.Context, input model.NewThread) (*model.Thread, error) {
	res, err := r.thread.Create(ctx, thread.CreateRequest{
		Title: input.Title,
	})
	if err != nil {
		return nil, err
	}
	return &model.Thread{
		ID:     res.ID(),
		Title:  res.Title(),
		Closed: res.Closed(),
	}, nil
}

func (r *mutationResolver) CloseThread(ctx context.Context, input model.CloseThread) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Threads(ctx context.Context, input model.GetThreadsInput) ([]*model.Thread, error) {
	thrds, err := r.thread.Get(ctx, thread.GetRequest{
		Limit:           input.Limit,
		LastEvaluatedID: input.LastEvaluatedID,
	})
	if err != nil {
		return nil, err
	}
	res := make([]*model.Thread, len(thrds))
	for i, thrd := range thrds {
		res[i] = &model.Thread{
			ID:     thrd.ID(),
			Title:  thrd.Title(),
			Closed: thrd.Closed(),
		}
	}
	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Nodes(ctx context.Context, ids []string) ([]model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}
