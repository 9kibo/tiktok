package utils

import (
	"encoding/base64"
	"os"
)

func ReadBase64File(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	dest := make([]byte, 0, 0)
	_, err = base64.StdEncoding.Decode(dest, bytes)
	return dest, err
}
