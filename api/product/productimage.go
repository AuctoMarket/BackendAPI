package product

import (
	"BackendAPI/data"
	"BackendAPI/store"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
Creates ID's for each image and uses the id for the filename of the image. Stores the
Stores the image in the s3 bucket and returns a array of the urls to access the image.
*/
func CreateProductImages(db *sql.DB, client *s3.Client, productId string, images []io.Reader) (data.ProductImageCreateData, *utils.ErrorHandler) {
	var response data.ProductImageCreateData

	validateErr := validateCreateProductImages(db, productId, images)

	if validateErr != nil {
		return response, validateErr
	}

	query := `INSERT INTO product_images(product_id, image_no) VALUES `

	for i := 0; i < len(images); i++ {
		query += `('` + productId + "'," + strconv.Itoa(i+1) + `)`
		if i < len(images)-1 {
			query += `,`
		}
	}

	fmt.Println(query)

	query += ` RETURNING product_image_id;`

	rows, err := db.QueryContext(context.Background(), query)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(err)
		utils.LogError(err, "Error in Inserting Product Image Rows")
		return response, errResp
	}

	for rows.Next() {
		var id string
		rows.Scan(&id)
		response.Images = append(response.Images, id)
	}

	err = store.UploadImages(client, response.Images, images)

	if err != nil {
		errResp := utils.InternalServerError(err)
		return response, errResp
	}

	query = `UPDATE products
	SET image_count = $1 WHERE product_id = $2`

	_, err = db.ExecContext(context.Background(), query, len(images), productId)

	if err != nil {
		errResp := utils.InternalServerError(err)
		return response, errResp
	}

	response.ProductId = productId
	return response, nil
}

/*
Validates the insertion of new product images for a product
*/
func validateCreateProductImages(db *sql.DB, productId string, images []io.Reader) *utils.ErrorHandler {
	if !DoesProductExist(db, productId) {
		return utils.BadRequestError("Product with given id does not exist")
	}

	if doProductImagesExist(db, productId) {
		return utils.BadRequestError("Product with given id already has images")
	}

	if len(images) > 5 {
		return utils.BadRequestError("Too many images uploaded, at most 5 images per post")
	}

	if len(images) == 0 {
		return utils.BadRequestError("No images attached, at least 1 image per post")
	}

	return nil
}

/*
Checks to see if a product has images already attached to it
*/
func doProductImagesExist(db *sql.DB, productId string) bool {
	var productImageExists bool
	query := `SELECT EXISTS(SELECT * FROM product_images WHERE product_id = $1);`
	err := db.QueryRowContext(context.Background(), query, productId).Scan(&productImageExists)

	if err != nil {
		return false
	}

	return productImageExists
}

/*
Transforms an image to an image path
*/
func makeImagePath(imageId string) (string, *utils.ErrorHandler) {
	api_env, envExists := os.LookupEnv("API_ENV")

	if !envExists {
		errResp := utils.InternalServerError(nil)
		utils.LogError(errors.New("Error in loading. env"), "Error in getting product by id")
		return "", errResp
	}

	if imageId == "" {
		errResp := utils.InternalServerError(nil)
		utils.LogError(errors.New("Error in creating image path: no image id"), "Error in creating image path: no image id")
		return "", errResp
	}

	if api_env == "local" {
		imageId = os.Getenv("S3_LOCAL_URL") + "/products/images/" + imageId
	} else if api_env == "dev" {
		imageId = os.Getenv("S3_DEV_URL") + "/products/images/" + imageId
	} else {
		imageId = os.Getenv("S3_PROD_URL") + "/products/images/" + imageId
	}

	return imageId, nil
}
