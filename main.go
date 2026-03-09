package main

import (
	"log"
	"net/http"

	"fractal-adhoc/graph"
	"fractal-adhoc/server"
)

func main() {
	const order = 3

	g := graph.Build(order)

	log.Printf(
		"Предфрактальный граф K4 порядка %d: %d узлов, %d рёбер",
		order, len(g.Nodes), len(g.Edges),
	)

	mux := http.NewServeMux()
	server.RegisterHandlers(mux, g)

	log.Println("Откройте в браузере: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
