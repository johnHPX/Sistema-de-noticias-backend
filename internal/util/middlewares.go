package util

import (
	"net/http"
)

func Authenticate(nextFunction http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := validateToken(r)
		if err != nil {
			CreateHttpErrorResponse(http.StatusUnauthorized, 01, err, "token is invalid")
			return
		}
		nextFunction(w, r)
	}
}
