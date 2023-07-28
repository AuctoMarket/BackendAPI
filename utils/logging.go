package utils

import (
<<<<<<< HEAD
=======
	"errors"
>>>>>>> 005bc68 (Add login and signup API)
	"log"
)

/*
Logs internal server errors and returns sanitized error response
*/
<<<<<<< HEAD
func LogError(err error, msg string) {
	log.Println(msg)
	log.Println(err)
=======
func LogError(msg string, err error) error {
	log.Println(msg)
	log.Println(err)
	return errors.New("Something went wrong")
>>>>>>> 005bc68 (Add login and signup API)
}
