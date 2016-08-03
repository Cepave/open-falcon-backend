package tail

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Cepave/open-falcon-backend/g"
	"github.com/mitchellh/cli"
)

// Command is a Command implementation that runs a Consul agent.
// The command will not end unless a shutdown message is sent on the
// ShutdownCh. If two messages are sent on the ShutdownCh it will forcibly
// exit.
type Command struct {
	Revision          string
	Version           string
	VersionPrerelease string
	Ui                cli.Ui
}

func (c *Command) Run(args []string) int {
	if len(args) != 1 {
		return cli.RunResultHelp
	}
	moduleName := args[0]
	err := g.ModuleExists(moduleName)
	if err != nil {
		fmt.Println(err)
		fmt.Println("** start failed **")
		return g.Command_EX_ERR
	}

	logPath := "./" + moduleName + "/" + g.LogDir + "/" + moduleName + ".log"
	cmd := exec.Command("tail", "-f", logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	dir, _ := os.Getwd()
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		fmt.Println("** tail failed **")
		return g.Command_EX_ERR
	}
	return 0
}

func (c *Command) Synopsis() string {
	return "Display an Open-Falcon module's log"
}

func (c *Command) Help() string {
	helpText := `
Usage: open-falcon tail [Module]

  Display the log of the specified Open-Falcon module.
  A module represents a single node in a cluster.

Modules:

  ` + strings.Join(g.AllModulesInOrder, " ")
	return strings.TrimSpace(helpText)
}
