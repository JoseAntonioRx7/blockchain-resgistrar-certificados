package utils

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

var PrivateKey ed25519.PrivateKey
var PublicKey ed25519.PublicKey

// GenerateInstitutionkeys cria um par de chaves para a instituição ed25519
func GenerateInstitutionKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		return nil, nil, fmt.Errorf("falha ao gerar chaves criptograficas: %v", err)
	}

	PublicKey = publicKey
	PrivateKey = privateKey
	return publicKey, privateKey, nil

}

// SignData assina  hash do certificado usando a chave privada da instituição
func SignData(privateKey ed25519.PrivateKey, data []byte) []byte {
	return ed25519.Sign(privateKey, data)
}
