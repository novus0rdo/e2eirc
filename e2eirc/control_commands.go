package e2eirc

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func (p *peer) recieveControlCommand(msg string) {
	fmt.Println("Control command recieved!", msg)

	components := strings.SplitN(msg, " ", 3)
	command := components[1]
	payload := components[2]

	switch command {
	case "HANDSHAKE":
		p.recieveControlCommandHandshake(payload)
	case "SETNICK":
		p.recieveControlCommandSetNick(payload)
	case "SETKEY":
		p.recieveControlCommandSetKey(payload)
	}
}

func (p *peer) recieveControlCommandSetNick(payload string) {
	p.name = payload
}

func (p *peer) recieveControlCommandHandshake(payload string) {
	p.setPublicKey(payload)

	if p.canTrust() {
		p.sendEncryptionKey()
		return
	}

	if !p.handshaking {
		p.beginHandshake()
	}

	p.requestTrust()
}

func (p *peer) recieveControlCommandSetKey(payload string) {
	key, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return
	}

	p.decryptionKey = []byte(key)
	p.decryptUnencryptedMessages()
}
