package middlewares

import (
    env "AuthInGo/config/env"
    repo "AuthInGo/db/repositories"
    "context"
    dbConfig "AuthInGo/config/db"
    "fmt"
    
    "net/http"
    "strconv"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            authHeader = r.Header.Get("Authorisation")
        }

        if authHeader == "" {
            http.Error(w, "Authorization Header is Required", http.StatusUnauthorized)
            return
        }

        if !strings.HasPrefix(authHeader, "Bearer ") && !strings.HasPrefix(authHeader, "Bearer") {
            http.Error(w, "Authorization Header must start with Bearer", http.StatusUnauthorized)
            return
        }

        token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
        if token == "" {
            http.Error(w, "Token is required", http.StatusUnauthorized)
            return
        }

        claims := jwt.MapClaims{}
        parsedToken, err := jwt.ParseWithClaims(token, &claims, func(jwtToken *jwt.Token) (interface{}, error) {
            return []byte(env.GetString("JWT_SECRET", "TOKEN")), nil
        })

        if err != nil || !parsedToken.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        userId, okId := claims["id"].(float64)
        email, okEmail := claims["email"].(string)

        if !okId || !okEmail {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "userID", strconv.FormatFloat(userId, 'f', 0, 64))
        ctx = context.WithValue(ctx, "email", email)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}


func RequireAllRoles(urr repo.UserRoleRepository, roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            val := r.Context().Value("userID")
            if val == nil {
                http.Error(w, "Unauthorized: Missing user context", http.StatusUnauthorized)
                return
            }

            userIdStr, ok := val.(string)
            if !ok {
                http.Error(w, "Unauthorized: Invalid user context", http.StatusUnauthorized)
                return
            }

            userId, err := strconv.ParseInt(userIdStr, 10, 64)
            if err != nil {
                http.Error(w, "Invalid user ID", http.StatusUnauthorized)
                return
            }

            hasAllRoles, hasAllRolesErr := urr.HasAllRoles(userId, roles)

            fmt.Println("userid", userId, "roles", roles, "hasAllRoles", hasAllRoles)

            if hasAllRolesErr != nil {
                http.Error(w, "Error checking user roles: "+hasAllRolesErr.Error(), http.StatusInternalServerError)
                return
            }



            if !hasAllRoles {
                http.Error(w, "Forbidden: You do not have the required roles", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userIdStr := r.Context().Value("userID").(string)
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusUnauthorized)
				return
			}

			dbConn, dbErr := dbConfig.SetupDB()
			if dbErr != nil {
				http.Error(w, "Database connection error: "+dbErr.Error(), http.StatusInternalServerError)
				return
			}

			urr := repo.NewUserRoleRepository(dbConn)

			hasAnyRole, hasAnyRolesErr := urr.HasAnyRole(userId, roles)
			fmt.Println("userid", userId, "roles", roles, "hasAnyRole", hasAnyRole)
			if hasAnyRolesErr != nil {
				http.Error(w, "Error checking user roles: "+hasAnyRolesErr.Error(), http.StatusInternalServerError)
				return
			}

			if !hasAnyRole {
				http.Error(w, "Forbidden: You do not have the required roles", http.StatusForbidden)
				return
			}

			fmt.Println("User has all required roles:", roles)

			next.ServeHTTP(w, r)
		})
	}
}

