package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"example.com/ness-api-function/graph/generated"
	"example.com/ness-api-function/graph/model"
)

func (r *mutationResolver) CreateThread(ctx context.Context, input model.NewThread) (*model.Thread, error) {
	return &model.Thread{
		ID:          "threadID",
		Title:       "title",
		Description: "description",
		Closed:      false,
	}, nil
}

func (r *queryResolver) Threads(ctx context.Context) ([]*model.Thread, error) {
	return []*model.Thread{
		{
			ID:          "threadID1",
			Title:       "title1",
			Description: "description1",
			Closed:      false,
		},
		{
			ID:          "threadID2",
			Title:       "title2",
			Description: "description2",
			Closed:      false,
		},
		{
			ID:          "threadID3",
			Title:       "title3",
			Description: "description3",
			Closed:      false,
		},
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
