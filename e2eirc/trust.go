package e2eirc

const trustUser = "$E2ECtrl"

type trustRequest struct {
	peer   *peer
	result chan bool
}

func newTrustRequest(p *peer) *trustRequest {
	request := trustRequest{}
	request.peer = p
	request.result = make(chan bool, 1)
	return &request
}

func (p *peer) canTrust() bool {
	if p.trusted {
		return true
	}

	if p.publicKey != nil {
		return isTrustedKey(p.publicKeyFingerprint())
	}

	return false
}

func (p *peer) requestTrust() {
	go p.conn.addTrustRequest(p)
}

func (p *peer) trust() {
	p.trusted = true
	addTrustedKey(p.publicKeyFingerprint())
}

func (c *connection) addTrustRequest(p *peer) {
	if !c.trustListening {
		c.trustListening = true
		go c.trustListen()
	}
	c.trustQueue <- newTrustRequest(p)
}

func (c *connection) trustListen() {
	for {
		request := <-c.trustQueue
		request.send()
		<-request.result
	}
}

func (r *trustRequest) send() {
	connection := r.peer.conn
	connection.pendingTrustRequest = r

	connection.sendTrustMessage(r.peer.displayName() + " would like to exchange keys with you. Do you trust this user? (Yes / No)")
}

func (r *trustRequest) setResult(result bool) {
	r.peer.trusted = result

	if r.peer.trusted {
		r.peer.trust()
		r.peer.sendEncryptionKey()
		r.peer.conn.sendTrustMessage(r.peer.displayName() + " has been trusted!")
	} else {
		r.peer.conn.sendTrustMessage(r.peer.displayName() + " has not been trusted. They will be unable to decrypt messages from you.")
		removeTrustedKey(r.peer.publicKeyFingerprint())
	}

	r.peer.conn.pendingTrustRequest = nil

	r.result <- result
}

func (c *connection) sendTrustMessage(payload string) {
	//msg := ""
	msg := ":" + trustUser + "!~trust@internal PRIVMSG " + c.name + " :" + payload + "\n"
	c.writeClient([]byte(msg))
}

func (c *connection) handleTrustMessage(message string) {

	switch message {
	case "REGENKEY":
		c.handleRegenKey(message)
	case "LISTTRUST":
		c.handleListTrust(message)
	default:
		c.handleTrustResultMessage(message)
	}
}

func (c *connection) handleListTrust(message string) {
	for _, p := range c.peerRegistry {
		trusted := "TRUSTED"

		if !p.canTrust() {
			trusted = "UNTRUSTED"
		}

		c.sendTrustMessage(p.displayName() + " " + trusted)
	}
}

func (c *connection) handleRegenKey(message string) {
	c.generateRandomKey()
	c.sendTrustMessage("Key Regenerated")

	for _, peer := range c.peerRegistry {
		if peer.canTrust() {
			peer.sendControlCommand("SETKEY", c.encryptionKeyBase64())
			c.sendTrustMessage("Notifying " + peer.displayName() + " of change")
		}
	}
}

func (c *connection) handleTrustResultMessage(message string) {
	accept := message == "Y" || message == "y" || message == "yes" || message == "Yes" || message == "YES"
	decline := message == "N" || message == "n" || message == "no" || message == "No" || message == "NO"

	if c.pendingTrustRequest == nil && (accept || decline) {
		c.sendTrustMessage("No pending requests. Type HELP for more info.")
		return
	}

	if accept {
		c.pendingTrustRequest.setResult(true)
		return
	} else if decline {
		c.pendingTrustRequest.setResult(false)
		return
	}

	c.sendTrustMessage("Command not recognized: '" + message + "'. Type HELP for more info.")

	if c.pendingTrustRequest != nil {
		c.pendingTrustRequest.send()
	}
}
