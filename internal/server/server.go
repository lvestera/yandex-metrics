package server

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/server/adapters"
	"github.com/lvestera/yandex-metrics/internal/server/compressor"
	"github.com/lvestera/yandex-metrics/internal/server/handlers"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	Cfg *Config
}

func (s *Server) Run() error {
	ms := storage.NewMemStorage()
	if err := ms.Init(s.Cfg.Restore, s.Cfg.FileStoragePath); err != nil {
		return err
	}
	if err := logger.Initialize(); err != nil {
		return err
	}

	db, err := sql.Open("pgx", s.Cfg.DBConfig)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	go ms.Save(s.Cfg.StorageInterval)

	quit := make(chan os.Signal)

	go func() {
		<-quit
		logger.Log.Info("Receive interrupt signal. Server Close")
	}()

	logger.Log.Info("Server starts at " + s.Cfg.Addr)
	return http.ListenAndServe(s.Cfg.Addr, MetricRouter(ms, db))
}

func MetricRouter(metric storage.Repository, db *sql.DB) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger)
	r.Use(compressor.RequestCompress)
	r.Use(compressor.ResponseCompress)

	r.Method(http.MethodPost, "/update/{mtype}/{name}/{value}", handlers.UpdateHandler{Ms: metric, Format: adapters.HTTP{}})
	r.Method(http.MethodGet, "/value/{mtype}/{name}", handlers.ViewHandler{Ms: metric, Format: adapters.HTTP{}})
	r.Method(http.MethodGet, "/ping", handlers.PingHandler{Db: db})
	r.Method(http.MethodGet, "/", handlers.ListHandler{Ms: metric})

	r.Method(http.MethodPost, "/update/", handlers.UpdateHandler{Ms: metric, Format: adapters.JSON{}})
	r.Method(http.MethodPost, "/value/", handlers.ViewHandler{Ms: metric, Format: adapters.JSON{}})

	return r
}
