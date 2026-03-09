package api

import (
	"net/http"
)

func RegisterRoutes() {
	// 0. Debug (Para verificar usuários no banco)
	// URL: http://localhost:8080/api/debug/users
	http.HandleFunc("/api/debug/users", DebugUsersHandler)

	// 1. Administrativo (Para criar a primeira instituição)
	// URL: http://localhost:8080/api/admin/register-institution
	http.HandleFunc("/api/admin/register-institution", RegisterInstitutionHandler)

	// 2. Autenticação e Verificação Pública
	// URL: http://localhost:8080/login
	http.HandleFunc("/login", LoginHandler)
	// URL: http://localhost:8080/verify
	http.HandleFunc("/verify", VerifyHandler)

	// 3. Operações Protegidas (Exigem Token JWT)
	// URL: http://localhost:8080/list
	http.HandleFunc("/list", JWTMiddleware(ListCertificatesHandler))
	// URL: http://localhost:8080/register
	http.HandleFunc("/register", JWTMiddleware(RegisterHandler))
}
