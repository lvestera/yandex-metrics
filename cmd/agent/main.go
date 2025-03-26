package main

import (
	"log"

	"github.com/lvestera/yandex-metrics/internal/agent"
	"github.com/lvestera/yandex-metrics/internal/agent/config"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	agent := agent.NewAgent(cfg)

	if err := agent.Run(); err != nil {
		logger.Log.Fatal("Can't start agent: " + err.Error())
	}
}

// func main1() {

// 	mainCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
// 	defer cancel()

// 	cfg, err := config.NewConfig()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	metric, err := storage.NewMemStorage(false, "")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := logger.Initialize(); err != nil {
// 		log.Fatal(err)
// 	}

// 	logger.Log.Info(fmt.Sprint("Client starts at ", cfg.Addr, " with pollInterval ", cfg.PollInterval, " and report interval ", cfg.ReportInterval))

// 	c := &agent.MetricClient{
// 		Host: cfg.Addr,
// 	}

// 	doneCh := make(chan os.Signal, 1)

// 	go agent.UpdateMainMetrics(metric, cfg.PollInterval)
// 	go agent.UpdateAdditionalMetrics(metric, cfg.PollInterval)

// 	jobs := make(chan struct{})
// 	for w := 1; w <= cfg.RateLimit; w++ {
// 		go agent.Send(jobs, metric, c, cfg.Key)
// 	}

// 	go agent.PrepareSend(jobs, cfg.ReportInterval)

// 	go func() {
// 		<-doneCh
// 		logger.Log.Info("Receive interrupt signal. Server Close")
// 	}()

// 	signal.Notify(doneCh, os.Interrupt, syscall.SIGTERM)
// 	<-doneCh

// 	_, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	fmt.Println("Agent stoped")
// }
