package main

import (
	"net/http"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/{mtype}/{name}/{value}", updateHandler)

	return http.ListenAndServe(`:8080`, mux)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var ok bool

	switch r.PathValue("mtype") {
	case "gauge":
		ok = updateGauge(r.PathValue("name"), r.PathValue("value"))
	case "counter":
		ok = updateCounter(r.PathValue("name"), r.PathValue("value"))
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateGauge(name string, mvalue string) bool {
	_, err := strconv.ParseFloat(mvalue, 64)
	return err == nil
}

func updateCounter(name string, mvalue string) bool {
	_, err := strconv.ParseInt(mvalue, 10, 64)
	return err == nil
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

func (ms MemStorage) AddGauge(name string, value float64) {

}

func (ms MemStorage) AddCounter(name string, value int64) {

}
