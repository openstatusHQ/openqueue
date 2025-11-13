package task_runner

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/openstatushq/openqueue/pkg/database"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
	_ "modernc.org/sqlite"
)

func TestTask_Success(t *testing.T) {
	// Setup a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if _, err = db.Exec(database.Schema); err != nil {
		t.Fatal("cannot create table")
	}

	t.Run("GET Method", func(t *testing.T) {
		task := &v1.Task{
			Url:    server.URL,
			Method: "GET",
		}

		queueOpts := QueueOpts{
			Retry: 1,
			Db:    db,
		}

		err := Task(context.Background(), queueOpts, task, "0")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Get Headers Method", func(t *testing.T) {
		task := &v1.Task{
			Url:    server.URL,
			Method: "GET",
			Headers: map[string]string{
				"test": "test",
			},
		}

		queueOpts := QueueOpts{
			Retry: 1,
			Db:    db,
		}

		err := Task(context.Background(), queueOpts, task, "1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestTask_Failure(t *testing.T) {
	// Setup a test HTTP server that returns 400
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	task := &v1.Task{
		Url:    server.URL,
		Method: "GET",
	}

	queueOpts := QueueOpts{
		Retry: 1,
		Db:    nil,
	}

	err := Task(context.Background(), queueOpts, task, "2")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
