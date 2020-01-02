package esi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
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

func (c *Crypt) Encrypt(text string) (string, error) {
	hasher := sha256.New()
	hasher.Write([]byte(c.key))
	key := base64.URLEncoding.EncodeToString(hasher.Sum(nil))[:32]
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize + len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	result := fmt.Sprintf("%s:%s", hex.EncodeToString(iv), hex.EncodeToString(ciphertext[aes.BlockSize:]))
	return result, nil
}