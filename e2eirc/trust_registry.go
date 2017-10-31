package e2eirc

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var trustRegistry = map[string]bool{}
var trustRegistryLoaded = false

func loadTrustRegistry() {
	if !trustRegistryLoaded {
		trustRegistryLoaded = true
		path := registryPath()
		file, err := ioutil.ReadFile(path)
		if err == nil {
			reg := rsaDecrypt(sharedPrivateKey, string(file))
			if reg != "" {
				lines := strings.Split(string(reg), "\n")
				for _, key := range lines {
					trustRegistry[key] = true
				}
			}
		}
	}
}

func isTrustedKey(fingerprint string) bool {
	loadTrustRegistry()
	return trustRegistry[fingerprint]
}

func addTrustedKey(fingerprint string) {
	loadTrustRegistry()
	trustRegistry[fingerprint] = true
	saveTrustRegistry()
}

func removeTrustedKey(fingerprint string) {
	loadTrustRegistry()
	delete(trustRegistry, fingerprint)
	saveTrustRegistry()
}

func saveTrustRegistry() {
	lines := make([]string, len(trustRegistry))
	index := 0
	for key := range trustRegistry {
		lines[index] = key
		index++
	}

	linesString := strings.Join(lines, "\n")

	data := rsaEncrypt(&sharedPrivateKey.PublicKey, linesString)

	err := ioutil.WriteFile(registryPath(), []byte(data), perms)
	fmt.Println(err, "<ER")
}
