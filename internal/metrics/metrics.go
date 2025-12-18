package metrics

import (
	"fmt"
	"log"
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

type Metrics struct {
	fileServerHits atomic.Int32
	logger         *log.Logger
}

func NewMetrics(logger *log.Logger) *Metrics {
	return &Metrics{
		logger: logger,
	}
}

func (m *Metrics) Reset() {
	m.fileServerHits.Store(0)
}

func (m *Metrics) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	payload := fmt.Sprintf(metricsTemplate, m.fileServerHits.Load())
	w.Write([]byte(payload))
}

func (m *Metrics) ServerHitCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
