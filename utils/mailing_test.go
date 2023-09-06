package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendOtpMail(t *testing.T) {
	testEmail := "test@aucto.io"
	testOtp := "000000"

	err := LoadDotEnv("../.env")
	assert.Empty(t, err)

	err = SendOtpMail(testEmail, testOtp)
	assert.Empty(t, err)
}
