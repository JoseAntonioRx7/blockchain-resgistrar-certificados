package api

import (
	"cert-chain/blockchain"
	"cert-chain/database"
	"cert-chain/utils"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
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

	// Gera o Hash uma única vez (removemos a duplicata)
	hash, err := utils.HashFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashBytes := []byte(hash)
	signature := utils.SignData(utils.PrivateKey, hashBytes)

	// Cria o objeto de transação
	tx := blockchain.CertificateTransaction{
		ID:          	utils.GenerateID(),
		StudentName: 	r.FormValue("student_name"),
		Institution: 	r.FormValue("institution"),
		Course:      	r.FormValue("course"),
		FileHash:    	hash,
		Signature:		fmt.Sprintf("%x", signature), // Simulando assinatura (na prática, seria a assinatura digital da instituição)
		Timestamp:   	time.Now().Unix(),
	}

	// Salva no banco de dados
	query := `INSERT INTO certificates (id, student_name, institution, course, file_hash, timestamp) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = database.DB.Exec(query, tx.ID, tx.StudentName, tx.Institution, tx.Course, tx.FileHash, tx.Timestamp)
	if err != nil {
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	// Adiciona na Blockchain
	Chain.AddBlock([]blockchain.CertificateTransaction{tx})

	// Retorna sucesso
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

	for _, block := range Chain.Blocks {
		for _, tx := range block.Transactions {
			allCerts = append(allCerts, tx)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allCerts)
}