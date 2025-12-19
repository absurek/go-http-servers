package response

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	ErrorText string `json:"error"`
}

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

func JSON[T any](w http.ResponseWriter, status int, resp T) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	payload, err := json.Marshal(resp)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.Write(payload)
}

func NotFound(w http.ResponseWriter) {
	JSON(w, http.StatusNotFound, errorResponse{
		ErrorText: "not found",
	})
}

func BadRequest(w http.ResponseWriter, errorText string) {
	JSON(w, http.StatusBadRequest, errorResponse{
		ErrorText: errorText,
	})
}

func InvalidRequestBody(w http.ResponseWriter) {
	BadRequest(w, "invalid request body")
}

func Unauthorized(w http.ResponseWriter) {
	JSON(w, http.StatusUnauthorized, errorResponse{
		ErrorText: "unauthorized",
	})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
