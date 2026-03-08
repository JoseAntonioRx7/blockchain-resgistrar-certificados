package main

import (
	"cert-chain/api"
	"cert-chain/blockchain"
	"cert-chain/database"
	"cert-chain/utils"
	"fmt"
	"net/http"
)

func main() {
	// 1. Configurações Iniciais e Criptografia
	publicKey, privateKey, err := utils.GenerateInstitutionKeys()
	if err != nil {
		fmt.Printf("Erro ao gerar chaves: %v\n", err)
	} else {
		fmt.Printf("Chaves geradas com sucesso!\nPublic Key: %x\nPrivate Key: %x\n", publicKey, privateKey)
	}

	// 2. Inicialização de Infraestrutura (DB e Blockchain)
	database.InitDB()
	api.Chain = blockchain.LoadBlockchain()
	fmt.Println("Blockchain carregada com sucesso!")

	// 3. Definição de Rotas (Devem ser declaradas ANTES do ListenAndServe)
	api.RegisterRoutes()

	// Servidor de arquivos estáticos (Frontend na pasta /web)
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	// Servidor de PDFs (Consolidado para a pasta correta de saída)
	// Certifique-se de que o RegisterHandler salva nesta mesma pasta
	pdfFs := http.FileServer(http.Dir("./pdfs"))
	http.Handle("/pdfs/", http.StripPrefix("/pdfs/", pdfFs))

	// 4. Inicialização Única do Servidor
	port := ":8080"
	fmt.Printf("TTLedger Online: http://localhost%s\n", port)
	
	serverErr := http.ListenAndServe(port, nil)
	if serverErr != nil {
		fmt.Printf("Falha crítica no servidor: %v\n", serverErr)
	}
}