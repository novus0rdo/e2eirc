package e2eirc

import (
	"fmt"
	"strings"
)

func PrintBanner() {
	fmt.Print(`
███████╗██████╗ ███████╗██╗██████╗  ██████╗
██╔════╝╚════██╗██╔════╝██║██╔══██╗██╔════╝
█████╗   █████╔╝█████╗  ██║██████╔╝██║     
██╔══╝  ██╔═══╝ ██╔══╝  ██║██╔══██╗██║     
███████╗███████╗███████╗██║██║  ██║╚██████╗
╚══════╝╚══════╝╚══════╝╚═╝╚═╝  ╚═╝ ╚═════╝
`)

	fmt.Println("Version: " + version())

	if beta() {
		fmt.Println(`
WARNING: YOU ARE ON A BETA RELEASE!
THE PURPOSE OF THIS RELEASE IS FOR EVALUATION,
SECURITY RESEARCH, AND DEVELOPMENT! THE SECURTIY
OF THIS RELEASE IS NOT GUARENTEED IN ANY WAY.
DO NOT USE THIS VERSION FOR MISSION CRITICAL
COMMUNICATIONS. YOU MAY NOT BE SAFE.
		`)
	}
}

func version() string {
	return "0.0.1-Beta"
}

func beta() bool {
	return strings.HasSuffix(version(), "Beta")
}
