package blockchain

type CertificateTransaction struct {
	ID          	string `json:"id"`
	StudentName 	string `json:"student_name"`
	Institution 	string `json:"institution"`
	InstitutionID 	int   `json:"institution_id"` // Para vincular ao banco
	Course      	string `json:"course"`
	FileHash    	string `json:"file_hash"`
	Signature   	string `json:"signature"`
	Timestamp   	int64  `json:"timestamp"`

}
