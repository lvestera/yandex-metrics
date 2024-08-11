package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lvestera/yandex-metrics/internal/storage"
)

type UpdateHandler struct {
	Ms storage.Repository
}

func (uh UpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ok bool

	mType := chi.URLParam(r, "mtype")
	mName := chi.URLParam(r, "name")
	mValue := chi.URLParam(r, "value")

	w.Header().Add("Content-Type", "text/plain")

	switch mType {
	case "gauge":
		ok = uh.updateGauge(mName, mValue)
	case "counter":
		ok = uh.updateCounter(mName, mValue)
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Println(uh.Ms.GetAllMetrics())
}

func (uh UpdateHandler) updateGauge(name string, mvalue string) bool {
	value, err := strconv.ParseFloat(mvalue, 64)

	uh.Ms.AddGauge(name, value)

	return err == nil
}

func (uh UpdateHandler) updateCounter(name string, mvalue string) bool {
	value, err := strconv.ParseInt(mvalue, 10, 64)

	uh.Ms.AddCounter(name, value)
	return err == nil
}
