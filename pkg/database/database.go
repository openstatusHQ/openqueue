package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"
)

var schema = `
CREATE TABLE IF NOT EXISTS task (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  method TEXT NOT NULL,
  headers TEXT NOT NULL,
  body TEXT,
  url TEXT,
  created_at INTEGER NOT NULL DEFAULT (unixepoch()),
  scheduled_at INTEGER NOT NULL DEFAULT (unixepoch()),
  updated_at INTEGER NOT NULL DEFAULT (unixepoch())
);

CREATE TABLE IF NOT EXISTS task_execution (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  status text NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (unixepoch()),
  updated_at INTEGER NOT NULL DEFAULT (unixepoch()),
  FOREIGN KEY (task_id) REFERENCES task(id)
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

func CreateTask(ctx context.Context, db *sqlx.DB, task *Task) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("failed to get database")
	}
    tx := db.MustBegin()

    r, err := tx.NamedExec("INSERT INTO task (method, headers, body, url, created_at, scheduled_at, updated_at) VALUES (:method, :headers, :body, :url, :created_at, :scheduled_at, :updated_at)", task)
    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("failed to create task: %w", err)
    }
    id, err := r.LastInsertId()
    if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("failed to get last insert id: %w", err)
    }

    jr := &TaskExecution{
		TaskId:     id,
		Status:    StatusPending,

	}

	r,err = tx.NamedExec("INSERT INTO task_execution (task_id, status) VALUES (:task_id, :status)", jr)
 	if err != nil {
        tx.Rollback()
        return 0, fmt.Errorf("failed to insert task_execution %w", err)
    }

    tx.Commit()

    log.Ctx(ctx).Info().Msgf("Created task with id %d", id)

	return id, nil
}

func GetTask(ctx context.Context, db *sqlx.DB, id int64) (*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("failed to get database")
	}
	job := &Task{}
	err := db.Get(job, "SELECT * FROM task WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return  job, nil
}
