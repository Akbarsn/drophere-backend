package jwt_tools

import (
	"context"
	"net/http"
	"strings"

	"github.com/bccfilkom/drophere-go/domain"
)

func (j *JWTAuthenticator) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// cast inner function to HandlerFunc
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			spaceIdx := strings.IndexByte(authHeader, ' ')
			if spaceIdx < 0 {
				writeGqlError(w, "Invalid Authorization header")
				return
			}
			authHeaderPrefix := authHeader[:spaceIdx]
			authToken := authHeader[spaceIdx+1:]

			if authHeaderPrefix != "bearer" && authHeaderPrefix != "Bearer" {
				writeGqlError(w, "Invalid Authorization header")
				return
			}

			userID, err := j.validateAndGetUserID(authToken)
			if err != nil {
				writeGqlError(w, "Invalid or expired token")
				return
			}

			// get the user from the database
			user, err := j.userRepo.FindByID(userID)
			if err != nil {
				writeGqlError(w, "Server Error")
				return
			}

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// GetAuthenticatedUser finds the user from the context. REQUIRES Middleware to have run.
func (j *JWTAuthenticator) GetAuthenticatedUser(ctx context.Context) *domain.User {
	raw, _ := ctx.Value(userCtxKey).(*domain.User)
	return raw
}
