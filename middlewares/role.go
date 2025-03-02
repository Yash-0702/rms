package middlewares

import (
	"net/http"
	"rms/utils"
)

func ShouldHaveRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := UserContext(r).Role
			if userRole != role {
				utils.ResponseError(w, http.StatusUnauthorized, nil, "unauthorized")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
