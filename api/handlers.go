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

// Placeholder para teste - No próximo passo, buscaremos isso do Postgres baseado no Login
const TEMP_PRIVATE_KEY = "0000000000000000000000000000000000000000000000000000000000000000"

func setupCORS(w *http.ResponseWriter, r *http.Request) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

	// 1. Gera o Hash do arquivo
	hash, err := utils.HashFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. ASSINATURA DIGITAL (Corrigido para o novo utils/crypto.go)
	// O SignData agora espera o hash em string e a chave privada em Hex
	signature, err := utils.SignData(hash, TEMP_PRIVATE_KEY)
	if err != nil {
		http.Error(w, "Erro ao assinar digitalmente", http.StatusInternalServerError)
		return
	}

	tx := blockchain.CertificateTransaction{
		ID:           utils.GenerateID(),
		StudentName:  r.FormValue("student_name"),
		Institution:  r.FormValue("institution"),
		Course:       r.FormValue("course"),
		FileHash:     hash,
		Signature:    signature, // Já vem em formato Hex do utils
		Timestamp:    time.Now().Unix(),
	}

	// 3. Persistência no PostgreSQL
	query := `INSERT INTO certificates (id, student_name, institution, course, file_hash, signature, timestamp) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = database.DB.Exec(query, 
		tx.ID, tx.StudentName, tx.Institution, tx.Course, tx.FileHash, tx.Signature, tx.Timestamp,
	)
	if err != nil {
		http.Error(w, "Erro ao salvar no banco de dados", http.StatusInternalServerError)
		return
	}

	// 4. Mineração na Blockchain
	Chain.AddBlock([]blockchain.CertificateTransaction{tx})

	// 5. Geração Automática do PDF
	pdfPath, err := utils.GenerateCertificatePDF(tx)
	if err != nil {
		fmt.Println("Alerta: Erro ao gerar PDF:", err)
	}

	// 6. Resposta Única de Sucesso
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Certificado registrado e PDF gerado!",
		"hash":     hash,
		"id":       tx.ID,
		"pdf_path": pdfPath,
	})
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
	json.NewEncoder(w).Encode(map[string]string{"error": "Certificado nao encontrado"})
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) { return }

	var creds LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Requisicao invalida", http.StatusBadRequest)
		return
	}

	// MVP: Admin hardcoded. Próximo passo: SELECT no Postgres.
	if creds.Username != "admin" || creds.Password != "123456" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Credenciais invalidas"})
		return
	}

	tokenString, err := utils.GenerateJWT(creds.Username)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}