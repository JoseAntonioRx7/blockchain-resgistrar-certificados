package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

var DB *sql.DB
func Connect(dsn string) {
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}
}

var db *sql.DB

func InitDB() {
	connStr := "host=localhost port=5432 user=postgres password=tedcrypto1239 dbname=cert_chain sslmode=disable"

	// 2. Prepara a conexão
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao preparar o banco: ", err)
	}

	// 3. O "Ping" é o que realmente bate na porta do banco para testar a senha
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar no banco (Senha errada ou banco offline): ", err)
	}

	fmt.Println("Banco de Dados conectado com sucesso!")

	// 4. Chama o algoritmo de criação da tabela
	createTable()
}

// Algoritmo 2: Estruturar os Dados
func createTable() {
	// Este é o comando SQL. "IF NOT EXISTS" impede que dê erro se a tabela já existir.
	query := `
	CREATE TABLE IF NOT EXISTS certificates (
		id VARCHAR(64) PRIMARY KEY,
		student_name VARCHAR(255) NOT NULL,
		institution VARCHAR(255) NOT NULL,
		course VARCHAR(255) NOT NULL,
		file_hash TEXT UNIQUE NOT NULL,
		timestamp BIGINT NOT NULL
	);`

	// O 'Exec' apenas executa o comando SQL sem esperar uma resposta de dados de volta
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Erro ao criar a tabela: ", err)
	}
	fmt.Println("Tabela 'certificates' pronta para receber dados!")
}

