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
	query := `SELECT products.seller_id, sellers.seller_name, title, description, condition, price, 
		product_type, language, expansion, posted_date::TEXT, product_quantity, sold_quantity, product_image_id,image_no, 
		COALESCE(preorder_information.order_by::TEXT, ''), COALESCE(preorder_information.releases_on::TEXT, ''), 
		COALESCE(product_discounts.discount, 0) 
		FROM ((((
			products INNER JOIN product_images ON products.product_id = product_images.product_id)
				INNER JOIN sellers ON products.seller_id = sellers.seller_id)
					LEFT OUTER JOIN preorder_information ON products.product_id = preorder_information.product_id)
						LEFT OUTER JOIN product_discounts ON product_discounts.product_id = products.product_id)
		WHERE products.product_id = $1;`
	rows, err := db.QueryContext(context.Background(), query, productId)
	defer rows.Close()

	for rows.Next() {
		var image string
		var imageNo int
		rows.Scan(
			&response.SellerInfo.SellerId, &response.SellerInfo.SellerName,
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
	err := validateCreateProduct(db, product)
	var response data.CreateProductResponseData
	if err != nil {
		var response data.CreateProductResponseData
		return response, err
	}

	if product.ProductType == "Buy-Now" {
		response, err = CreateBuyNow(db, product)

	} else {
		response, err = CreatePreOrder(db, product)
	}

	if err != nil {
		return response, err
	}

	if product.Discount > 0 {
		err = addDiscount(db, response.ProductId, product.Discount)
	}

	return response, err

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
		product_id, order_by, releases_on) 
		VALUES ($1,$2::timestamptz,$3::timestamptz);`

	_, err = db.ExecContext(context.Background(), query, response.ProductId, product.OrderBy, product.ReleasesOn)

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
Adds discount for a product in terms of the product amount in cents
*/
func addDiscount(db *sql.DB, productId string, discountAmount int) *utils.ErrorHandler {
	query := `INSERT INTO product_discounts(
		product_id, discount) 
		VALUES ($1,$2);`

	_, err := db.ExecContext(context.Background(), query, productId, discountAmount)
	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Inserting Product Discounts rows")
		return errResp
	}

	return nil

}

/*
Gets a list of products specified by the parameters. If no such products exist returns an empty list and if params are incorrect,
returns a 400 bad request
*/

func GetProductList(db *sql.DB, request data.GetProductListRequestData) (data.GetProductListResponseData, *utils.ErrorHandler) {
	var response data.GetProductListResponseData
	var products []data.GetProductResponseData
	productMap := make(map[string]int)

	query := `
	SELECT
		products.product_id,
		seller_id,
		seller_name,
		title,
		description,
		condition,
		price,
		product_type,
		language,
		expansion,
		posted_date::TEXT,
		product_quantity,
		sold_quantity,
		product_image_id,
		image_no,
		COALESCE(order_by::TEXT, '') AS order_by,
		COALESCE(releases_on::TEXT, '') AS releases_on,
		COALESCE(discount, 0) AS discount
	FROM
		product_images
		RIGHT JOIN (
			SELECT
				products.product_id,
				sellers.seller_id,
				seller_name,
				title,
				description,
				condition,
				price,
				product_type,
				language,
				expansion,
				posted_date,
				product_quantity,
				sold_quantity,
				order_by,
				releases_on,
				discount
			FROM (((products
						INNER JOIN sellers ON products.seller_id = sellers.seller_id)
					LEFT OUTER JOIN preorder_information ON products.product_id = preorder_information.product_id)
				LEFT OUTER JOIN product_discounts ON product_discounts.product_id = products.product_id)`

	query = AddProductFiltering(query, request.Prices, request.Languages, request.ProductTypes, request.Expansions)
	query = AddProductSorting(query, request.SortBy)
	query = AddPagesProduct(query, request.Anchor, request.Limit)
	query += `) products ON products.product_id = product_images.product_id`
	query = AddProductSorting(query, request.SortBy)
	query += `, image_no ASC`

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
			&product.Title, &product.Description, &product.Condition, &product.Price, &product.ProductType,
			&product.Language, &product.Expansion, &product.PostedDate, &product.Quantity,
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
			products = append(products, product)
			index := len(products)
			productMap[product.ProductId] = index
		} else {
			p := products[productMap[product.ProductId]-1]
			p.ProductImages = append(
				p.ProductImages,
				data.ProductImageData{ProductImagePath: imagePath, ProductImageNo: imageNo})
			products[productMap[product.ProductId]-1] = p

		}

	}

	response.Products = products
	response.ProductCount = getProductCount(db, request.Prices, request.Languages, request.ProductTypes, request.Expansions)

	return response, nil
}

/*
Gets the total number of products after applying filters
*/
func getProductCount(db *sql.DB, prices []string, languages []string, productTypes []string, expansions []string) int {
	var count int
	query := `SELECT COUNT(*) FROM products`
	query = AddProductFiltering(query, prices, languages, productTypes, expansions)

	db.QueryRowContext(context.Background(), query).Scan(&count)

	return count
}

/*
Adds the sorting to the query to determine the order of the products
*/
func AddProductSorting(query string, sortBy string) string {
	if sortBy == "price-low" {
		query += ` ORDER BY products.price ASC, discount DESC`
	} else if sortBy == "price-high" {
		query += ` ORDER BY products.price DESC, discount ASC`
	} else if sortBy == "name-asc" {
		query += ` ORDER BY products.title ASC`
	} else if sortBy == "name-desc" {
		query += ` ORDER BY products.title DESC`
	} else {
		query += ` ORDER BY products.posted_date DESC`
	}

	return query
}

/*
Adds the filtering to the query to filter out certain products
*/
func AddProductFiltering(query string, prices []string, languages []string, productTypes []string, expansions []string) string {
	var hasFiltered bool = false

	if len(productTypes) > 0 {
		for i := 0; i < len(productTypes); i++ {
			if !hasFiltered {
				filter := ` WHERE products.product_type =`
				if productTypes[i] == "Pre-Order" {
					filter += ` 'Pre-Order'`
					query += filter
				}
				if productTypes[i] == "Buy-Now" {
					filter += ` 'Buy-Now'`
					query += filter
				}
				hasFiltered = true
			} else {
				filter := ` OR products.product_type =`
				if productTypes[i] == "Pre-Order" {
					filter += ` 'Pre-Order'`
					query += filter
				}
				if productTypes[i] == "Buy-Now" {
					filter += ` 'Buy-Now'`
					query += filter
				}
			}
		}
	}

	if len(languages) > 0 {
		for i := 0; i < len(languages); i++ {
			if !hasFiltered {
				filter := ` WHERE products.language =`
				if languages[i] == "Eng" {
					filter += ` 'Eng'`
					query += filter
				}
				if languages[i] == "Jap" {
					filter += ` 'Jap'`
					query += filter
				}
				hasFiltered = true
			} else {
				var filter string
				if i == 0 {
					filter += ` AND products.language =`
				} else {
					filter += ` OR products.language =`
				}

				if languages[i] == "Eng" {
					filter += ` 'Eng'`
					query += filter
				}
				if languages[i] == "Jap" {
					filter += ` 'Jap'`
					query += filter
				}
			}
		}
	}

	if len(expansions) > 0 {
		for i := 0; i < len(expansions); i++ {
			if !hasFiltered {
				query += ` WHERE products.expansion = '` + expansions[i] + `'`
				hasFiltered = true
			} else {
				if i == 0 {
					query += ` AND products.expansion = '` + expansions[i] + `'`
				} else {
					query += ` OR products.expansion = '` + expansions[i] + `'`
				}
			}
		}
	}

	if len(prices) > 0 {
		for i := 0; i < len(prices); i++ {
			var filter string
			if !hasFiltered {
				filter += ` WHERE products.price`
				hasFiltered = true
			} else {
				if i == 0 {
					filter += ` AND products.price`
				} else {
					filter += ` OR products.price`
				}
			}

			if prices[i] == "0-20" {
				filter += ` BETWEEN 0 AND 2000`
				query += filter
			}
			if prices[i] == "20-50" {
				filter += ` BETWEEN 2000 AND 5000`
				query += filter
			}
			if prices[i] == "50-100" {
				filter += ` BETWEEN 5000 AND 10000`
				query += filter
			}
			if prices[i] == "100-200" {
				filter += ` BETWEEN 10000 AND 20000`
				query += filter
			}
			if prices[i] == "200" {
				filter += ` >= 20000`
				query += filter
			}
		}
	}

	return query
}

/*
Adds pages to the query to allow for pagination
*/
func AddPagesProduct(query string, anchor int, limit int) string {
	query += ` OFFSET ` + strconv.Itoa(anchor) + ` LIMIT ` + strconv.Itoa(limit)
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

	if product.Discount < 0 {
		utils.LogMessage("Discount is less than 0")
		return utils.BadRequestError("Bad discount data")
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
