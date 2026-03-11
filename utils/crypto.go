package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateInstitutionKeys gera o par de chaves Ed25519.
// A pública vai para o banco, a privada é mostrada uma única vez ao admin.
func GenerateInstitutionKeys() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("falha ao gerar chaves criptograficas: %v", err)
	}

	// Retornamos em Hex para facilitar o "copiar e colar" do usuário e o armazenamento no Postgres
	return hex.EncodeToString(pub), hex.EncodeToString(priv), nil
}

// SignData assina o hash do certificado. 
// A privateKeyHex deve vir do formulário preenchido pela instituição (volátil).
func SignData(hash string, privateKeyHex string) (string, error) {
	// 1. Decodifica a chave privada Hex para o formato original (bytes)
	privBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("chave privada invalida: %v", err)
	}

	// O Ed25519 do Go espera que a chave privada tenha 64 bytes (seed + pub)
	if len(privBytes) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("tamanho da chave privada incorreto")
	}

	// 2. Gera a assinatura digital sobre o hash do arquivo
	signature := ed25519.Sign(privBytes, []byte(hash))

	// 3. Retorna a assinatura em Hex para o registro na Blockchain
	return hex.EncodeToString(signature), nil
}