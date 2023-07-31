package product

import (
	"BackendAPI/data"
	"BackendAPI/store"
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateProduct(t *testing.T) {

	//Test 1: No errors, product is valid
	testCreateProduct1 := data.ProductCreateData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 3}
	err := validateCreateProduct(testCreateProduct1)
	assert.Empty(t, err)

	//Test 2: Error, product price is invalid
	testCreateProduct2 := data.ProductCreateData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-Now", Price: -5, Condition: 3}
	err = validateCreateProduct(testCreateProduct2)
	assert.Error(t, err)
	assert.Equal(t, "Bad price data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 3: Error, product condition is invalid
	testCreateProduct3 := data.ProductCreateData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: 13}
	err = validateCreateProduct(testCreateProduct3)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 4: Error, product condition is invalid
	testCreateProduct4 := data.ProductCreateData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-Now", Price: 10, Condition: -1}
	err = validateCreateProduct(testCreateProduct4)
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 5: No errors, product is valid
	testCreateProduct5 := data.ProductCreateData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Pre-Order", Price: 0, Condition: 5}
	err = validateCreateProduct(testCreateProduct5)
	assert.Empty(t, err)

	//Test 6: No errors, product is valid
	testCreateProduct6 := data.ProductCreateData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-Now", Price: 0, Condition: 0}
	err = validateCreateProduct(testCreateProduct6)
	assert.Empty(t, err)

	//Test 7:Error product type is wrong
	testCreateProduct7 := data.ProductCreateData{
		Title: "test", SellerId: "test", Description: "This is a test description",
		ProductType: "Buy-It-Now", Price: 0, Condition: 0}
	err = validateCreateProduct(testCreateProduct7)
	assert.Error(t, err)
	assert.Equal(t, "Bad product_type data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())
}

func TestDoesProductExist(t *testing.T) {
	db, err := store.SetupTestDB()
	assert.NoError(t, err)

	productIds := createDummyProducts(db)

	//Test 1: Product Id Exists
	doesExist := doesProductExist(db, productIds[0])
	assert.Equal(t, true, doesExist)

	//Test 2: Product Id Exists
	doesExist = doesProductExist(db, productIds[1])
	assert.Equal(t, true, doesExist)

	//Test 3: Product Id Exists
	doesExist = doesProductExist(db, productIds[2])
	assert.Equal(t, true, doesExist)

	//Test 4: Product Id does not exist
	doesExist = doesProductExist(db, "wrong id")
	assert.Equal(t, false, doesExist)

	//Test 4: Product Id does not exist
	doesExist = doesProductExist(db, productIds[0]+"1")
	assert.Equal(t, false, doesExist)

	store.CloseDB(db)
}

func TestGetProductById(t *testing.T) {
	db, err := store.SetupTestDB()
	assert.NoError(t, err)

	productIds := createDummyProducts(db)

	//Test 1: Product Id exists
	response, err := GetProductById(db, productIds[0])
	assert.Empty(t, err)
	assert.Equal(t, "Test", response.Title)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 10, response.Price)
	assert.Equal(t, int8(3), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)

	//Test 2: Product Id exists
	response, err = GetProductById(db, productIds[1])
	assert.Empty(t, err)
	assert.Equal(t, "Test", response.Title)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 100, response.Price)
	assert.Equal(t, int8(5), response.Condition)
	assert.Equal(t, "Buy-Now", response.ProductType)

	//Test 3: Product Id exists
	response, err = GetProductById(db, productIds[2])
	assert.Empty(t, err)
	assert.Equal(t, "Test", response.Title)
	assert.Equal(t, "This is a test description", response.Description)
	assert.Equal(t, 90, response.Price)
	assert.Equal(t, int8(4), response.Condition)
	assert.Equal(t, "Pre-Order", response.ProductType)

	//Test 4: Product Id exists
	response, err = GetProductById(db, "wrong id")
	assert.Error(t, err)

}

func TestCreateProduct(t *testing.T) {
	db, startupErr := store.SetupTestDB()
	assert.NoError(t, startupErr)

	sellerId := createDummySeller(db)
	var dummyCreateProducts []data.ProductCreateData = []data.ProductCreateData{
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10, Condition: 0},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 0, Condition: 5},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 100, Condition: -1},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 100, Condition: 6},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: -100, Condition: 5}}

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

	//Test 3: Error, Condition is greater than 5
	response, err = CreateProduct(db, dummyCreateProducts[3])
	assert.Error(t, err)
	assert.Equal(t, "Bad condition data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())

	//Test 3: Error, Price is less than 0
	response, err = CreateProduct(db, dummyCreateProducts[4])
	assert.Error(t, err)
	assert.Equal(t, "Bad price data", err.Error())
	assert.Equal(t, 400, err.ErrorCode())
}

func createDummySeller(db *sql.DB) string {
	var sellerId string
	query := `INSERT INTO sellers(email, password) VALUES ('test@gmail.com','test') RETURNING seller_id`
	db.QueryRowContext(context.Background(), query).Scan(&sellerId)

	return sellerId
}

func createDummyProducts(db *sql.DB) []string {
	sellerId := createDummySeller(db)

	var dummyCreateProducts []data.ProductCreateData = []data.ProductCreateData{
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10, Condition: 3},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 100, Condition: 5},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 90, Condition: 4}}
	var productIds []string

	for i := 0; i < len(dummyCreateProducts); i++ {
		query := `INSERT INTO products(
			title, seller_id, description, product_type, posted_date, price, condition) 
			VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING product_id;`
		postedDate := time.Now()
		var productId string
		db.QueryRowContext(
			context.Background(), query,
			dummyCreateProducts[i].Title, dummyCreateProducts[i].SellerId, dummyCreateProducts[i].Description,
			dummyCreateProducts[i].ProductType, postedDate, dummyCreateProducts[i].Price, dummyCreateProducts[i].Condition).Scan(&productId)
		productIds = append(productIds, productId)
	}

	return productIds
}
