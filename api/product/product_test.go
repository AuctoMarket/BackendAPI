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
		ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 3}
	err := validateCreateProduct(db, testCreateProduct1)
	assert.Empty(t, err)

	//Test 2: Error, product price is invalid
	testCreateProduct2 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: -5, Condition: 3, Quantity: 3}
	err = validateCreateProduct(db, testCreateProduct2)
	assert.Error(t, err)
	assert.Equal(t, "Bad price data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 3: Error, product condition is invalid
	testCreateProduct3 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 13, Quantity: 3}
	err = validateCreateProduct(db, testCreateProduct3)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 4: Error, product condition is invalid
	testCreateProduct4 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: -1, Quantity: 3}
	err = validateCreateProduct(db, testCreateProduct4)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 5: No errors, product is valid
	testCreateProduct5 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3, ReleasesOn: "test", OrderBy: "test", Discount: 10}
	err = validateCreateProduct(db, testCreateProduct5)
	assert.Empty(t, err)

	//Test 6: No errors, product is valid
	testCreateProduct6 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 0, Condition: 0, Quantity: 3}
	err = validateCreateProduct(db, testCreateProduct6)
	assert.Empty(t, err)

	//Test 7:Error product type is wrong
	testCreateProduct7 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-It-Now", Price: 0, Condition: 0, Quantity: 3}
	err = validateCreateProduct(db, testCreateProduct7)
	assert.Error(t, err)
	assert.Equal(t, "Bad product_type data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 8:Error seller id does not exist
	testCreateProduct8 := data.CreateProductData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 3}
	err = validateCreateProduct(db, testCreateProduct8)
	assert.Error(t, err)
	assert.Equal(t, "Bad seller_id data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 9: No errors, product is valid
	testCreateProduct9 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 0}
	err = validateCreateProduct(db, testCreateProduct9)
	assert.Error(t, err)
	assert.Equal(t, "Bad quantity data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 10: Missing release date
	testCreateProduct10 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3, OrderBy: "test", Discount: 10}
	err = validateCreateProduct(db, testCreateProduct10)
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 11: Missing order by
	testCreateProduct11 := data.CreateProductData{
		Title: "test", SellerId: sellerId, Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3, ReleasesOn: "test", Discount: 10}
	err = validateCreateProduct(db, testCreateProduct11)
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())
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
	assert.Equal(t, 10, response.Price)
	assert.Equal(t, int8(3), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)
	assert.Equal(t, 3, response.Quantity)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com"+productImageIds[0], response.ProductImages[0].ProductImagePath)

	//Test 2: Product Id exists
	response, err = GetProductById(db, productIds[1])
	assert.Empty(t, err)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 100, response.Price)
	assert.Equal(t, int8(5), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)
	assert.Equal(t, 3, response.Quantity)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com"+productImageIds[1], response.ProductImages[0].ProductImagePath)

	//Test 3: Product Id exists
	response, err = GetProductById(db, productIds[2])
	assert.Empty(t, err)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 90, response.Price)
	assert.Equal(t, int8(4), response.Condition)
	assert.Equal(t, "Pre-Order", response.ProductType)
	assert.Equal(t, 3, response.Quantity)
	assert.Equal(t, "https://aucto-s3-local.s3.ap-southeast-1.amazonaws.com"+productImageIds[2], response.ProductImages[0].ProductImagePath)

	//Test 4: Product Id does not exist
	response, err = GetProductById(db, "wrong id")
	assert.Error(t, err)

	store.CloseDB(db)
}

func TestCreateProduct(t *testing.T) {
	db, startupErr := store.SetupTestDB("../../.env")
	assert.NoError(t, startupErr)

	sellerId, dataErr := createDummySeller(db)
	assert.NoError(t, dataErr)

	var dummyCreateProducts []data.CreateProductData = []data.CreateProductData{
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10, Condition: 0, Quantity: 3},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
			ReleasesOn: "2023-10-02 04:44:17.170135", OrderBy: "2023-10-02 04:44:17.170135",
			Discount: 10},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 100, Condition: -1, Quantity: 3},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 100, Condition: 6, Quantity: 3},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: -100, Condition: 5, Quantity: 3},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 100, Condition: 5, Quantity: 0},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
			OrderBy: "2023-10-02 04:44:17.170135", Discount: 10},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 0, Condition: 5, Quantity: 3,
			ReleasesOn: "2023-10-02 04:44:17.170135", Discount: 10}}

	//Test 1: No error, product is created
	response, err := CreateProduct(db, dummyCreateProducts[0])
	assert.Empty(t, err)
	assert.Equal(t, "Test", response.Title)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 10, response.Price)
	assert.Equal(t, int8(0), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)

	//Test 2: No error, product is created
	response, err = CreateProduct(db, dummyCreateProducts[1])
	assert.Empty(t, err)
	assert.Equal(t, "Test", response.Title)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 0, response.Price)
	assert.Equal(t, int8(5), response.Condition)
	assert.Equal(t, "Pre-Order", response.ProductType)

	//Test 3: Error, Condition is negative
	response, err = CreateProduct(db, dummyCreateProducts[2])
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 4: Error, Condition is greater than 5
	response, err = CreateProduct(db, dummyCreateProducts[3])
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 5: Error, Price is less than 0
	response, err = CreateProduct(db, dummyCreateProducts[4])
	assert.Error(t, err)
	assert.Equal(t, "Bad price data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 6: Error, Quantity is 0
	response, err = CreateProduct(db, dummyCreateProducts[5])
	assert.Error(t, err)
	assert.Equal(t, "Bad quantity data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 7: Error, Pre order release date does not exist
	response, err = CreateProduct(db, dummyCreateProducts[6])
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 8: Error, Pre order order by does not exist
	response, err = CreateProduct(db, dummyCreateProducts[7])
	assert.Error(t, err)
	assert.Equal(t, "Bad pre-order data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	store.CloseDB(db)
}

func TestGetBuyNowList(t *testing.T) {
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
	req := data.GetProductListData{SortBy: "None", MinPrice: 0, MaxPrice: 0, ProductType: "None"}
	res, err := GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(res))

	//Test 2: Get Product List Sorted by Price (Low-High)
	req = data.GetProductListData{SortBy: "price-low", MinPrice: 0, MaxPrice: 0, ProductType: "None"}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, 10, res[0].Price)
	assert.Equal(t, 90, res[1].Price)

	//Test 3: Get Product List Sorted by Price (High-Low)
	req = data.GetProductListData{SortBy: "price-high", MinPrice: 0, MaxPrice: 0, ProductType: "None"}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, 100, res[0].Price)
	assert.Equal(t, 90, res[1].Price)

	//Test 4: Get Product List Sorted by Name (A-Z)
	req = data.GetProductListData{SortBy: "name-asc", MinPrice: 0, MaxPrice: 0, ProductType: "None"}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "Test", res[0].Title)
	assert.Equal(t, "Test1", res[1].Title)

	//Test 4: Get Product List Sorted by Name (Z-A)
	req = data.GetProductListData{SortBy: "name-desc", MinPrice: 0, MaxPrice: 0, ProductType: "None"}
	res, err = GetProductList(db, req)
	assert.Empty(t, err)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, "Test2", res[0].Title)
	assert.Equal(t, "Test1", res[1].Title)

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
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10, Condition: 3, Quantity: 3},
		{Title: "Test1", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 100, Condition: 5, Quantity: 3},
		{Title: "Test2", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 90, Condition: 4, Quantity: 3}}
	var productIds []string

	for i := 0; i < len(dummyCreateProducts); i++ {
		query := `INSERT INTO products(
			title, seller_id, description, product_type, posted_date, price, condition, product_quantity) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING product_id;`
		postedDate := time.Now()
		var productId string
		err := db.QueryRowContext(
			context.Background(), query,
			dummyCreateProducts[i].Title, dummyCreateProducts[i].SellerId, dummyCreateProducts[i].Description,
			dummyCreateProducts[i].ProductType, postedDate, dummyCreateProducts[i].Price, dummyCreateProducts[i].Condition,
			dummyCreateProducts[i].Quantity).Scan(&productId)
		if err != nil {
			return nil, err
		}
		productIds = append(productIds, productId)
	}

	return productIds, nil
}
