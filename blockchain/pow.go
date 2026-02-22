package blockchain

import "bytes"

const Difficulty = 3

func (b *Block) Mine() {
	var hash []byte
	for {
		hash = b.calculateHash()
		if bytes.HasPrefix(hash, []byte{0, 0, 0}) {
			break
		}
		b.Nonce++
	}
	b.Hash = hash
}
