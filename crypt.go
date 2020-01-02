package esi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

type Crypt struct {
	key string
}

func (c *Crypt) Decrypt(cipherText string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(c.key))
	key := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:32]
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	s := strings.Split(cipherText, ":")
	iv, err := hex.DecodeString(s[0])

	if err != nil {
		return "", fmt.Errorf("Unable to decode iv: %v", err)
	}

	cs, _ := hex.DecodeString(s[1])

	//if len(s[1])%aes.BlockSize != 0 {
	//	return "", fmt.Errorf("cipherText is not a multiple of the block size")
	//}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cs, cs)

	result := string(cs[:])
	return result, nil
}

