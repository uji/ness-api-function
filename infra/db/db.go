package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB() *sqlx.DB {
	opt := fmt.Sprintf(
		"host=%s dbname=%s port=%s user=%s password=%s sslmode= disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
	)
	db, err := sqlx.Open("postgres", opt)
	if err != nil {
		panic(err)
	}

	if db.Ping() != nil {
		panic(err)
	}
	return db
}
