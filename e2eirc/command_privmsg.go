package e2eirc

import (
	"fmt"
	"strings"
)

func commandPrivMSG() *command {
	c := newCommand("PRIVMSG")
	c.encrypt = encryptPrivMSG
	c.decrypt = decryptPrivMSG

	return c
}

func encryptPrivMSG(m *message) {
	msgIndex := m.offset + 2

	msg := strings.Join(m.components[msgIndex:], " ")

	start := 0
	if msg[0] == ':' {
		start = 1
	}

	msg = msg[start : len(msg)-2]

	if m.components[m.offset+1] == trustUser {
		m.conn.handleTrustMessage(msg)
		return
	}

	trimmedComponents := m.components[:msgIndex]

	msg = m.conn.aesEncryptString(msg)

	result := strings.Join(trimmedComponents, " ") + " :" + msg + "\n"
	fmt.Print("<<<", result)

	m.conn.writeServer([]byte(result))
}

func decryptPrivMSG(m *message) {
	// Extract Sender
	sender := m.components[0][1:]
	m.senderComponents = strings.SplitN(sender, "!~", 2)

	sender = m.senderComponents[0]

	// If message is coming from server and has trust user's username
	// just block it
	if sender == trustUser {
		return
	}

	peer := m.conn.peerWithID(sender)
	peer.decryptPrivMSG(m)
}
