package apiv1

import (
	"context"
	"errors"
	"fmt"

	"github.com/openstatushq/openqueue/pkg/database"
	"github.com/openstatushq/openqueue/pkg/task_runner"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
	"github.com/rs/zerolog/log"

	"connectrpc.com/connect"
)

func (api *APIv1) CreateTask(ctx context.Context, req *connect.Request[v1.CreateTaskRequest]) (*connect.Response[v1.CreateTaskResponse], error) {

	log.Ctx(ctx).Debug().Msg("Creating task")

	//  Use Protovalidate ?
	if req.Msg.QueueName == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("queue name is required"))
	}

	i, ok := api.queues[req.Msg.QueueName]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("queue not found"))
	}

	task := &database.Task{
		Method:  req.Msg.GetTask().Method.String(),
		Headers: fmt.Sprintf("%v", req.Msg.GetTask().Headers),
		Body:    req.Msg.GetTask().Body,
		URL:     req.Msg.GetTask().Url,
	}
	id, err := database.CreateTask(ctx, i.Db, task)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	go task_runner.Task(ctx, i, req.Msg.Task, id)

	return &connect.Response[v1.CreateTaskResponse]{
		Msg: &v1.CreateTaskResponse{
			TaskId: id,
		},
	}, nil
}
