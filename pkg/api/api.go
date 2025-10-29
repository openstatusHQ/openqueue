package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/openstatushq/openqueue/pkg/api/apiv1"
)


func RegisterAPIs() *chi.Mux {

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Healty"))
	})

	apiv1.RegisterAPIv1(r)

	return r
}
