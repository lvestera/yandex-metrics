package handlers

import (
	"html/template"
	"net/http"

	"github.com/lvestera/yandex-metrics/internal/storage"
)

type ListHandler struct {
	Ms storage.Repository
}

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Metrics</title>
	</head>
	<body>
		<h2>Gauges</h2>
		<ul>
		{{range $key, $value := .Gauges}}<li>{{ $key }} - {{ $value }}</li>{{else}}<div><strong>no rows</strong></div>{{end}}
		</ul>

		<h2>Counters</h2>
		<ul>
		{{range $key, $value := .Counters}}<li>{{ $key }} - {{ $value }}</li>{{else}}<div><strong>no rows</strong></div>{{end}}
		</ul>
	</body>
</html>`

type ViewData struct {
	Gauges   map[string]string
	Counters map[string]string
}

func (h ListHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	gauges := make(map[string]string)
	counters := make(map[string]string)
	var mValue string

	for _, m := range h.Ms.GetMetrics() {
		mValue, _ = m.GetValue()
		if m.MType == "gauge" {
			gauges[m.ID] = mValue
		}
		if m.MType == "counter" {
			counters[m.ID] = mValue
		}
	}

	data := ViewData{
		Gauges:   gauges,
		Counters: counters,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
