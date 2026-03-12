package ai

import (
	"context"
	"fmt"
	"encoding/json"
	"github.com/google/generative-ai-go/genai"
)

// AnalyzeCertificates envia os dados para o Gemini e retorna a análise
func (core *AICore) AnalyzeCertificates(ctx context.Context, certData interface{}) (string, error) {
	// 1. Transformamos os dados do banco em uma String JSON bonita para a IA
	jsonData, _ := json.MarshalIndent(certData, "", "  ")

	// 2. Criamos o "Prompt" (A instrução mestre)
	prompt := fmt.Sprintf(`
		Você é um Auditor de Segurança de Rede Blockchain especializado em registros acadêmicos.
		Analise os seguintes registros de certificados emitidos recentemente em nossa rede:

		%s

		Sua tarefa:
		1. Identifique emissões em massa suspeitas (mesmo curso/aluno em segundos).
		2. Verifique se nomes de alunos ou cursos parecem gerados por bots ou são incoerentes.
		3. Dê um nível de risco (Baixo, Médio ou Alto).
		4. Retorne um resumo conciso e direto para o administrador.
	`, string(jsonData))

	// 3. Enviamos para o Gemini
	resp, err := core.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	// 4. Extraímos o texto da resposta
	if len(resp.Candidates) > 0 {
		c := resp.Candidates[0]
		if c.Content != nil {
			return fmt.Sprintf("%v", c.Content.Parts[0]), nil
		}
	}

	return "IA não conseguiu gerar uma análise no momento.", nil
}