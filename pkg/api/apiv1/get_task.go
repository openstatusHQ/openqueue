package apiv1

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/openstatushq/openqueue/pkg/database"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
	"github.com/rs/zerolog/log"
)

func (api *APIv1) GetTask(ctx context.Context, req *connect.Request[v1.GetTaskRequest]) (*connect.Response[v1.GetTaskResponse], error) {
	log.Ctx(ctx).Debug().Msg("Get task")

	//  Use Protovalidate ?
	if req.Msg.TaskId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("task id is required"))
	}

	i, ok := api.queues[req.Msg.QueueName]
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("queue not found"))
	}

	task, err := database.GetTask(ctx, i.Db, req.Msg.TaskId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s, err := database.GetTaskStatus(ctx, i.Db, req.Msg.TaskId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var status v1.TaskStatus

	switch s {
	case database.StatusCompleted:
		status = v1.TaskStatus_TASK_STATUS_COMPLETED
	case database.StatusPending:
		status = v1.TaskStatus_TASK_STATUS_PENDING
	case database.StatusFailed:
		status = v1.TaskStatus_TASK_STATUS_FAILED
	case database.StatusRunning:
		status = v1.TaskStatus_TASK_STATUS_IN_PROGRESS
	default:
		status = v1.TaskStatus_TASK_STATUS_PENDING
	}

	return connect.NewResponse(&v1.GetTaskResponse{
		Task: &v1.Task{
			Url:    task.URL,
			Body:   task.Body,
			Method: task.Method,
		},
		TaskId: task.ID,
		Status: status,
	}), nil
}
