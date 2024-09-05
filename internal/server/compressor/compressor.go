package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/lvestera/yandex-metrics/internal/server/logger"
)

func ResponseCompress(h http.Handler) http.Handler {
	zipFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") &&
			!(strings.Contains(r.Header.Get("Content-Type"), "text/html") || strings.Contains(r.Header.Get("Content-Type"), "html/text")) {
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

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (gz gzipReader) Read(p []byte) (n int, err error) {
	return gz.zr.Read(p)
}
func (gz gzipReader) Close() (err error) {
	if err := gz.r.Close(); err != nil {
		return err
	}
	return gz.zr.Close()
}

func RequestCompress(h http.Handler) http.Handler {
	unzipFn := func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		newBody := new(gzipReader)
		newBody.r = r.Body

		newReader, err := gzip.NewReader(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newBody.zr = newReader

		r.Body = newBody

		defer newBody.Close()

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(unzipFn)
}
