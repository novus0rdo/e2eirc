package e2eirc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func (c *connection) generateRandomKey() {
	key := make([]byte, 32)
	_, err := rand.Read(key)

	if err != nil {
		fmt.Println("Fatal: Could not generate random key")
		os.Exit(0)
	}

	c.encryptionKey = key
}

func (c *connection) encryptionKeyBase64() string {
	return base64.StdEncoding.EncodeToString(c.encryptionKey)
}

func (c *connection) aesEncryptString(str string) string {
	in := []byte(str)

	block, err := aes.NewCipher(c.encryptionKey)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(in))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], in)

	msg := base64.StdEncoding.EncodeToString(ciphertext)
	return msg
}

func (p *peer) aesDecryptString(str string) string {
	ciphertext, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		return str
	}

	block, err := aes.NewCipher(p.decryptionKey)
	if err != nil {
		return str
	}

	if len(ciphertext) < aes.BlockSize {
		return str
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func (p *peer) sendEncryptionKey() {
	p.sendControlCommandRSA("SETKEY", p.conn.encryptionKeyBase64())
}
