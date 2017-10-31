package e2eirc

import (
	"crypto/rsa"
	"fmt"
	"strings"
)

type peer struct {
	id   string
	name string

	conn                *connection
	unencryptedMessages []*message

	decryptionKey []byte
	publicKey     *rsa.PublicKey

	trusted     bool
	handshaking bool
}

func (c *connection) loadRegistry() {
	c.peerRegistry = map[string]*peer{}
}

func (c *connection) peerWithID(id string) *peer {
	peer := c.peerRegistry[id]

	if peer == nil {
		peer = c.newPeerWithID(id)
	}

	return peer
}

func (c *connection) newPeerWithID(id string) *peer {
	p := peer{id: id, conn: c}
	p.name = id
	p.unencryptedMessages = []*message{}
	c.peerRegistry[id] = &p
	return &p
}

func (p *peer) decryptPrivMSG(m *message) {
	msgIndex := m.offset + 2

	msg := strings.Join(m.components[msgIndex:], " ")
	trimmedComponents := m.components[:msgIndex]

	plaintext := true

	// Extract Message
	msg = msg[1 : len(msg)-2]

	// Check if RSA Message
	if strings.HasPrefix(msg, "RSA ") {
		msg = p.conn.rsaDecrypt(msg)
		plaintext = false
	}

	// Check for non encrypted control command
	if strings.HasPrefix(msg, "CTRL ") {
		p.recieveControlCommand(msg)
		return
	}

	// Determine if we can decrypt or if we have to request the decryption key
	if !p.canDecrypt() {
		p.beginHandshake()
		p.unencryptedMessages = append(p.unencryptedMessages, m)
		return
	}

	// Decrypt Message
	decMsg := p.aesDecryptString(msg)
	if decMsg != msg {
		plaintext = false
	}

	msg = decMsg

	// Check for encrypted control command
	if strings.HasPrefix(msg, "CTRL ") {
		p.recieveControlCommand(msg)
		return
	}

	// Properly Name Sender
	m.senderComponents[0] = p.name

	sender := strings.Join(m.senderComponents, "!~")
	trimmedComponents[0] = ":" + sender

	// Add cleartext warning
	if plaintext {
		msg += " [SENT AS CLEARTEXT]"
	}

	// Rebuild Message string
	result := strings.Join(trimmedComponents, " ") + " :" + msg + "\n"

	p.conn.writeClient([]byte(result))
}

func (p *peer) canDecrypt() bool {
	return p.decryptionKey != nil
}

//
// Handshakes
//

func (p *peer) beginHandshake() {
	p.sendControlCommand("HANDSHAKE", p.conn.publicKeyBase64())
	p.handshaking = true
}

//
// Control Commands
//

// Basics
func (p *peer) sendControlCommand(command string, payload string) {
	msg := "CTRL " + command

	if payload != "" {
		msg += " " + payload
	}

	p.sendPrivateMessage(msg)
}

func (p *peer) sendControlCommandRSA(command string, payload string) {
	msg := "CTRL " + command

	if payload != "" {
		msg += " " + payload
	}

	msg = p.rsaEncrypt(msg)

	p.sendPrivateMessage(msg)
}

func (p *peer) sendPrivateMessage(message string) {
	msg := "PRIVMSG " + p.id + " :" + message
	p.conn.writeServer([]byte(msg + "\n"))
}

//
// After the handshake, decrypt the messages that were sent previously
//
func (p *peer) decryptUnencryptedMessages() {
	if p.canDecrypt() {
		for _, message := range p.unencryptedMessages {
			message.decrypt()
		}

		p.unencryptedMessages = []*message{}
	}
}

func (p *peer) displayName() string {
	fingerprintHex := p.publicKeyFingerprintSha1()
	fingerprintChunks := make([]string, len(fingerprintHex)/2)
	for i := range fingerprintChunks {
		fingerprintChunks[i] = string(fingerprintHex[(i*2)]) + string(fingerprintHex[(i*2)+1])
	}

	fingerprint := strings.Join(fingerprintChunks, ":")

	fmt.Println(fingerprint, fingerprintHex)

	return p.name + " (" + fingerprint + ")"
}
