package main

import (
	"github.com/lvestera/yandex-metrics/internal/server"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

func main() {

	parseFlags()
	server := &server.Server{
		Addr: addr,
	}
	if err := server.Run(); err != nil {
		logger.Log.Fatal("Can't start server: " + err.Error())
	}
}
