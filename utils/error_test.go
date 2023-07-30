package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadRequestError(t *testing.T) {
	//Test 1: Bad Request Error created with error code 400 and message
	testBadRequestError1 := BadRequestError("Test Error 1")
	assert.Equal(t, "Test Error 1", testBadRequestError1.Error())
	assert.Equal(t, 400, testBadRequestError1.Code)

	//Test 1: Bad Request Error created with error code 400 and standard message
	testBadRequestError2 := BadRequestError("")
	assert.Equal(t, "Bad request", testBadRequestError2.Error())
	assert.Equal(t, 400, testBadRequestError2.Code)
}

func TestInternalServerError(t *testing.T) {
	//Test 1: Internal Server Error created with error code 500 and standard message
	testInternalServerError1 := InternalServerError()
	assert.Equal(t, "Something went wrong", testInternalServerError1.Error())
	assert.Equal(t, 500, testInternalServerError1.Code)
}

func TestUnauthorizedError(t *testing.T) {
	//Test 1: Unauthorized Error created with error code 401 and message
	testUnauthorizedError1 := UnauthorizedError("Test Error 1")
	assert.Equal(t, "Test Error 1", testUnauthorizedError1.Error())
	assert.Equal(t, 401, testUnauthorizedError1.Code)

	//Test 1: Unauthorized Error created with error code 401 and standard message
	testUnauthorizedError2 := UnauthorizedError("")
	assert.Equal(t, "Unauthorized user", testUnauthorizedError2.Error())
	assert.Equal(t, 401, testUnauthorizedError1.Code)
}
