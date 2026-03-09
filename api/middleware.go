package api

import(
	"cert-chain/utils"
	"net/http"
	"strings"
	"context"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware "envolve" outra função HTTP
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token nao fornecido", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Formato de token invalido", http.StatusUnauthorized)
			return
		}

		token, err := utils.ValidateJWT(bearerToken[1])
		if err != nil || !token.Valid {
			http.Error(w, "Token invalido ou expirado", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			username := claims["username"].(string)

			ctx := context.WithValue(r.Context(), "username", username)

			next(w, r.WithContext(ctx))
			return
		}

		http.Error(w, "Token invalido", http.StatusUnauthorized)
	}
}
