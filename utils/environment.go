package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

/*
Function to load a .env file into the OS environment
*/
func LoadDotEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		wd, _ := os.Getwd()
		LogError(err, fmt.Sprintf("Could not load env file stated in path, PWD:%s", wd))
	}

	return err
}

/*
Function to get environment variables from the .Env file if they are a
int
*/
func GetDotEnvInt(key string) (int, error) {
	num, err := strconv.ParseInt(os.Getenv(key), 10, 32)

	if err != nil {
		return 0, err
	}

	return int(num), nil
}
