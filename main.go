package main

import (
	"fmt"
	"github.com/Cepave/open-falcon/commands/src/agent"
	"github.com/Cepave/open-falcon/commands/src/aggregator"
	"github.com/Cepave/open-falcon/commands/src/api"
	"github.com/Cepave/open-falcon/commands/src/graph"
	"github.com/Cepave/open-falcon/commands/src/hbs"
	"github.com/Cepave/open-falcon/commands/src/judge"
	"github.com/Cepave/open-falcon/commands/src/nodata"
	"github.com/Cepave/open-falcon/commands/src/query"
	"github.com/Cepave/open-falcon/commands/src/sender"
	"github.com/Cepave/open-falcon/commands/src/task"
	"github.com/Cepave/open-falcon/commands/src/transfer"
	"github.com/mitchellh/cli"
	"io/ioutil"
	"log"
	"os"
)

// Commands is the mapping of all the available Consul commands.
var Commands map[string]cli.CommandFactory

func init() {
	//ui := &cli.BasicUi{Writer: os.Stdout}

	Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &agent.Command{}, nil
		},

		"aggregator": func() (cli.Command, error) {
			return &aggregator.Command{}, nil
		},

		"api": func() (cli.Command, error) {
			return &api.Command{}, nil
		},

		"graph": func() (cli.Command, error) {
			return &graph.Command{}, nil
		},

		"hbs": func() (cli.Command, error) {
			return &hbs.Command{}, nil
		},

		"judge": func() (cli.Command, error) {
			return &judge.Command{}, nil
		},

		"nodata": func() (cli.Command, error) {
			return &nodata.Command{}, nil
		},

		"query": func() (cli.Command, error) {
			return &query.Command{}, nil
		},

		"sender": func() (cli.Command, error) {
			return &sender.Command{}, nil
		},

		"task": func() (cli.Command, error) {
			return &task.Command{}, nil
		},

		"transfer": func() (cli.Command, error) {
			return &transfer.Command{}, nil
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
