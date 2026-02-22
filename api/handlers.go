package api

import (
	"cert-chain/blockchain"
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

	tx := blockchain.CertificateTransaction{
		ID:          utils.GenerateID(),
		StudentName: r.FormValue("student_name"),
		Institution: r.FormValue("institution"),
		Course:      r.FormValue("course"),
		FileHash:    hash,
		Timestamp:   time.Now().Unix(),
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
