package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}

	//think about Migration

	stmt1, err := db.Prepare(`create table IF NOT EXISTS users(
									id BIGSERIAL PRIMARY KEY,
									user_name TEXT NOT NULL ,
									created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt1.Close()

	_, err = stmt1.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt2, err := db.Prepare(`create table IF NOT EXISTS notes(
									id BIGSERIAL PRIMARY KEY,
									user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE cascade,
									title TEXT NOT NULL,
									content TEXT NOT NULL,
									created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
									updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt2.Close()

	_, err = stmt2.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil

}
