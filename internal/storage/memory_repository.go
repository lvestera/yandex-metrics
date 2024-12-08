package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
	rwm      sync.RWMutex
	filepath string
}

func NewMemStorage(restore bool, filepath string) (*MemStorage, error) {
	ms := new(MemStorage)
	err := ms.Init(restore, filepath)
	return ms, err
}

func (ms *MemStorage) Init(restore bool, filepath string) error {
	ms.Gauges = make(map[string]float64)
	ms.Counters = make(map[string]int64)
	ms.filepath = filepath

	if restore {
		file, err := os.OpenFile(ms.filepath, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		if len(data) != 0 {
			var result []models.Metric

			jsonErr := json.Unmarshal(data, &result)
			if jsonErr != nil {
				return err
			}

			for _, elem := range result {
				if elem.MType == "gauge" {
					ms.AddGauge(elem.ID, *elem.Value)
				} else {
					ms.AddCounter(elem.ID, *elem.Delta)
				}
			}
		}
	}

	return nil
}

func (ms *MemStorage) SetGauges(gauges map[string]float64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()
	for name, value := range gauges {
		ms.Gauges[name] = value
	}
}

func (ms *MemStorage) AddGauge(name string, value float64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()
	ms.Gauges[name] = value
}

func (ms *MemStorage) AddCounter(name string, value int64) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()
	ms.Counters[name] += value
}

func (ms *MemStorage) AddMetrics(metrics []models.Metric) (int, error) {
	ms.rwm.Lock()
	defer ms.rwm.Unlock()

	count := 0

	for _, m := range metrics {
		ok, err := ms.AddMetric(m)

		if ok {
			count = count + 1
		}

		if err != nil {
			return count, err
		}
	}

	return count, nil
}

func (ms *MemStorage) AddMetric(m models.Metric) (bool, error) {
	switch m.MType {
	case "gauge":
		ms.AddGauge(m.ID, *m.Value)
		return true, nil
	case "counter":
		ms.AddCounter(m.ID, *m.Delta)
		return true, nil
	default:
		return false, errors.New("incorrect metric type")
	}
}

func (ms *MemStorage) GetMetric(mtype string, name string) (m models.Metric, err error) {
	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	m = models.Metric{ID: name, MType: mtype}
	switch mtype {
	case "gauge":
		val, ok := ms.Gauges[name]
		if !ok {
			return m, errors.New("not found")
		}

		m.Value = &val
		return m, nil
	case "counter":
		val, ok := ms.Counters[name]
		if !ok {
			return m, errors.New("not found")
		}

		m.Delta = &val
		return m, nil
	default:
		return m, errors.New("incorrect metric type")
	}
}

func (ms *MemStorage) GetMetrics() ([]models.Metric, error) {
	metrics := []models.Metric{}

	ms.rwm.RLock()
	defer ms.rwm.RUnlock()
	for name, elem := range ms.Gauges {
		m := models.Metric{ID: name, MType: "gauge", Value: &elem}
		metrics = append(metrics, m)
	}
	for name, elem := range ms.Counters {
		m := models.Metric{ID: name, MType: "counter", Delta: &elem}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (ms *MemStorage) Save(interval int) error {
	for {
		runtime.Gosched()

		data, err := ms.GetMetrics()
		if err != nil {
			return err
		}

		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return err
		}

		if len(data) > 0 {
			err := func() error {
				file, err := os.OpenFile(ms.filepath, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					return err
				}
				defer file.Close()

				_, err = file.Write(jsonData)
				if err != nil {
					return err
				}
				return nil
			}()

			// file, err := os.OpenFile(ms.filepath, os.O_WRONLY|os.O_CREATE, 0666)
			// if err != nil {
			// 	return err
			// }
			// defer file.Close()

			// _, err = file.Write(jsonData)
			if err != nil {
				logger.Log.Info("Can't save into file")
				return err
			}

			logger.Log.Info("Save into file")
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
