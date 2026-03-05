package main

import (
	"cert-chain/api"
	"cert-chain/blockchain"
	"cert-chain/database"
	"cert-chain/utils"
	"net/http"
	"fmt"
	
)


func main() {

	// Gerar chaves para a instituição
	publicKey, privateKey, err := utils.GenerateInstitutionKeys()

	if err != nil {
		fmt.Printf("Erro ao gerar chaves: %v\n", err)
	 }else{
		fmt.Printf("Chaves geradas com sucesso!\nPublic Key: %x\nPrivate Key: %x\n", publicKey, privateKey)
	 }

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

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Erro ao iniciar o servidor: %v\n", err)
	}
}