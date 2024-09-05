package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/server/adapters"
	"github.com/lvestera/yandex-metrics/internal/server/handlers"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
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
	if err := logger.Initialize(); err != nil {
		return err
	}
	return http.ListenAndServe(s.Addr, MetricRouter(s.ms))
}

func MetricRouter(metric storage.Repository) chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/update/{mtype}/{name}/{value}", logger.RequestLogger(handlers.UpdateHandler{Ms: metric, Format: adapters.Http{}}))
	r.Method(http.MethodGet, "/value/{mtype}/{name}", logger.RequestLogger(handlers.ViewHandler{Ms: metric, Format: adapters.Http{}}))
	r.Method(http.MethodGet, "/", logger.RequestLogger(handlers.ListHandler{Ms: metric}))

	r.Method(http.MethodPost, "/update/", logger.RequestLogger(handlers.UpdateHandler{Ms: metric, Format: adapters.Json{}}))
	r.Method(http.MethodPost, "/value/", logger.RequestLogger(handlers.ViewHandler{Ms: metric, Format: adapters.Json{}}))

	return r
}
