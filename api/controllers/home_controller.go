package controllers

import (
	"net/http"

	"github.com/Tchayo/gotickets/api/responses"
)

// Home : description
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")
}
