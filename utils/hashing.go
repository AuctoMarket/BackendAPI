package utils

import (
<<<<<<< HEAD
<<<<<<< HEAD
	"log"
=======
	"fmt"
>>>>>>> 005bc68 (Add login and signup API)
=======
	"log"
>>>>>>> e5d2750 (Add Tests for Login/Signup)

	"golang.org/x/crypto/bcrypt"
)

/*
Uses the bcrypt encryption algorithm to hash the password
Salt is stored as part of the result of the bcrypt algorithm
MinCost is the minimum cost of running the algorithm and is
a constant found in the bcrypt library
*/
func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

/*
Compare passwords takes in the string hashed pwd and converts
it to a []byte and compares the plaintext and the ciphertext to
see wether they are equivalent.
*/
<<<<<<< HEAD
<<<<<<< HEAD
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlaintext := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlaintext)
	if err != nil {
		log.Println(err)
=======
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
=======
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlaintext := []byte(plainPwd)
>>>>>>> e5d2750 (Add Tests for Login/Signup)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlaintext)
	if err != nil {
<<<<<<< HEAD
		fmt.Println(err)
>>>>>>> 005bc68 (Add login and signup API)
=======
		log.Println(err)
>>>>>>> e5d2750 (Add Tests for Login/Signup)
		return false
	}

	return true
}
