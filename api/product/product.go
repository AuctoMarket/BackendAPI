package product

import (
	"context"
	"database/sql"
)

func GetProductById(db *sql.DB, productId string) {

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
