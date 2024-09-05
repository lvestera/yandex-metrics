package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/lvestera/yandex-metrics/internal/models"
	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

type MClient interface {
	SendUpdate(mtype string, name string, value string) error
}
type MetricClient struct {
	Host string
}

func (c *MetricClient) SendUpdate(mtype string, name string, value string) error {

	var err error
	var body []byte
	m := models.Metric{ID: name, MType: mtype}

	m.SetValue(value)

	body, err = json.Marshal(m)
	if err == nil {

		url := fmt.Sprint("http://", c.Host, "/update/")
		client := resty.New()

		if body, err = Compress(body); err != nil {
			logger.Log.Error(err.Error())
			return err
		}

		_, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(url)

		logger.Log.Info(fmt.Sprint("Send the", mtype, "metric", name, "to server"))
	}

	if err != nil {
		logger.Log.Error(err.Error())
	}

	return err
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
