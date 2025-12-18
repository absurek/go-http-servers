package response

import "net/http"

func Unauthorized(w http.ResponseWriter) {
	JSON(w, http.StatusUnauthorized, errorResponse{
		ErrorText: "unauthorized",
	})
}
