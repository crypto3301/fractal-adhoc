package server

import (
	"encoding/json"
	"net/http"

	"fractal-adhoc/graph"
	"fractal-adhoc/web"
)

func RegisterHandlers(mux *http.ServeMux, g graph.GraphData) {
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/graph", graphHandler(g))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(web.Page))
}

func graphHandler(g graph.GraphData) http.HandlerFunc {
	payload, err := json.Marshal(g)
	if err != nil {
		panic("graph: failed to marshal graph data: " + err.Error())
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}
}
