package rill

// The Command interface represent the entire command
// description.
type Command interface {
	CommandDescription() string
	CommandLabel() string
	CommandAliases() []string
	Clone() interface{}
}

type Commands struct {
	commands map[string]*NamespaceCommand
}

func Rill() *Commands {
	return &Commands{make(map[string]*NamespaceCommand)}
}

func (c *Commands) Add(command Command, shortDesc string, desc string) *NamespaceCommand {
	namespace := newNamespaceCommand(c, command, shortDesc, desc)
	c.registerCommand(command, namespace)
	return namespace
}

func (c *Commands) registerCommand(command Command, namespace *NamespaceCommand) {
	c.commands[command.CommandLabel()] = namespace
	for _, alias := range command.CommandAliases() {
		c.commands[alias] = namespace
	}
}
