package utils

import (
	"cert-chain/blockchain"
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GenerateCertificatePDF centraliza a lógica de criação e salvamento do documento
func GenerateCertificatePDF(tx blockchain.CertificateTransaction) (string, error) {
	// 1. Inicialização do Documento (A4 Paisagem)
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// 2. Identidade Visual e Moldura
	setupLayout(pdf)

	// 3. Conteúdo Principal (Título e Texto)
	addCertificateContent(pdf, tx)

	// 4. Rodapé Técnico (Dados da Blockchain)
	addBlockchainFooter(pdf, tx)

	// 5. Gestão de Arquivos e Salvamento
	return savePDF(pdf, tx.ID)
}

// setupLayout define as cores e a borda do certificado
func setupLayout(pdf *gofpdf.Fpdf) {
	pdf.SetLineWidth(2)
	pdf.SetDrawColor(94, 0, 63) // Cor institucional
	pdf.Rect(10, 10, 277, 190, "D")
}

// addCertificateContent insere os dados do aluno e curso
func addCertificateContent(pdf *gofpdf.Fpdf, tx blockchain.CertificateTransaction) {
	// Título
	pdf.SetFont("Arial", "B", 35)
	pdf.SetTextColor(43, 0, 29)
	pdf.Ln(30)
	pdf.CellFormat(0, 15, "CERTIFICADO DE CONCLUSAO", "0", 1, "C", false, 0, "")

	// Corpo do texto
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 18)
	pdf.SetTextColor(0, 0, 0)
	
	texto := fmt.Sprintf("Certificamos que o(a) aluno(a) %s concluiu com exito o curso de %s, ministrado pela instituição %s.", 
		tx.StudentName, tx.Course, tx.Institution)
	
	pdf.MultiCell(0, 10, texto, "0", "C", false)
}

// addBlockchainFooter adiciona os metadados de segurança no rodapé
func addBlockchainFooter(pdf *gofpdf.Fpdf, tx blockchain.CertificateTransaction) {
	pdf.Ln(30)
	pdf.SetFont("Courier", "I", 10)
	pdf.SetTextColor(100, 100, 100)
	
	dataStr := time.Unix(tx.Timestamp, 0).Format("02/01/2006 15:04:05")
	info := fmt.Sprintf("Autenticidade garantida via Blockchain\nID: %s\nHash: %s\nData: %s", 
		tx.ID, tx.FileHash, dataStr)
	
	pdf.MultiCell(0, 5, info, "0", "C", false)
}

// savePDF garante que a pasta exista e salva o arquivo final
func savePDF(pdf *gofpdf.Fpdf, id string) (string, error) {
	// Sincronizado com a rota definida em api/routes.go
	outputDir := "pdfs" 
	
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s/cert_%s.pdf", outputDir, id)
	err := pdf.OutputFileAndClose(fileName)
	
	if err != nil {
		return "", err
	}

	return fileName, nil
}