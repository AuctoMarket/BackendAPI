package utils

import (
	"errors"
	"log"
)

/*
Logs internal server errors and returns sanitized error response
*/
func LogError(msg string, err error) error {
	log.Println(msg)
	log.Println(err)
	return errors.New("Something went wrong")
}
