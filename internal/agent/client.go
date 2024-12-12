package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

const maxRetries = 3

type MClient interface {
	SendUpdate(m models.Metric) error
	SendBatchUpdate(metrics []models.Metric) error
}
type MetricClient struct {
	Host string
}

func (c *MetricClient) SendUpdate(m models.Metric) error {

	var err error
	var body []byte

	if body, err = json.Marshal(m); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	url := fmt.Sprint("http://", c.Host, "/update/")
	client := resty.New()

	if body, err = Compress(body); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(body).
		Post(url)

	if err != nil {
		logger.Log.Error(err.Error())
	}

	logger.Log.Info(fmt.Sprint("Send the ", m.MType, " metric ", m.ID, " to server"))

	return err
}

func (c *MetricClient) SendBatchUpdate(metrics []models.Metric) error {
	var err error
	var body []byte

	if body, err = json.Marshal(metrics); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	url := fmt.Sprint("http://", c.Host, "/updates/")
	client := resty.New()

	if body, err = Compress(body); err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	delay := 1

	for i := 0; i < maxRetries; i++ {
		_, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(body).
			Post(url)

		if err != nil {
			logger.Log.Error(err.Error())
		} else {
			logger.Log.Info(fmt.Sprint("Send ", len(metrics), " metrics to server"))
			return nil
		}

		time.Sleep(time.Duration(delay))
		delay += 2
	}

	// _, err = client.R().
	// 	SetHeader("Content-Type", "application/json").
	// 	SetHeader("Content-Encoding", "gzip").
	// 	SetBody(body).
	// 	Post(url)

	// if err != nil {
	// 	logger.Log.Error(err.Error())
	// }

	return errors.New(fmt.Sprint("Failed to send ", len(metrics), " metrics to server after ", maxRetries, " attempts"))
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	// создаём переменную w — в неё будут записываться входящие данные,
	// которые будут сжиматься и сохраняться в bytes.Buffer
	w := gzip.NewWriter(&b)

	// запись данных
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	// обязательно нужно вызвать метод Close() — в противном случае часть данных
	// может не записаться в буфер b; если нужно выгрузить все упакованные данные
	// в какой-то момент сжатия, используйте метод Flush()
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	// переменная b содержит сжатые данные
	return b.Bytes(), nil
}
