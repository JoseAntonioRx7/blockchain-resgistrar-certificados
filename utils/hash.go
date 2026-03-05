package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
)

func HashFile(file multipart.File) (string, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
