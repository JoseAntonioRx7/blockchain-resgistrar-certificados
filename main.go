package main

import (
	"cert-chain/api"
	"cert-chain/blockchain"
	"fmt"
	"net/http"
	"cert-chain/database"
)

func main() {
	// 1. INICIA O BANCO DE DADOS PRIMEIRO
	database.InitDB()

	// 2. Carrega a Blockchain (O código que já fizemos)
	api.Chain = blockchain.LoadBlockchain()
	fmt.Println("Blockchain funcionando!")
	
	api.RegisterRoutes()
	
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}