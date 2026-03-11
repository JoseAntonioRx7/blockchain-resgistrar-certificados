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


// --- AUXILIARES ---

func sendJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

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

// --- HANDLERS PÚBLICOS ---

// Corrigido para bater com o routes.go
func RegisterInstitutionHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) { return }

	var data struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}

	// Gera chaves para a nova instituição
	pubKey, privKey, err := utils.GenerateInstitutionKeys()
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao gerar chaves"})
		return
	}

	query := `INSERT INTO institutions (name, username, password_hash, public_key) VALUES ($1, $2, $3, $4)`
	_, err = database.DB.Exec(query, data.Name, data.Username, data.Password, pubKey)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Username já existe ou erro no banco"})
		return
	}

	// Retorna a chave privada apenas uma vez para o administrador
	sendJSON(w, http.StatusCreated, map[string]string{
		"message":     "Instituição criada com sucesso!",
		"private_key": privKey,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) { return }

	var data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}

	var storedPassword string
	err := database.DB.QueryRow("SELECT password_hash FROM institutions WHERE username = $1", data.Username).Scan(&storedPassword)

	// O Erro 401 acontece se esta verificação falhar
	if err != nil || data.Password != storedPassword {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário ou senha incorretos"})
		return
	}

	token, _ := utils.GenerateJWT(data.Username)
	sendJSON(w, http.StatusOK, map[string]string{"token": token})
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if setupCORS(&w, r) { return }

	hash := r.URL.Query().Get("hash")
	if hash == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Hash vazio"})
		return
	}

	if Chain != nil {
		for _, block := range Chain.Blocks {
			for _, tx := range block.Transactions {
				if tx.FileHash == hash {
					sendJSON(w, http.StatusOK, tx)
					return
				}
			}
		}
	}
	sendJSON(w, http.StatusNotFound, map[string]string{"error": "Certificado não encontrado"})
}

// --- HANDLERS PROTEGIDOS ---

func ListCertificatesHandler(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Erro de contexto"})
		return
	}

	results := []blockchain.CertificateTransaction{}
	query := `SELECT c.id, c.student_name, c.course, c.file_hash, c.blockchain_hash, c.timestamp
			  FROM certificates c JOIN institutions i ON c.institution_id = i.id
			  WHERE i.username = $1 ORDER BY c.timestamp DESC`

	rows, err := database.DB.Query(query, username)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tx blockchain.CertificateTransaction
			var bHash string
			rows.Scan(&tx.ID, &tx.StudentName, &tx.Course, &tx.FileHash, &bHash, &tx.Timestamp)
			tx.Signature = bHash
			results = append(results, tx)
		}
	}
	sendJSON(w, http.StatusOK, results)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Limite de upload (32MB)
    if err := r.ParseMultipartForm(32 << 20); err != nil {
        sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Erro no formulário"})
        return
    }

    // Identificação da Instituição
    username := r.Context().Value("username").(string)
    var instID int
    var instName string
    err := database.DB.QueryRow("SELECT id, name FROM institutions WHERE username = $1", username).Scan(&instID, &instName)
    if err != nil {
        sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Instituição não encontrada"})
        return
    }

    // 2. Recuperar o arquivo (Corrigido: _ substitui handler para evitar erro de compilação)
    file, _, err := r.FormFile("file") //
    if err != nil {
        sendJSON(w, http.StatusBadRequest, map[string]string{"error": "Arquivo ausente no formulário"})
        return
    }
    defer file.Close()

    // 3. Gerar Hash do arquivo
    hash, err := utils.HashFile(file)
    if err != nil {
        sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao gerar hash"})
        return
    }
    
    // IMPORTANTE: Resetar ponteiro do arquivo para leitura do salvamento
    file.Seek(0, 0)

    // 4. Salvar o arquivo fisicamente usando a lógica centralizada no utils (pdf.go)
    // Isso garante que o nome do arquivo seja cert_[HASH].pdf para bater com o frontend
    _, err = utils.SaveUploadedCertificate(file, hash) 
    if err != nil {
        fmt.Println("Aviso: Falha ao salvar arquivo físico:", err)
    }

    // 5. Preparar Transação da Blockchain
    tx := blockchain.CertificateTransaction{
        ID:            utils.GenerateID(),
        StudentName:   r.FormValue("student_name"), 
        Institution:   instName,
        InstitutionID: instID,
        Course:        r.FormValue("course"),
        FileHash:      hash,
        Signature:     r.FormValue("private_key"), 
        Timestamp:     time.Now().Unix(),
    }

    // 6. Registro no Banco de Dados
    _, err = database.DB.Exec(`INSERT INTO certificates (id, institution_id, student_name, course, file_hash, blockchain_hash, timestamp) 
                      VALUES ($1, $2, $3, $4, $5, $6, $7)`, tx.ID, instID, tx.StudentName, tx.Course, tx.FileHash, tx.Signature, tx.Timestamp)
    
    if err != nil {
        sendJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro ao salvar no banco: " + err.Error()})
        return
    }

    // 7. Minerar Bloco
    if Chain != nil {
        Chain.AddBlock([]blockchain.CertificateTransaction{tx})
    }
    
    // Resposta de Sucesso
    sendJSON(w, http.StatusOK, map[string]string{
        "message": "Certificado minerado com sucesso!",
        "hash":    hash,
    })
}