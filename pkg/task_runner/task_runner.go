package task_runner

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/jmoiron/sqlx"
	"github.com/openstatushq/openqueue/pkg/database"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
)

type QueueOpts struct {
	Retry int
	Db    *sqlx.DB
}

func isSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func Task(ctx context.Context, queueOpts QueueOpts, task *v1.Task, taskId int64) error {

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	var retry = 0
	operation := func() (int, error) {
		retry += 1
		req, err := http.NewRequestWithContext(ctx, task.Method.String(), task.Url, nil)

		if err != nil {
			return 0, err
		}

		for key, value := range task.Headers {
			req.Header.Add(key, value)
		}

		// An example request that may fail.
		resp, err := httpClient.Do(req)

		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		// In case on non-retriable error, return Permanent error to stop retrying.
		// For this HTTP example, client errors are non-retriable.
		if !isSuccessful(resp.StatusCode) {
			_, err := database.AddTaskExecution(ctx, queueOpts.Db, &database.TaskExecution{
				Status: database.StatusFailed,
				TaskId: (taskId),
			})
			if err != nil {
				return 0, err
			}

			return 0, fmt.Errorf("bad request, status code: %d", resp.StatusCode)
		}
		_, err = database.AddTaskExecution(ctx, queueOpts.Db, &database.TaskExecution{
			Status: database.StatusFailed,
			TaskId: (taskId),
		})
		if err != nil {
			return 0, err
		}

		// Return successfulresponse.
		return 1, nil
	}

	result, err := backoff.Retry(ctx, operation, backoff.WithMaxTries(uint(queueOpts.Retry)))
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Operation is successful.

	fmt.Println(result)
	return nil
}
