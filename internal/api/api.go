package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/absurek/go-http-servers/internal/chirps"
	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/metrics"
	"github.com/absurek/go-http-servers/internal/users"
)

type Api struct {
	db        *sql.DB
	dbQueries *database.Queries
	metrics   *metrics.Metrics
	logger    *log.Logger

	usersHandler  *users.UsersHandler
	chirpsHandler *chirps.ChirpsHandler
}

func NewApi(db *sql.DB, dbQueries *database.Queries, metrics *metrics.Metrics, logger *log.Logger) *Api {
	usersHandler := users.NewUsersHandler(db, dbQueries, logger)
	chirpsHandler := chirps.NewChirpsHandler(db, dbQueries, logger)

	return &Api{
		db:        db,
		dbQueries: dbQueries,
		metrics:   metrics,
		logger:    logger,

		usersHandler:  usersHandler,
		chirpsHandler: chirpsHandler,
	}
}

func (a *Api) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/healthz", a.GetHealthz)

	mux.HandleFunc("POST /api/users", a.usersHandler.CreateUser)

	mux.HandleFunc("GET /api/chirps", a.chirpsHandler.GetAllChirps)
	mux.HandleFunc("POST /api/chirps", a.chirpsHandler.CreateChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", a.chirpsHandler.GetChirp)
}
