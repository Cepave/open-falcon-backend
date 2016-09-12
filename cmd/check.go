package cmd

import (
	"fmt"
	"strings"

	"github.com/Cepave/open-falcon-backend/g"
	"github.com/spf13/cobra"
)

var Check = &cobra.Command{
	Use:   "check [Module ...]",
	Short: "Check the status of Open-Falcon modules",
	Long: `
Check if the specified Open-Falcon modules are running.

Modules:

  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: check,
}

func check(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return c.Usage()
	}

	if len(args) == 1 && args[0] == "all" {
		for moduleName := range g.Modules {
			if g.IsRunning(moduleName) {
				fmt.Printf("%20s %10s %15s \n", g.ModuleApps[moduleName], "UP", g.Pid(moduleName))
			} else {
				fmt.Printf("%20s %10s %15s \n", g.ModuleApps[moduleName], "DOWN", "-")
			}
		}
	} else {
		for _, moduleName := range args {
			if !g.HasModule(moduleName) {
				return fmt.Errorf("%s doesn't exist\n", moduleName)
			}

			if g.IsRunning(moduleName) {
				fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
			} else {
				fmt.Print("[", g.ModuleApps[moduleName], "] down\n")
			}
		}
	}
	return nil
}
