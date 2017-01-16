package rpc

import (
	"fmt"
	log "github.com/Cepave/open-falcon-backend/common/logruslog"
)

var logger = log.NewDefaultLogger("WARN")

// Builds defer function for handling panic,
// and sets the value of error into error object
func HandleError(err *error) func() {
	return func() {
		p := recover()
		if p != nil {
			logger.Errorf("Panic in RPC(GoLang): %v", p)

			switch errObject := p.(type) {
			case error:
				*err = errObject
			default:
				newErr := fmt.Errorf("Has error on RPC: %v", errObject)
				*err = newErr
			}
		}
	}
}
