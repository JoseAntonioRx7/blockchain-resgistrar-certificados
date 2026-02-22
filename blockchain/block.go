package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
)

type Block struct {
	Timestamp    int64
	Transactions []CertificateTransaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
}

func (b *Block) calculateHash() []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(strconv.FormatInt(b.Timestamp, 10)),
			[]byte(strconv.Itoa(b.Nonce)),
			b.PrevHash,
			[]byte(fmtTransactions(b.Transactions)),
		},
		[]byte{},
	)

	hash := sha256.Sum256(data)
	return hash[:]
}

func fmtTransactions(txs []CertificateTransaction) string {
	var txStr string
	for _, tx := range txs {
		txStr += tx.ID + tx.FileHash + strconv.FormatInt(tx.Timestamp, 10)
	}
	return txStr
}
