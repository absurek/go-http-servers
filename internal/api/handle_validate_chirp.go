package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/absurek/go-http-servers/internal/response"
)

var profanes = [...]string{
	"kerfuffle",
	"sharbert",
	"fornax",
}

type validateChirpRequest struct {
	Body string `json:"body"`
}

type validateChirpResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	var req validateChirpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.InvalidRequestBody(w)
		return
	}

	if len(req.Body) > 140 {
		response.BadRequest(w, "Chirp is too long")
		return
	}

	words := strings.Split(req.Body, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		for _, profane := range profanes {
			if lower == profane {
				words[i] = "****"
			}
		}
	}

	response.JSON(w, validateChirpResponse{
		CleanedBody: strings.Join(words, " "),
	})
}
