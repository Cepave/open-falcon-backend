package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Cepave/open-falcon-backend/commands/reload"
	"github.com/Cepave/open-falcon-backend/commands/restart"
	"github.com/Cepave/open-falcon-backend/commands/start"
	"github.com/Cepave/open-falcon-backend/commands/status"
	"github.com/Cepave/open-falcon-backend/commands/stop"
	"github.com/Cepave/open-falcon-backend/commands/tail"
	"github.com/mitchellh/cli"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory
var startCommand cli.Command

func init() {
	//ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"start": func() (cli.Command, error) {
			return &start.Command{}, nil
		},
		"stop": func() (cli.Command, error) {
			return &stop.Command{}, nil
		},
		"restart": func() (cli.Command, error) {
			return &restart.Command{}, nil
		},
		"status": func() (cli.Command, error) {
			return &status.Command{}, nil
		},
		"tail": func() (cli.Command, error) {
			return &tail.Command{}, nil
		},
		"reload": func() (cli.Command, error) {
			return &reload.Command{}, nil
		},
	}
}

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.SetOutput(ioutil.Discard)

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		//	HelpFunc: cli.BasicHelpFunc("consul"),
	}
	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
