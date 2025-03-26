package server

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/server/adapters"
	"github.com/lvestera/yandex-metrics/internal/server/compressor"
	"github.com/lvestera/yandex-metrics/internal/server/config"
	"github.com/lvestera/yandex-metrics/internal/server/handlers"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/server/sign"
	"github.com/lvestera/yandex-metrics/internal/storage"
	"golang.org/x/sync/errgroup"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	Cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Cfg: cfg,
	}
}

func (s *Server) Run() error {

	ctx, cancel := context.WithCancel(context.Background())
	//init logger
	if err := logger.Initialize(); err != nil {
		cancel()
		return err
	}
	//init storage repository
	repository, err := storage.NewStorageRepository(s.Cfg)
	if err != nil {
		cancel()
		return err
	}

	sign.NewSign(s.Cfg.Key)

	go repository.Save(ctx, s.Cfg.StorageInterval)

	quit := make(chan os.Signal)
	go func() {
		<-quit
		logger.Log.Info("Receive interrupt signal. Server Close")
		cancel()
	}()

	httpServer := &http.Server{
		Addr:    s.Cfg.Addr,
		Handler: MetricRouter(repository),
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Log.Info("Server starts at " + s.Cfg.Addr)
		return httpServer.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func MetricRouter(metric storage.Repository) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger)
	r.Use(compressor.RequestCompress)
	r.Use(compressor.ResponseCompress)
	r.Use(sign.RequestHashCheck)

	r.Method(http.MethodPost, "/update/{mtype}/{name}/{value}", handlers.UpdateHandler{Ms: metric, Format: adapters.HTTP{}})
	r.Method(http.MethodGet, "/value/{mtype}/{name}", handlers.ViewHandler{Ms: metric, Format: adapters.HTTP{}})
	r.Method(http.MethodGet, "/ping", handlers.PingHandler{Ms: metric})
	r.Method(http.MethodGet, "/", handlers.ListHandler{Ms: metric})

	r.Method(http.MethodPost, "/updates/", handlers.UpdateBatchHandler{Ms: metric, Format: adapters.JSON{}})
	r.Method(http.MethodPost, "/update/", handlers.UpdateHandler{Ms: metric, Format: adapters.JSON{}})
	r.Method(http.MethodPost, "/value/", handlers.ViewHandler{Ms: metric, Format: adapters.JSON{}})

	return r
}
