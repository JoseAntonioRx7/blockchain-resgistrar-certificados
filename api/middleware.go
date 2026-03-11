package api

import (
	"cert-chain/utils"
	"context"
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware centraliza a segurança e o CORS das rotas protegidas
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Resolve CORS antes de qualquer validação. 
		// Se for OPTIONS, o setupCORS retorna true e encerra a requisição aqui.
		if setupCORS(&w, r) {
			return
		}

		// 2. Extração do Token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token não fornecido"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Formato de token inválido"})
			return
		}

		// 3. Validação
		token, err := utils.ValidateJWT(parts[1])
		if err != nil || !token.Valid {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token inválido ou expirado"})
			return
		}

		// 4. Injeção no Contexto
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			username, ok := claims["username"].(string)
			if !ok {
				sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Payload inválido"})
				return
			}

			// Passa o username adiante através do contexto da requisição
			ctx := context.WithValue(r.Context(), "username", username)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Falha na autorização"})
	}
}