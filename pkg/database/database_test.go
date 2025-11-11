package database

import (
	"context"
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


    if _, err = db.Exec(schema); err != nil {
        t.Fatal("cannot create table")
	}

    job := &Job{
        Method:      "POST",
        Headers:     "{\"Content-Type\":\"application/json\"}",
        Body:        "{\"foo\":\"bar\"}",
        URL:      	 "https://openstatus.dev",
        CreatedAt:   time.Now().Unix(),
        ScheduledAt: time.Now().Unix(),
        UpdatedAt:   time.Now().Unix(),
    }

    id, err := CreateTask(context.Background(), db, job)
    if err != nil {
        t.Fatalf("CreateTask failed: %v", err)
    }
    if id == 0 {
        t.Fatalf("CreateTask failed: id is zero")
    }

}
