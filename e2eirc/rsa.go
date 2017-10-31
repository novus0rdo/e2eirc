package e2eirc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var sharedPrivateKey *rsa.PrivateKey

func getPassword(msg string, err bool) string {
	if !err && sharedConfig.keyPassword != "" {
		return sharedConfig.keyPassword
	}

	fmt.Print(msg)
	fmt.Print(": ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))

	pass := string(bytePassword)

	fmt.Println("\n")

	return strings.TrimSpace(pass)
}

func Unlock() {
	loadOrGenerateRSAKey(false)
}

func generateRSAKey() *rsa.PrivateKey {
	reader := rand.Reader
	bitSize := 2048

	fmt.Print("Generating " + strconv.Itoa(bitSize) + " bit private key. Please wait...")
	key, _ := rsa.GenerateKey(reader, bitSize)

	blockType := "RSA PRIVATE KEY"

	cipherType := x509.PEMCipherAES256

	fmt.Println(" Done")

	password := getPassword("Enter a new password for your private key. If you lose it you won't be able to confirm your identity on chat", false)

	encryptedPEMBlock, _ := x509.EncryptPEMBlock(rand.Reader,
		blockType,
		x509.MarshalPKCS1PrivateKey(key),
		[]byte(password),
		cipherType)

	pemdata := pem.EncodeToMemory(encryptedPEMBlock)
	ioutil.WriteFile(dirPath()+"/private_key.pem", []byte(pemdata), perms)

	return key
}

func loadOrGenerateRSAKey(retry bool) *rsa.PrivateKey {
	path := dirPath() + "/private_key.pem"

	_, err := os.Stat(path)
	notExist := os.IsNotExist(err)

	if !notExist {
		if !retry {
			fmt.Print("Key found! Loading...")
		}

		file, err := ioutil.ReadFile(path)
		if err == nil {
			block, _ := pem.Decode(file)

			if !retry {
				fmt.Println(" Done")
			}

			var password string

			if retry {
				password = getPassword("Password Incorrect. Try Again", true)
			} else {
				password = getPassword("Enter Private Key Password", false)
			}

			b, err := x509.DecryptPEMBlock(block, []byte(password))

			if err != nil {
				return loadOrGenerateRSAKey(true)
			}

			key, err := x509.ParsePKCS1PrivateKey(b)

			if err == nil {
				sharedPrivateKey = key
				return sharedPrivateKey
			}
		}
	}

	sharedPrivateKey = generateRSAKey()
	return sharedPrivateKey
}

func (c *connection) loadRSAKey() {
	if sharedPrivateKey == nil {
		loadOrGenerateRSAKey(false)
	}

	c.privateKey = sharedPrivateKey
	c.publicKey = &sharedPrivateKey.PublicKey
}

func (p *peer) rsaEncrypt(buf string) string {
	return rsaEncrypt(p.publicKey, buf)
}

func rsaEncrypt(key *rsa.PublicKey, buf string) string {
	secretMessage := []byte(buf)

	rng := rand.Reader

	ciphertext, err := rsa.EncryptPKCS1v15(rng, key, secretMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return ""
	}

	return "RSA " + base64.StdEncoding.EncodeToString(ciphertext)
}

func (c *connection) rsaDecrypt(buf string) string {
	return rsaDecrypt(c.privateKey, buf)
}

func rsaDecrypt(key *rsa.PrivateKey, buf string) string {
	// Remove RSA Prefix
	buf = buf[4:]

	b, err := base64.StdEncoding.DecodeString(buf)

	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader

	plaintext, err := rsa.DecryptPKCS1v15(rng, key, b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		return ""
	}

	return string(plaintext)
}

func (c *connection) publicKeyBase64() string {
	bytes, _ := x509.MarshalPKIXPublicKey(c.publicKey)
	b64 := base64.StdEncoding.EncodeToString(bytes)

	return b64
}

func (p *peer) setPublicKey(b64String string) {
	b, err := base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		return
	}

	publicKeyI, err := x509.ParsePKIXPublicKey(b)

	if err != nil {
		return
	}

	publicKey := publicKeyI.(*rsa.PublicKey)
	p.publicKey = publicKey
}

func (p *peer) publicKeyFingerprint() string {
	h := sha256.New()

	bytes, _ := x509.MarshalPKIXPublicKey(p.publicKey)
	h.Write(bytes)

	return hex.EncodeToString(h.Sum(nil))
}

func (p *peer) publicKeyFingerprintSha1() string {
	h := sha1.New()

	bytes, _ := x509.MarshalPKIXPublicKey(p.publicKey)
	h.Write(bytes)

	return hex.EncodeToString(h.Sum(nil))
}
