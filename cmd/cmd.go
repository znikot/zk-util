package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"sort"
	"strings"
)

var exeName string

func init() {
	// get the executable name
	exeName = path.Base(os.Args[0])
}

// Command interface. example code:
//
//		package test
//		import (
//			"fmt"
//			"flag"
//			"github.com/znikot/zk-util/cmd"
//		)
//
//		func init() *testCommand{
//	 		cmd.AddCommand(new(testCommand).init())
//		}
//
//		type testCommand struct{
//			Name string
//			fs *flag.FlagSet
//		}
//
//		func (c *testCommand) init() *testCommand {
//			c.fs = flag.NewFlagSet(c.Name(), flag.ContinueOnError)
//			return c
//		}
//
//		func (c *testCommand) Name() string {
//			return "test"
//		}
//
//		func (c *testCommand) Description() string {
//			return "test command"
//		}
//
//		func (c *testCommand) Usage(){
//			c.fs.Usage()
//		}
//
//		func (c *testCommand) Exec(args ...string){
//			// check arguments and print usage to the console
//			if len(args) == 0 {
//				c.Usage()
//				return
//			}
//
//			// parse flags from args
//			c.fs.Parse(args)
//			fmt.Printf("name: %s\n",c.Name)
//		}
type Command interface {
	// name of command
	Name() string
	// print usage of command
	Usage()
	// execute command
	Exec(args ...string) error
	// description of command
	Description() string
}

var commands = map[string]Command{}

func AddCommand(command Command) {
	commands[command.Name()] = command
}

// get all commands, sorted by name
func AllCommands() []Command {
	cmds := make([]Command, 0)

	for _, v := range commands {
		cmds = append(cmds, v)
	}

	// sort
	sort.SliceStable(cmds, func(i, j int) bool {
		return strings.Compare(cmds[i].Name(), cmds[j].Name()) < 0
	})

	return cmds
}

// get command by the name
func GetCommand(name string) (Command, bool) {
	command, ok := commands[name]
	return command, ok
}

// run command with arguments
func Run(args ...string) error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	// get command
	command, ok := GetCommand(os.Args[1])
	if !ok {
		return fmt.Errorf("invalid command %s\n", os.Args[1])
	}

	// exec command
	return command.Exec(os.Args[2:]...)
}

// invoke this function will block until os signal received
//
// Example:
//
//	func main(){
//		cmd.WaitSignal(func(sig os.Signal) {
//			// do something while signal received
//
//			os.Exit(0)
//
//		}, syscall.SIGINT, syscall.SIGTERM)
//	}
func WaitSignal(handler func(s os.Signal), signals ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, signals...)

	go func() {
		sig := <-sigs
		if handler != nil {
			handler(sig)
		}
		done <- true
	}()
	<-done

	close(sigs)
	close(done)
}

// print usage
func printUsage() {
	fmt.Printf("usage: %s <command> [options...]\n", exeName)
	fmt.Println()
	fmt.Println("commands:")
	for _, k := range AllCommands() {
		fmt.Printf("\t%s%s\n", fixLen(k.Name(), 20), k.Description())
	}
}

// fix string len
func fixLen(str string, l int) string {
	if len(str) >= l {
		return str[:l]
	}
	left := l - len(str)
	for i := 0; i < left; i++ {
		str += " "
	}
	return str
}
