package blockchain

import (
	"cert-chain/database"
	"encoding/json"
	"encoding/hex"
	"log"
	"time"
)

// Estrutura que representa a Blockchain
type Blockchain struct {
	Blocks []*Block
}

// CreateGenesisBlock isola a lógica de criação do primeiro bloco
func CreateGenesisBlock() *Block {
	genesis := &Block{
		Index:        0, // Garante que o ID no banco será 0
		Timestamp:    time.Now().Unix(),
		Transactions: []CertificateTransaction{},
		PrevHash:     []byte{},
	}
	genesis.Mine()
	return genesis
}

// --- INTEGRAÇÃO COM POSTGRESQL ---

// SaveBlockToDB salva um único bloco recém-minerado no PostgreSQL
func SaveBlockToDB(b *Block) {
	txsJSON, err := json.Marshal(b.Transactions)
	if err != nil {
		log.Println("Erro ao converter transacoes para JSON:", err)
		return
	}

	query := `INSERT INTO blocks (index, timestamp, prev_hash, hash, nonce, transactions) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = database.DB.Exec(query,
		b.Index,
		b.Timestamp,
		hex.EncodeToString(b.PrevHash),
		hex.EncodeToString(b.Hash),
		b.Nonce,
		string(txsJSON),
	)

	if err != nil {
		log.Println("Erro ao salvar bloco no banco de dados:", err)
	} else {
		log.Printf("Bloco %d salvo no banco de dados com sucesso!\n", b.Index)
	}
}

// LoadBlockchain carrega todo o histórico de blocos do PostgreSQL
func LoadBlockchain() *Blockchain {
    chain := &Blockchain{}

    query := `SELECT index, timestamp, prev_hash, hash, nonce, transactions FROM blocks ORDER BY index ASC`
    rows, err := database.DB.Query(query)

    if err != nil {
        log.Fatal("Erro ao buscar blocos no banco: ", err)
    }
    defer rows.Close()

    for rows.Next() {
        // Criamos um ponteiro para Block, exatamente como a sua struct exige
        b := &Block{}
        var txsJSON string

        err := rows.Scan(&b.Index, &b.Timestamp, &b.PrevHash, &b.Hash, &b.Nonce, &txsJSON)
        if err != nil {
            log.Fatal("Erro ao ler dados do bloco: ", err)
        }

        // Remonta as transações
        json.Unmarshal([]byte(txsJSON), &b.Transactions)

        chain.Blocks = append(chain.Blocks, b)
    }

    // Inteligência de inicialização: Se o banco estiver vazio, cria o Gênese
    if len(chain.Blocks) == 0 {
        log.Println("Nenhum bloco encontrado. Criando Bloco Genese no Postgres")
        genesisBlock := CreateGenesisBlock()
        chain.Blocks = append(chain.Blocks, genesisBlock)
        SaveBlockToDB(genesisBlock)
    } else {
        log.Printf("Blockchain funcionando com sucesso! Total de blocos: %d\n", len(chain.Blocks))
    }

    return chain
}

// AddBlock adiciona um bloco e salva no PostgreSQL automaticamente
func (bc *Blockchain) AddBlock(txs []CertificateTransaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	
	newBlock := &Block{
		Index:        prevBlock.Index + 1, // Fundamental para a tabela 'blocks' do SQL
		Timestamp:    time.Now().Unix(),
		Transactions: txs,
		PrevHash:     prevBlock.Hash,
	}
	
	newBlock.Mine()
	bc.Blocks = append(bc.Blocks, newBlock)

	// Esta linha garante a persistência no banco
	SaveBlockToDB(newBlock)
}