package main

import (
	"cert-chain/api"
	"cert-chain/blockchain"
	"cert-chain/database"
	"fmt"
	"net/http"
)

func main() {
	// 1. Inicializa a infraestrutura
	database.InitDB() // Conecta ao banco que vimos no pgAdmin
	
	// 2. Carrega a Blockchain global
	api.Chain = blockchain.LoadBlockchain()
	fmt.Println("TTLedger: Corrente de blocos sincronizada.")

	// 3. Define as rotas (incluindo as novas que vamos criar)
	api.RegisterRoutes()

	// 4. Servidores de arquivos
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.Handle("/pdfs/", http.StripPrefix("/pdfs/", http.FileServer(http.Dir("./pdfs"))))

	fmt.Println("Servidor Multi-Tenant rodando em http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro fatal: %v\n", err)
	}
}