package utils

import (
	"log"
)

/*
Logs internal server errors and returns sanitized error response
*/
func LogInternalServerError(msg string, err error) *ErrorHandler {
	log.Println(msg)
	log.Println(err)
	return InternalServerError()
}
