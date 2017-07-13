package db

import (
	"github.com/Cepave/open-falcon-backend/common/utils"
)

// Defines the type of database error
type DbError struct {
	*utils.StackError
}

// Panic with database error if the error is vialbe
func PanicIfError(err error) {
	if !utils.IsViable(err) {
		return
	}

	panic(NewDatabaseError(err))
}

// Constructs an error of database
func NewDatabaseError(err error) *DbError {
	stackError, ok := err.(*utils.StackError)
	if ok {
		return &DbError{stackError}
	}

	return &DbError{utils.BuildErrorWithCaller(err)}
}
