package agent

import (
	"testing"
	"time"

	. "github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdate(t *testing.T) {

	metric := &MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}
	var pollCount string

	_, ok := metric.GetMetric("counter", "PollCount")
	assert.False(t, ok)
	for _, name := range MetricsName {
		_, ok := metric.GetMetric("gauge", name)
		assert.False(t, ok)
	}

	go Update(metric, 2)
	time.Sleep(2 * time.Second)

	pollCount = "1"
	val, ok := metric.GetMetric("counter", "PollCount")
	assert.True(t, ok)
	assert.Equal(t, pollCount, val)

	for _, name := range MetricsName {
		_, ok := metric.GetMetric("gauge", name)
		assert.True(t, ok)
	}

	time.Sleep(2 * time.Second)

	pollCount = "2"
	val, ok = metric.GetMetric("counter", "PollCount")
	assert.True(t, ok)
	assert.Equal(t, pollCount, val)
}

type fakeClient struct {
	mock.Mock
}

func (c *fakeClient) SendUpdate(mtype string, name string, value string) error {
	c.Called(mtype, name, value)
	return nil
}

func TestSend(t *testing.T) {

	metric := &MemStorage{
		Counters: make(map[string]int64),
		Gauges:   make(map[string]float64),
	}
	metric.AddGauge("mg1", 1)
	metric.AddGauge("mg2", 2)
	metric.AddCounter("mc1", 1)

	var mockClient = &fakeClient{}

	mockClient.On("SendUpdate", "gauge", "mg1", "1").Return(nil)
	mockClient.On("SendUpdate", "gauge", "mg2", "2").Return(nil)
	mockClient.On("SendUpdate", "counter", "mc1", "1").Return(nil)

	go Send(metric, mockClient, 10)
	time.Sleep(10 * time.Second)

	mockClient.AssertNumberOfCalls(t, "SendUpdate", 3)
}
