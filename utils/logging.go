package utils

import (
	"log"
)

/*
Logs internal server errors and returns sanitized error response
*/
func LogError(err error, msg string) {
	log.Println(msg)
	log.Println(err)
}
