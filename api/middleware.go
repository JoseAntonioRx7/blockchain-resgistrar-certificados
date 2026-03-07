package api

import(
	"cert-chain/utils"
	"net/http"
	"strings"
)

// JWTMiddleware "envolve" outra função HTTP
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
        // ESSA PARTE É CRUCIAL: Liberar o "Pre-flight" do navegador
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

		

		// 2. Pega o cabeçalho "Authorization" que o frontend vai enviar
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token nao fornecido", http.StatusUnauthorized)
			return
		}

		// O formato padrão do mercado é enviar: "Bearer <token_aqui>"
		// Então precisamos separar a palavra "Bearer" do token real
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			http.Error(w, "Formato de token invalido", http.StatusUnauthorized)
			return
		}

		// 3. Valida o token usando a nossa função do utils
		token, err := utils.ValidateJWT(bearerToken[1])
		if err != nil || !token.Valid {
			http.Error(w, "Token invalido ou expirado", http.StatusUnauthorized)
			return
		}

		// 4. Se chegou até aqui, o crachá é válido!
		// Chama a próxima função (que será o nosso RegisterHandler)
		next(w, r)
	}
}