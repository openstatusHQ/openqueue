package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"
)

var schema = `
CREATE TABLE IF NOT EXISTS job (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  method TEXT NOT NULL,
  headers TEXT NOT NULL,
  body TEXT,
  url TEXT,
  created_at INTEGER NOT NULL,
  scheduled_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS job_run (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  job_id INTEGER NOT NULL,
  status text NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  FOREIGN KEY (job_id) REFERENCES job(id)
);
`

func GetDatabase(ctx context.Context, name string) *sqlx.DB {

	log.Ctx(ctx).Debug().Msgf("Creating database %s", name)
	dbName := fmt.Sprintf("file:./data/%s.db", name)
	db, err := sqlx.Open("sqlite", dbName)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
		return nil
	}

	if _, err = db.Exec(schema); err != nil {
		log.Fatal().Err(err).Msg("failed to execute schema")

		return nil
	}
	return db

}

func CreateTask(ctx context.Context, db *sqlx.DB, job *Job) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("failed to get database")
	}
    tx := db.MustBegin()

    r, err := tx.NamedExec("INSERT INTO job (method, headers, body, url, created_at, scheduled_at, updated_at) VALUES (:method, :headers, :body, :url, :created_at, :scheduled_at, :updated_at)", job)
    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("failed to create task: %w", err)
    }
    id, err := r.LastInsertId()
    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("failed to get last insert id: %w", err)
    }
    tx.Commit()

    log.Ctx(ctx).Info().Msgf("Created task with id %d", id)


	return id, nil
}
