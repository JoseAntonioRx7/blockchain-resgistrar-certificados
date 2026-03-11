package api

import (
	"net/http"
)

func RegisterRoutes() {
	// --- ROTAS PÚBLICAS ---
	// Estas rotas chamam setupCORS() internamente no Handler.
	http.HandleFunc("/api/login", LoginHandler)
	http.HandleFunc("/api/verify", VerifyHandler)
	http.HandleFunc("/api/admin/register-institution", RegisterInstitutionHandler)

	// --- ROTAS PROTEGIDAS (JWT) ---
	// Estas rotas NÃO chamam setupCORS() internamente, pois o JWTMiddleware já faz isso.
	http.HandleFunc("/api/list", JWTMiddleware(ListCertificatesHandler))
	http.HandleFunc("/api/register", JWTMiddleware(RegisterHandler))

	// --- ARQUIVOS ESTÁTICOS ---
	fs := http.FileServer(http.Dir("./pdfs"))
	http.Handle("/pdfs/", http.StripPrefix("/pdfs/", fs))
}