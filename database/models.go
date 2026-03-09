package database

import "time"

// Institution representa a universidade ou escola (ex: UFPE, USP, etc.)
type Institution struct {
    ID           int       `json:"id"`
    Name         string    `json:"name"`           // Nome completo da instituição
    Username     string    `json:"username"`       // Login único
    PasswordHash string    `json:"-"`              // Oculto no JSON por segurança
    PublicKey    []byte    `json:"public_key"`     // Usado para validar que a instituição X gerou o bloco
    CreatedAt    time.Time `json:"created_at"`
}

// Certificate representa o diploma registrado na rede
type Certificate struct {
    ID             int       `json:"id"`
    InstitutionID  int       `json:"institution_id"` // A mágica acontece aqui (Chave Estrangeira)
    StudentName    string    `json:"student_name"`
    Course         string    `json:"course"`
    FileHash       string    `json:"file_hash"`      // Hash do PDF original
    BlockchainHash string    `json:"blockchain_hash"`// Hash do bloco onde foi minerado
    IssueDate      time.Time `json:"issue_date"`
}	