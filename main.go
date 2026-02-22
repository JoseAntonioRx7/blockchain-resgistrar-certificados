package main

import (
	"cert-chain/api"
	"cert-chain/blockchain" // Importe o pacote da blockchain
	"fmt"
	"net/http"
)

func main() {
	// 1. INICIALIZE A VARIAVEL GLOBAL DO PACOTE API
	// Isso cria o bloco gênese e evita que 'Chain' seja nil
	api.Chain = blockchain.NewBlockchain()

	fmt.Println("Servidor rodando em http://localhost:8080")

	api.RegisterRoutes()

	// 2. ADICIONE UM HANDLER PARA SERVIR O FRONTEND (Opcional, mas ajuda)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Erro ao abrir servidor: %v\n", err)
	}
}
