package e2eirc

type cryptoFunc func(m *message)

type command struct {
	name    string
	encrypt cryptoFunc
	decrypt cryptoFunc
}

var commands = map[string]*command{}

func newCommand(name string) *command {
	c := command{name: name}
	return &c
}

func (c *command) register() {
	commands[c.name] = c
}

func RegisterCommands() {
	commandPrivMSG().register()
	commandNICK().register()
}
