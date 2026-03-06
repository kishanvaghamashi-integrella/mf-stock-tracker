package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kishanvaghamashi-integrella/mf-stock-tracker/internal/util"
)

func JWTAuth(next http.Handler) http.Handler {
	// Exclude public endpoints from JWT validation
	excludedPaths := []struct {
		method string
		path   string
		prefix bool
	}{
		{"", "/swagger", true},              // "" matching any method
		{"POST", "/api/users/", false},      // user creation
		{"POST", "/api/users/login", false}, // user login
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		for _, ep := range excludedPaths {
			if ep.method != "" && ep.method != r.Method {
				continue
			}

			if (ep.prefix && strings.HasPrefix(path, ep.path)) || (!ep.prefix && path == ep.path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.SendErrorResponse(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			util.SendErrorResponse(w, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		tokenString := parts[1]

		userID, err := util.ValidateToken(tokenString)
		if err != nil {
			util.SendErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// Set the user_id in the request context
		ctx := context.WithValue(r.Context(), util.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
