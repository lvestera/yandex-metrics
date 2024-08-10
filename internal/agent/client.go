package agent

import (
	"bytes"
	"fmt"
	"net/http"
)

type MClient interface {
	SendUpdate(mtype string, name string, value string) error
}
type MetricClient struct {
	Host string
}

func (c *MetricClient) SendUpdate(mtype string, name string, value string) error {
	url := fmt.Sprint(c.Host, "/update/", mtype, "/", name, "/", value)
	resp, err := http.Post(url, "text/plain", bytes.NewReader([]byte("")))
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

	return err
}
