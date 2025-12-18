package admin

import (
	"net/http"

	"github.com/absurek/go-http-servers/internal/response"
)

type resetResponse struct {
	Status string `json:"status"`
}

func (a *Admin) Reset(w http.ResponseWriter, r *http.Request) {
	a.metrics.Reset()

	err := a.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		response.InternalServerError(w)
		return
	}

	response.JSON(w, http.StatusOK, resetResponse{
		Status: "OK",
	})
}
