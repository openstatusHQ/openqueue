package server

import (
	"context"

	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"

	"connectrpc.com/connect"
)

func (s TaskServer) CreateTask(ctx context.Context, req *connect.Request[v1.CreateTaskRequest]) (*connect.Response[v1.CreateTaskResponse], error) {

	return nil, nil
}
