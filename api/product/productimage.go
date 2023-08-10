package product

import (
	"BackendAPI/data"
	"BackendAPI/store"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

/*
Creates ID's for each image and uses the id for the filename of the image. Stores the
Stores the image in the s3 bucket and returns a array of the urls to access the image.
*/
func CreateProductImages(db *sql.DB, client *s3.Client, productId string, images []io.Reader) (data.ProductImageCreateData, *utils.ErrorHandler) {
	var response data.ProductImageCreateData

	if !doesProductExist(db, productId) {
		return response, utils.BadRequestError("Product with given id does not exist")
	}

	if len(images) > 5 {
		return response, utils.BadRequestError("Too many images uploaded, at most 5 images per post")
	}

	if len(images) == 0 {
		return response, utils.BadRequestError("No images attached, at least 1 image per post")
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
