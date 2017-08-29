package http

import (
	"flag"
	"fmt"
	"github.com/Cepave/open-falcon-backend/common/logruslog"
)

var logger = logruslog.NewDefaultLogger("INFO")

var webHost = flag.String("test.web_host", "0.0.0.0", "Listening Host(0.0.0.0)")
var webPort = flag.Uint("test.web_port", 0, "Listening port of web")

func GetWebHost() string {
	if *webHost == "0.0.0.0" {
		return "127.0.0.1"
	}

	return *webHost
}
func GetWebPort() uint16 {
	return uint16(*webPort)
}
func GetWebUrl() string {
	return fmt.Sprintf("http://%s:%d", GetWebHost(), GetWebPort())
}
