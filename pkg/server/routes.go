package server

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/openstatushq/openqueue/pkg/api/apiv1"
	"github.com/rs/zerolog/log"
)


func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(LogMiddleware)
	r.Use(middleware.Recoverer)
	r.Get("/health", s.healthHandler)

	a:= apiv1.NewAPIv1(s.dbs)
	apiv1.RegisterAPIv1(r, a)

	return r
}

// healthHandler responds with the health status of the server.
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {

	render.JSON(w, r, map[string]any{
		"status": "ok",
	})
	render.Status(r, http.StatusOK)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		t1 := time.Now()
		defer func() {
			t2 := time.Now()

			// Recover and record stack traces in case of a panic
			if rec := recover(); rec != nil {
				log.Ctx(r.Context()).Error().
					Str("type", "error").
					Timestamp().
					Interface("recover_info", rec).
					Bytes("debug_stack", debug.Stack()).
					Msg("log system error")
				http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			// log end request
			log.Ctx(r.Context()).Info().
				Str("type", "access").
				Timestamp().
				Fields(map[string]any{
					"request_id": middleware.GetReqID(r.Context()),
					"remote_ip":  r.RemoteAddr,
					"url":        r.URL.Path,
					"proto":      r.Proto,
					"method":     r.Method,
					"user_agent": r.Header.Get("User-Agent"),
					"status":     ww.Status(),
					"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
					"bytes_in":   r.Header.Get("Content-Length"),
					"bytes_out":  ww.BytesWritten(),
				}).
				Msg("incoming_request")
		}()

		next.ServeHTTP(ww, r)
	})
}
