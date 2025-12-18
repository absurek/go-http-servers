package admin

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/metrics"
)

type Admin struct {
	db        *sql.DB
	dbQueries *database.Queries
	logger    *log.Logger
	metrics   *metrics.Metrics
}

func NewAdmin(db *sql.DB, dbQueries *database.Queries, metrics *metrics.Metrics, logger *log.Logger) *Admin {
	return &Admin{
		db:        db,
		dbQueries: dbQueries,
		metrics:   metrics,
		logger:    logger,
	}
}

func (a *Admin) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/reset", a.Reset)
	mux.HandleFunc("/admin/metrics", a.metrics.GetMetrics)
}
