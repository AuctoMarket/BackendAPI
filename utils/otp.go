package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GetOtp(length int) string {
	rand.Seed(time.Now().UnixNano())
	var otp string
	for i := 0; i < length; i++ {
		// generate a random integer between 0 and 9
		randomInt := strconv.Itoa(rand.Intn(10))
		otp += randomInt
	}

	return otp
}
