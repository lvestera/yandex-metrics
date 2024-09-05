package agent

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/lvestera/yandex-metrics/internal/models"
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

		_, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(url)
	}

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (c *MetricClient) SendUpdate1(mtype string, name string, value string) error {

	url := fmt.Sprint("http://", c.Host, "/update/", mtype, "/", name, "/", value)

	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		fmt.Println(err)
	}

	return err
}
