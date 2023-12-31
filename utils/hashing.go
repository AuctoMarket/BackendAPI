package utils

import (
	"log"

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
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlaintext := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlaintext)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
