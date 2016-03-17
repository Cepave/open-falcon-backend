package restart

import (
	"github.com/Cepave/open-falcon/commands/start"
	"github.com/Cepave/open-falcon/commands/stop"
	"github.com/Cepave/open-falcon/g"
	"github.com/mitchellh/cli"
	"strings"
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
	if len(args) == 0 {
		return cli.RunResultHelp
	}
	var stopCmd stop.Command
	var startCmd start.Command
	stopCmd.Run(args)
	startCmd.Run(args)
	return 0
}

func (c *Command) Synopsis() string {
	return "Restart Open-Falcon modules"
}

func (c *Command) Help() string {
	helpText := `
Usage: open-falcon restart [Module ...]

  Restart the specified Open-Falcon modules and run until a stop command is received.
  A module represents a single node in a cluster.

Modules:

  ` + "all " + strings.Join(g.GetAllModuleArgs(), " ")
	return strings.TrimSpace(helpText)
}
