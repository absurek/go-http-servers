package response

import (
	"encoding/json"
	"net/http"
)

func BadRequest(w http.ResponseWriter, errorText string) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	payload, err := json.Marshal(errorResponse{
		ErrorText: errorText,
	})
	if err != nil {
		InternalServerError(w)
		return
	}

	w.Write(payload)
}

func InvalidRequestBody(w http.ResponseWriter) {
	BadRequest(w, "invalid request body")
}
