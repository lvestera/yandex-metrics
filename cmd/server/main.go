package main

import (
	"log"

	"github.com/lvestera/yandex-metrics/internal/server"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

func main() {

	cfg, err := server.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	server := &server.Server{
		Cfg: cfg,
	}
	if err := server.Run(); err != nil {
		logger.Log.Fatal("Can't start server: " + err.Error())
	}
}
