package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDotEnv(t *testing.T) {
	testEnvPath := "/Users/ekam/Desktop/AuctoCode/BackendAPI/test.env"
	testEnvPathFake := "/Users/ekam/Desktop/AuctoCode/BackendAPI/fake.env"

	//Test 1: Get .env variable that exists
	res, err := GetDotEnv("TEST_ENV", testEnvPath)
	assert.NoError(t, err)
	assert.Equal(t, "test", res)

	//Test 2: Get .env variable that does not exist
	res, err = GetDotEnv("TEST_ENV_FAKE", testEnvPath)
	assert.Empty(t, res)

	//Test 3: Path to .env is incorrect
	res, err = GetDotEnv("TEST_ENV_FAKE", testEnvPathFake)
	assert.Error(t, err)
}

func TestGetDotEnvInt(t *testing.T) {
	testEnvPath := "/Users/ekam/Desktop/AuctoCode/BackendAPI/test.env"
	testEnvPathFake := "/Users/ekam/Desktop/AuctoCode/BackendAPI/fake.env"

	//Test 1: Get .env variable that exists
	res, err := GetDotEnvInt("TEST_ENV_INT", testEnvPath)
	assert.NoError(t, err)
	assert.Equal(t, 1, res)

	//Test 2: Get .env variable that does not exist
	res, err = GetDotEnvInt("TEST_ENV_FAKE", testEnvPath)
	assert.Empty(t, res)

	//Test 3: Path to .env is incorrect
	res, err = GetDotEnvInt("TEST_ENV", testEnvPathFake)
	assert.Error(t, err)
}
