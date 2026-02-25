package blockchain

import (
	"bytes"
)

// Subimos para 4 para aumentar a segurança da startup
const Difficulty = 4 

func (b *Block) Mine() {
	var hash []byte
    
    // Criamos o prefixo de zeros dinamicamente com base na Difficulty
    target := bytes.Repeat([]byte{0}, Difficulty)

	for {
		hash = b.calculateHash()
		if bytes.HasPrefix(hash, target) {
			break
		}
		b.Nonce++
	}
	b.Hash = hash
}