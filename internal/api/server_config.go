package api

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

const metricsTemplate = `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`

type config struct {
	fileServerHits atomic.Int32
}

func NewConfig() *config {
	return &config{}
}

func (ac *config) ServerHitCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (ac *config) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	payload := fmt.Sprintf(metricsTemplate, ac.fileServerHits.Load())
	w.Write([]byte(payload))
}

func (ac *config) ResetServerHits(w http.ResponseWriter, r *http.Request) {
	ac.fileServerHits.Store(0)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("OK\n"))
}
