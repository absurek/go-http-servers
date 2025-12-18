package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/absurek/go-http-servers/internal/auth"
	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/response"
)

const maxExpiresIn = 1 * time.Hour

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type loginResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     string    `json:"token"`
}

type UsersHandler struct {
	jwtSecret string
	db        *sql.DB
	dbQueries *database.Queries
	logger    *log.Logger
}

func NewUsersHandler(jwtSecret string, db *sql.DB, dbQueries *database.Queries, logger *log.Logger) *UsersHandler {
	return &UsersHandler{
		jwtSecret: jwtSecret,
		db:        db,
		dbQueries: dbQueries,
		logger:    logger,
	}
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.InvalidRequestBody(w)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.Printf("ERROR(CreateUser): hash password: %v", err)
		response.InternalServerError(w)
		return
	}

	user, err := h.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		h.logger.Printf("ERROR(CreateUser): db create user: %v", err)
		response.InternalServerError(w)
		return
	}

	response.JSON(w, http.StatusCreated, userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	})
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.BadRequest(w, "invalid request schema")
		return
	}

	user, err := h.dbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.Unauthorized(w)
		default:
			h.logger.Printf("Error(Login): db get user by email: %v", err)
			response.InternalServerError(w)
		}

		return
	}

	isValidPassword, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil {
		h.logger.Printf("Error(Login): check password: %v", err)
		response.InternalServerError(w)
		return
	}

	if !isValidPassword {
		response.Unauthorized(w)
		return
	}

	expiresIn := max(time.Duration(req.ExpiresInSeconds)*time.Second, maxExpiresIn)
	jwt, err := auth.MakeJWT(user.ID, h.jwtSecret, expiresIn)
	if err != nil {
		h.logger.Printf("Error(Login): make jwt (user_id=%s): %v", user.ID, err)
		response.InternalServerError(w)
	}

	response.JSON(w, http.StatusOK, loginResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Token:     jwt,
	})
}
