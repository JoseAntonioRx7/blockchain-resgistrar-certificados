package api

import (
	"cert-chain/blockchain"
	"cert-chain/database"
	"cert-chain/utils"
	"encoding/json"	
	"net/http"
	"time"
)

var Chain *blockchain.Blockchain

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Arquivo inválido", http.StatusBadRequest)
		return
	}
	defer file.Close()

	hash, err := utils.HashFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hash, err = utils.HashFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx := blockchain.CertificateTransaction{
		ID:          utils.GenerateID(),
		StudentName: r.FormValue("student_name"),
		Institution: r.FormValue("institution"),
		Course:      r.FormValue("course"),
		FileHash:    hash,
		Timestamp:   time.Now().Unix(),
	}

	// 1. Defina o comando SQL com placeholders
	query := `INSERT INTO certificates (id, student_name, institution, course, file_hash, timestamp) 
		  VALUES ($1, $2, $3, $4, $5, $6)`

	// 2. Execute a query passando os dados do seu objeto 'tx' (CertificateTransaction)
	// O banco de dados valida os tipos automaticamente
	_, err = database.DB.Exec(query, tx.ID, tx.StudentName, tx.Institution, tx.Course, tx.FileHash, tx.Timestamp)

	// 3. Tratamento de erro profissional: Se falhar aqui, o registro não prossegue
	if err != nil {
		// Retorna erro 500 para o frontend entender que algo falhou
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	Chain.AddBlock([]blockchain.CertificateTransaction{tx})

	resp := map[string]interface{}{
		"message": "Certificado registrado com sucesso",
		"hash":    hash,
		"id":      tx.ID,
	}

	json.NewEncoder(w).Encode(resp)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")

	for _, block := range Chain.Blocks {
		for _, tx := range block.Transactions {
			if tx.FileHash == hash {
				json.NewEncoder(w).Encode(tx)
				return
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]string{
		"error": "Certificado não encontrado",
	})
}

func ListCertificatesHandler(w http.ResponseWriter, r *http.Request) {
	var allCerts []blockchain.CertificateTransaction

	// Percorre todos os blocos (pulando o gênese se quiser)
	for _, block := range Chain.Blocks {
		for _, tx := range block.Transactions {
			allCerts = append(allCerts, tx)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allCerts)
}
