package e2eirc

import (
	"crypto/rsa"
	"fmt"
	"io"
	"net"
	"strconv"
)

type connection struct {
	client *net.Conn
	server *net.Conn

	name string

	peerRegistry map[string]*peer

	encryptionKey []byte

	trustListening      bool
	trustQueue          chan *trustRequest
	pendingTrustRequest *trustRequest

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

const bufferSize = 32768

func serverAddr() string {
	return sharedConfig.ircHost + ":" + strconv.Itoa(sharedConfig.ircPort)
}

func connectionWithNetworkConnection(conn *net.Conn) *connection {
	c := connection{}
	c.client = conn
	c.generateRandomKey()
	c.loadRSAKey()
	c.loadRegistry()

	c.trustListening = false
	c.trustQueue = make(chan *trustRequest, 100)

	serverConn, _ := net.Dial("tcp", serverAddr())
	c.server = &serverConn

	go c.listenReadClient()
	go c.listenReadServer()

	fmt.Println("Connection established to " + serverAddr() + " <-> " + (*conn).LocalAddr().String())

	return &c
}

func (c *connection) close() {
	(*c.client).Close()
	(*c.server).Close()

	fmt.Println("Connection closed")
}

func readNetConnection(conn *net.Conn) ([]byte, error) {
	buf := make([]byte, bufferSize)

	l, err := (*conn).Read(buf)

	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading:", err.Error())
		}

		return buf, err
	}

	return buf[:l], nil
}

func (c *connection) listenReadClient() {
	for {
		buf, err := readNetConnection(c.client)
		if err != nil {
			c.close()
			break
		}

		c.readClient(buf)
	}
}

func (c *connection) listenReadServer() {
	for {
		buf, err := readNetConnection(c.server)
		if err != nil {
			c.close()
			break
		}

		c.readServer(buf)
	}
}

func (c *connection) readClient(buf []byte) {
	m := messageWithBytesAndConnection(buf, c)
	m.encrypt()
}

func (c *connection) readServer(buf []byte) {
	m := messageWithBytesAndConnection(buf, c)
	m.decrypt()
}

func (c connection) writeClient(buf []byte) {
	(*c.client).Write(buf)
}

func (c connection) writeServer(buf []byte) {
	(*c.server).Write(buf)
}
