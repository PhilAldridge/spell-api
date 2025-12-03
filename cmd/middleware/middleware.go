package middleware

import (
	"net/http"
	"time"

	"github.com/PhilAldridge/spell-api/ent"
	"github.com/PhilAldridge/spell-api/internal/auth"
)

func AuthMiddleware(getUser func(r *http.Request, id int) (*ent.User, error)) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authz := r.Header.Get("Authorization")
            if authz == "" {
                http.Error(w, "missing Authorization header", http.StatusUnauthorized)
                return
            }

            const prefix = "Bearer "
            if len(authz) <= len(prefix) || authz[:len(prefix)] != prefix {
                http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
                return
            }

            token := authz[len(prefix):]

            userID,claims, err := auth.ParseAccessToken(token)
			if exp,err:= claims.GetExpirationTime(); err!=nil || exp.Unix() < time.Now().Unix() {
				http.Error(w, "token expired", http.StatusUnauthorized)
			}

            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }

            user, err := getUser(r, userID)
            if err != nil {
                http.Error(w, "user not found", http.StatusUnauthorized)
                return
            }

            ctx := auth.NewContextWithUser(r.Context(), user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func RequireRoles(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := r.Context().Value("user").(*auth.Claims)
            if !ok {
                http.Error(w, "unauthorised", http.StatusUnauthorized)
                return
            }

            allowed := false
            for _, role := range roles {
                if user.AccountType == role {
                    allowed = true
                    break
                }
            }

            if !allowed {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
