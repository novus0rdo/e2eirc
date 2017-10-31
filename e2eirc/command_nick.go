package e2eirc

func commandNICK() *command {
	c := newCommand("NICK")
	c.encrypt = encryptNICK
	c.decrypt = decryptNICK

	return c
}

func encryptNICK(m *message) {
	name := m.components[m.offset+1]
	m.conn.name = name[:len(name)-2]

	m.conn.writeServer(m.buf)
}

func decryptNICK(m *message) {
	m.conn.writeClient(m.buf)
}
