package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOtp(t *testing.T) {
	res := GetOtp(6)
	assert.NotEmpty(t, res)
	assert.Equal(t, 6, len(res))

	res2 := GetOtp(10)
	assert.NotEmpty(t, res2)
	assert.Equal(t, 10, len(res2))
}
