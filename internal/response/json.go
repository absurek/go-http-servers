package response

import (
	"encoding/json"
	"net/http"
)

func JSON[T any](w http.ResponseWriter, resp T) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	payload, err := json.Marshal(resp)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.Write(payload)
}
