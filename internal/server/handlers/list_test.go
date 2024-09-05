package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListHandler(t *testing.T) {

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
	router.Method(http.MethodGet, "/", ListHandler{Ms: metrics})

	ts := httptest.NewServer(router)
	defer ts.Close()

	type want struct {
		statusCode  int
		contentType string
		response    []string
	}

	var testTable = []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Successfully getting list of metric",
			method: http.MethodGet,
			url:    "/",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/html",
				response: []string{
					"<!DOCTYPE html>",
					"counter1 - 1",
					"counter2 - 2",
					"gauge1 - 1.1",
					"gauge2 - 3",
				},
			},
		},
		{
			name:   "Failed getting metric: wrong http method",
			method: http.MethodPost,
			url:    "/",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
				contentType: "",
				response:    []string{},
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
			respStringBody := string(respBody)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			for _, v := range tt.want.response {
				assert.True(t, strings.Contains(respStringBody, v))
			}

		})
	}
}
