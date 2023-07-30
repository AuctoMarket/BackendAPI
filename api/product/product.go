package product

import (
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
)

func GetProductById(db *sql.DB, productId string) (data.ProductData, *utils.ErrorHandler) {
	var product data.ProductData

	productExists := doesProductExist(db, productId)
	if !productExists {
		return product, utils.NotFoundError("Product with given id does not exist")
	}
	query := `SELECT title, description, condition, price, product_type, posted_date from products WHERE product_id = $1;`
	err := db.QueryRowContext(context.Background(), query, productId).Scan(
		&product.Title, &product.Description, &product.Condition, &product.Price, &product.ProductType)

	if err != nil {
		errResp := utils.InternalServerError()
		utils.LogError(err, "Error in Selecting Product rows")
		return product, errResp
	}

	return product, nil
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
