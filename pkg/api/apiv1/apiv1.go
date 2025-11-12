package apiv1

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/openstatushq/openqueue/proto/gen/api/v1/apiconnect"
)

type QueueOpts struct {
	Retry int
	Db  *sqlx.DB
}


type APIv1 struct {
	queues  map[string]QueueOpts
}

func NewAPIv1(queues map[string]QueueOpts) *APIv1 {
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
