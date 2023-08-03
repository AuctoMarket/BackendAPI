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
func GetProductById(db *sql.DB, productId string) (data.ProductResponseData, *utils.ErrorHandler) {
	var response data.ProductResponseData

	productExists := doesProductExist(db, productId)
	if !productExists {
		return response, utils.NotFoundError("Product with given id does not exist")
	}
	query := `SELECT seller_id, title, description, condition, price, 
		product_type, posted_date from products WHERE product_id = $1;`
	err := db.QueryRowContext(context.Background(), query, productId).Scan(
		&response.SellerId, &response.Title, &response.Description,
		&response.Condition, &response.Price, &response.ProductType, &response.PostedDate)

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
func CreateProduct(db *sql.DB, product data.ProductCreateData) (data.ProductResponseData, *utils.ErrorHandler) {
	var response data.ProductResponseData

	validationErr := validateCreateProduct(db, product)
	if validationErr != nil {
		return response, validationErr
	}

	postedDate := time.Now()
	query := `INSERT INTO products(
			title, seller_id, description, product_type, posted_date, price, condition) 
			VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING product_id;`

	err := db.QueryRowContext(
		context.Background(), query,
		product.Title, product.SellerId, product.Description,
		product.ProductType, postedDate, product.Price, product.Condition).Scan(&response.ProductId)

	if err != nil {
		errResp := utils.InternalServerError(err)
		utils.LogError(err, "Error in Inserting Product rows")
		return response, errResp
	}

	product.CreateResponseFromRequest(&response)
	response.PostedDate = postedDate.String()

	return response, nil
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

	return nil
}
