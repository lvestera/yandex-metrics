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

	metric, _ := NewMemStorage(false, "")
	var pollCount int64

	_, err := metric.GetMetric("counter", "PollCount")
	assert.NotNil(t, err)
	for _, name := range MetricsName {
		_, err := metric.GetMetric("gauge", name)
		assert.NotNil(t, err)
	}

	go Update(metric, 2)
	time.Sleep(2 * time.Second)

	pollCount = 1
	pollMetric := models.Metric{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}
	val, err := metric.GetMetric("counter", "PollCount")
	assert.Equal(t, nil, err)
	assert.Equal(t, pollMetric, val)

	for _, name := range MetricsName {
		_, err := metric.GetMetric("gauge", name)
		assert.Equal(t, nil, err)
	}

	time.Sleep(2 * time.Second)

	pollCount = 2
	pollMetric = models.Metric{
		ID:    "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}
	val, err = metric.GetMetric("counter", "PollCount")
	assert.Equal(t, nil, err)
	assert.Equal(t, pollMetric, val)
}

type fakeClient struct {
	mock.Mock
}

func (c *fakeClient) SendUpdate(m models.Metric) error {
	c.Called(m)
	return nil
}

func TestSend(t *testing.T) {

	metric, _ := NewMemStorage(false, "")

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
