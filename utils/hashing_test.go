package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashAndSalt(t *testing.T) {
	testPwd := []byte("Test1234")
	testPwd2 := []byte("Test1234jnsbyhanshr hshbr  shhfbfbskjje")
	testPwd3 := []byte("Test1234")

	//Test 1: Hashing works and does not result in an error and
	//returns a non empty hash password
	response, err := HashAndSalt(testPwd)
	assert.NoError(t, err)
	assert.NotEqual(t, response, "")

	//Test 2: Hashing works and does not result in an error and
	//returns a non empty hash password
	response, err = HashAndSalt(testPwd2)
	assert.NoError(t, err)
	assert.NotEqual(t, response, "")

	//Test 3: Hashing 2 of the same password but ensuring the hash
	//returned is different
	response1, err := HashAndSalt(testPwd)
	assert.NoError(t, err)
	response2, err := HashAndSalt(testPwd3)
	assert.NoError(t, err)
	assert.NotEqual(t, response1, response2)

}

func TestComparePasswords(t *testing.T) {
	testPwd := "Test1234"
	testPwd2 := "ehjjcbasdjvbkesndklsncklns"
	testPwd3 := "test1234"
	hashedPwd, _ := HashAndSalt([]byte(testPwd))

	//Test 1: Passwords are the same and output is true
	response := ComparePasswords(hashedPwd, testPwd)
	assert.Equal(t, response, true)

	//Test 2: Passwords are not the same and output is false
	response = ComparePasswords(hashedPwd, testPwd2)
	assert.Equal(t, response, false)

	//Test 3: Passwords are not the same but are similar and output is false
	response = ComparePasswords(hashedPwd, testPwd3)
	assert.Equal(t, response, false)
}
