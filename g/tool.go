package g

import (
	"fmt"
	"github.com/toolkits/file"
	"os/exec"
	"strings"
)

func configExists(cfg string) bool {
	if !file.IsExist(cfg) {
		return false
	}
	return true
}

func GetConfFileArgs(cfg string) ([]string, error) {
	if !configExists(cfg) {
		return nil, fmt.Errorf("expect config file: %s\n", cfg)
	}
	return []string{"-c", cfg}, nil
}

func ModuleExists(name string) error {
	if Modules[name] {
		return nil
	}
	return fmt.Errorf("This module doesn't exist: %s", name)
}

func CheckModulePid(name string) (string, error) {
	output, err := exec.Command("pgrep", name).Output()
	if err != nil {
		return "", err
	}
	pidStr := strings.TrimSpace(string(output))
	return pidStr, nil
}

func GetAllModuleArgs() []string {
	var allModules []string
	for _, name := range Order {
		if Modules[name] {
			allModules = append(allModules, name)
		}
	}
	return allModules
}

func CheckModuleStatus(name string) int {
	err := ModuleExists(name)
	if err != nil {
		fmt.Println(err)
		return ModuleNonexistent
	}

	fmt.Print("Checking status [", ModuleApps[name], "]...")

	pidStr, err := CheckModulePid(ModuleApps[name])
	if err != nil {
		fmt.Println("not running!!")
		return ModuleExistentNotRunning
	}

	fmt.Println("running with PID [", pidStr, "]!!")
	return ModuleRunning
}
