package blockchain

import (
	"encoding/json"
	"os"
	"time"
)

// O nome do arquivo onde os dados serão salvos
const dbFile = "blockchain.json"

// Estrutura que representa a Blockchain (Já existia no seu código)
type Blockchain struct {
	Blocks []*Block
}

// Função para criar uma nova Blockchain (Já existia no seu código)
func NewBlockchain() *Blockchain {
	genesis := &Block{
		Timestamp:    time.Now().Unix(),
		Transactions: []CertificateTransaction{},
		PrevHash:     []byte{},
	}
	genesis.Mine()
	return &Blockchain{Blocks: []*Block{genesis}}
}

// --- NOVAS FUNÇÕES DE PERSISTÊNCIA ---

// Salva a blockchain inteira no arquivo JSON
func (bc *Blockchain) Save() {
	data, _ := json.MarshalIndent(bc, "", "  ")
	os.WriteFile(dbFile, data, 0644)
}

// Tenta carregar do arquivo, se não existir, cria uma nova
func LoadBlockchain() *Blockchain {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		bc := NewBlockchain()
		bc.Save() // Salva o bloco gênese inicial
		return bc
	}

	file, _ := os.ReadFile(dbFile)
	var bc Blockchain
	err := json.Unmarshal(file, &bc)
	if err != nil {
		return NewBlockchain()
	}
	return &bc
}

// Adiciona um bloco e salva no arquivo automaticamente
func (bc *Blockchain) AddBlock(txs []CertificateTransaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{
		Timestamp:    time.Now().Unix(),
		Transactions: txs,
		PrevHash:     prevBlock.Hash,
	}
	newBlock.Mine()
	bc.Blocks = append(bc.Blocks, newBlock)
	
	// Esta linha garante a persistência na startup
	bc.Save() 
}