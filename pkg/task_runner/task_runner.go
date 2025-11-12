package task_runner

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/jmoiron/sqlx"
)

type QueueOpts struct {
	Retry int
	Db    *sqlx.DB
}

func isSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func Task(ctx context.Context, queueOpts QueueOpts) {

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	operation := func() (string, error) {

		req, err := http.NewRequestWithContext(ctx, "GET", "http://httpbin.org/get", nil)

		// An example request that may fail.
		resp, err := httpClient.Do(req)

		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		// In case on non-retriable error, return Permanent error to stop retrying.
		// For this HTTP example, client errors are non-retriable.
		if !isSuccessful(resp.StatusCode) {
			return "", fmt.Errorf("bad request, status code: %d", resp.StatusCode)
		}
		// Return successful response.
		return "hello", nil
	}

	result, err := backoff.Retry(context.TODO(), operation, backoff.WithMaxTries(uint(queueOpts.Retry)))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Operation is successful.

	fmt.Println(result)

}
