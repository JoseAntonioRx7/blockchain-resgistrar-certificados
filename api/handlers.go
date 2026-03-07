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

// Função auxiliar para configurar CORS e Headers de Autorização
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

	// 1. Obtém o username do contexto (injetado pelo JWTMiddleware)
	usernameCtx := r.Context().Value("username")
	if usernameCtx == nil {
		http.Error(w, "Erro de autenticação: usuário não identificado", http.StatusUnauthorized)
		return
	}
	username := usernameCtx.(string)

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

	// 2. Busca os dados reais da instituição no PostgreSQL
	var privateKeyHex, institutionName string
	queryDB := `SELECT private_key, name FROM institutions WHERE username = $1`
	err = database.DB.QueryRow(queryDB, username).Scan(&privateKeyHex, &institutionName)
	if err != nil {
		http.Error(w, "Instituição não autorizada no banco de dados", http.StatusUnauthorized)
		return
	}

	// 3. Gera o Hash do arquivo
	hash, err := utils.HashFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. ASSINATURA DIGITAL REAL (Utilizando a chave privada da instituição logada)
	signature, err := utils.SignData(hash, privateKeyHex)
	if err != nil {
		http.Error(w, "Erro ao realizar assinatura digital institucional", http.StatusInternalServerError)
		return
	}

	tx := blockchain.CertificateTransaction{
		ID:           utils.GenerateID(),
		StudentName:  r.FormValue("student_name"),
		Institution:  institutionName, // Nome oficial vindo do banco (impede fraude no front)
		Course:       r.FormValue("course"),
		FileHash:     hash,
		Signature:    signature,
		Timestamp:    time.Now().Unix(),
	}

	// 5. Persistência no PostgreSQL
	queryInsert := `INSERT INTO certificates (id, student_name, institution, course, file_hash, signature, timestamp) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = database.DB.Exec(queryInsert, 
		tx.ID, tx.StudentName, tx.Institution, tx.Course, tx.FileHash, tx.Signature, tx.Timestamp,
	)
	if err != nil {
		http.Error(w, "Erro ao salvar registro no banco de dados", http.StatusInternalServerError)
		return
	}

	// 6. Mineração na Blockchain
	Chain.AddBlock([]blockchain.CertificateTransaction{tx})

	// 7. Geração Automática do PDF
	pdfPath, err := utils.GenerateCertificatePDF(tx)
	if err != nil {
		fmt.Println("Alerta: Falha na geração do PDF:", err)
	}

	// 8. Resposta de Sucesso
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Certificado registrado com sucesso e assinado digitalmente!",
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

	var storedPassword string
	query := `SELECT password FROM institutions WHERE username = $1`
	
	err := database.DB.QueryRow(query, creds.Username).Scan(&storedPassword)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario nao encontrado"})
		return
	}

	if creds.Password != storedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Senha incorreta"})
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