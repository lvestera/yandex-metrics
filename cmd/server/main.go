package main

import (
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/handlers"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {

	metric := &storage.MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}

	mh := handlers.MetricsHandlers{Ms: metric}

	mux := http.NewServeMux()
	mux.Handle("POST /update/{mtype}/{name}/{value}", mh)

	return http.ListenAndServe(`:8080`, mux)
}
