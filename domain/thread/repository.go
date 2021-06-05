//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package thread

import (
	"context"
	"time"
)

type repository struct {
	dnm DynamoDB
	es  ElasticSearch
}

type DynamoDBThreadRow struct {
	Id        string
	TeamID    string
	CreaterID string
	Title     string
	Closed    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t DynamoDBThreadRow) toThread() Thread {
	return &thread{
		id:        t.Id,
		teamID:    TeamID(t.TeamID),
		createrID: UserID(t.CreaterID),
		title:     t.Title,
		closed:    t.Closed,
		createdAt: t.CreatedAt,
		updatedAt: t.UpdatedAt,
	}
}

type DynamoDB interface {
	GetThreadsByIDs(ctx context.Context, ids []string) (map[string]DynamoDBThreadRow, error)
	Find(ctx context.Context, id string) (DynamoDBThreadRow, error)
	Create(ctx context.Context, thread Thread) error
	Update(ctx context.Context, thread Thread) error
}

type SearchThreadIDsRequest struct {
	Size int
	From int
	Word string
}

type SearchThreadIDsOption int

const (
	SearchThreadIDsOptionOnlyClosed SearchThreadIDsOption = iota + 1
	SearchThreadIDsOptionOnlyOpened
)

type ElasticSearch interface {
	SearchThreadIDs(context.Context, SearchThreadIDsRequest, ...SearchThreadIDsOption) ([]string, error)
	PutThread(context.Context, Thread) error
}

var _ Repository = &repository{}

func NewDynamoRepository(
	dnm DynamoDB,
	es ElasticSearch,
) *repository {
	return &repository{dnm, es}
}

func (d *repository) get(ctx context.Context, req repositoryGetRequest) ([]Thread, error) {
	size := req.size
	if size > 100 {
		size = 100
	}

	esreq := SearchThreadIDsRequest{
		Size: size,
		From: req.from,
		Word: req.word,
	}

	opts := make([]SearchThreadIDsOption, 0, 1)
	if req.closed.Valid {
		if req.closed.Bool {
			opts = append(opts, SearchThreadIDsOptionOnlyClosed)
		} else {
			opts = append(opts, SearchThreadIDsOptionOnlyOpened)
		}
	}

	ids, err := d.es.SearchThreadIDs(ctx, esreq, opts...)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return []Thread{}, nil
	}

	rslt, err := d.dnm.GetThreadsByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	res := make([]Thread, len(ids))
	for i, id := range ids {
		t := rslt[id]
		res[i] = &thread{
			id:        t.Id,
			teamID:    TeamID(t.TeamID),
			createrID: UserID(t.CreaterID),
			title:     t.Title,
			closed:    t.Closed,
			createdAt: t.CreatedAt,
			updatedAt: t.UpdatedAt,
		}
	}

	return res, nil
}

func (d *repository) find(ctx context.Context, req repositoryFindRequest) (Thread, error) {
	thrd, err := d.dnm.Find(ctx, req.threadID)
	if err != nil {
		return nil, err
	}
	return thrd.toThread(), nil
}

func (d *repository) create(ctx context.Context, req repositoryCreateRequest) error {
	if err := d.dnm.Create(ctx, req.thread); err != nil {
		return err
	}
	return d.es.PutThread(ctx, req.thread)
}

func (d *repository) update(ctx context.Context, req repositoryUpdateRequest) error {
	if err := d.dnm.Update(ctx, req.thread); err != nil {
		return err
	}
	return d.es.PutThread(ctx, req.thread)
}
