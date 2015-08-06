package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

//MD5 encription
func MD5(cad string) string {
	hash := sha1.New()
	hash.Write([]byte(cad))
	return hex.EncodeToString(hash.Sum(nil))
}

//GenerateToken a random
func GenerateToken(size int) string {
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}
	return base64.URLEncoding.EncodeToString(rb)
}
