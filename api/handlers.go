package api

import (
	"cert-chain/blockchain"
	"cert-chain/database"
	"cert-chain/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Chain *blockchain.Blockchain

// 1. Utilitário de CORS (Centralizado)
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

// 2. Registro de Certificados (Multi-Tenant)
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}

	// Extrai o token do header Authorization PRIMEIRO (antes de ParseMultipartForm)
	authHeader := r.Header.Get("Authorization")
	var tokenStr string
	if authHeader != "" && len(authHeader) > 7 {
		tokenStr = authHeader[7:] // Remove "Bearer "
	}

	// Parse o multipart form primeiro - sem limite de tamanho
	if err := r.ParseMultipartForm(0); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao processar formulario: " + err.Error()})
		return
	}

	// Se token não veio do header, tenta do FormData
	if tokenStr == "" {
		tokenStr = r.FormValue("token")
	}

	if tokenStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token nao fornecido"})
		return
	}

	// Valida o token JWT manualmente
	token, err := utils.ValidateJWT(tokenStr)
	if err != nil || !token.Valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token invalido ou expirado: " + err.Error()})
		return
	}

	// Extrai o username do token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token invalido"})
		return
	}
	username := claims["username"].(string)

	// Busca ID, Nome e Chave Privada da Instituição logada
	var instID int
	var instName, privateKeyHex string
	queryInst := `SELECT id, name, private_key FROM institutions WHERE username = $1`
	err = database.DB.QueryRow(queryInst, username).Scan(&instID, &instName, &privateKeyHex)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Instituicao nao encontrada: " + err.Error()})
		return
	}

	// Processamento do Arquivo
	// Primeiro obtemos os valores do form
	studentName := r.FormValue("student_name")
	course := r.FormValue("course")

	// Agora obtemos o arquivo
	file, _, err := r.FormFile("file")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao ler arquivo: " + err.Error()})
		return
	}
	defer file.Close()

	hash, err := utils.HashFile(file)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao processar hash do arquivo: " + err.Error()})
		return
	}

	file.Seek(0, 0) // Rebobina para o passo de salvar o PDF

	// Assinatura com a chave privada única da instituição
	signature, err := utils.SignData(hash, privateKeyHex)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao assinar: " + err.Error()})
		return
	}

	tx := blockchain.CertificateTransaction{
		ID:            utils.GenerateID(),
		StudentName:   studentName,
		Institution:   instName,
		InstitutionID: instID, // Vincula ao ID do banco
		Course:        course,
		FileHash:      hash,
		Signature:     signature,
		Timestamp:     time.Now().Unix(),
	}

	// 5. Persistência no PostgreSQL
	queryInsert := `INSERT INTO certificates (id, institution_id, student_name, course, file_hash, blockchain_hash, timestamp) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = database.DB.Exec(queryInsert, tx.ID, instID, tx.StudentName, tx.Course, tx.FileHash, tx.Signature, tx.Timestamp)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao salvar no banco: " + err.Error()})
		return
	}

	// Adiciona na Blockchain e salva o PDF
	Chain.AddBlock([]blockchain.CertificateTransaction{tx})
	pdfPath, _ := utils.SaveUploadedCertificate(file, tx.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Certificado assinado pela " + instName,
		"id":      tx.ID,
		"hash":    hash,
		"pdf":     pdfPath,
	})
}

// 3. Listagem Filtrada (Onde a mágica acontece)
func ListCertificatesHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}

	// Obtém o token do header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 {
		http.Error(w, `{"error": "Token nao fornecido"}`, http.StatusUnauthorized)
		return
	}
	tokenStr := authHeader[7:]

	// Valida o token
	token, err := utils.ValidateJWT(tokenStr)
	if err != nil || !token.Valid {
		http.Error(w, `{"error": "Token invalido ou expirado"}`, http.StatusUnauthorized)
		return
	}

	// Extrai o username
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, `{"error": "Token invalido"}`, http.StatusUnauthorized)
		return
	}
	username := claims["username"].(string)

	// Filtra os certificados no banco para mostrar apenas os desta instituição
	// Incluímos o ID para os links de download do PDF
	query := `SELECT c.id, c.student_name, c.course, c.file_hash, c.blockchain_hash, c.timestamp
			  FROM certificates c
			  JOIN institutions i ON c.institution_id = i.id
			  WHERE i.username = $1
			  ORDER BY c.timestamp DESC`

	rows, err := database.DB.Query(query, username)
	if err != nil {
		http.Error(w, `{"error": "Erro ao buscar dados"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []blockchain.CertificateTransaction
	for rows.Next() {
		var tx blockchain.CertificateTransaction
		var blockchainHash string
		// Scan into correct fields - blockchain_hash goes to blockchainHash, not Signature
		err := rows.Scan(&tx.ID, &tx.StudentName, &tx.Course, &tx.FileHash, &blockchainHash, &tx.Timestamp)
		if err != nil {
			log.Printf("Erro ao scanear linha: %v", err)
			continue
		}
		tx.Signature = blockchainHash // Use blockchain_hash as signature for compatibility
		results = append(results, tx)
	}

	// Retorna array vazio se não houver resultados
	if results == nil {
		results = []blockchain.CertificateTransaction{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// 4. Verificação Universal (Pública)
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}
	hash := r.URL.Query().Get("hash")

	for _, block := range Chain.Blocks {
		for _, tx := range block.Transactions {
			if tx.FileHash == hash {
				json.NewEncoder(w).Encode(tx)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Hash nao encontrado na rede"})
}

// RegisterInstitutionHandler cria uma nova instituição com seu próprio par de chaves
func RegisterInstitutionHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}

	var data struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Dados invalidos", http.StatusBadRequest)
		return
	}

	// 1. Gera o par de chaves exclusivo para esta instituição
	pubKey, privKey, err := utils.GenerateInstitutionKeys()
	if err != nil {
		http.Error(w, "Erro ao gerar chaves criptograficas", http.StatusInternalServerError)
		return
	}

	// 2. Insere na tabela 'institutions' que você criou no pgAdmin
	query := `INSERT INTO institutions (name, username, password_hash, public_key, private_key) 
			  VALUES ($1, $2, $3, $4, $5)`

	// Nota: Em produção, use bcrypt para o password_hash
	_, err = database.DB.Exec(query, data.Name, data.Username, data.Password, pubKey, privKey)
	if err != nil {
		http.Error(w, "Erro ao salvar instituicao: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Retorna a chave privada (A instituição deve guardar isso com segurança!)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Instituicao cadastrada com sucesso!",
		"private_key": fmt.Sprintf("%x", privKey),
	})
}

// DebugUsersHandler lista todos os usuários (apenas para debug)
func DebugUsersHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}

	query := `SELECT id, name, username, password_hash FROM institutions`
	rows, err := database.DB.Query(query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, username, passwordHash string
		if err := rows.Scan(&id, &name, &username, &passwordHash); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{
			"id":            id,
			"name":          name,
			"username":      username,
			"password_hash": passwordHash,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// LoginHandler autentica a instituição e retorna um token JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Metodo nao permitido",
		})
		return
	}

	defer r.Body.Close()

	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Dados invalidos",
		})
		return
	}

	if data.Username == "" || data.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Username e password sao obrigatorios",
		})
		return
	}

	var storedPassword string
	query := `SELECT password_hash FROM institutions WHERE username = $1`
	err := database.DB.QueryRow(query, data.Username).Scan(&storedPassword)

	if err != nil || data.Password != storedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Credenciais invalidas",
		})
		return
	}

	token, err := utils.GenerateJWT(data.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao gerar token",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login realizado com sucesso",
		"token":   token,
	})
}
