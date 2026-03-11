package utils

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
)

// SaveUploadedCertificate agora usa o hash para o nome do arquivo
func SaveUploadedCertificate(file multipart.File, hash string) (string, error) {
    outputDir := "pdfs"

    // 1. Garante que a pasta pdfs existe
    if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
        return "", err
    }

    // 2. Define o nome usando o HASH (para bater com o link do Frontend /api/pdfs/cert_HASH.pdf)
    fileName := fmt.Sprintf("%s/cert_%s.pdf", outputDir, hash)

    // 3. Cria o arquivo no disco
    out, err := os.Create(fileName)
    if err != nil {
        return "", err
    }
    defer out.Close()

    // 4. Copia o conteúdo
    _, err = io.Copy(out, file)
    if err != nil {
        return "", err
    }

    return fileName, nil
}