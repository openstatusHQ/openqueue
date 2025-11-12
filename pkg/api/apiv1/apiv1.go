package apiv1

import (
	"github.com/go-chi/chi/v5"
	"github.com/openstatushq/openqueue/pkg/task_runner"
	"github.com/openstatushq/openqueue/proto/gen/api/v1/apiconnect"
)

type APIv1 struct {
	queues map[string]task_runner.QueueOpts
}

func NewAPIv1(queues map[string]task_runner.QueueOpts) *APIv1 {
	return &APIv1{
		queues: queues,
	}
}

func RegisterAPIv1(r *chi.Mux, a *APIv1) {
	path, handler := apiconnect.NewQueueServiceHandler(a)
	r.Group(func(r chi.Router) {
		r.Mount(path, handler)
	})
}
