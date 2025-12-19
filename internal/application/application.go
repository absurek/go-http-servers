package application

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/absurek/go-http-servers/internal/admin"
	"github.com/absurek/go-http-servers/internal/api"
	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/metrics"
	"github.com/absurek/go-http-servers/internal/website"
	_ "github.com/lib/pq"
)

type Application struct {
	jwtSecret string
	db        *sql.DB
	server    *http.Server
	logger    *log.Logger

	metrics *metrics.Metrics
	website *website.Website
	admin   *admin.Admin
	api     *api.Api
}

func NewApplication(addr string, logger *log.Logger) (*Application, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
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

	api := api.NewApi(jwtSecret, db, dbQueries, metr, logger)
	api.SetupRoutes(mux)

	server := &http.Server{
		Handler: mux,
		Addr:    addr,
	}

	return &Application{
		jwtSecret: jwtSecret,
		db:        db,
		server:    server,
		logger:    logger,
		metrics:   metr,
		website:   website,
		admin:     admin,
		api:       api,
	}, nil
}

func (a *Application) ListenAndServe() error {
	return a.server.ListenAndServe()
}

func (a *Application) Close() error {
	return a.db.Close()
}

func (a *Application) ListenForInterrupt() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	sig := <-sigCh
	a.logger.Printf("Received signal %s. Shutting down...\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Fatal("Graceful shutdown failed:", err)
	}
}
