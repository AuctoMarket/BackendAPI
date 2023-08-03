package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDotEnvInt(t *testing.T) {
	testEnvPath := "/Users/ekam/Desktop/AuctoCode/BackendAPI/test.env"

	err := LoadDotEnv(testEnvPath)
	assert.NoError(t, err)

	//Test 1: Get .env variable that exists
	res, err := GetDotEnvInt("TEST_ENV_INT")
	assert.NoError(t, err)
	assert.Equal(t, 1, res)

	//Test 2: Get .env variable that does not exist
	res, err = GetDotEnvInt("TEST_ENV_FAKE")
	assert.Empty(t, res)

}
