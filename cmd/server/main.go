package main

import (
	"github.com/lvestera/yandex-metrics/internal/server"
)

func main() {
	server := &server.Server{}
	if err := server.Run(); err != nil {
		panic(err)
	}
}
