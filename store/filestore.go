package store

import (
	"BackendAPI/utils"
	"context"
	"errors"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
Create a S3 client that can be used to interact with AWS s3
*/
func CreateNewS3() (*s3.Client, error) {

	var (
		region, hasRegion = os.LookupEnv("AWS_REGION")
		secret, hasSecret = os.LookupEnv("AWS_SECRET")
		key, hasKey       = os.LookupEnv("AWS_KEY")
		token             = ""
	)

	if !(hasKey && hasRegion && hasSecret) {
		return nil, errors.New("Error in loading environment variables for S3")
	}

	creds := credentials.NewStaticCredentialsProvider(key, secret, token)
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(creds), config.WithRegion(region))

	if err != nil {
		return nil, errors.New("There was an error loading config files:" + err.Error())
	}

	awsS3Client := s3.NewFromConfig(cfg)

	return awsS3Client, nil
}

/*
Upload a list of images to the S3 Bucket specified by the environment variables
*/
func UploadImages(client *s3.Client, keys []string, files []io.Reader) error {
	var bucket string
	var hasBucket bool

	bucket, hasBucket = os.LookupEnv("S3_BUCKET_NAME")

	fileType := "image/png"

	if !hasBucket {
		return errors.New("Error in loading environment variables, Bucket name does not exist:")
	}

	for i := 0; i < len(keys); i++ {
		image := files[i]
		err := uploadFile(client, keys[i], image, bucket, fileType)

		if err != nil {
			return err
		}
	}

	return nil
}

/*
Upload a an individiual image to the S3 Bucket specified in the bucket argument
*/
func uploadFile(client *s3.Client, key string, file io.Reader, bucket string, fileType string) error {
	uploader := manager.NewUploader(client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileType),
	})

	if err != nil {
		utils.LogError(err, "Error uploading to S3 Bucket")
	}

	return err
}
