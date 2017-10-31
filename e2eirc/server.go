package e2eirc

import (
	"fmt"
	"strconv"
)

func Start() {
	fmt.Println("Starting TCP server on port " + sharedConfig.localHost + ":" + strconv.Itoa(sharedConfig.localPort))
	fmt.Println("Connect your IRC client to " + sharedConfig.localHost + "/" + strconv.Itoa(sharedConfig.localPort))

	fmt.Println("Relaying to E2E Messages to " + serverAddr())

	fmt.Print("\n\n---------------------------------------------------\n\n")

	startTCP(sharedConfig.localHost, sharedConfig.localPort)
}
