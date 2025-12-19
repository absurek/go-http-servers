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

const jwtExpiresIn = 1 * time.Hour
const refreshTokenExpiresIn = 60 * 24 * time.Hour

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type loginResponse struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

type refreshResponse struct {
	Token string `json:"token"`
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

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		response.Unauthorized(w)
		return
	}

	userID, err := auth.ValidateJWT(jwt, h.jwtSecret)
	if err != nil {
		response.Unauthorized(w)
		return
	}

	var req userRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.InvalidRequestBody(w)
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		h.logger.Panicf("Error(UpdateUser): hash password (user_id=%s): %v", userID, err)
		response.InternalServerError(w)
		return
	}

	user, err := h.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		h.logger.Printf("Error(UpdateUser): update user (user_id=%s): %v", userID, err)
		response.InternalServerError(w)
		return
	}

	response.JSON(w, http.StatusOK, userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	})
}

func (h *UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req userRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.InvalidRequestBody(w)
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

	jwt, err := auth.MakeJWT(user.ID, h.jwtSecret, jwtExpiresIn)
	if err != nil {
		h.logger.Printf("Error(Login): make jwt (user_id=%s): %v", user.ID, err)
		response.InternalServerError(w)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		h.logger.Printf("Error(Login): make refresh token (user_id=%s): %v", user.ID, err)
		response.InternalServerError(w)
		return
	}

	_, err = h.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refreshTokenExpiresIn),
	})
	if err != nil {
		h.logger.Printf("Error(Login): create refresh token (user_id=%s): %v", user.ID, err)
		response.InternalServerError(w)
		return
	}

	response.JSON(w, http.StatusOK, loginResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		Token:        jwt,
		RefreshToken: refreshToken,
	})
}

func (h *UsersHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		h.logger.Printf("Error(Refresh): get bearer token: %v", err)
		response.Unauthorized(w)
		return
	}

	refreshToken, err := h.dbQueries.GetRefreshToken(r.Context(), bearerToken)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.Unauthorized(w)
		default:
			h.logger.Printf("Error(Refresh): get refresh token: %v", err)
			response.InternalServerError(w)
		}

		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		response.Unauthorized(w)
		return
	}

	jwt, err := auth.MakeJWT(refreshToken.UserID, h.jwtSecret, jwtExpiresIn)
	if err != nil {
		h.logger.Printf("Error(Refres): make jwt (user_id=%s): %v", refreshToken.UserID, err)
		response.InternalServerError(w)
		return
	}

	response.JSON(w, http.StatusOK, refreshResponse{
		Token: jwt,
	})
}

func (h *UsersHandler) Revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		h.logger.Printf("Error(Revoke): get bearer token: %v", err)
		response.Unauthorized(w)
		return
	}

	err = h.dbQueries.RevokeToken(r.Context(), refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.Unauthorized(w)
		default:
			h.logger.Printf("Error(Revoke): revoke token: %v", err)
			response.InternalServerError(w)
		}

		return
	}

	response.NoContent(w)
}
