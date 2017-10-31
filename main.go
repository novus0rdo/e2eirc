package main

import "github.com/Novus0rdo/e2eirc/e2eirc"

func main() {
	e2eirc.PrintBanner()
	e2eirc.ParseFlags()
	e2eirc.RegisterCommands()
	e2eirc.Unlock()
	e2eirc.Start()
}
