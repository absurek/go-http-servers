package polka

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/absurek/go-http-servers/internal/auth"
	"github.com/absurek/go-http-servers/internal/database"
	"github.com/absurek/go-http-servers/internal/response"
	"github.com/absurek/go-http-servers/internal/settings"
	"github.com/google/uuid"
)

type PolkaHandler struct {
	settings  settings.Settings
	db        *sql.DB
	dbQueries *database.Queries
	logger    *log.Logger
}

func NewPolkaHandler(s settings.Settings, db *sql.DB, dbQueries *database.Queries, logger *log.Logger) *PolkaHandler {
	return &PolkaHandler{
		settings:  s,
		db:        db,
		dbQueries: dbQueries,
		logger:    logger,
	}
}

type webhooksRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserId string `json:"user_id"`
	} `json:"data"`
}

func (h *PolkaHandler) Webhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != h.settings.PolkaKey {
		response.Unauthorized(w)
		return
	}

	var req webhooksRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.InvalidRequestBody(w)
		return
	}

	if req.Event != "user.upgraded" {
		response.NoContent(w)
		return
	}

	userID, err := uuid.Parse(req.Data.UserId)
	if err != nil {
		response.NotFound(w)
		return
	}

	err = h.dbQueries.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			response.NotFound(w)
		default:
			h.logger.Printf("Error(Webhooks): upgrade user to chirpy red (user_id=%s): %v", userID, err)
			response.InternalServerError(w)
		}

		return
	}

	response.NoContent(w)
}
