package utils

import (
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"errors"
>>>>>>> 005bc68 (Add login and signup API)
=======
>>>>>>> c9cf9e0 (Update Error Handling)
=======
>>>>>>> 3a877cd (Update Error Handling)
	"log"
)

/*
Logs internal server errors and returns sanitized error response
*/
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
=======
func LogInternalServerError(msg string, err error) *ErrorHandler {
	log.Println(msg)
	log.Println(err)
	return InternalServerError()
>>>>>>> c9cf9e0 (Update Error Handling)
=======
func LogError(err error, msg string) {
	log.Println(msg)
	log.Println(err)
>>>>>>> a806d09 (Update logging)
=======
func LogInternalServerError(msg string, err error) *ErrorHandler {
	log.Println(msg)
	log.Println(err)
	return InternalServerError()
>>>>>>> 3a877cd (Update Error Handling)
}
