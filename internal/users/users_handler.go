package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/response"
)

type createUserRequest struct {
	Email string `json:"email"`
}

type createUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersHandler struct {
	db        *sql.DB
	dbQueries *database.Queries
	logger    *log.Logger
}

func NewUsersHandler(db *sql.DB, dbQueries *database.Queries, logger *log.Logger) *UsersHandler {
	return &UsersHandler{
		db:        db,
		dbQueries: dbQueries,
		logger:    logger,
	}
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.InvalidRequestBody(w)
		return
	}

	user, err := h.dbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		h.logger.Printf("ERROR(CreateUser): db create user: %v", err)
		response.InternalServerError(w)
		return
	}

	response.JSON(w, http.StatusCreated, createUserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	})
}
