package response

import "net/http"

func NotFound(w http.ResponseWriter) {
	JSON(w, http.StatusNotFound, errorResponse{
		ErrorText: "not found",
	})
}
