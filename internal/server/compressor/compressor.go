package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

func RequestCompress(h http.Handler) http.Handler {
	zipFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") && !strings.Contains(r.Header.Get("Content-Type"), "text/html") {
			h.ServeHTTP(w, r)
			return
		}

		logger.Log.Info("Using gzip")

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}

	return http.HandlerFunc(zipFn)
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
