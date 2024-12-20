package middleware

import (
	"advancedGo/configs"
	"advancedGo/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextPhoneKey key = "ContextPhoneKey"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuth(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthed(w)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.NewJWT(config.Auth.Token).Parse(token)
		if !isValid {
			writeUnauthed(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
