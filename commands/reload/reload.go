package reload

import (
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
	return g.Command_EX_OK
}

func (c *Command) Synopsis() string {
	return "Reload an Open-Falcon module's configuration file"
}

func (c *Command) Help() string {
	helpText := `
Usage: open-falcon reload [Module]

  Reload the configuration file of the specified Open-Falcon module.
  A module represents a single node in a cluster.

Modules:

  `
	return strings.TrimSpace(helpText)
}
