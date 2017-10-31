package e2eirc

import "flag"

type config struct {
	ircPort int
	ircHost string

	localPort int
	localHost string

	keyPassword string
	dirPath     string
}

var sharedConfig = config{}

func ParseFlags() {
	flag.IntVar(&sharedConfig.ircPort, "port", 6667, "The port of the remote IRC server")
	flag.StringVar(&sharedConfig.ircHost, "host", "0.0.0.0", "The host of the remote IRC server")

	flag.IntVar(&sharedConfig.localPort, "local_port", 6666, "The local IRC port to connect to")
	flag.StringVar(&sharedConfig.localHost, "local_host", "0.0.0.0", "The address to listen on")

	flag.StringVar(&sharedConfig.keyPassword, "key", "", "The password to decrypt your private key, leave blank for STDIN")
	flag.StringVar(&sharedConfig.dirPath, "dir", "", "The path to the directory containing your private and trusted keys, leave blank for ~/.e2eirc")

	flag.Parse()
}
