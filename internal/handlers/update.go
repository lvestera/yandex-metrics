package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lvestera/yandex-metrics/internal/storage"
)

type MetricsHandlers struct {
	Ms storage.Repository
}

func (mh MetricsHandlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ok bool
	w.Header().Add("Content-Type", "text/plain")

	switch r.PathValue("mtype") {
	case "gauge":
		ok = mh.updateGauge(r.PathValue("name"), r.PathValue("value"))
	case "counter":
		ok = mh.updateCounter(r.PathValue("name"), r.PathValue("value"))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Println(mh.Ms.GetAllMetrics())
}

func (mh MetricsHandlers) updateGauge(name string, mvalue string) bool {
	value, err := strconv.ParseFloat(mvalue, 64)

	mh.Ms.AddGauge(name, value)

	return err == nil
}

func (mh MetricsHandlers) updateCounter(name string, mvalue string) bool {
	value, err := strconv.ParseInt(mvalue, 10, 64)

	mh.Ms.AddCounter(name, value)
	return err == nil
}
