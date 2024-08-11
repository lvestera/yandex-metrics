package agent

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type MClient interface {
	SendUpdate(mtype string, name string, value string) error
}
type MetricClient struct {
	Host string
}

func (c *MetricClient) SendUpdate(mtype string, name string, value string) error {

	url := fmt.Sprint("http://", c.Host, "/update/", mtype, "/", name, "/", value)

	client := resty.New()

	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(url)

	if err != nil {
		fmt.Println(err)
	}

	// resp, err := http.Post(url, "text/plain", bytes.NewReader([]byte("")))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer resp.Body.Close()

	return err
}
