package utils

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// Chave secreta super protegida. (Em produção, isso ficaria num arquivo .env)
// É ela que garante que ninguém consiga falsificar o seu token.
var jwtSecretKey = []byte("super_senha_secreta_blockchain_2026")

// GenerateJWT cria o token para a instituição que fez login com sucesso
func GenerateJWT(username string) (string, error) {
	// 1. Cria o payload (as informações que vão dentro do token)
	// Aqui dizemos quem é o usuário e que o token expira em 2 horas
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	}

	// 2. Monta o token usando o algoritmo de assinatura HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Assina o token com a sua chave secreta e retorna a string final
	return token.SignedString(jwtSecretKey)
}

// ValidateJWT pega o token recebido e verifica se ele foi assinado pela sua chave
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// A função de parse lê o token e usa uma função anônima para fornecer a chave secreta correta
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o algoritmo de assinatura é o esperado (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metodo de assinatura inesperado: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})
}