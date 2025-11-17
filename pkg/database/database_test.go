package database

import (
	"testing"
	"time"

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
	t.Run("Create Task failed ", func(t *testing.T) {

		job := &Task{
			Method:  "POST",
			Headers: "{\"Content-Type\":\"application/json\"}",
			Body:    "{\"foo\":\"bar\"}",
			URL:     "https://openstatus.dev",
		}

		r, err := CreateTask(t.Context(), nil, job)

		if err == nil {
			t.Fatalf("We are expecting error")
		}

		if r != "" {
			t.Fatalf("Should not return id")
		}

	})
	t.Run("create task ok ", func(t *testing.T) {

		job := &Task{
			Method:  "POST",
			Headers: "{\"Content-Type\":\"application/json\"}",
			Body:    "{\"foo\":\"bar\"}",
			URL:     "https://openstatus.dev",
		}

		id, err := CreateTask(t.Context(), db, job)
		if err != nil {
			t.Fatalf("CreateTask failed: %v", err)
		}
		if id == "" {
			t.Fatalf("CreateTask failed: id is empty")
		}

		gotJob, err := GetTask(t.Context(), db, id)
		if err != nil {
			t.Fatalf("CreateTask failed: %v", err)
		}
		if gotJob.URL != job.URL {
			t.Fatalf("GetTask failed: expected URL %s, got %s", job.URL, gotJob.URL)
		}
	})

	t.Run("Get Task execution ", func(t *testing.T) {

		job := &Task{
			Method:  "POST",
			Headers: "{\"Content-Type\":\"application/json\"}",
			Body:    "{\"foo\":\"bar\"}",
			URL:     "https://openstatus.dev",
		}

		id, err := CreateTask(t.Context(), db, job)

		if err != nil {
			t.Fatalf("could not set up test")
		}

		_, err = AddTaskExecution(t.Context(), db, &TaskExecution{
			TaskId: id,
			Status: StatusPending,
		})

		if err != nil {
			t.Fatalf("could not create task execution %v", err)
		}

		status, err := GetTaskStatus(t.Context(), db, id)
		if err != nil {

			t.Fatalf("could not get Task %v", err)
		}
		if status != StatusPending {
			t.Fatalf("Status should be equal to pending but got %s", status)
		}
	})
	t.Run("Get Task status with multiple execution  ", func(t *testing.T) {

		job := &Task{
			Method:  "POST",
			Headers: "{\"Content-Type\":\"application/json\"}",
			Body:    "{\"foo\":\"bar\"}",
			URL:     "https://openstatus.dev",
		}

		id, err := CreateTask(t.Context(), db, job)

		if err != nil {
			t.Fatalf("could not set up test")
		}

		_, err = AddTaskExecution(t.Context(), db, &TaskExecution{
			TaskId: id,
			Status: StatusPending,
		})

		if err != nil {
			t.Fatalf("could not set up test")
		}

		time.Sleep(100 * time.Millisecond)
		_, err = AddTaskExecution(t.Context(), db, &TaskExecution{
			TaskId: id,
			Status: StatusFailed,
		})
		if err != nil {
			t.Fatalf("could not set up test")
		}

		time.Sleep(100 * time.Millisecond)

		_, err = AddTaskExecution(t.Context(), db, &TaskExecution{
			TaskId: id,
			Status: StatusCompleted,
		})

		if err != nil {
			t.Fatalf("could not create task execution %v", err)
		}

		status, err := GetTaskStatus(t.Context(), db, id)
		if err != nil {

			t.Fatalf("could not get Task %v", err)
		}
		if status != StatusCompleted {
			t.Fatalf("Status should be equal to completed but got %s", status)
		}
	})

}
