package database

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func TestCreateTask(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if _, err = db.Exec(Schema); err != nil {
		t.Fatal("cannot create table")
	}

	newId, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}

	job := &Task{
		ID:      newId.String(),
		Method:  "POST",
		Headers: "{\"Content-Type\":\"application/json\"}",
		Body:    "{\"foo\":\"bar\"}",
		URL:     "https://openstatus.dev",
	}

	id, err := CreateTask(context.Background(), db, job)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if id == "" {
		t.Fatalf("CreateTask failed: id is nil")
	}

	gotJob, err := GetTask(context.Background(), db, id)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if gotJob.URL != job.URL {
		t.Fatalf("GetTask failed: expected URL %s, got %s", job.URL, gotJob.URL)
	}

}
