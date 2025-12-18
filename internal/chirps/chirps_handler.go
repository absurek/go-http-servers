package chirps

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/response"
	"github.com/google/uuid"
)

var profanes = [...]string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

type createChirpRequest struct {
	UserID string `json:"user_id"`
	Body   string `json:"body"`
}

type chirpResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChirpsHandler struct {
	db        *sql.DB
	dbQueries *database.Queries
	logger    *log.Logger
}

func NewChirpsHandler(db *sql.DB, dbQueries *database.Queries, logger *log.Logger) *ChirpsHandler {
	return &ChirpsHandler{
		db:        db,
		dbQueries: dbQueries,
		logger:    logger,
	}
}

func cleanBody(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		for _, profane := range profanes {
			if lower == profane {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}

func (h *ChirpsHandler) CreateChirp(w http.ResponseWriter, r *http.Request) {
	var req createChirpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.InvalidRequestBody(w)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.Unauthorized(w)
		return
	}

	if len(req.Body) > 140 {
		response.BadRequest(w, "Chirp is too long")
		return
	}

	cleanedBody := cleanBody(req.Body)
	chirp, err := h.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: userID,
		Body:   cleanedBody,
	})
	if err != nil {
		h.logger.Printf("ERROR(CreateChirp): db create chirp: %v", err)
		response.InternalServerError(w)
		return
	}

	response.JSON(w, http.StatusCreated, chirpResponse{
		ID:        chirp.ID.String(),
		UserID:    chirp.UserID.String(),
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
	})
}

func (h *ChirpsHandler) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := h.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		h.logger.Printf("ERROR(GetAllChirps): db get all chirps: %v", err)
		response.InternalServerError(w)
		return
	}

	var resp []chirpResponse
	for _, chirp := range chirps {
		resp = append(resp, chirpResponse{
			ID:        chirp.ID.String(),
			UserID:    chirp.UserID.String(),
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
		})
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *ChirpsHandler) GetChirp(w http.ResponseWriter, r *http.Request) {
	pathChirpID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(pathChirpID)
	if err != nil {
		h.logger.Printf("ERROR(GetChirp): invalid id: %s", pathChirpID)
		response.BadRequest(w, "invalid chirp id")
		return
	}

	chirp, err := h.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(w)
		default:
			h.logger.Printf("Error(GetChirp): db get chirp by id: %v", err)
			response.InternalServerError(w)
		}

		return
	}

	response.JSON(w, http.StatusOK, chirpResponse{
		ID:        chirp.ID.String(),
		UserID:    chirp.UserID.String(),
		Body:      chirp.Body,
		CreatedAt: chirp.CreatedAt.Time,
		UpdatedAt: chirp.UpdatedAt.Time,
	})
}
