package apiv1

import (
	"context"
	"errors"

	"github.com/openstatushq/openqueue/pkg/database"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
	"github.com/rs/zerolog/log"

	"connectrpc.com/connect"
)

func (api *APIv1)  CreateTask(ctx context.Context, req *connect.Request[v1.CreateTaskRequest]) (*connect.Response[v1.CreateTaskResponse], error) {


	log.Ctx(ctx).Debug().Msg("Creating task")

	//  Use Protovalidate ?
	if req.Msg.QueueName == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("queue name is required"))
	}

	i, ok := api.dbs[req.Msg.QueueName]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("queue not found"))
	}

	id, err := database.CreateTask(ctx, i, &database.Job{

	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.CreateTaskResponse]{
		Msg: &v1.CreateTaskResponse{
			TaskId: id,
		},
	}, nil
}
