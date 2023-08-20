package product

import (
	"BackendAPI/api/seller"
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"time"
)

/*
Gets a new product by its product id. If the product with the given product id does not exist, a
404 Not found error is returned
*/
func GetProductById(db *sql.DB, productId string) (data.GetProductResponseData, *utils.ErrorHandler) {
	var response data.GetProductResponseData

	productExists := doesProductExist(db, productId)
	if !productExists {
		return response, utils.NotFoundError("Product with given id does not exist")
	}
	query := `SELECT products.seller_id, sellers.seller_name, sellers.followers, title, description, condition, price, 
		product_type, posted_date, product_quantity, sold_quantity, product_image_id,image_no 
		FROM (
			(products INNER JOIN product_images
				 ON products.product_id = product_images.product_id)
			INNER JOIN sellers
				ON products.seller_id = sellers.seller_id) 
		WHERE products.product_id = $1;`
	rows, err := db.QueryContext(context.Background(), query, productId)
	defer rows.Close()

	for rows.Next() {
		var image string
		var imageNo int
		rows.Scan(
			&response.SellerInfo.SellerId, &response.SellerInfo.SellerName, &response.SellerInfo.Followers,
			&response.Title, &response.Description, &response.Condition, &response.Price,
			&response.ProductType, &response.PostedDate, &response.Quantity, &response.SoldQuantity,
			&image, &imageNo)

		image, pathErr := makeImagePath(image)
		if err != nil {
			return response, pathErr
		}

		response.ProductImages = append(response.ProductImages, data.ProductImageData{ProductImagePath: image, ProductImageNo: imageNo})
	}

	if err != nil {
		errResp := utils.InternalServerError(err)
		utils.LogError(err, "Error in Selecting Product rows")
		return response, errResp
	}

	response.ProductId = productId

	return response, nil
}

/*
Creates a product given product information. If there is an issue with the inputed data, it returns a
400 bad request.
*/
func CreateProduct(db *sql.DB, product data.ProductCreateData) (data.ProductCreateResponseData, *utils.ErrorHandler) {
	var response data.ProductCreateResponseData

	validationErr := validateCreateProduct(db, product)
	if validationErr != nil {
		return response, validationErr
	}

	postedDate := time.Now()
	query := `INSERT INTO products(
			title, seller_id, description, product_type, posted_date, price, condition, product_quantity) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING product_id;`

	err := db.QueryRowContext(
		context.Background(), query,
		product.Title, product.SellerId, product.Description,
		product.ProductType, postedDate, product.Price, product.Condition, product.Quantity).Scan(&response.ProductId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Inserting Product rows")
		return response, errResp
	}

	product.CreateProductResponseFromRequest(&response)
	response.PostedDate = postedDate.String()

	return response, nil
}

/*
Gets a list of products specified by the parameters. If no such products exist returns an empty list and if params are incorrect,
returns a 400 bad request
*/

func GetProductList(db *sql.DB, sellerId string, sortBy string) ([]data.GetProductResponseData, *utils.ErrorHandler) {
	var products []data.GetProductResponseData
	productMap := make(map[string]data.GetProductResponseData)

	validErr := validateGetProductList(db, sellerId, sortBy)

	if validErr != nil {
		return products, validErr
	}

	query := `SELECT products.product_id, products.seller_id, sellers.seller_name, sellers.followers, title, description, condition, price, 
	product_type, posted_date, product_quantity, sold_quantity, product_image_id,image_no 
	FROM (
		(products INNER JOIN product_images
			 ON products.product_id = product_images.product_id)
		INNER JOIN sellers
			ON products.seller_id = sellers.seller_id)`

	//Add query params to the query
	if sellerId != "None" {
		query = query + ` WHERE sellers.seller_id = '` + sellerId + `'`
	}

	if sortBy == "None" {
		query = query + ` ORDER BY posted_date DESC`
	}

	query = query + `;`

	rows, err := db.QueryContext(context.Background(), query)

	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting product rows")
		return products, errResp
	}

	for rows.Next() {
		var product data.GetProductResponseData
		var imagePath string
		var imageNo int
		// scan the product
		err = rows.Scan(&product.ProductId, &product.SellerInfo.SellerId, &product.SellerInfo.SellerName,
			&product.SellerInfo.Followers, &product.Title, &product.Description, &product.Condition,
			&product.Price, &product.ProductType, &product.PostedDate, &product.Quantity, &product.SoldQuantity,
			&imagePath, &imageNo)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in selecting product rows")
			return products, errResp
		}

		//Convert image id to path
		imagePath, pathErr := makeImagePath(imagePath)
		if err != nil {
			return products, pathErr
		}

		// Add the product to the map, if it exists add just the image to the product images array
		if productMap[product.ProductId].ProductId == "" {
			product.ProductImages = append(product.ProductImages, data.ProductImageData{ProductImagePath: imagePath, ProductImageNo: imageNo})
		} else {
			product = productMap[product.ProductId]
			product.ProductImages = append(product.ProductImages, data.ProductImageData{ProductImagePath: imagePath, ProductImageNo: imageNo})
		}

		productMap[product.ProductId] = product
	}

	// iterate through map and return the list of unique products
	for _, v := range productMap {
		products = append(products, v)
	}

	return products, nil
}

/*
Checks wether a Product with a given product id already exists in the database
and returns true if it does false otherwise.
*/
func doesProductExist(db *sql.DB, productId string) bool {
	var productExists bool
	query := `SELECT EXISTS(SELECT * FROM products WHERE product_id = $1);`
	err := db.QueryRowContext(context.Background(), query, productId).Scan(&productExists)

	if err != nil {
		return false
	}

	return productExists
}

/*
Validates the various fields in the create product request body to ensure they are valid.
Returns error if request body is not valid
*/
func validateCreateProduct(db *sql.DB, product data.ProductCreateData) *utils.ErrorHandler {
	if product.Condition < 0 || product.Condition > 5 {
		utils.LogMessage("Condition is less than 0 or greater than 5")
		return utils.BadRequestError("Bad condition data")
	}

	if product.Price < 0 {
		utils.LogMessage("Price is less than 0")
		return utils.BadRequestError("Bad price data")
	}

	if product.ProductType != "Buy-Now" && product.ProductType != "Pre-Order" {
		utils.LogMessage("Product type is not recognised")
		return utils.BadRequestError("Bad product_type data")
	}

	if !seller.DoesSellerExist(db, product.SellerId) {
		utils.LogMessage("Seller Id provided does not exist")
		return utils.BadRequestError("Bad seller_id data")
	}

	if product.Quantity <= 0 {
		utils.LogMessage("Quantity cannot be less than 1")
		return utils.BadRequestError("Bad quantity data")
	}

	return nil
}

/*
Validates the various query params in the get products request and returns an error if they are
invalid.
*/
func validateGetProductList(db *sql.DB, sellerId string, sortBy string) *utils.ErrorHandler {
	if sortBy != "None" && sortBy != "posted_date" {
		utils.LogMessage("Field to sort by does not exist")
		return utils.BadRequestError("Bad sort_by param")
	}

	if sellerId != "None" && !seller.DoesSellerExist(db, sellerId) {
		utils.LogMessage("Seller Id Does not exist")
		return utils.BadRequestError("Bad seller_id param")
	}

	return nil
}
