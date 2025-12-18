package website

import (
	"log"
	"net/http"

	"github.com/absurek/go-http-servers/internal/metrics"
)

type Website struct {
	metrics *metrics.Metrics
	logger  *log.Logger
}

func NewWebsite(metrics *metrics.Metrics, logger *log.Logger) *Website {
	return &Website{
		metrics: metrics,
		logger:  logger,
	}
}

func (w *Website) SetupRoutes(mux *http.ServeMux) {
	fileServerHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", w.metrics.ServerHitCounter(fileServerHandler))
}
