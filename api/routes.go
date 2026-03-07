package api

import "net/http"

// RegisterRoutes configura os endpoints da API
func RegisterRoutes() {
	// 1. Rotas Públicas (Qualquer pessoa ou instituição pode acessar)
	http.HandleFunc("/login", LoginHandler)          // A recepção para pegar o crachá
	http.HandleFunc("/verify", VerifyHandler)        // Para validar diplomas
	http.HandleFunc("/list", ListCertificatesHandler) // Para ver o dashboard

	// 2. Rota Privada (Apenas instituições com o "Crachá" JWT podem acessar)
	// O JWTMiddleware intercepta a chamada antes de chegar no RegisterHandler
	http.HandleFunc("/register", JWTMiddleware(RegisterHandler))
}