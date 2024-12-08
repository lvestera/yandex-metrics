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

func (c *fakeClient) SendBatchUpdate(metrics []models.Metric) error {
	c.Called(metrics)
	return nil
}

func (c *fakeClient) SendUpdate(m models.Metric) error {
	c.Called(m)
	return nil
}
