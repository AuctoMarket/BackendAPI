package order

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

func TestDoesOrderExist(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	buyerIds := createDummyBuyers(db)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)
	orderIds, err := createDummyOrders(db, productIds[0], buyerIds[0], 10000)
	assert.NoError(t, err)

	//Test 1: OrderId exists
	orderExists := DoesOrderExist(db, orderIds[0])
	assert.Equal(t, true, orderExists)

	//Test 2: OrderId exists
	orderExists = DoesOrderExist(db, orderIds[1])
	assert.Equal(t, true, orderExists)

	//Test 3: OrderId does not exist
	orderExists = DoesOrderExist(db, "wrong id")
	assert.Equal(t, false, orderExists)

	store.CloseDB(db)
}

func TestDoesGuestOrderExist(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)
	orderIds, err := createDummyGuestOrders(db, productIds[0], "Test@gmail.com", 10000)
	assert.NoError(t, err)

	//Test 1: OrderId exists
	orderExists := DoesGuestOrderExist(db, orderIds[0])
	assert.Equal(t, true, orderExists)

	//Test 2: OrderId exists
	orderExists = DoesGuestOrderExist(db, orderIds[1])
	assert.Equal(t, true, orderExists)

	//Test 3: OrderId does not exist
	orderExists = DoesGuestOrderExist(db, "wrong id")
	assert.Equal(t, false, orderExists)

	store.CloseDB(db)
}

func TestCalculatePaymentAmount(t *testing.T) {
	//Test 1: No additional fees
	amount := calculatePaymentAmount(10000, 2, "paynow_online", "self-collection")
	assert.Equal(t, 20000, amount)
	//Test 2: Minimum order fee only
	amount = calculatePaymentAmount(1000, 2, "paynow_online", "self-collection")
	assert.Equal(t, 2100, amount)
	//Test 3: Delivery fee only
	amount = calculatePaymentAmount(10000, 2, "paynow_online", "standard_delivery")
	assert.Equal(t, 20400, amount)
	//Test 4: Delivery fee and minumum order fee
	amount = calculatePaymentAmount(1000, 2, "paynow_online", "standard_delivery")
	assert.Equal(t, 2500, amount)
	//Test 5: Card fee only
	amount = calculatePaymentAmount(10000, 2, "card", "self_collection")
	assert.Equal(t, 20400, amount)
	//Test 5: Card fee and delivery fee
	amount = calculatePaymentAmount(10000, 2, "card", "standard_delivery")
	assert.Equal(t, 20808, amount)
	//Test 5: Card fee and delivery fee and minimum order fee
	amount = calculatePaymentAmount(1000, 2, "card", "standard_delivery")
	assert.Equal(t, 2550, amount)
}

func TestValidateCreateGuestOrderRequest(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)

	dummyGuestOrderRequests := []data.CreateGuestOrderRequestData{
		{ProductId: productIds[0], Email: "test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[1], Email: "test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10400},
		{ProductId: productIds[0], Email: "test@gmail.com", OrderQuantity: 10,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 100400},
		{ProductId: productIds[0], Email: "test@gmail.com", OrderQuantity: -1,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 0},
		{ProductId: "wrong_id", Email: "test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], Email: "test@gmail.com", OrderQuantity: 1,
			PaymentType: "card", DeliveryType: "self-collect", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], Email: "test@gmail.com", OrderQuantity: 1,
			PaymentType: "card_payment", DeliveryType: "self-collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], Email: "test@gmail.com", OrderQuantity: 1,
			PaymentType: "card", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 12000},
	}

	//Test 1: No errors, no fees
	valErr := validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[0])
	assert.Empty(t, valErr)

	//Test 2: No errors, delivery fee
	valErr = validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[1])
	assert.Empty(t, valErr)

	//Test 3: To high an order quantity
	valErr = validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[2])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 4: To low an order quantity
	valErr = validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[3])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 5:  Product Id does not exist
	valErr = validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[4])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 6: Wrong delivery type
	valErr = validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[5])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 7: Wrong payment type
	valErr = validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[6])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 8: Wrong amount
	valErr = validateCreateGuestOrderRequest(db, dummyGuestOrderRequests[7])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	store.CloseDB(db)
}

func TestValidateCreateOrderRequest(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	buyerIds := createDummyBuyers(db)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)

	dummyOrderRequests := []data.CreateOrderRequestData{
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[1], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10400},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 10,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 100400},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: -1,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 0},
		{ProductId: "wrong_id", BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "standard_delivery", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "card", DeliveryType: "self-collect", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "card_payment", DeliveryType: "self-collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "card", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 12000},
		{ProductId: productIds[0], BuyerId: "wrong id", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
	}

	//Test 1: No errors, no fees
	valErr := validateCreateOrderRequest(db, dummyOrderRequests[0])
	assert.Empty(t, valErr)

	//Test 2: No errors, delivery fee
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[1])
	assert.Empty(t, valErr)

	//Test 3: To high an order quantity
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[2])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 4: To low an order quantity
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[3])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 5:  Product Id does not exist
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[4])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 6: Wrong delivery type
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[5])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 7: Wrong payment type
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[6])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 8: Wrong amount
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[7])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 9: Buyer id does not exist
	valErr = validateCreateOrderRequest(db, dummyOrderRequests[8])
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	store.CloseDB(db)
}

func TestGetGuestOrderById(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)
	guestOrderIds, err := createDummyGuestOrders(db, productIds[0], "test@gmail.com", 10000)

	//Test 1: Order Id exists
	guestOrder, getErr := GetGuestOrderById(db, guestOrderIds[0])
	assert.Empty(t, getErr)
	assert.Equal(t, productIds[0], guestOrder.ProductId)
	assert.Equal(t, "test@gmail.com", guestOrder.Email)
	assert.Equal(t, "Test", guestOrder.AddressLine1)
	assert.Equal(t, "123456", guestOrder.PostalCode)

	//Test 2: Order Id does not exist
	guestOrder, getErr = GetGuestOrderById(db, "wrong_id")
	assert.NotEmpty(t, getErr)
	assert.Equal(t, 404, getErr.ErrorCode())

	store.CloseDB(db)
}

func TestGetOrderById(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	buyerIds := createDummyBuyers(db)
	assert.NoError(t, err)
	orderIds, err := createDummyOrders(db, productIds[0], buyerIds[0], 10000)

	//Test 1: Order Id exists
	order, getErr := GetOrderById(db, orderIds[0])
	assert.Empty(t, getErr)
	assert.Equal(t, productIds[0], order.ProductId)
	assert.Equal(t, buyerIds[0], order.BuyerId)
	assert.Equal(t, "Test", order.AddressLine1)
	assert.Equal(t, "123456", order.PostalCode)

	//Test 2: Order Id does not exist
	order, getErr = GetOrderById(db, "wrong_id")
	assert.NotEmpty(t, getErr)
	assert.Equal(t, 404, getErr.ErrorCode())

	store.CloseDB(db)
}

func TestGuestCreateOrder(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)

	dummyGuestOrders := []data.CreateGuestOrderRequestData{
		{ProductId: productIds[0], Email: "Test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], Email: "Test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], Email: "Test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10003},
		{ProductId: productIds[0], Email: "Test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_qr-code", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10003},
		{ProductId: productIds[0], Email: "Test@gmail.com", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
	}

	//Test 1: No errors in guest order
	guestOrder, orderErr := CreateGuestOrder(db, dummyGuestOrders[0])
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, guestOrder)

	//Test 2: No errors in guest order
	guestOrder, orderErr = CreateGuestOrder(db, dummyGuestOrders[1])
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, guestOrder)

	//Test 3: Incorrect amount
	guestOrder, orderErr = CreateGuestOrder(db, dummyGuestOrders[2])
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 4: Incorrect Payment Type
	guestOrder, orderErr = CreateGuestOrder(db, dummyGuestOrders[3])
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 5: Incorrect Delivery Type
	guestOrder, orderErr = CreateGuestOrder(db, dummyGuestOrders[4])
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	store.CloseDB(db)

}

func TestCreateOrder(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)
	buyerIds := createDummyBuyers(db)

	dummyOrders := []data.CreateOrderRequestData{
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10003},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_qr-code", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", PostalCode: "123456",
			Amount: 10003},
		{ProductId: productIds[0], BuyerId: buyerIds[0], OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
		{ProductId: productIds[0], BuyerId: "wrong_id", OrderQuantity: 1,
			PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
			AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
			Amount: 10000},
	}

	//Test 1: No errors in guest order
	guestOrder, orderErr := CreateOrder(db, dummyOrders[0])
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, guestOrder)

	//Test 2: No errors in guest order
	guestOrder, orderErr = CreateOrder(db, dummyOrders[1])
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, guestOrder)

	//Test 3: Incorrect amount
	guestOrder, orderErr = CreateOrder(db, dummyOrders[2])
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 4: Incorrect Payment Type
	guestOrder, orderErr = CreateOrder(db, dummyOrders[3])
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 5: Incorrect Delivery Type
	guestOrder, orderErr = CreateOrder(db, dummyOrders[4])
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 6: Buyer Id does not exist
	guestOrder, orderErr = CreateOrder(db, dummyOrders[5])
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

}

func createDummyBuyers(db *sql.DB) []string {
	var dummyAccounts []data.BuyerSignUpData = []data.BuyerSignUpData{{Email: "test@gmail.com", Password: "Test1234"},
		{Email: "test2@gmail.com", Password: "Test1234"}, {Email: "test3@gmail.com", Password: "Test1234"}}

	var buyerIds []string
	for i := 0; i < len(dummyAccounts); i++ {
		var buyerId string
		query := `INSERT INTO buyers(email, password) VALUES ($1,$2) RETURNING buyer_id;`
		hashedPwd, _ := utils.HashAndSalt([]byte(dummyAccounts[i].Password))
		db.QueryRowContext(context.Background(), query, dummyAccounts[i].Email, hashedPwd).Scan(&buyerId)
		buyerIds = append(buyerIds, buyerId)
	}

	return buyerIds
}

func createDummySeller(db *sql.DB) (string, error) {
	var sellerId string
	query := `INSERT INTO sellers(email, seller_name, password) VALUES ('test@gmail.com','test','test') RETURNING seller_id`
	err := db.QueryRowContext(context.Background(), query).Scan(&sellerId)

	return sellerId, err
}

func createDummyProducts(db *sql.DB, sellerId string) ([]string, error) {
	var dummyCreateProducts []data.ProductCreateData = []data.ProductCreateData{
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10000, Condition: 3, Quantity: 3},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10000, Condition: 5, Quantity: 3},
		{Title: "Test", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Pre-Order", Price: 9000, Condition: 4, Quantity: 3}}
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

func createDummyOrders(db *sql.DB, productId string, buyerId string, price int) ([]string, error) {
	var orderIds []string
	dummyOrders := []data.CreateOrderRequestData{{ProductId: productId, BuyerId: buyerId, OrderQuantity: 1,
		PaymentType: "card", DeliveryType: "self-collect", PhoneNumber: "12345678",
		AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
		Amount: price}, {ProductId: productId, BuyerId: buyerId, OrderQuantity: 1,
		PaymentType: "paynow_online", DeliveryType: "self-collect", PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456",
		Amount: price}}
	//SQL Query to insert new order
	query := `INSERT INTO orders(
		product_id, buyer_id, delivery_type, order_quantity, payment_type, 
		payment_status, phone_number, order_date, address_line_1, 
		address_line_2, postal_code) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) 
		RETURNING order_id;`

	for i := 0; i < len(dummyOrders); i++ {
		var orderId string
		request := dummyOrders[i]
		orderDate := time.Now()
		paymentStatus := "pending"
		err := db.QueryRowContext(context.Background(), query, request.ProductId, request.BuyerId, request.DeliveryType,
			request.OrderQuantity, request.PaymentType, paymentStatus, request.PhoneNumber, orderDate,
			request.AddressLine1, utils.NewNullableString(request.AddressLine2), request.PostalCode).Scan(&orderId)

		if err != nil {
			return nil, err
		}

		orderIds = append(orderIds, orderId)
	}

	return orderIds, nil
}

func createDummyGuestOrders(db *sql.DB, productId string, email string, price int) ([]string, error) {
	var guestOrderIds []string
	dummyOrders := []data.CreateGuestOrderRequestData{{ProductId: productId, Email: email, OrderQuantity: 1,
		PaymentType: "card", DeliveryType: "self_collection", PhoneNumber: "12345678",
		AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
		Amount: price}, {ProductId: productId, Email: email, OrderQuantity: 1,
		PaymentType: "paynow_online", DeliveryType: "self_collection", PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456",
		Amount: price}}
	//SQL Query to insert new order
	query := `INSERT INTO guest_orders(
		product_id, email, delivery_type, order_quantity, payment_type, 
		payment_status, phone_number, order_date, address_line_1, 
		address_line_2, postal_code) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) 
		RETURNING guest_order_id;`

	for i := 0; i < len(dummyOrders); i++ {
		var guestOrderId string
		request := dummyOrders[i]
		orderDate := time.Now()
		paymentStatus := "pending"
		err := db.QueryRowContext(context.Background(), query, request.ProductId, request.Email, request.DeliveryType,
			request.OrderQuantity, request.PaymentType, paymentStatus, request.PhoneNumber, orderDate,
			request.AddressLine1, utils.NewNullableString(request.AddressLine2), request.PostalCode).Scan(&guestOrderId)

		if err != nil {
			return nil, err
		}

		guestOrderIds = append(guestOrderIds, guestOrderId)
	}

	return guestOrderIds, nil
}
