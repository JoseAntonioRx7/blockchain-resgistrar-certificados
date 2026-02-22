package blockchain

import "time"

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	genesis := &Block{
		Timestamp:    time.Now().Unix(),
		Transactions: []CertificateTransaction{},
		PrevHash:     []byte{},
	}

	genesis.Mine()

	return &Blockchain{Blocks: []*Block{genesis}}
}

func (bc *Blockchain) AddBlock(txs []CertificateTransaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := &Block{
		Timestamp:    time.Now().Unix(),
		Transactions: txs,
		PrevHash:     prevBlock.Hash,
	}

	newBlock.Mine()
	bc.Blocks = append(bc.Blocks, newBlock)
}
