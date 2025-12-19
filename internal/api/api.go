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
	jwtSecret string
	db        *sql.DB
	dbQueries *database.Queries
	metrics   *metrics.Metrics
	logger    *log.Logger

	usersHandler  *users.UsersHandler
	chirpsHandler *chirps.ChirpsHandler
}

func NewApi(jwtSecret string, db *sql.DB, dbQueries *database.Queries, metrics *metrics.Metrics, logger *log.Logger) *Api {
	usersHandler := users.NewUsersHandler(jwtSecret, db, dbQueries, logger)
	chirpsHandler := chirps.NewChirpsHandler(jwtSecret, db, dbQueries, logger)

	return &Api{
		jwtSecret: jwtSecret,
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
	mux.HandleFunc("PUT /api/users", a.usersHandler.UpdateUser)
	mux.HandleFunc("POST /api/login", a.usersHandler.Login)
	mux.HandleFunc("POST /api/refresh", a.usersHandler.Refresh)
	mux.HandleFunc("POST /api/revoke", a.usersHandler.Revoke)

	mux.HandleFunc("GET /api/chirps", a.chirpsHandler.GetAllChirps)
	mux.HandleFunc("POST /api/chirps", a.chirpsHandler.CreateChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", a.chirpsHandler.GetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", a.chirpsHandler.DeleteChirp)
}
