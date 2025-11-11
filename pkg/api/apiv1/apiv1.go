package apiv1

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/openstatushq/openqueue/proto/gen/api/v1/apiconnect"
)

type APIv1 struct {
	dbs  map[string]*sqlx.DB
}

func NewAPIv1(dbs map[string]*sqlx.DB) *APIv1 {
	return &APIv1{
		dbs: dbs,
	}
}



func RegisterAPIv1(r *chi.Mux, a *APIv1) {
	path, handler := apiconnect.NewQueueServiceHandler(a)
	r.Group(func(r chi.Router) {
		r.Mount(path, handler)
	})
}
