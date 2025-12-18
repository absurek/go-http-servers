package application

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/absurek/go-http-servers/internal/admin"
	"github.com/absurek/go-http-servers/internal/api"
	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/metrics"
	"github.com/absurek/go-http-servers/internal/website"
	_ "github.com/lib/pq"
)

type Application struct {
	db     *sql.DB
	server *http.Server
	logger *log.Logger

	metrics *metrics.Metrics
	website *website.Website
	admin   *admin.Admin
	api     *api.Api
}

func NewApplication(addr string, logger *log.Logger) (*Application, error) {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	dbQueries := database.New(db)

	mux := &http.ServeMux{}
	metr := metrics.NewMetrics(logger)

	website := website.NewWebsite(metr, logger)
	website.SetupRoutes(mux)

	admin := admin.NewAdmin(db, dbQueries, metr, logger)
	admin.SetupRoutes(mux)

	api := api.NewApi(db, dbQueries, metr, logger)
	api.SetupRoutes(mux)

	server := &http.Server{
		Handler: mux,
		Addr:    addr,
	}

	return &Application{
		db:      db,
		server:  server,
		logger:  logger,
		metrics: metr,
		website: website,
		admin:   admin,
		api:     api,
	}, nil
}

func (a *Application) ListenAndServe() error {
	return a.server.ListenAndServe()
}

func (a *Application) Close() error {
	return a.db.Close()
}
