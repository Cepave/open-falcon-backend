package rpc

import (
	"fmt"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"

	"github.com/Cepave/open-falcon-backend/common/utils"
)

var logger = log.NewDefaultLogger("WARN")

// Builds defer function for handling panic,
// and sets the value of error into error object
func HandleError(err *error) func() {
	return utils.PanicToError(
		err,
		func(p interface{}) error {
			logger.Errorf("Panic in RPC(GoLang): %v", p)

			if errObject, ok := p.(error); ok {
				return errObject
			}

			return fmt.Errorf("Has error on RPC: %v", p)
		},
	)
}
