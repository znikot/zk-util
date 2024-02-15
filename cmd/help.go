package cmd

import "fmt"

func init() {
	AddCommand(&helpCommand{})
}

type helpCommand struct{}

func (c *helpCommand) Name() string {
	return "help"
}

func (c *helpCommand) Description() string {
	return "show command help infomation."
}

func (c *helpCommand) Usage() {
	fmt.Printf("usage: %s help <command>\n", exeName)
}

func (c *helpCommand) Exec(args ...string) error {
	if len(args) < 1 {
		c.Usage()
		return nil
	}

	// get command
	command, ok := GetCommand(args[0])
	if !ok {
		return fmt.Errorf("invalid command %s\n", args[0])
	}
	command.Usage()
	return nil
}
