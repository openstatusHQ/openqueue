package task_runner

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/openstatushq/openqueue/pkg/database"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
	"github.com/rs/zerolog/log"
)

type QueueOpts struct {
	Retry int
	Db    *sqlx.DB
}

func isSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func Task(ctx context.Context, queueOpts QueueOpts, task *v1.Task, taskId string) error {

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	log.Ctx(ctx).Debug().Msgf("Starting task %s", taskId)

	operation := func() (string, error) {

		execId, err := uuid.NewV7()
		if err != nil {

			return "", err
		}
		id := execId.String()
		_, err = database.AddTaskExecution(ctx, queueOpts.Db, &database.TaskExecution{
			Status: database.StatusPending,
			TaskId: (taskId),
			ID:     execId.String(),
		})

		req, err := http.NewRequestWithContext(ctx, task.Method, task.Url, bytes.NewReader([]byte(task.Body)))

		if err != nil {
			return id, err
		}

		for key, value := range task.Headers {
			req.Header.Add(key, value)
		}

		resp, err := httpClient.Do(req)

		if err != nil {
			return id, err
		}
		defer resp.Body.Close()

		// For this HTTP example, client errors are non-retriable.
		if !isSuccessful(resp.StatusCode) {
			err := database.UpdateTaskStatus(ctx, queueOpts.Db, id, database.StatusFailed)
			if err != nil {
				return id, err
			}

			return id, fmt.Errorf("bad request, status code: %d", resp.StatusCode)
		}
		err = database.UpdateTaskStatus(ctx, queueOpts.Db, id, database.StatusCompleted)

		if err != nil {
			return id, err
		}

		// Return successfulresponse.
		return id, nil
	}

	_, err := backoff.Retry(ctx, operation, backoff.WithMaxTries(uint(queueOpts.Retry)))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Operation is successful.
	log.Ctx(ctx).Debug().Msgf("Task %s completed successfully", taskId)
	return nil
}
