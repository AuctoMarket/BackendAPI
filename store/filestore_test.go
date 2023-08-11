package store

import (
	"BackendAPI/utils"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewS3(t *testing.T) {
	os.Clearenv()
	//Test 1: No .env file
	s3, err := CreateNewS3()
	assert.Error(t, err)
	assert.Empty(t, s3)

	//Load wrong .env file for tests
	loadErr := utils.LoadDotEnv("../test.env")
	assert.NoError(t, loadErr)

	//Test 2: Wrong .env file
	s3, err = CreateNewS3()
	assert.Error(t, err)
	assert.Empty(t, s3)

	//Load correct .env file for tests
	loadErr = utils.LoadDotEnv("../.env")
	assert.NoError(t, loadErr)

	//Test 3: Correct .env file
	s3, err = CreateNewS3()
	assert.NoError(t, err)
	assert.NotEmpty(t, s3)
}

func TestUploadImages(t *testing.T) {
	utils.LoadDotEnv("../.env")
	s3Client, s3Error := CreateNewS3()
	assert.NoError(t, s3Error)

	os.Clearenv()

	//Test 1: Bad env variables
	var files []io.Reader
	buf := bytes.NewBufferString("hello\n")
	files = append(files, buf)
	keys := []string{"test"}

	err := UploadImages(s3Client, keys, files)
	assert.Error(t, err)
	assert.Equal(t, "Error in loading environment variables, Bucket name does not exist:", err.Error())

	//Test 2: Upload successful
	utils.LoadDotEnv("../.env")
	err = UploadImages(s3Client, keys, files)
	assert.NoError(t, err)
}
