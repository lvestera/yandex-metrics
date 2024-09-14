package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/server/adapters"
	"github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestViewHandler(t *testing.T) {

	metrics := &storage.MemStorage{
		Counters: map[string]int64{
			"counter1": 1,
			"counter2": 2,
		},
		Gauges: map[string]float64{
			"gauge1": 1.1,
			"gauge2": 3,
		},
	}

	router := chi.NewRouter()
	router.Method(http.MethodGet, "/value/{mtype}/{name}", ViewHandler{Ms: metrics, Format: adapters.HTTP{}})

	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		statusCode  int
		contentType string
		response    string
	}

	var testTable = []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Successfully getting counter metric",
			method: http.MethodGet,
			url:    "/value/counter/counter1",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain",
				response:    "1",
			},
		},
		{
			name:   "Successfully getting gauge metric",
			method: http.MethodGet,
			url:    "/value/gauge/gauge1",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/plain",
				response:    "1.1",
			},
		},
		{
			name:   "Failed getting metric: wrong http method",
			method: http.MethodPost,
			url:    "/value/gauge/gauge1",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "",
				response:    "",
			},
		},
		{
			name:   "Failed getting metric: wrong metric type",
			method: http.MethodGet,
			url:    "/value/unexpecting/gauge1",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    http.StatusText(http.StatusNotFound),
			},
		},
		{
			name:   "Failed getting metric: unexisting metrics",
			method: http.MethodGet,
			url:    "/value/gauge/someOtherMetric",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    http.StatusText(http.StatusNotFound),
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(tt.method, ts.URL+tt.url, nil)
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, strings.TrimSuffix(string(respBody), "\n"))

		})
	}
}

func TestViewHandlerJson(t *testing.T) {

	metrics := &storage.MemStorage{
		Counters: map[string]int64{
			"counter1": 1,
			"counter2": 2,
		},
		Gauges: map[string]float64{
			"gauge1": 1.1,
			"gauge2": 3,
		},
	}

	router := chi.NewRouter()
	router.Method(http.MethodPost, "/value/", ViewHandler{Ms: metrics, Format: adapters.JSON{}})

	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		statusCode  int
		contentType string
		response    string
	}

	var testTable = []struct {
		name   string
		method string
		url    string
		body   string
		want   want
	}{
		{
			name:   "Successfully getting counter metric",
			method: http.MethodPost,
			url:    "/value/",
			body:   "{\"id\":\"counter1\",\"type\":\"counter\"}",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"id\":\"counter1\",\"type\":\"counter\",\"delta\":1}",
			},
		},
		{
			name:   "Successfully getting gauge metric",
			method: http.MethodPost,
			url:    "/value/",
			body:   "{\"id\":\"gauge1\",\"type\":\"gauge\"}",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "application/json",
				response:    "{\"id\":\"gauge1\",\"type\":\"gauge\",\"value\":1.1}",
			},
		},
		{
			name:   "Failed getting metric: wrong http method",
			method: http.MethodGet,
			url:    "/value/",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "",
				response:    "",
			},
		},
		{
			name:   "Failed getting metric: wrong metric type",
			method: http.MethodPost,
			url:    "/value/",
			body:   "{\"id\":\"gauge1\",\"type\":\"unexpecting\"}",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    http.StatusText(http.StatusNotFound),
			},
		},
		{
			name:   "Failed getting metric: unexisting metrics",
			method: http.MethodPost,
			url:    "/value/",
			body:   "{\"id\":\"someOtherMetric\",\"type\":\"gauge\"}",
			want: want{
				statusCode:  http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				response:    http.StatusText(http.StatusNotFound),
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest(tt.method, ts.URL+tt.url, strings.NewReader(tt.body))
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, strings.TrimSuffix(string(respBody), "\n"))

		})
	}
}
