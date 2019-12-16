package middlewares

import (
	"github.com/Tchayo/gotickets/gotickets/api/responses"
	"errors"
	"net/http"

	
	"github.com/Tchayo/gotickets/api/auth"
	"github.com/Tchayo/gotickets/api/responses"
)

func SetMiddlewareJSON(next http.HandleFunc) http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandleFunc) http.HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}