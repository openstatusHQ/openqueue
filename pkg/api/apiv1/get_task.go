package apiv1

import (
	"context"

	"connectrpc.com/connect"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
)

func (api *APIv1) GetTask(ctx context.Context, req *connect.Request[v1.GetTaskRequest]) (*connect.Response[v1.GetTaskResponse], error) {
	return nil, nil
}
