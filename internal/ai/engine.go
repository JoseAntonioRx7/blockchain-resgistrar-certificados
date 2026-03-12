package ai

import (
	"context"
	"log"
	"os"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Client estruturado para ser reutilizado
type AICore struct {
	Client *genai.Client
	Model  *genai.GenerativeModel
}

func NewAICore(ctx context.Context) *AICore {
	apiKey := os.Getenv("GEMINI_API_KEY") // Pegaremos do seu .env ou sistema
	if apiKey == "" {
		log.Fatal("ERRO: Variável GEMINI_API_KEY não configurada")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Erro ao conectar no Gemini: %v", err)
	}

	// Usaremos o modelo Flash por ser gratuito e rápido
	model := client.GenerativeModel("gemini-1.5-flash")
	
	return &AICore{
		Client: client,
		Model:  model,
	}
}