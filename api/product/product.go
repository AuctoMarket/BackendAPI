package product

import (
	"BackendAPI/api/seller"
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"strconv"
	"time"
)

/*
Gets a new product by its product id. If the product with the given product id does not exist, a
404 Not found error is returned
*/
func GetProductById(db *sql.DB, productId string) (data.GetProductResponseData, *utils.ErrorHandler) {
	var response data.GetProductResponseData

	productExists := DoesProductExist(db, productId)
	if !productExists {
		return response, utils.NotFoundError("Product with given id does not exist")
	}
	query := `SELECT products.seller_id, sellers.seller_name, sellers.followers, title, description, condition, price, 
		product_type, language, expansion, posted_date::TEXT, product_quantity, sold_quantity, product_image_id,image_no, 
		COALESCE(preorder_information.order_by::TEXT, ''), COALESCE(preorder_information.releases_on::TEXT, ''), 
		COALESCE(preorder_information.discount, 0) 
		FROM (((
			products INNER JOIN product_images ON products.product_id = product_images.product_id)
				INNER JOIN sellers ON products.seller_id = sellers.seller_id)
					LEFT OUTER JOIN preorder_information ON products.product_id = preorder_information.product_id)
		WHERE products.product_id = $1;`
	rows, err := db.QueryContext(context.Background(), query, productId)
	defer rows.Close()

	for rows.Next() {
		var image string
		var imageNo int
		rows.Scan(
			&response.SellerInfo.SellerId, &response.SellerInfo.SellerName, &response.SellerInfo.Followers,
			&response.Title, &response.Description, &response.Condition, &response.Price,
			&response.ProductType, &response.Language, &response.Expansion, &response.PostedDate, &response.Quantity,
			&response.SoldQuantity, &image, &imageNo, &response.OrderBy, &response.ReleasesOn, &response.Discount)

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
func CreateProduct(db *sql.DB, product data.CreateProductData) (data.CreateProductResponseData, *utils.ErrorHandler) {
	validationErr := validateCreateProduct(db, product)
	if validationErr != nil {
		var response data.CreateProductResponseData
		return response, validationErr
	}

	if product.ProductType == "Buy-Now" {
		return CreateBuyNow(db, product)
	} else {
		return CreatePreOrder(db, product)
	}
}

/*
Handles the creation of a buynow product listing
*/
func CreateBuyNow(db *sql.DB, product data.CreateProductData) (data.CreateProductResponseData, *utils.ErrorHandler) {
	var response data.CreateProductResponseData
	postedDate := time.Now()
	query := `INSERT INTO products(
			title, seller_id, description, product_type, language, expansion, 
			posted_date, price, condition, product_quantity) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING product_id, posted_date::TEXT;`

	err := db.QueryRowContext(
		context.Background(), query,
		product.Title, product.SellerId, product.Description,
		product.ProductType, product.Language, product.Expansion, postedDate, product.Price,
		product.Condition, product.Quantity).Scan(&response.ProductId, &response.PostedDate)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Inserting Product rows")
		return response, errResp
	}

	product.ProductCreateResponseFromRequest(&response)

	return response, nil
}

/*
Handles the creation of a preorder product listing
*/
func CreatePreOrder(db *sql.DB, product data.CreateProductData) (data.CreateProductResponseData, *utils.ErrorHandler) {
	var response data.CreateProductResponseData
	postedDate := time.Now()
	query := `INSERT INTO products(
		title, seller_id, description, product_type, language, expansion, 
		posted_date, price, condition, product_quantity) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING product_id, posted_date::TEXT;`

	err := db.QueryRowContext(
		context.Background(), query,
		product.Title, product.SellerId, product.Description,
		product.ProductType, product.Language, product.Expansion, postedDate, product.Price,
		product.Condition, product.Quantity).Scan(&response.ProductId, &response.PostedDate)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Inserting Product rows")
		return response, errResp
	}

	query = `INSERT INTO preorder_information(
		product_id, order_by, releases_on, discount) 
		VALUES ($1,$2::timestamptz,$3::timestamptz,$4);`

	_, err = db.ExecContext(context.Background(), query, response.ProductId, product.OrderBy, product.ReleasesOn, product.Discount)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Inserting Preorder Information rows")
		return response, errResp
	}

	product.ProductCreateResponseFromRequest(&response)
	response.PostedDate = postedDate.String()

	return response, nil
}

/*
Gets a list of products specified by the parameters. If no such products exist returns an empty list and if params are incorrect,
returns a 400 bad request
*/

func GetProductList(db *sql.DB, request data.GetProductListData) ([]data.GetProductResponseData, *utils.ErrorHandler) {
	var response []data.GetProductResponseData
	productMap := make(map[string]int)

	query := `SELECT products.product_id, products.seller_id, sellers.seller_name, sellers.followers, title, description, condition, price, 
	product_type, language, expansion, posted_date::TEXT, product_quantity, sold_quantity, product_image_id,image_no, 
	COALESCE(preorder_information.order_by::TEXT, ''),COALESCE(preorder_information.releases_on::TEXT, ''), 
	COALESCE(preorder_information.discount, 0) 
	FROM (((
		products INNER JOIN product_images ON products.product_id = product_images.product_id)
			INNER JOIN sellers ON products.seller_id = sellers.seller_id)
				LEFT OUTER JOIN preorder_information ON products.product_id = preorder_information.product_id)`

	query = AddProductFiltering(query, request.MinPrice, request.MaxPrice, request.Language, request.ProductType, request.Expansion)
	query = AddProductSorting(query, request.SortBy)
	rows, err := db.QueryContext(context.Background(), query)

	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting product rows")
		return response, errResp
	}

	for rows.Next() {
		var product data.GetProductResponseData
		var imagePath string
		var imageNo int
		// scan the product
		err = rows.Scan(&product.ProductId, &product.SellerInfo.SellerId, &product.SellerInfo.SellerName,
			&product.SellerInfo.Followers, &product.Title, &product.Description, &product.Condition,
			&product.Price, &product.ProductType, &product.Language, &product.Expansion, &product.PostedDate, &product.Quantity,
			&product.SoldQuantity, &imagePath, &imageNo, &product.OrderBy, &product.ReleasesOn, &product.Discount)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in selecting product rows")
			return response, errResp
		}

		//Convert image id to path
		imagePath, pathErr := makeImagePath(imagePath)
		if err != nil {
			return response, pathErr
		}

		// If product is already in response array, add the image to that product, otherwise add the entire product
		if productMap[product.ProductId] == 0 {
			product.ProductImages = append(product.ProductImages, data.ProductImageData{ProductImagePath: imagePath, ProductImageNo: imageNo})
			response = append(response, product)
			index := len(response)
			productMap[product.ProductId] = index
		} else {
			p := response[productMap[product.ProductId]-1]
			p.ProductImages = append(
				p.ProductImages,
				data.ProductImageData{ProductImagePath: imagePath, ProductImageNo: imageNo})
			response[productMap[product.ProductId]-1] = p

		}

	}

	return response, nil
}

/*
Adds the sorting to the query to determine the order of the products
*/
func AddProductSorting(query string, sortBy string) string {
	if sortBy == "price-low" {
		query += ` ORDER BY products.price ASC, preorder_information.discount DESC`
	} else if sortBy == "price-high" {
		query += ` ORDER BY products.price DESC, preorder_information.discount ASC`
	} else if sortBy == "name-asc" {
		query += ` ORDER BY products.title ASC`
	} else if sortBy == "name-desc" {
		query += ` ORDER BY products.title DESC`
	} else {
		query += ` ORDER BY products.posted_date DESC`
	}

	query += `, product_images.image_no ASC`
	return query
}

/*
Adds the filtering to the query to filter out certain products
*/
func AddProductFiltering(query string, minPrice int, maxPrice int, language string, productType string, expansion string) string {
	var hasFiltered bool = false

	if productType != "None" {
		if !hasFiltered {
			query += ` WHERE products.product_type =`
			if productType == "Pre-Order" {
				query += ` 'Pre-Order'`
			} else {
				query += ` 'Buy-Now'`
			}
			hasFiltered = true
		}
	}

	if language != "None" {
		if !hasFiltered {
			query += ` WHERE products.language =`
			if language == "Eng" {
				query += ` 'Eng'`
			} else {
				query += ` 'Jap'`
			}
			hasFiltered = true
		} else {
			query += ` AND products.language =`
			if language == "Eng" {
				query += ` 'Eng'`
			} else {
				query += ` 'Jap'`
			}
		}
	}

	if minPrice > 0 {
		if !hasFiltered {
			query += ` WHERE products.price >= ` + strconv.Itoa(minPrice)
			hasFiltered = true
		} else {
			query += ` AND products.price >= ` + strconv.Itoa(minPrice)
		}
	}

	if maxPrice > 0 {
		if !hasFiltered {
			query += ` WHERE products.price <= ` + strconv.Itoa(maxPrice)
			hasFiltered = true
		} else {
			query += ` AND products.price <= ` + strconv.Itoa(maxPrice)
		}
	}

	if expansion != "None" {
		if !hasFiltered {
			query += ` WHERE products.expansion = '` + expansion + `'`
			hasFiltered = true
		} else {
			query += ` AND products.expansion = '` + expansion + `'`
		}
	}
	return query
}

/*
Adds pages to the query to allow for pagination
*/
func AddPagesProduct(query string, page int) string {
	return query
}

/*
Checks wether a Product with a given product id already exists in the database
and returns true if it does false otherwise.
*/
func DoesProductExist(db *sql.DB, productId string) bool {
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
func validateCreateProduct(db *sql.DB, product data.CreateProductData) *utils.ErrorHandler {
	if product.Condition < 0 || product.Condition > 5 {
		utils.LogMessage("Condition is less than 0 or greater than 5")
		return utils.BadRequestError("Bad condition data")
	}

	if product.Language != "Eng" && product.Language != "Jap" {
		utils.LogMessage("Product language is not recognised")
		return utils.BadRequestError("Bad language data")
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

	if product.ProductType == "Pre-Order" && product.ReleasesOn == "" {
		utils.LogMessage("Pre orders need a release date")
		return utils.BadRequestError("Bad pre-order data")
	}

	if product.ProductType == "Pre-Order" && product.OrderBy == "" {
		utils.LogMessage("Pre orders need a order by date")
		return utils.BadRequestError("Bad pre-order data")
	}

	return nil
}
