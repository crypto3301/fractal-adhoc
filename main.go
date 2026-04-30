package main

import (
	"flag"
	"log"
	"net/http"

	"fractal-adhoc/graph"
	"fractal-adhoc/server"
)

func main() {
	n := flag.Int("n", 64, "Начальное количество узлов (≥ 64)")
	addr := flag.String("addr", ":8080", "Адрес сервера")
	flag.Parse()

	if *n < 4 {
		*n = 64
	}

	g := graph.Build(*n)
	log.Printf("Предфрактальный граф K4 порядка %d: %d узлов", g.Order, g.N)

	mux := http.NewServeMux()
	server.RegisterHandlers(mux, g)

	log.Printf("Откройте: http://localhost%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
