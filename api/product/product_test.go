package product

import (
	"BackendAPI/data"
	"BackendAPI/store"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateProduct(t *testing.T) {
	db, dbErr := store.SetupTestDB("../../.env")
	assert.NoError(t, dbErr)

	sellerId, dataErr := createDummySeller(db)
	assert.NoError(t, dataErr)

	//Test 1: No errors, product is valid
	testCreateProduct1 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 3, Language: "Eng", Expansion: "Test"}
	err := validateCreateProduct(db, testCreateProduct1)
	assert.Empty(t, err)

	//Test 2: Error, product price is invalid
	testCreateProduct2 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: -5, Condition: 3, Quantity: 3, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct2)
	assert.Error(t, err)
	assert.Equal(t, "Bad price data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 3: Error, product condition is invalid
	testCreateProduct3 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 13, Quantity: 3, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct3)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 4: Error, product condition is invalid
	testCreateProduct4 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: -1, Quantity: 3, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct4)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 5: No errors, product is valid
	testCreateProduct5 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3, ReleasesOn: "2023-10-02 15:59:59.170135+00",
		OrderBy: "2023-10-02 15:59:59.170135+00", Discount: 10, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct5)
	assert.Empty(t, err)

	//Test 6: No errors, product is valid
	testCreateProduct6 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 0, Condition: 0, Quantity: 3, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct6)
	assert.Empty(t, err)

	//Test 7:Error product type is wrong
	testCreateProduct7 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-It-Now", Price: 0, Condition: 0, Quantity: 3, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct7)
	assert.Error(t, err)
	assert.Equal(t, "Bad product_type data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 8:Error seller id does not exist
	testCreateProduct8 := data.CreateProductData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 3, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct8)
	assert.Error(t, err)
	assert.Equal(t, "Bad seller_id data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 9: No errors, product is valid
	testCreateProduct9 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 0, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct9)
	assert.Error(t, err)
	assert.Equal(t, "Bad quantity data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 10: Missing release date
	testCreateProduct10 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
		OrderBy: "2023-10-02 15:59:59.170135+00", Discount: 10, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct10)
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 11: Missing order by
	testCreateProduct11 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
		ReleasesOn: "2023-10-02 15:59:59.170135+00", Discount: 10, Language: "Eng", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct11)
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 12: Wrong Language
	testCreateProduct12 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 3, Language: "Wrong", Expansion: "Test"}
	err = validateCreateProduct(db, testCreateProduct12)
	assert.Error(t, err)
	assert.Equal(t, "Bad language data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	store.CloseDB(db)
}

func TestDoesProductExist(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)

	sellerid, sellerErr := createDummySeller(db)
	assert.NoError(t, sellerErr)
	productIds, productErr := createDummyProducts(db, sellerid)
	assert.NoError(t, productErr)

	//Test 1: Product Id Exists
	doesExist := DoesProductExist(db, productIds[0])
	assert.Equal(t, true, doesExist)

	//Test 2: Product Id Exists
	doesExist = DoesProductExist(db, productIds[1])
	assert.Equal(t, true, doesExist)

	//Test 3: Product Id Exists
	doesExist = DoesProductExist(db, productIds[2])
	assert.Equal(t, true, doesExist)

	//Test 4: Product Id does not exist
	doesExist = DoesProductExist(db, "wrong id")
	assert.Equal(t, false, doesExist)

	//Test 4: Product Id does not exist
	doesExist = DoesProductExist(db, productIds[0]+"1")
	assert.Equal(t, false, doesExist)

	store.CloseDB(db)
}

func TestGetProductById(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)

	sellerid, sellerErr := createDummySeller(db)
	assert.NoError(t, sellerErr)
	productIds, productErr := createDummyProducts(db, sellerid)
	assert.NoError(t, productErr)
	productImageIds, imageErr := createDummyProductImages(db, productIds)
	assert.NoError(t, imageErr)

	//Test 1: Product Id exists
	response, err := GetProductById(db, productIds[0])
	assert.Empty(t, err)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 10000, response.Price)
	assert.Equal(t, int8(3), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)
	assert.Equal(t, 3, response.Quantity)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com"+productImageIds[0], response.ProductImages[0].ProductImagePath)
	assert.Equal(t, "Eng", response.Language)
	assert.Equal(t, "Test", response.Expansion)

	//Test 2: Product Id exists, Japanese language
	response, err = GetProductById(db, productIds[1])
	assert.Empty(t, err)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 20000, response.Price)
	assert.Equal(t, int8(5), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)
	assert.Equal(t, 3, response.Quantity)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com"+productImageIds[1], response.ProductImages[0].ProductImagePath)
	assert.Equal(t, "Jap", response.Language)
	assert.Equal(t, "Test", response.Expansion)

	//Test 3: Pre-Order, Product id exists
	response, err = GetProductById(db, productIds[3])
	assert.Empty(t, err)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 10000, response.Price)
	assert.Equal(t, int8(4), response.Condition)
	assert.Equal(t, "Pre-Order", response.ProductType)
	assert.Equal(t, 3, response.Quantity)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com"+productImageIds[3], response.ProductImages[0].ProductImagePath)
	assert.Equal(t, "Eng", response.Language)
	assert.Equal(t, "Test", response.Expansion)

	//Test 4: Buy-Now but with discount
	response, err = GetProductById(db, productIds[4])
	assert.Empty(t, err)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 10000, response.Price)
	assert.Equal(t, int8(4), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)
	assert.Equal(t, 3, response.Quantity)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com"+productImageIds[4], response.ProductImages[0].ProductImagePath)
	assert.Equal(t, "Eng", response.Language)
	assert.Equal(t, "Test", response.Expansion)
	assert.Equal(t, 1000, response.Discount)

	//Test 5: Product Id does not exist
	response, err = GetProductById(db, "wrong id")
	assert.Error(t, err)

	store.CloseDB(db)
}

func TestCreateProduct(t *testing.T) {
	db, startupErr := store.SetupTestDB("../../.env")
	assert.NoError(t, startupErr)

	sellerId, dataErr := createDummySeller(db)
	assert.NoError(t, dataErr)

	//Test 1: No error, product is created
	product := data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 0, Quantity: 3, Language: "Eng", Expansion: "Test"}
	response, err := CreateProduct(db, product)
	assert.Empty(t, err)
	assert.Equal(t, "Test", response.Title)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 10, response.Price)
	assert.Equal(t, int8(0), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)
	assert.Equal(t, "Test", response.Expansion)
	assert.Equal(t, "Eng", response.Language)

	//Test 2: No error, product is created
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
		ReleasesOn: "2024-10-02 04:44:17.170135", OrderBy: "2024-10-02 04:44:17.170135",
		Discount: 10, Language: "Eng", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Empty(t, err)
	assert.Equal(t, "Test", response.Title)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 0, response.Price)
	assert.Equal(t, int8(5), response.Condition)
	assert.Equal(t, "Pre-Order", response.ProductType)
	assert.Equal(t, "Test", response.Expansion)
	assert.Equal(t, "Eng", response.Language)

	//Test 3: Error, Condition is negative
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 100, Condition: -1, Quantity: 3, Language: "Eng", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 4: Error, Condition is greater than 5
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 100, Condition: 6, Quantity: 3, Language: "Eng", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 5: Error, Price is less than 0
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: -100, Condition: 5, Quantity: 3, Language: "Eng", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Error(t, err)
	assert.Equal(t, "Bad price data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 6: Error, Quantity is 0
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 100, Condition: 5, Quantity: 0, Language: "Eng", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Error(t, err)
	assert.Equal(t, "Bad quantity data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 7: Error, Pre order release date does not exist
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
		OrderBy: "2023-10-02 04:44:17.170135", Discount: 10, Language: "Eng", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 8: Error, Pre order order by does not exist
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
		ReleasesOn: "2023-10-02 04:44:17.170135", Discount: 10, Language: "Eng", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 9: Error Language does not exist
	product = data.CreateProductData{
		Title: "Test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 100, Condition: 5, Quantity: 0, Language: "wrong", Expansion: "Test"}
	response, err = CreateProduct(db, product)
	assert.Error(t, err)
	assert.Equal(t, "Bad language data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	store.CloseDB(db)
}
func TestAddProductSorting(t *testing.T) {
	//Test 1: Sort by price low to high
	query := AddProductSorting("", "price-low")
	assert.Equal(t, ` ORDER BY products.price ASC, discount DESC`, query)

	//Test 2: Sort by price high to low
	query = AddProductSorting("", "price-high")
	assert.Equal(t, ` ORDER BY products.price DESC, discount ASC`, query)

	//Test 3: Name ascending
	query = AddProductSorting("", "name-asc")
	assert.Equal(t, ` ORDER BY products.title ASC`, query)

	//Test 4: Name descending
	query = AddProductSorting("", "name-desc")
	assert.Equal(t, ` ORDER BY products.title DESC`, query)

	//Test 5: Default
	query = AddProductSorting("", "None")
	assert.Equal(t, ` ORDER BY products.posted_date DESC`, query)

	//Test 6: Random string
	query = AddProductSorting("", "sdjknvjk")
	assert.Equal(t, ` ORDER BY products.posted_date DESC`, query)

}

func TestAddProductFiltering(t *testing.T) {

	//Test 1: No filters
	query := AddProductFiltering("", 0, 0, "None", "None", "None")
	assert.Equal(t, "", query)

	//Test 2: Preorders
	query = AddProductFiltering("", 0, 0, "None", "Pre-Order", "None")
	assert.Equal(t, query, ` WHERE products.product_type = 'Pre-Order'`)

	//Test 3: Buy-Now
	query = AddProductFiltering("", 0, 0, "None", "Buy-Now", "None")
	assert.Equal(t, query, ` WHERE products.product_type = 'Buy-Now'`)

	//Test 4: English
	query = AddProductFiltering("", 0, 0, "Eng", "None", "None")
	assert.Equal(t, query, ` WHERE products.language = 'Eng'`)

	//Test 5: Preorders and Japanese
	query = AddProductFiltering("", 0, 0, "Jap", "Pre-Order", "None")
	assert.Equal(t, query, ` WHERE products.product_type = 'Pre-Order' AND products.language = 'Jap'`)

	//Test 6: Min price
	query = AddProductFiltering("", 10, 0, "None", "None", "None")
	assert.Equal(t, query, ` WHERE products.price >= 10`)

	//Test 7: Max price in japanese
	query = AddProductFiltering("", 0, 100, "Jap", "None", "None")
	assert.Equal(t, query, ` WHERE products.language = 'Jap' AND products.price <= 100`)

	//Test 8: Max price & min price in japanese for buy-now
	query = AddProductFiltering("", 10, 100, "Jap", "Buy-Now", "None")
	assert.Equal(t, query, ` WHERE products.product_type = 'Buy-Now' AND products.language = 'Jap' AND products.price >= 10 AND products.price <= 100`)
}

func TestGetProductList(t *testing.T) {
	utils.LoadDotEnv("../../.env")
	db, startupErr := store.SetupTestDB("../../.env")
	assert.NoError(t, startupErr)

	sellerid, sellerErr := createDummySeller(db)
	assert.NoError(t, sellerErr)
	productIds, productErr := createDummyProducts(db, sellerid)
	assert.NoError(t, productErr)
	_, productImageErr := createDummyProductImages(db, productIds)
	assert.NoError(t, productImageErr)

	//Test 1: Get Product List default options
	req := data.GetProductListRequestData{SortBy: "None", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err := GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)

	//Test 2: Get Product List Sorted by Price (Low-High)
	req = data.GetProductListRequestData{SortBy: "price-low", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, 10000, res.Products[0].Price)
	assert.Equal(t, 10000, res.Products[1].Price)

	//Test 3: Get Product List Sorted by Price (High-Low)
	req = data.GetProductListRequestData{SortBy: "price-high", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, 20000, res.Products[0].Price)
	assert.Equal(t, 10000, res.Products[1].Price)

	//Test 4: Get Product List Sorted by Name (A-Z)
	req = data.GetProductListRequestData{SortBy: "name-asc", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, "Test", res.Products[0].Title)
	assert.Equal(t, "Test1", res.Products[1].Title)

	//Test 4: Get Product List Sorted by Name (Z-A)
	req = data.GetProductListRequestData{SortBy: "name-desc", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, "Test4", res.Products[0].Title)
	assert.Equal(t, "Test3", res.Products[1].Title)

	//Test 5: Get Product List Buy-Now
	req = data.GetProductListRequestData{SortBy: "name-desc", MinPrice: 0, MaxPrice: 0, ProductType: "Buy-Now",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 4, res.ProductCount)
	assert.Equal(t, "Test4", res.Products[0].Title)
	assert.Equal(t, "Test2", res.Products[1].Title)

	//Test 6: Get Product List Pre-Order
	req = data.GetProductListRequestData{SortBy: "name-desc", MinPrice: 0, MaxPrice: 0, ProductType: "Pre-Order",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 1, res.ProductCount)
	assert.Equal(t, "Test3", res.Products[0].Title)

	//Test 7: Get Product Min Price 15000
	req = data.GetProductListRequestData{SortBy: "name-asc", MinPrice: 15000, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 1, res.ProductCount)
	assert.Equal(t, "Test1", res.Products[0].Title)

	//Test 8: Get Product Max price 15000
	req = data.GetProductListRequestData{SortBy: "name-asc", MinPrice: 0, MaxPrice: 15000, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 4, res.ProductCount)
	assert.Equal(t, "Test", res.Products[0].Title)
	assert.Equal(t, "Test2", res.Products[1].Title)

	//Test 9: Get Product expansion Test
	req = data.GetProductListRequestData{SortBy: "name-asc", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "Test", Anchor: 0, Limit: 10}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 4, res.ProductCount)
	assert.Equal(t, "Test", res.Products[0].Title)
	assert.Equal(t, "Test1", res.Products[1].Title)
	assert.Equal(t, "Test3", res.Products[2].Title)

	//Test 10: Get Product List Low Limit
	req = data.GetProductListRequestData{SortBy: "None", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 1}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, 1, len(res.Products))

	//Test 11: Get Product List Limit same as count
	req = data.GetProductListRequestData{SortBy: "None", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 0, Limit: 3}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, 3, len(res.Products))

	//Test 12: Get Product List Anchor added
	req = data.GetProductListRequestData{SortBy: "None", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 1, Limit: 3}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, 3, len(res.Products))

	//Test 13: Anchor higher than count
	req = data.GetProductListRequestData{SortBy: "None", MinPrice: 0, MaxPrice: 0, ProductType: "None",
		Language: "None", Expansion: "None", Anchor: 10, Limit: 3}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 5, res.ProductCount)
	assert.Equal(t, 0, len(res.Products))
	store.CloseDB(db)
}

func createDummySeller(db *sql.DB) (string, error) {
	var sellerId string
	query := `INSERT INTO sellers(email, seller_name, password) VALUES ('test@aucto.io','test','test') RETURNING seller_id`
	err := db.QueryRowContext(context.Background(), query).Scan(&sellerId)

	return sellerId, err
}

func createDummyProducts(db *sql.DB, sellerId string) ([]string, error) {
	var dummyCreateProducts []data.CreateProductData = []data.CreateProductData{
		//Product 1: Test Buy-Now Product English
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10000, Condition: 3, Quantity: 3, Language: "Eng",
			Expansion: "Test"},
		//Product 2: Test Buy-Now Product Japanese
		{Title: "Test1", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 20000, Condition: 5, Quantity: 3, Language: "Jap",
			Expansion: "Test"},
		//Product 3: Test Buy-Now Product expansion 'Test2'
		{Title: "Test2", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10000, Condition: 4, Quantity: 3, Language: "Eng",
			Expansion: "Test2"},
		//Product 4: Test Pre-Order Product
		{Title: "Test3", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 10000, Condition: 4, Quantity: 3, Language: "Eng",
			Expansion: "Test"},
		//Product 5: Test Buy-Now Product with discount
		{Title: "Test4", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10000, Condition: 4, Quantity: 3, Language: "Eng",
			Expansion: "Test", Discount: 1000}}
	var productIds []string

	for i := 0; i < len(dummyCreateProducts); i++ {
		query := `INSERT INTO products(
			title, seller_id, description, product_type, language, expansion, posted_date, price, condition, product_quantity) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING product_id;`
		postedDate := time.Now()
		var productId string
		err := db.QueryRowContext(
			context.Background(), query,
			dummyCreateProducts[i].Title, dummyCreateProducts[i].SellerId, dummyCreateProducts[i].Description,
			dummyCreateProducts[i].ProductType, dummyCreateProducts[i].Language, dummyCreateProducts[i].Expansion, postedDate,
			dummyCreateProducts[i].Price, dummyCreateProducts[i].Condition, dummyCreateProducts[i].Quantity).Scan(&productId)
		if err != nil {
			return nil, err
		}
		productIds = append(productIds, productId)
	}

	query := `INSERT INTO product_discounts(product_id, discount) VALUES ($1,$2);`
	_, err := db.ExecContext(context.Background(), query, productIds[4], dummyCreateProducts[4].Discount)
	if err != nil {
		return nil, err
	}

	return productIds, nil
}
