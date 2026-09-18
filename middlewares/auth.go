package middlewares

import (
    env "AuthInGo/config/env"
    "fmt"
    "net/http"
    "strings"
	"strconv"
	"context"

    "github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorisation")

        if authHeader == "" {
            http.Error(w, "Authorisation Header is Required", http.StatusUnauthorized)
        }

        if !strings.HasPrefix(authHeader, "Bearer") {
            http.Error(w, "Authorised Header must start with Header", http.StatusUnauthorized)
        }

        token := strings.TrimPrefix(authHeader, "Bearer")
        if token == "" {
            http.Error(w, "Token is required", http.StatusUnauthorized)
        }

        claims := jwt.MapClaims{}
        _, err := jwt.ParseWithClaims(token, &claims, func(jwtToken *jwt.Token) (interface{}, error) {
            return []byte(env.GetString("JWT_SECRET", "TOKEN")), nil
        })
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        userId, okId := claims["id"].(float64)
        email, okEmail := claims["email"].(string)

        if !okId || !okEmail {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        fmt.Println("Aunthicated user Id:", int64(userId), "Email:", email)

        
		ctx := context.WithValue(r.Context(), "userID", strconv.FormatFloat(userId, 'f', 0, 64))
		ctx = context.WithValue(ctx, "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
    })
}