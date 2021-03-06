package utils

import (
	"context"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// CtxUserString type for using with context
type CtxUserString string

// Claims jwt claim
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// LoginRequired Middleware to protect endpoints
func LoginRequired(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Bearer")
		if tokenStr == "" {
			http.Error(w, "Empty GET request", 400)
			LOG.Println("Empty Get request")
			return
		}

		jwtKey := []byte(os.Getenv("JWT_SECRET_KEY"))

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Empty GET request", http.StatusUnauthorized)
			LOG.Println(err)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			LOG.Println("Invalid Token")
			return
		}

		ctx := r.Context()
		newctx := context.WithValue(ctx, CtxUserString("user"), claims.Username)
		req := r.WithContext(newctx)

		next(w, req)
	}
}
