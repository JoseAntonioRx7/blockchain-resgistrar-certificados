package blockchain

type CertificateTransaction struct {
	ID          string `json:"id"`
	StudentName string `json:"student_name"`
	Institution string `json:"institution"`
	Course      string `json:"course"`
	FileHash    string `json:"file_hash"`
	Signature   string `json:"signature"`
	Timestamp   int64  `json:"timestamp"`

}
