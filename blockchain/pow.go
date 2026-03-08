package blockchain

import (
	"bytes"
	"fmt"
	"time"
)

// Difficulty = 3 é o equilíbrio perfeito para a TTLedger:
// Garante o conceito de PoW sem destruir a experiência do usuário (UX).
const Difficulty = 3 

func (b *Block) Mine() {
	start := time.Now()
	var hash []byte
	
	// Calculamos o alvo uma única vez fora do loop para performance
	target := bytes.Repeat([]byte{0}, Difficulty)

	fmt.Printf("TTLedger: Iniciando mineração do bloco %d...\n", b.Index)

	for {
		hash = b.calculateHash()
		
		// Verificação de prefixo (Proof of Work)
		if bytes.HasPrefix(hash, target) {
			break
		}
		
		b.Nonce++

		// Batimento cardíaco: Log a cada 1 milhão de tentativas para 
		// mostrar que o sistema não travou no backend.
		if b.Nonce%1000000 == 0 {
			fmt.Printf("Minerando... Nonce atual: %d\n", b.Nonce)
		}
	}

	b.Hash = hash
	duration := time.Since(start)
	
	fmt.Printf("Bloco Minerado! Nonce: %d | Tempo: %v | Hash: %x\n", b.Nonce, duration, b.Hash)
}