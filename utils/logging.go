package utils

import (
	"log"
)

/*
Logs internal server errors
*/
func LogError(err error, msg string) {
	log.Println(msg)
	log.Println(err)
}

/*
Logs custom error messages
*/
func LogMessage(msg string) {
	log.Println(msg)
}
