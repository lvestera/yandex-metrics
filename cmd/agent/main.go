package main

import (
	"fmt"
	"log"

	"github.com/lvestera/yandex-metrics/internal/agent"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

func main() {

	err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	metric := storage.NewMemStorage()
	if err := metric.Init(false, ""); err != nil {
		log.Fatal(err)
	}
	if err := logger.Initialize(); err != nil {
		log.Fatal(err)
	}

	logger.Log.Info(fmt.Sprint("Client starts at ", addr, " with pollInterval ", pollInterval, " and report interval ", reportInterval))

	c := &agent.MetricClient{
		Host: addr,
	}

	go agent.Update(metric, pollInterval)
	agent.Send(metric, c, reportInterval)
}
