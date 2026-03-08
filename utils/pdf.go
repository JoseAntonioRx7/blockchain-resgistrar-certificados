package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// SaveUploadedCertificate pega o arquivo enviado no formulário e salva com o ID da transação
func SaveUploadedCertificate(file multipart.File, id string) (string, error) {
	outputDir := "pdfs"

	// 1. Garante que a pasta pdfs existe
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", err
	}

	// 2. Define o nome do arquivo com o ID da transação (para bater com o Frontend)
	fileName := fmt.Sprintf("%s/cert_%s.pdf", outputDir, id)

	// 3. Cria o arquivo vazio no disco
	out, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 4. Copia o conteúdo do arquivo que o usuário enviou para o arquivo no disco
	_, err = io.Copy(out, file)
	if err != nil {
		return "", err
	}

	return fileName, nil
}