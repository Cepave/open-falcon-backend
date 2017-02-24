package db

import (
	"fmt"
	"path"
	"runtime"
)

// Defines the type of database error
type DbError struct {
	errorMessage string
}

// Implements interface error
func (dbError *DbError) Error() string {
	return dbError.errorMessage
}

// Panic with database error if the error is vialbe
func PanicIfError(err error) {
	if err == nil {
		return
	}

	callerPc, file, line, _ := runtime.Caller(1)
	callerFunc := runtime.FuncForPC(callerPc)

	finalFileName := path.Base(file)

	panic(NewDatabaseError(
		fmt.Errorf(
			"[RDB Error] %s @ File:%s (line:%d)\n\t%s",
			callerFunc.Name(), finalFileName, line,
			err,
		),
	))
}

// Constructs an error of database
func NewDatabaseError(err error) *DbError {
	return &DbError { err.Error() }
}
