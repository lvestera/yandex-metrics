package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/server/adapters"
	. "github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler(t *testing.T) {

	type want struct {
		statusCode  int
		contentType string
		allMetrics  map[string]map[string]string
	}

	emptyMap := map[string]string{}

	tests := []struct {
		name       string
		requestURL string
		want       want
	}{
		{
			name:       "OK gauge test",
			requestURL: "/update/gauge/metric/1",
			want: want{
				statusCode:  200,
				contentType: "text/plain",
				allMetrics: map[string]map[string]string{
					"gauge": map[string]string{
						"metric": "1",
					},
					"counter": emptyMap,
				},
			},
		},
		{
			name:       "OK counter test",
			requestURL: "/update/counter/metric/1",
			want: want{
				statusCode:  200,
				contentType: "text/plain",
				allMetrics: map[string]map[string]string{
					"gauge": emptyMap,
					"counter": map[string]string{
						"metric": "1",
					},
				},
			},
		},
		{
			name:       "Fail, no metric name",
			requestURL: "/update/counter/",
			want: want{
				statusCode:  404,
				contentType: "text/plain; charset=utf-8",
				allMetrics: map[string]map[string]string{
					"gauge":   emptyMap,
					"counter": emptyMap,
				},
			},
		},
		{
			name:       "Fail, incorrect metric type",
			requestURL: "/update/other/metric/1",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				allMetrics: map[string]map[string]string{
					"gauge":   emptyMap,
					"counter": emptyMap,
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
			uh := UpdateHandler{
				Ms:     m,
				Format: adapters.Http{},
			}

			r := chi.NewRouter()
			r.Method(http.MethodPost, "/update/{mtype}/{name}/{value}", uh)

			ts := httptest.NewServer(r)
			defer ts.Close()

			result, err := ts.Client().Post(fmt.Sprintf("%v%v", ts.URL, tt.requestURL), "text/plain", nil)
			require.NoError(t, err)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			assert.Equal(t, tt.want.allMetrics, uh.Ms.GetAllMetrics())
		})
	}
}

func TestUpdateHandlerJson(t *testing.T) {

	type want struct {
		statusCode  int
		contentType string
		allMetrics  map[string]map[string]string
	}

	emptyMap := map[string]string{}

	tests := []struct {
		name       string
		requestURL string
		body       string
		want       want
	}{
		{
			name:       "OK gauge test",
			requestURL: "/update/",
			body:       "{\"id\":\"metric\",\"type\":\"gauge\",\"value\":1}",
			want: want{
				statusCode:  200,
				contentType: "application/json",
				allMetrics: map[string]map[string]string{
					"gauge": map[string]string{
						"metric": "1",
					},
					"counter": emptyMap,
				},
			},
		},
		{
			name:       "OK counter test",
			requestURL: "/update/",
			body:       "{\"id\":\"metric\",\"type\":\"counter\",\"delta\":1}",
			want: want{
				statusCode:  200,
				contentType: "application/json",
				allMetrics: map[string]map[string]string{
					"gauge": emptyMap,
					"counter": map[string]string{
						"metric": "1",
					},
				},
			},
		},
		{
			name:       "Fail, no metric name",
			requestURL: "/update/",
			body:       "{\"type\":\"counter\",\"delta\":1}",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				allMetrics: map[string]map[string]string{
					"gauge":   emptyMap,
					"counter": emptyMap,
				},
			},
		},
		{
			name:       "Fail, incorrect metric type",
			requestURL: "/update/",
			body:       "{\"id\":\"metric\",\"type\":\"other\",\"value\":1}",
			want: want{
				statusCode:  400,
				contentType: "text/plain; charset=utf-8",
				allMetrics: map[string]map[string]string{
					"gauge":   emptyMap,
					"counter": emptyMap,
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
			uh := UpdateHandler{
				Ms:     m,
				Format: adapters.Json{},
			}

			r := chi.NewRouter()
			r.Method(http.MethodPost, "/update/", uh)

			ts := httptest.NewServer(r)
			defer ts.Close()

			result, err := ts.Client().Post(fmt.Sprintf("%v%v", ts.URL, tt.requestURL), "application/json", strings.NewReader(tt.body))
			require.NoError(t, err)
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			assert.Equal(t, tt.want.allMetrics, uh.Ms.GetAllMetrics())
		})
	}
}
