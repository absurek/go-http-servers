package response

import (
	"encoding/json"
	"net/http"
)

func InternalServerError(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	payload, err := json.Marshal(errorResponse{
		ErrorText: "internal server error",
	})
	if err != nil {
		return
	}

	w.Write(payload)
}
