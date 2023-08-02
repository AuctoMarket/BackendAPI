package utils

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

/*
Function to get environment variables from the .Env file if they are a
string
*/
func GetDotEnv(key string, path string) (string, error) {
	err := godotenv.Load(path)

	if os.Getenv("API_ENV") == "lambda" {
		return os.Getenv(key), nil
	}

	if err != nil {
		return "", err
	}

	return os.Getenv(key), nil
}

/*
Function to get environment variables from the .Env file if they are a
int
*/
func GetDotEnvInt(key string, path string) (int, error) {
	err := godotenv.Load(path)

	if os.Getenv("API_ENV") == "lambda" {
		num, err := strconv.ParseInt(os.Getenv(key), 10, 32)
		return int(num), err
	}

	if err != nil {
		return 0, err
	}

	num, err := strconv.ParseInt(os.Getenv(key), 10, 32)

	if err != nil {
		return 0, err
	}

	return int(num), nil
}
