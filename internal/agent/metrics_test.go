package agent

import (
	"testing"
	"time"

	"github.com/lvestera/yandex-metrics/internal/models"
	. "github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdate(t *testing.T) {

	metric := NewMemStorage()
	metric.Init(false, "")
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

func (c *fakeClient) SendUpdate(m models.Metric) error {
	c.Called(m)
	return nil
}

func TestSend(t *testing.T) {

	metric := NewMemStorage()
	metric.Init(false, "")
	metric.AddGauge("mg1", 1)
	metric.AddGauge("mg2", 2)
	metric.AddCounter("mc1", 1)

	var mockClient = &fakeClient{}

	var x1, x2 float64 = 1, 2
	var y int64 = 1
	gaugeMetric1 := models.Metric{ID: "mg1", MType: "gauge", Value: &x1}
	gaugeMetric2 := models.Metric{ID: "mg2", MType: "gauge", Value: &x2}
	counterMetric1 := models.Metric{ID: "mc1", MType: "counter", Delta: &y}

	mockClient.On("SendUpdate", gaugeMetric1).Return(nil)
	mockClient.On("SendUpdate", gaugeMetric2).Return(nil)
	mockClient.On("SendUpdate", counterMetric1).Return(nil)

	go Send(metric, mockClient, 10)
	time.Sleep(15 * time.Second)

	mockClient.AssertNumberOfCalls(t, "SendUpdate", 6)
}
