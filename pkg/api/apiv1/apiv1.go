package apiv1

import (
	"context"

	"connectrpc.com/connect"
	"github.com/go-chi/chi/v5"
	v1 "github.com/openstatushq/openqueue/proto/gen/api/v1"
	"github.com/openstatushq/openqueue/proto/gen/api/v1/apiconnect"
)

type APIv1 struct {
}

func NewAPIv1() *APIv1 {
	return &APIv1{}
}

func (s *APIv1) Push(context.Context, *connect.Request[v1.PushRequest]) (*connect.Response[v1.PushResponse], error) {
	// Implement your Push logic here
	return connect.NewResponse(&v1.PushResponse{}), nil
}

func RegisterAPIv1(r *chi.Mux) {
	server := NewAPIv1()
	path, handler := apiconnect.NewApiHandler(server)
	r.Group(func(r chi.Router) {
		r.Mount(path, handler)
	})
}
