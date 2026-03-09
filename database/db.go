package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Apenas UMA variável global, com letra maiúscula para ser exportada
var DB *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=tedcrypto1239 dbname=cert_chain sslmode=disable"

	var err error
	// Atribuímos a conexão diretamente à variável global DB (maiúscula)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao preparar o banco: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar no banco (Senha errada ou banco offline): ", err)
	}

	fmt.Println("Banco de Dados conectado com sucesso!")

	createTable()
}

func createTable() {
	// 1. CRIAR INSTITUTIONS PRIMEIRO (nada depende dela)
	queryInstitutions := `
	CREATE TABLE IF NOT EXISTS institutions (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		username VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		public_key TEXT NOT NULL,
		private_key TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := DB.Exec(queryInstitutions)
	if err != nil {
		log.Fatal("Erro ao criar a tabela institutions: ", err)
	}
	fmt.Println("Tabela 'institutions' pronta para participar da rede multi-tenant!")

	// 2. CRIAR CERTIFICATES (depende de institutions)
	queryCertificates := `
	CREATE TABLE IF NOT EXISTS certificates (
		id VARCHAR(64) PRIMARY KEY,
		institution_id INTEGER NOT NULL,
		student_name VARCHAR(255) NOT NULL,
		course VARCHAR(255) NOT NULL,
		file_hash TEXT UNIQUE NOT NULL,
		blockchain_hash TEXT,
		timestamp BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (institution_id) REFERENCES institutions(id)
	);`

	_, err = DB.Exec(queryCertificates)
	if err != nil {
		log.Fatal("Erro ao criar a tabela certificates: ", err)
	}
	fmt.Println("Tabela 'certificates' pronta para receber dados!")

	// 3. CRIAR BLOCKS (nada depende dela)
	queryBlocks := `
	CREATE TABLE IF NOT EXISTS blocks (
		index INTEGER PRIMARY KEY,
		timestamp BIGINT NOT NULL,
		prev_hash TEXT NOT NULL,
		hash TEXT NOT NULL,
		nonce BIGINT NOT NULL,
		transactions JSONB NOT NULL
	);`

	_, err = DB.Exec(queryBlocks)
	if err != nil {
		log.Fatal("Erro ao criar a tabela blocks: ", err)
	}
	fmt.Println("Tabela 'blocks' pronta para registrar a blockchain!")
}
