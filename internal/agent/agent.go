package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/lvestera/yandex-metrics/internal/agent/config"
	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
	"github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/shirou/gopsutil/v4/mem"
)

type Agent struct {
	Cfg    *config.Config
	Client *MetricClient
}

func NewAgent(cfg *config.Config) *Agent {
	return &Agent{
		Cfg: cfg,
		Client: &MetricClient{
			Host: cfg.Addr,
		},
	}
}

func (a *Agent) Run() error {
	mainCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	repository, err := storage.NewMemStorage(false, "")
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.Initialize(); err != nil {
		log.Fatal(err)
	}

	logger.Log.Info(fmt.Sprint("Client starts at ", a.Cfg.Addr, " with pollInterval ", a.Cfg.PollInterval, " and report interval ", a.Cfg.ReportInterval))

	wg := &sync.WaitGroup{}

	//update main metrics
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		ticker := time.NewTicker(time.Duration(a.Cfg.PollInterval) * time.Second)
		var rtm runtime.MemStats
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.updateMainMetrics(ctx, repository, rtm)
			}
		}
	}(mainCtx)

	//update additional metrics
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		ticker := time.NewTicker(time.Duration(a.Cfg.PollInterval) * time.Second)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				a.updateAdditionalMetrics(ctx, repository)
			}
		}
	}(mainCtx)

	jobs := make(chan struct{})
	for w := 1; w <= a.Cfg.RateLimit; w++ {
		wg.Add(1)
		go func(ctx context.Context, jobs <-chan struct{}) {
			defer wg.Done()
			a.send(ctx, jobs, repository)
		}(mainCtx, jobs)
	}

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		ticker := time.NewTicker(time.Duration(a.Cfg.ReportInterval) * time.Second)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				jobs <- struct{}{}
			}
		}
	}(mainCtx)

	wg.Wait()
	fmt.Println("Agent stoped")

	return nil
}

func (a *Agent) updateMainMetrics(ctx context.Context, m storage.Repository, rtm runtime.MemStats) {
	runtime.ReadMemStats(&rtm)

	m.AddMetrics(collectMetrics(&rtm))

	var pollMetricValue int64 = 1

	m.AddMetric(models.Metric{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollMetricValue,
	})
}

func (a *Agent) updateAdditionalMetrics(ctx context.Context, m storage.Repository) {

	var val float64

	v, _ := mem.VirtualMemory()
	metrics := make([]models.Metric, 0)

	val = float64(v.Total)
	metrics = append(metrics, models.Metric{
		ID:    "TotalMemory",
		MType: "gauge",
		Value: &val,
	})

	val = float64(v.Free)
	metrics = append(metrics, models.Metric{
		ID:    "FreeMemory",
		MType: "gauge",
		Value: &val,
	})
	val = float64(v.UsedPercent)
	metrics = append(metrics, models.Metric{
		ID:    "CPUutilization1",
		MType: "gauge",
		Value: &val,
	})

	m.AddMetrics(metrics)
}

func (a *Agent) send(ctx context.Context, jobs <-chan struct{}, m storage.Repository) {

	for {
		select {
		case <-ctx.Done():
			return
		case <-jobs:
			metrics, err := m.GetMetrics(ctx)
			if err != nil {
				logger.Log.Info("Get metrics failed")
			}

			err = a.Client.SendBatchUpdate(metrics, a.Cfg.Key)
			if err != nil {
				logger.Log.Info(fmt.Sprint("Sending the batch of ", len(metrics), "metrics failed: ", err.Error()))
			}
		}
	}
}
