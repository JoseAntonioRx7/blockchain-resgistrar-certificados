package main

import (
	"cert-chain/api"
	"cert-chain/blockchain"
	"cert-chain/database"
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"cert-chain/internal/ai"
	"context"
)


func main() {

	godotenv.Load() // Carrega as variáveis do .env para o ambiente

	// 1. Inicializa a infraestrutura (Banco de Dados)
	database.InitDB()

	// 2. Garante que a pasta de PDFs exista
	if _, err := os.Stat("./pdfs"); os.IsNotExist(err) {
		os.Mkdir("./pdfs", 0755)
		fmt.Println("Diretório ./pdfs criado com sucesso.")
	}

	// 3. Carrega a Blockchain global e injeta no pacote API
	api.Chain = blockchain.LoadBlockchain()
	fmt.Println("TTLedger: Corrente de blocos sincronizada.")

	// --- ORDEM DAS ROTAS É CRUCIAL ---

	// 4. Rota de Arquivos Estáticos (PDFs)
	// IMPORTANTE: Esta rota deve vir ANTES das rotas genéricas
	fs := http.FileServer(http.Dir("./pdfs"))
	http.Handle("/api/pdfs/", http.StripPrefix("/api/pdfs/", fs))

	// 5. Registra as Rotas da API (Login, Register, List, Verify)
	api.RegisterRoutes()

	// 6. Servir o Frontend (Arquivos da pasta 'web')
	// O FileServer para "/" deve ser o ÚLTIMO, pois ele aceita qualquer caminho
	http.Handle("/", http.FileServer(http.Dir("./web")))

	fmt.Println("================================================")
	fmt.Println("TTLedger Node ativo em http://localhost:8080")
	fmt.Println("Segurança: Non-Custodial ativa")
	fmt.Println("================================================")

	// Inicializa o Core da IA (Precisa do contexto)
	ctx := context.Background()
	aiCore := ai.NewAICore(ctx)

	// Nova rota de Auditoria
	http.HandleFunc("/api/admin/audit-network", api.AuditHandler(database.DB, aiCore))

	// Dentro da func main(), antes do ListenAndServe:
	ctx = context.Background()
	aiCore = ai.NewAICore(ctx) // 'ai' é o nome que está no seu import

	// Registra a rota de auditoria usando o Handler que acabamos de criar
	http.HandleFunc("/api/admin/audit-network", api.AuditNetworkHandler(database.DB, aiCore))

	// Inicia o servidor usando o DefaultServeMux (nil)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Erro fatal no servidor: %v\n", err)
	}
}			