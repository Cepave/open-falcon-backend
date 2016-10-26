package db

import (
	"fmt"
)

// Defines the type of database error
type DbError struct {
	errorMessage string
}

// Implements interface error
func (dbError *DbError) Error() string {
	return dbError.errorMessage
}

// Constructs a error of database
func NewDatabaseError(err error) *DbError {
	return &DbError {
		fmt.Sprintf("Database Error. %v", err),
	}
}
