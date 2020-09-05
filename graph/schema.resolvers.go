package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

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

func (r *queryResolver) Threads(ctx context.Context) ([]*model.Thread, error) {
	thrds, err := r.thread.Get(ctx, thread.GetRequest{
		Limit:  100,
		Offset: 0,
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
