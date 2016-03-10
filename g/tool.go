package g

import (
	"github.com/toolkits/file"
//	"log"
    "fmt"
)

func configExists(cfg string) bool {
    if !file.IsExist(cfg) {
        return false
    }
	return true
}

func ConfigArgs(cfg string) []string {
    if !configExists(cfg) {
        fmt.Println("expect config file:", cfg)
        return nil
    }
    return []string {"-c", cfg}
}
