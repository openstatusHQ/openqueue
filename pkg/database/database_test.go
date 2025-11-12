package database

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func TestCreateTask(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if _, err = db.Exec(schema); err != nil {
		t.Fatal("cannot create table")
	}

	job := &Task{
		Method:  "POST",
		Headers: "{\"Content-Type\":\"application/json\"}",
		Body:    "{\"foo\":\"bar\"}",
		URL:     "https://openstatus.dev",
	}

	id, err := CreateTask(context.Background(), db, job)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if id == 0 {
		t.Fatalf("CreateTask failed: id is zero")
	}

}

func TestCreateAndGetTask(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if _, err = db.Exec(schema); err != nil {
		t.Fatal("cannot create table")
	}

	job := &Task{
		Method:  "POST",
		Headers: "{\"Content-Type\":\"application/json\"}",
		Body:    "{\"foo\":\"bar\"}",
		URL:     "https://openstatus.dev",
	}

	id, err := CreateTask(context.Background(), db, job)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if id == 0 {
		t.Fatalf("CreateTask failed: id is zero")
	}

	gotJob, err := GetTask(context.Background(), db, id)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}
	if gotJob.URL != job.URL {
		t.Fatalf("GetTask failed: expected URL %s, got %s", job.URL, gotJob.URL)
	}

}
