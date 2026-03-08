package api

import "net/http"

func RegisterRoutes() {
    http.HandleFunc("/login", LoginHandler)
    http.HandleFunc("/verify", VerifyHandler)
    http.HandleFunc("/list", ListCertificatesHandler)
    http.HandleFunc("/register", JWTMiddleware(RegisterHandler))
    
    // Removido o http.Handle("/pdfs/") daqui para centralizar no main.go
}