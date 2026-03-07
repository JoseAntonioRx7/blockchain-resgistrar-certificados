package utils

import (
	"cert-chain/blockchain"
	"fmt"
	"os"
	"time"
	"github.com/jung-kurt/gofpdf"
)

// GenerateCertificatePDF cria um arquivo PDF formatado para o aluno
func GenerateCertificatePDF(tx blockchain.CertificateTransaction) (string, error) {
	pdf := gofpdf.New("L", "mm", "A4", "") // "L" para Paisagem (Landscape)
	pdf.AddPage()

	// 1. Configuração de Borda e Fundo
	pdf.SetLineWidth(2)
	pdf.SetDrawColor(94, 0, 63) // Cor primária que definimos no CSS
	pdf.Rect(10, 10, 277, 190, "D")

	// 2. Título do Certificado
	pdf.SetFont("Arial", "B", 35)
	pdf.SetTextColor(43, 0, 29)
	pdf.Ln(30)
	pdf.CellFormat(0, 15, "CERTIFICADO DE CONCLUSAO", "0", 1, "C", false, 0, "")

	// 3. Texto do Certificado
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 18)
	pdf.SetTextColor(0, 0, 0)
	
	texto := fmt.Sprintf("Certificamos que o(a) aluno(a) %s concluiu com êxito o curso de %s, ministrado pela instituição %s.", 
		tx.StudentName, tx.Course, tx.Institution)
	
	pdf.MultiCell(0, 10, texto, "0", "C", false)

	// 4. Detalhes da Blockchain (Onde a mágica acontece)
	pdf.Ln(30)
	pdf.SetFont("Courier", "I", 10)
	pdf.SetTextColor(100, 100, 100)
	
	dataStr := time.Unix(tx.Timestamp, 0).Format("02/01/2006 15:04:05")
	infoBlockchain := fmt.Sprintf("Autenticidade garantida via Blockchain\nID da Transacao: %s\nHash do Arquivo: %s\nData de Registro: %s", 
		tx.ID, tx.FileHash, dataStr)
	
	pdf.MultiCell(0, 5, infoBlockchain, "0", "C", false)

	// 5. Salvar o arquivo
	// Criamos uma pasta para os PDFs se ela não existir
	os.MkdirAll("generated_pdfs", os.ModePerm)
	
	fileName := fmt.Sprintf("generated_pdfs/cert_%s.pdf", tx.ID)
	err := pdf.OutputFileAndClose(fileName)
	
	if err != nil {
		return "", err
	}

	return fileName, nil
}