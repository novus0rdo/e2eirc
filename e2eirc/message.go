package e2eirc

import (
	"strings"
)

type message struct {
	buf  []byte
	conn *connection

	components       []string
	senderComponents []string
	offset           int
}

func messageWithBytesAndConnection(buf []byte, conn *connection) message {
	return message{buf: buf, conn: conn}
}

func (m *message) encrypt() {
	str := string(m.buf)
	m.components = strings.Split(str, " ")

	commandName := m.components[0]
	command := commands[commandName]
	if command != nil {
		command.encrypt(m)
		return
	}

	m.conn.writeServer(m.buf)
}

func (m *message) decrypt() {
	str := string(m.buf)

	m.components = strings.Split(str, " ")
	m.offset = 1

	commandName := m.components[1]
	command := commands[commandName]
	if command != nil {
		command.decrypt(m)
		return
	}

	m.conn.writeClient(m.buf)
}
