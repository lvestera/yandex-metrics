package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMetricsHandlers(t *testing.T) {

	type want struct {
		statusCode  int
		contentType string
		allMetrics  map[string]map[string]string
	}

	tests := []struct {
		name       string
		requestUrl string
		want       want
	}{
		{
			name:       "OK gauge test",
			requestUrl: "/update/gauge/metric/1",
			want: want{
				statusCode:  200,
				contentType: "text/plain",
				allMetrics: map[string]map[string]string{
					"gauge": map[string]string{
						"metric": "1",
					},
					"counter": map[string]string{},
				},
			},
		},
		{
			name:       "OK counter test",
			requestUrl: "/update/counter/metric/1",
			want: want{
				statusCode:  200,
				contentType: "text/plain",
				allMetrics: map[string]map[string]string{
					"gauge": map[string]string{},
					"counter": map[string]string{
						"metric": "1",
					},
				},
			},
		},
		{
			name:       "OK counter test2",
			requestUrl: "/update/counter/metric/1",
			want: want{
				statusCode:  200,
				contentType: "text/plain",
				allMetrics: map[string]map[string]string{
					"gauge": map[string]string{},
					"counter": map[string]string{
						"metric": "1",
					},
				},
			},
		},
		{
			name:       "Fail, no metric name",
			requestUrl: "/update/counter/",
			want: want{
				statusCode:  404,
				contentType: "text/plain; charset=utf-8",
				allMetrics: map[string]map[string]string{
					"gauge":   map[string]string{},
					"counter": map[string]string{},
				},
			},
		},
		{
			name:       "Fail, incorrect metric type",
			requestUrl: "/update/other/metric/1",
			want: want{
				statusCode:  400,
				contentType: "text/plain",
				allMetrics: map[string]map[string]string{
					"gauge":   map[string]string{},
					"counter": map[string]string{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := &MemStorage{
				Counters: make(map[string]int64),
				Gauges:   make(map[string]float64),
			}
			mh := MetricsHandlers{
				Ms: m,
			}

			mux := http.NewServeMux()
			mux.Handle("POST /update/{mtype}/{name}/{value}", mh)

			ts := httptest.NewServer(mux)
			defer ts.Close()

			result, _ := ts.Client().Post(fmt.Sprintf("%v%v", ts.URL, tt.requestUrl), "text/plain", nil)

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			assert.Equal(t, tt.want.allMetrics, mh.Ms.GetAllMetrics())
		})
	}

}
