//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package thread

import (
	"context"
	"time"
)

type repository struct {
	es ElasticSearch
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

type ElasticSearchThreadRow struct {
	Id        string
	TeamID    string
	CreaterID string
	Title     string
	Closed    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t ElasticSearchThreadRow) toThread() Thread {
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

type ElasticSearch interface {
	FindThread(ctx context.Context, id string) (ElasticSearchThreadRow, error)
	SearchThreads(context.Context, SearchThreadIDsRequest, ...SearchThreadIDsOption) ([]ElasticSearchThreadRow, error)
	PutThread(context.Context, Thread) error
}

var _ Repository = &repository{}

func NewRepository(
	es ElasticSearch,
) *repository {
	return &repository{es}
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

	rslts, err := d.es.SearchThreads(ctx, esreq, opts...)
	if err != nil {
		return nil, err
	}

	res := make([]Thread, len(rslts))
	for i, r := range rslts {
		res[i] = &thread{
			id:        r.Id,
			teamID:    TeamID(r.TeamID),
			createrID: UserID(r.CreaterID),
			title:     r.Title,
			closed:    r.Closed,
			createdAt: r.CreatedAt,
			updatedAt: r.UpdatedAt,
		}
	}

	return res, nil
}

func (d *repository) find(ctx context.Context, req repositoryFindRequest) (Thread, error) {
	thrd, err := d.es.FindThread(ctx, req.threadID)
	if err != nil {
		return nil, err
	}
	return thrd.toThread(), nil
}

func (d *repository) create(ctx context.Context, req repositoryCreateRequest) error {
	return d.es.PutThread(ctx, req.thread)
}

func (d *repository) update(ctx context.Context, req repositoryUpdateRequest) error {
	return d.es.PutThread(ctx, req.thread)
}
