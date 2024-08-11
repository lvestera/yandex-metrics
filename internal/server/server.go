package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/server/handlers"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

type Server struct {
	ms   storage.Repository
	Addr string
}

func (s *Server) Run() error {
	s.ms = &storage.MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}
	return http.ListenAndServe(s.Addr, MetricRouter(s.ms))
}

func MetricRouter(metric storage.Repository) chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/update/{mtype}/{name}/{value}", handlers.UpdateHandler{Ms: metric})
	r.Method(http.MethodGet, "/value/{mtype}/{name}", handlers.ViewHandler{Ms: metric})
	r.Method(http.MethodGet, "/", handlers.ListHandler{Ms: metric})

	return r
}
