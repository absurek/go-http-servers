package main

import (
	"net/http"

	"github.com/absurek/go-http-servers/internal/api"
)

const addr = ":8080"

func main() {
	apiCfg := api.NewConfig()

	mux := &http.ServeMux{}
	server := http.Server{
		Handler: mux,
		Addr:    addr,
	}

	mux.HandleFunc("GET /api/healthz", api.GetHealthz)
	mux.HandleFunc("GET /admin/metrics", apiCfg.GetMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.ResetServerHits)
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirp)

	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.ServerHitCounter(fileServerHandler))

	server.ListenAndServe()
}
