package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"
)

var Schema = `
CREATE TABLE IF NOT EXISTS task (
	id TEXT NOT NULL,
	method TEXT NOT NULL,
	headers TEXT NOT NULL,
	body TEXT,
	url TEXT,
	created_at INTEGER NOT NULL DEFAULT (unixepoch()),
	scheduled_at INTEGER NOT NULL DEFAULT (unixepoch()),
	updated_at INTEGER NOT NULL DEFAULT (unixepoch())
);

CREATE TABLE IF NOT EXISTS task_execution (
	id TEXT NOT NULL,
	task_id TEXT NOT NULL,
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

	if _, err = db.Exec(Schema); err != nil {
		log.Fatal().Err(err).Msg("failed to execute schema")

		return nil
	}
	return db

}

func CreateTask(ctx context.Context, db *sqlx.DB, task *Task) (string, error) {
	if db == nil {
		return "", fmt.Errorf("failed to get database")
	}
	tx := db.MustBegin()

	id, err := uuid.NewV7()
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to generate uuid: %w", err)
	}
	//  Let's update task ID
	task.ID = id.String()

	_, err = tx.NamedExec("INSERT INTO task (id, method, headers, body, url) VALUES (:id, :method, :headers, :body, :url)", task)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	newId, err := uuid.NewV7()
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to generate uuid: %w", err)
	}
	jr := &TaskExecution{
		TaskId: newId.String(),
		Status: StatusPending,
	}

	_, err = tx.NamedExec("INSERT INTO task_execution (id, task_id, status) VALUES (:id, :task_id, :status)", jr)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to insert task_execution %w", err)
	}

	tx.Commit()

	log.Ctx(ctx).Info().Msgf("Created task with id %s", id)

	return id.String(), nil
}

func GetTask(ctx context.Context, db *sqlx.DB, id string) (*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("failed to get database")
	}
	job := &Task{}
	err := db.Get(job, "SELECT * FROM task WHERE id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return job, nil
}

func AddTaskExecution(ctx context.Context, db *sqlx.DB, execution *TaskExecution) (string, error) {
	if db == nil {
		return "", fmt.Errorf("failed to get database")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid: %w", err)
	}
	execution.ID = id.String()

	tx := db.MustBegin()
	_, err = tx.NamedExec("INSERT INTO task_execution (id, task_id, status) VALUES (:id, :task_id, :status)", execution)

	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to insert task_execution %w", err)
	}

	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to get last insert id: %w", err)
	}

	tx.Commit()
	return id.String(), nil
}
