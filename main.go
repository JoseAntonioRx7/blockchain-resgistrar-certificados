package main

import (
	"cert-chain/api"
	"cert-chain/blockchain"
	"fmt"
	"net/http"
)

func main() {
	// Em vez de NewBlockchain, usamos LoadBlockchain
	// Isso garante que os dados antigos apareçam na inicialização
	api.Chain = blockchain.LoadBlockchain()

	fmt.Println("Blockchain carregada e persistente!")
	
	api.RegisterRoutes()
	
	// Serve o frontend
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}