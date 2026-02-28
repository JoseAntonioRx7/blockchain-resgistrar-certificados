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
	query := `
	CREATE TABLE IF NOT EXISTS certificates (
		id VARCHAR(64) PRIMARY KEY,
		student_name VARCHAR(255) NOT NULL,
		institution VARCHAR(255) NOT NULL,
		course VARCHAR(255) NOT NULL,
		file_hash TEXT UNIQUE NOT NULL,
		timestamp BIGINT NOT NULL
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar a tabela: ", err)
	}
	fmt.Println("✅ Tabela 'certificates' pronta para receber dados!")
}