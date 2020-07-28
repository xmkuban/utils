package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math"
	"os"
)

const filechunk = 8192 // we settle for 8KB
func GetSha256Bytes(data []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func GetSha256BytesString(data []byte) (string, error) {
	resByte, err := GetSha256Bytes(data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(resByte), nil
}

func GetSha256(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, _ := file.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
	}

	checksum := hash.Sum(nil)

	return checksum, nil
}

func GetSha256String(filePath string) (string, error) {
	checksum, err := GetSha256(filePath)
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(checksum), nil
}
