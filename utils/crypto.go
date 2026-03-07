package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateInstitutionKeys agora retorna chaves em formato HEX para salvar no Postgres
func GenerateInstitutionKeys() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("falha ao gerar chaves criptograficas: %v", err)
	}

	// Convertemos para Hexadecimal para facilitar o armazenamento no banco
	return hex.EncodeToString(pub), hex.EncodeToString(priv), nil
}

// SignData assina o hash do certificado usando a chave privada (em HEX) da instituição
func SignData(hash string, privateKeyHex string) (string, error) {
	// 1. Converte a chave privada de Hex para Bytes
	privBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}

	// 2. Assina os dados (o hash do arquivo)
	signature := ed25519.Sign(privBytes, []byte(hash))

	// 3. Retorna a assinatura em Hex para ser gravada na Blockchain e no DB
	return hex.EncodeToString(signature), nil
}