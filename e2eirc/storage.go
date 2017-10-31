package e2eirc

import (
	"log"
	"os"
	"os/user"
)

const perms = 0600
const dirPerms = 0700

func dirPath() string {
	path := sharedConfig.dirPath

	if path == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		path = usr.HomeDir + "/.e2eirc"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, dirPerms)
	}

	return path
}

func registryPath() string {
	return dirPath() + "/trusted_keys"
}
