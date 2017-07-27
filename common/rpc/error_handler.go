package rpc

import (
	"fmt"

	log "github.com/Cepave/open-falcon-backend/common/logruslog"
	or "github.com/Cepave/open-falcon-backend/common/runtime"

	"github.com/Cepave/open-falcon-backend/common/utils"
)

var logger = log.NewDefaultLogger("WARN")

// Builds defer function for handling panic,
// and sets the value of error into error object
func HandleError(err *error) func() {
	return utils.PanicToError(
		err,
		func(p interface{}) error {
			stack := or.GetCallerInfoStack(2, 16).ConcatStringStack(" <- ")

			logger.Errorf("Panic in RPC(GoLang): %v", p)
			logger.Errorf("Panic Stack: %s", stack)

			if errObject, ok := p.(error); ok {
				return errObject
			}

			return fmt.Errorf("Has error on RPC: %v. Stack: %s", p, stack)
		},
	)
}
