package thread

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type psqlRepository struct {
	db *sqlx.DB
}

func NewPsqlRepository(db *sqlx.DB) *psqlRepository {
	return &psqlRepository{db}
}

var _ repository = &psqlRepository{}

func (p *psqlRepository) get(
	ctx context.Context,
	req repositoryGetRequest,
) ([]*Thread, error) {
	query := `
select
    uuid
  , title
  , closed
from threads
limit :limit
offset :offset
  `

	args := struct {
		Limit  int `db:"limit"`
		Offset int `db:"offset"`
	}{
		Limit:  req.limit,
		Offset: req.offset,
	}

	rslt := struct {
		UUID   string `db:"uuid"`
		Title  string `db:"title"`
		Closed bool   `db:"closed"`
	}{}

	rows, err := p.db.NamedQueryContext(ctx, query, &args)
	if err != nil {
		return nil, err
	}

	res := make([]*Thread, 0, 5)
	for rows.Next() {
		if err := rows.StructScan(&rslt); err != nil {
			return nil, err
		}
		res = append(res, &Thread{
			id:     rslt.UUID,
			title:  rslt.Title,
			closed: rslt.Closed,
		})
	}

	return res, nil
}

func (p *psqlRepository) create(ctx context.Context, req repositoryCreateRequest) (Thread, error) {
	query := `
insert into threads (
	title
)
values (
	:title
)
returning
		uuid
	, title
	, closed
	`

	args := struct {
		Title string
	}{
		Title: req.title,
	}

	rslt := struct {
		UUID   string `db:"uuid"`
		Title  string `db:"title"`
		Closed bool   `db:"closed"`
	}{}

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return Thread{}, err
	}
	defer tx.Rollback()

	rows, err := p.db.NamedQueryContext(ctx, query, &args)
	if err != nil {
		return Thread{}, err
	}
	if !rows.Next() {
		return Thread{}, errors.New("failed")
	}
	if err := rows.StructScan(&rslt); err != nil {
		return Thread{}, err
	}

	tx.Commit()

	return Thread{
		id:     rslt.UUID,
		title:  rslt.Title,
		closed: rslt.Closed,
	}, nil
}
