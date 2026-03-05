package api

import (
	"cert-chain/blockchain"
	"cert-chain/database"
	"cert-chain/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var Chain *blockchain.Blockchain

// Função auxiliar para configurar CORS e evitar repetição
func setupCORS(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		(*w).WriteHeader(http.StatusOK)
		return true
	}
	return false
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) { return }

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Arquivo invalido", http.StatusBadRequest)
		return
	}
	defer file.Close()

	hash, err := utils.HashFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashBytes := []byte(hash)
	signature := utils.SignData(utils.PrivateKey, hashBytes)

	tx := blockchain.CertificateTransaction{
		ID:           utils.GenerateID(),
		StudentName:  r.FormValue("student_name"),
		Institution:  r.FormValue("institution"),
		Course:       r.FormValue("course"),
		FileHash:     hash,
		Signature:    fmt.Sprintf("%x", signature),
		Timestamp:    time.Now().Unix(),
	}

	query := `INSERT INTO certificates (id, student_name, institution, course, file_hash, signature, timestamp) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = database.DB.Exec(query, 
		tx.ID, 
		tx.StudentName, 
		tx.Institution, 
		tx.Course, 
		tx.FileHash, 
		tx.Signature, 
		tx.Timestamp,
	)

	if err != nil {
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	// Adiciona na Blockchain e salva a persistência
	Chain.AddBlock([]blockchain.CertificateTransaction{tx})

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"message": "Certificado registrado com sucesso",
		"hash":    hash,
		"id":      tx.ID,
	}
	json.NewEncoder(w).Encode(resp)
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) { return }

	hash := r.URL.Query().Get("hash")
	w.Header().Set("Content-Type", "application/json")

	for _, block := range Chain.Blocks {
		for _, tx := range block.Transactions {
			if tx.FileHash == hash {
				json.NewEncoder(w).Encode(tx)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Certificado nao encontrado",
	})
}

func ListCertificatesHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) { return }

	var allCerts []blockchain.CertificateTransaction
	for _, block := range Chain.Blocks {
		for _, tx := range block.Transactions {
			allCerts = append(allCerts, tx)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allCerts)
}