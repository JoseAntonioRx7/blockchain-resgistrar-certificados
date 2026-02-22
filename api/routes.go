package api

import "net/http"

func RegisterRoutes() {
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/verify", VerifyHandler)
}
