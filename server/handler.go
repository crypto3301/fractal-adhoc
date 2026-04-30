package server

import (
	"encoding/json"
	"net/http"
	"sync"

	"fractal-adhoc/graph"
	"fractal-adhoc/web"
)

var (
	mu      sync.RWMutex
	payload []byte
)

// RegisterHandlers регистрирует маршруты:
//
//	GET  /        → HTML-страница
//	GET  /graph   → JSON с данными графа
//	POST /rebuild → пересобрать граф с новым числом узлов
func RegisterHandlers(mux *http.ServeMux, g graph.GraphData) {
	setPayload(g)
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/graph", graphHandler)
	mux.HandleFunc("/rebuild", rebuildHandler)
}

func setPayload(g graph.GraphData) {
	data, err := json.Marshal(g)
	if err != nil {
		panic("graph marshal error: " + err.Error())
	}
	mu.Lock()
	payload = data
	mu.Unlock()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(web.Page))
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.RLock()
	defer mu.RUnlock()
	w.Write(payload)
}

// rebuildHandler принимает POST /rebuild с телом {"n": 128}
// и пересобирает граф без перезапуска сервера.
func rebuildHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		N int `json:"n"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.N < 4 {
		http.Error(w, `{"error":"n must be >= 4"}`, http.StatusBadRequest)
		return
	}
	g := graph.Build(req.N)
	setPayload(g)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order": g.Order,
		"n":     g.N,
	})
}
