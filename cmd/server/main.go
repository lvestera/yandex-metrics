package main

import (
	"github.com/lvestera/yandex-metrics/internal/server"
)

func main() {

	parseFlags()
	server := &server.Server{
		Addr: addr,
	}
	if err := server.Run(); err != nil {
		panic(err)
	}
}
