package product

import (
	"BackendAPI/store"
	"BackendAPI/utils"
	"bytes"
	"context"
	"database/sql"
	"io"
	"os"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestDoProductImagesExist(t *testing.T) {
	utils.LoadDotEnv("../../.env")
	db, startupErr := store.SetupTestDB("../../.env")
	assert.NoError(t, startupErr)

	sellerid, sellerErr := createDummySeller(db)
	assert.NoError(t, sellerErr)
	productIds, productErr := createDummyProducts(db, sellerid)
	assert.NoError(t, productErr)

	//Test 1: No product Images
	res := doProductImagesExist(db, productIds[0])
	assert.Equal(t, false, res)

	//Test 2: No product Images
	res = doProductImagesExist(db, productIds[1])
	assert.Equal(t, false, res)

	//Test 3: No product Images
	createDummyProductImages(db, []string{productIds[2]})
	res = doProductImagesExist(db, productIds[2])
	assert.Equal(t, true, res)

	store.CloseDB(db)
}

func TestCreateProductImages(t *testing.T) {
	utils.LoadDotEnv("../../.env")
	db, startupErr := store.SetupTestDB("../../.env")
	assert.NoError(t, startupErr)
	s3Client, s3Error := store.CreateNewS3()
	assert.NoError(t, s3Error)

	sellerid, sellerErr := createDummySeller(db)
	assert.NoError(t, sellerErr)
	productIds, productErr := createDummyProducts(db, sellerid)
	assert.NoError(t, productErr)

	//Test 1: Creating image successfully
	var files []io.Reader
	buf := bytes.NewBufferString("hello\n")
	files = append(files, buf)

	res, err := CreateProductImages(db, s3Client, productIds[0], files)
	assert.Empty(t, err)
	assert.NotEmpty(t, res)

	//Test 2: Creating image successfully
	var files2 []io.Reader
	buf2 := bytes.NewBufferString("hello\n")
	files2 = append(files2, buf2)
	res, err = CreateProductImages(db, s3Client, productIds[1], files2)
	assert.Empty(t, err)
	assert.NotEmpty(t, res)

	//Test 3: Incorrect Product Id
	res, err = CreateProductImages(db, s3Client, "wrong id", files)
	assert.Error(t, err)
	assert.Equal(t, "Product with given id does not exist", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 4: Product already has images
	res, err = CreateProductImages(db, s3Client, productIds[0], files)
	assert.Error(t, err)
	assert.Equal(t, "Product with given id already has images", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 5: No files to submit
	res, err = CreateProductImages(db, s3Client, productIds[2], nil)
	assert.Error(t, err)
	assert.Equal(t, "No images attached, at least 1 image per post", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 6: To many files to submit
	var files3 []io.Reader
	buf31 := bytes.NewBufferString("hello\n")
	buf32 := bytes.NewBufferString("hello\n")
	buf33 := bytes.NewBufferString("hello\n")
	buf34 := bytes.NewBufferString("hello\n")
	buf35 := bytes.NewBufferString("hello\n")
	buf36 := bytes.NewBufferString("hello\n")
	files3 = append(files3, buf31)
	files3 = append(files3, buf32)
	files3 = append(files3, buf33)
	files3 = append(files3, buf34)
	files3 = append(files3, buf35)
	files3 = append(files3, buf36)
	res, err = CreateProductImages(db, s3Client, productIds[2], files3)
	assert.Error(t, err)
	assert.Equal(t, "Too many images uploaded, at most 5 images per post", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	store.CloseDB(db)
}

func TestValidateCreateProductImages(t *testing.T) {
	utils.LoadDotEnv("../../.env")
	db, startupErr := store.SetupTestDB("../../.env")
	assert.NoError(t, startupErr)

	sellerid, sellerErr := createDummySeller(db)
	assert.NoError(t, sellerErr)
	productIds, productErr := createDummyProducts(db, sellerid)
	assert.NoError(t, productErr)

	createDummyProductImages(db, []string{productIds[2]})

	//Test 1: No error
	var files []io.Reader
	buf := bytes.NewBufferString("hello\n")
	files = append(files, buf)
	testErr := validateCreateProductImages(db, productIds[0], files)
	assert.Empty(t, testErr)

	//Test 2: No Images
	var files2 []io.Reader = nil
	testErr = validateCreateProductImages(db, productIds[1], files2)
	assert.NotEmpty(t, testErr)
	assert.Equal(t, "No images attached, at least 1 image per post", testErr.Error())
	assert.Equal(t, 400, testErr.ErrorCode())

	//Test 3: To Many Images
	var files3 []io.Reader
	buf31 := bytes.NewBufferString("hello\n")
	buf32 := bytes.NewBufferString("hello\n")
	buf33 := bytes.NewBufferString("hello\n")
	buf34 := bytes.NewBufferString("hello\n")
	buf35 := bytes.NewBufferString("hello\n")
	buf36 := bytes.NewBufferString("hello\n")
	files3 = append(files3, buf31)
	files3 = append(files3, buf32)
	files3 = append(files3, buf33)
	files3 = append(files3, buf34)
	files3 = append(files3, buf35)
	files3 = append(files3, buf36)
	testErr = validateCreateProductImages(db, productIds[1], files3)
	assert.NotEmpty(t, testErr)
	assert.Equal(t, "Too many images uploaded, at most 5 images per post", testErr.Error())
	assert.Equal(t, 400, testErr.ErrorCode())

	//Test 4: Product Id does not exist
	var files4 []io.Reader
	buf = bytes.NewBufferString("hello\n")
	files4 = append(files4, buf)
	testErr = validateCreateProductImages(db, "wrong id", files4)
	assert.NotEmpty(t, testErr)
	assert.Equal(t, "Product with given id does not exist", testErr.Error())
	assert.Equal(t, 400, testErr.ErrorCode())

	//Test 5: Product Id already has images
	var files5 []io.Reader
	buf = bytes.NewBufferString("hello\n")
	files5 = append(files5, buf)
	createDummyProductImages(db, []string{productIds[0]})
	testErr = validateCreateProductImages(db, productIds[0], files5)
	assert.NotEmpty(t, testErr)
	assert.Equal(t, "Product with given id already has images", testErr.Error())
	assert.Equal(t, 400, testErr.ErrorCode())

	store.CloseDB(db)
}

func TestMakeImagePath(t *testing.T) {
	//Test 1: No env variables decalred
	os.Clearenv()
	_, err := makeImagePath("Test")
	assert.NotEmpty(t, err)

	utils.LoadDotEnv("../../.env")

	//Test 2: Environment variables present and image path is local
	res, err := makeImagePath("Test")
	assert.Empty(t, err)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com/products/images/Test", res)

	//Test 3: Empty String
	res, err = makeImagePath("")
	assert.NotEmpty(t, err)
	assert.Equal(t, 500, err.ErrorCode())

	os.Clearenv()
	utils.LoadDotEnv("../../.env")
}

func createDummyProductImage(db *sql.DB, productId string, imageNo int) (string, error) {
	var productImageId string

	utils.LoadDotEnv("../../.env")
	api_env, _ := os.LookupEnv("API_ENV")

	query2 := `INSERT INTO product_images (product_id, image_no) VALUES ($1, $2) RETURNING product_image_id;`
	err := db.QueryRowContext(context.Background(), query2, productId, strconv.Itoa(imageNo)).Scan(&productImageId)

	if err != nil {
		return productImageId, err
	}

	if api_env == "local" {
		productImageId = os.Getenv("S3_LOCAL_URL") + "/products/images/" + productImageId
	} else {
		productImageId = os.Getenv("S3_DEV_URL") + "/products/images/" + productImageId
	}

	return productImageId, nil
}

func createDummyProductImages(db *sql.DB, productIds []string) ([]string, error) {
	var productImageIds []string
	for i := 0; i < len(productIds); i++ {
		imageId, err := createDummyProductImage(db, productIds[i], i)
		productImageIds = append(productImageIds, imageId)

		if err != nil {
			return nil, err
		}
	}

	return productImageIds, nil
}
