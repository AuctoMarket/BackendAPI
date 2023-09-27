package order

import (
	"BackendAPI/data"
	"BackendAPI/store"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"strconv"
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
	orderIds, err := createDummyOrders(db, productIds, buyerIds[0])
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
	orderIds, err := createDummyGuestOrders(db, productIds, "test@aucto.io")
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

func TestValidatePaymentAmount(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)

	//Test 1: No additional fees
	var fees data.OrderFees = data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
		PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0, TotalPaid: 20000}
	var products []data.ProductOrder = []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}
	err = validatePaymentAmount(db, products, fees)
	assert.Empty(t, err)
	//Test 2: Minimum order fee only
	fees = data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
		PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 100, TotalPaid: 2100}
	products = []data.ProductOrder{{ProductId: productIds[5], OrderQuantity: 2}}
	err = validatePaymentAmount(db, products, fees)
	assert.Empty(t, err)
	//Test 3: Delivery fee only
	fees = data.OrderFees{PaymentType: "paynow_online", DeliveryType: "standard_delivery",
		PaymentFee: 0, DeliveryFee: 400, SmallOrderFee: 0, TotalPaid: 20400}
	products = []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}
	err = validatePaymentAmount(db, products, fees)
	assert.Empty(t, err)
	//Test 4: Delivery fee and minumum order fee
	fees = data.OrderFees{PaymentType: "paynow_online", DeliveryType: "standard_delivery",
		PaymentFee: 0, DeliveryFee: 400, SmallOrderFee: 100, TotalPaid: 2500}
	products = []data.ProductOrder{{ProductId: productIds[5], OrderQuantity: 2}}
	err = validatePaymentAmount(db, products, fees)
	assert.Empty(t, err)
	//Test 5: Card fee only
	fees = data.OrderFees{PaymentType: "card", DeliveryType: "self_collection",
		PaymentFee: 400, DeliveryFee: 0, SmallOrderFee: 0, TotalPaid: 20400}
	products = []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}
	err = validatePaymentAmount(db, products, fees)
	assert.Empty(t, err)
	//Test 5: Card fee and delivery fee
	fees = data.OrderFees{PaymentType: "card", DeliveryType: "standard_delivery",
		DeliveryFee: 400, PaymentFee: 408, TotalPaid: 20808}
	products = []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}
	err = validatePaymentAmount(db, products, fees)
	assert.Empty(t, err)
	//Test 5: Card fee and delivery fee and minimum order fee
	fees = data.OrderFees{PaymentType: "card", DeliveryType: "standard_delivery",
		DeliveryFee: 400, SmallOrderFee: 100, PaymentFee: 50, TotalPaid: 2550}
	products = []data.ProductOrder{{ProductId: productIds[5], OrderQuantity: 2}}
	err = validatePaymentAmount(db, products, fees)
	assert.Empty(t, err)
}

func TestValidateCreateGuestOrderRequest(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)

	//Test 1: No errors, no fees
	order := data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}},
		Email:    "test@aucto.io", PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 20000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}}
	valErr := validateCreateGuestOrderRequest(db, order)
	assert.Empty(t, valErr)

	//Test 2: No errors, delivery fee
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}},
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "standard_delivery",
			TotalPaid: 20400, PaymentFee: 0, DeliveryFee: 400, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.Empty(t, valErr)

	//Test 3: To high an order quantity
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 100}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 1000000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 4: To low an order quantity
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: -1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 100000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 5:  Product Id does not exist
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: "wrong id", OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 6: Wrong delivery type
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self-collect",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 7: Wrong payment type
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "card-payment", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 8: Correct amount Card Payment
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "card", DeliveryType: "self_collection",
			TotalPaid: 10200, PaymentFee: 200, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.Empty(t, valErr)

	//Test 9: Correct Amount Discount
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[4], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 9000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateGuestOrderRequest(db, order)
	assert.Empty(t, valErr)

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

	//Test 1: No errors, no fees
	order := data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 20000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr := validateCreateOrderRequest(db, order)
	assert.Empty(t, valErr)

	//Test 2: No errors, delivery fee
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "standard_delivery",
			TotalPaid: 20400, PaymentFee: 0, DeliveryFee: 400, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.Empty(t, valErr)

	//Test 3: To high an order quantity
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 100}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 100000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 4: To low an order quantity
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: -1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 100000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 5:  Product Id does not exist
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: "wrong id", OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 6: Wrong delivery type
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self-collect",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 7: Wrong payment type
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "card-payment", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.NotEmpty(t, valErr)
	assert.Equal(t, 400, valErr.ErrorCode())

	//Test 8: Correct amount Card Payment
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[1], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "card", DeliveryType: "self_collection",
			TotalPaid: 10200, PaymentFee: 200, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.Empty(t, valErr)

	//Test 9: Correct Amount Discount
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[4], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 9000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
	assert.Empty(t, valErr)

	//Test 11: Buyer Id does not exist
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[4], OrderQuantity: 1}}, BuyerId: "wrong id",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 9000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	valErr = validateCreateOrderRequest(db, order)
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
	guestOrderIds, err := createDummyGuestOrders(db, productIds, "test@aucto.io")

	//Test 1: Order Id exists
	guestOrder, getErr := GetGuestOrderById(db, guestOrderIds[0])
	assert.Empty(t, getErr)
	assert.Equal(t, productIds[0], guestOrder.Products[0].ProductId)
	assert.Equal(t, "test@aucto.io", guestOrder.Email)
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
	orderIds, err := createDummyOrders(db, productIds, buyerIds[0])

	//Test 1: Order Id exists
	order, getErr := GetOrderById(db, orderIds[0])
	assert.Empty(t, getErr)
	assert.Equal(t, productIds[0], order.Products[0].ProductId)
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

	//Test 1: No errors in guest order
	order := data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 20000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr := CreateGuestOrder(db, order)
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, response)

	//Test 2: No errors in guest order
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[1], OrderQuantity: 2}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 20000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456"}
	response, orderErr = CreateGuestOrder(db, order)
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, response)

	//Test 3: Incorrect amount
	order = data.CreateGuestOrderRequestData{Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 10003, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456"}
	response, orderErr = CreateGuestOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 4: Incorrect Payment Type
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_qr", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456"}
	response, orderErr = CreateGuestOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 5: Incorrect Delivery Type
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr = CreateGuestOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 6: Incorrect amount Discount
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[4], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr = CreateGuestOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 7: Correct Amount Discount
	order = data.CreateGuestOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[4], OrderQuantity: 1}}, Email: "test@aucto.io",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 9000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr = CreateGuestOrder(db, order)
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, response)

	store.CloseDB(db)
}

func TestCreateOrder(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	buyerIds := createDummyBuyers(db)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)

	//Test 1: No errors in order
	order := data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}, {ProductId: productIds[1], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 20000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr := CreateOrder(db, order)
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, response)

	//Test 2: No errors in order
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[1], OrderQuantity: 2}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 20000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456"}
	response, orderErr = CreateOrder(db, order)
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, response)

	//Test 3: Incorrect amount
	order = data.CreateOrderRequestData{Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 10003, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456"}
	response, orderErr = CreateOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 4: Incorrect Payment Type
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_qr", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", PostalCode: "123456"}
	response, orderErr = CreateOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 5: Incorrect Delivery Type
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0}, PhoneNumber: "12345678",
		AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr = CreateOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 6: Incorrect amount Discount
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[4], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr = CreateOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	//Test 7: Correct Amount Discount
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[4], OrderQuantity: 1}}, BuyerId: buyerIds[0],
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 9000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr = CreateOrder(db, order)
	assert.Empty(t, orderErr)
	assert.NotEmpty(t, response)

	//Test 8: Buyer Id does not exist
	order = data.CreateOrderRequestData{
		Products: []data.ProductOrder{{ProductId: productIds[0], OrderQuantity: 1}}, BuyerId: "wrong id",
		Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
			TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
		PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"}
	response, orderErr = CreateOrder(db, order)
	assert.NotEmpty(t, orderErr)
	assert.Equal(t, 400, orderErr.ErrorCode())

	store.CloseDB(db)

}

func TestUpdateOrderPaymentStatus(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	buyerIds := createDummyBuyers(db)
	assert.NoError(t, err)
	orderIds, err := createDummyOrders(db, productIds, buyerIds[0])

	//Test 1: Order status is completed
	testErr := UpdateOrderPaymentStatus(db, orderIds[0], data.PaymentValidationRequestData{Status: "completed"})
	assert.Empty(t, testErr)

	//Test 2: Order id does not exist
	testErr = UpdateOrderPaymentStatus(db, "wrong id", data.PaymentValidationRequestData{Status: "completed"})
	assert.NotEmpty(t, testErr)

	store.CloseDB(db)
}

func TestUpdateGuestOrderPaymentStatus(t *testing.T) {
	db, err := store.SetupTestDB("../../.env")
	assert.NoError(t, err)
	sellerId, err := createDummySeller(db)
	assert.NoError(t, err)
	productIds, err := createDummyProducts(db, sellerId)
	assert.NoError(t, err)
	guestOrderIds, err := createDummyGuestOrders(db, productIds, "test@aucto.io")

	//Test 1: Order status is completed
	testErr := UpdateGuestOrderPaymentStatus(db, guestOrderIds[0], data.PaymentValidationRequestData{Status: "completed"})
	assert.Empty(t, testErr)

	//Test 2: Order id does not exist
	testErr = UpdateGuestOrderPaymentStatus(db, "wrong id", data.PaymentValidationRequestData{Status: "completed"})
	assert.NotEmpty(t, testErr)

	store.CloseDB(db)
}

func createDummyBuyers(db *sql.DB) []string {
	var dummyAccounts []data.BuyerSignUpData = []data.BuyerSignUpData{{Email: "test@aucto.io", Password: "Test1234"},
		{Email: "test2@aucto.io", Password: "Test1234"}, {Email: "test3@aucto.io", Password: "Test1234"}}

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
			ProductType: "Buy-Now", Price: 10000, Condition: 5, Quantity: 3, Language: "Jap",
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
		{Title: "Test3", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 10000, Condition: 4, Quantity: 3, Language: "Eng",
			Expansion: "Test", Discount: 1000},
		//Product 6: Cheap products
		{Title: "Test6", SellerId: sellerId, Description: "This is a test description",
			ProductType: "Buy-Now", Price: 1000, Condition: 4, Quantity: 3, Language: "Eng",
			Expansion: "Test"}}
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

func createDummyOrders(db *sql.DB, products []string, buyerId string) ([]string, error) {
	var orderIds []string
	dummyOrders := []data.CreateOrderRequestData{
		//Order 1: Basic Order with Paynow_online, Self-Collection
		{Products: []data.ProductOrder{{ProductId: products[0], OrderQuantity: 1}}, BuyerId: buyerId,
			Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
				TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
			PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"},
		//Order 2: Basic Order with Paynow_online, Self-Collection and discount
		{Products: []data.ProductOrder{{ProductId: products[4], OrderQuantity: 1}}, BuyerId: buyerId,
			Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
				TotalPaid: 9000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
			PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"},
		//Order 3: Basic Pre-Order with Paynow_online, Self-Collection
		{Products: []data.ProductOrder{{ProductId: products[3], OrderQuantity: 1}}, BuyerId: buyerId,
			Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
				TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
			PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"},
	}
	//SQL Query to insert new order
	query := `INSERT INTO orders(
		buyer_id, delivery_type, delivery_fee, payment_type, payment_fee, small_order_fee, total_paid,
		phone_number, order_date, address_line_1, address_line_2, postal_code, telegram_handle) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) 
		RETURNING order_id;`

	for i := 0; i < len(dummyOrders); i++ {
		var orderId string
		request := dummyOrders[i]
		orderDate := time.Now()
		err := db.QueryRowContext(context.Background(), query, request.BuyerId, request.Fees.DeliveryType, request.Fees.DeliveryFee,
			request.Fees.PaymentType, request.Fees.PaymentFee, request.Fees.SmallOrderFee, request.Fees.TotalPaid,
			request.PhoneNumber, orderDate, request.AddressLine1, utils.NewNullableString(request.AddressLine2),
			request.PostalCode, utils.NewNullableString(request.TelegramHandle)).Scan(&orderId)

		if err != nil {
			return nil, err
		}

		query2 := `INSERT INTO order_products(product_id, order_id, quantity) VALUES`
		for i := 0; i < len(request.Products); i++ {
			query2 += `('` + request.Products[i].ProductId + `','` + orderId + `',` + strconv.Itoa(request.Products[i].OrderQuantity) + `)`
			if i < len(request.Products)-1 {
				query2 += `,`
			}
		}
		_, err = db.ExecContext(context.Background(), query2)

		if err != nil {
			return nil, err
		}

		orderIds = append(orderIds, orderId)
	}

	return orderIds, nil
}

func createDummyGuestOrders(db *sql.DB, products []string, email string) ([]string, error) {
	var guestOrderIds []string
	dummyOrders := []data.CreateGuestOrderRequestData{
		//Order 1: Basic Order with Paynow_online, Self-Collection
		{Products: []data.ProductOrder{{ProductId: products[0], OrderQuantity: 1}}, Email: email,
			Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
				TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
			PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"},
		//Order 2: Basic Order with Paynow_online, Self-Collection and discount
		{Products: []data.ProductOrder{{ProductId: products[4], OrderQuantity: 1}}, Email: email,
			Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
				TotalPaid: 9000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
			PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"},
		//Order 3: Basic Pre-Order with Paynow_online, Self-Collection
		{Products: []data.ProductOrder{{ProductId: products[3], OrderQuantity: 1}}, Email: email,
			Fees: data.OrderFees{PaymentType: "paynow_online", DeliveryType: "self_collection",
				TotalPaid: 10000, PaymentFee: 0, DeliveryFee: 0, SmallOrderFee: 0},
			PhoneNumber: "12345678", AddressLine1: "Test", AddressLine2: "Test", PostalCode: "123456"},
	}

	//SQL Query to insert new order
	query := `INSERT INTO guest_orders(
		email, delivery_type, delivery_fee, payment_type, payment_fee, small_order_fee, total_paid,
		phone_number, order_date, address_line_1, address_line_2, postal_code, telegram_handle) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) 
		RETURNING guest_order_id;`

	for i := 0; i < len(dummyOrders); i++ {
		var guestOrderId string
		request := dummyOrders[i]
		orderDate := time.Now()
		err := db.QueryRowContext(context.Background(), query, request.Email, request.Fees.DeliveryType,
			request.Fees.DeliveryFee, request.Fees.PaymentType, request.Fees.PaymentFee, request.Fees.SmallOrderFee,
			request.Fees.TotalPaid, request.PhoneNumber, orderDate, request.AddressLine1,
			utils.NewNullableString(request.AddressLine2), request.PostalCode, utils.NewNullableString(request.TelegramHandle)).Scan(&guestOrderId)

		if err != nil {
			return nil, err
		}

		query2 := `INSERT INTO guest_order_products(product_id, guest_order_id, quantity) VALUES`
		for i := 0; i < len(request.Products); i++ {
			query2 += `('` + request.Products[i].ProductId + `','` + guestOrderId + `',` + strconv.Itoa(request.Products[i].OrderQuantity) + `)`
			if i < len(request.Products)-1 {
				query2 += `,`
			}
		}
		_, err = db.ExecContext(context.Background(), query2)

		if err != nil {
			return nil, err
		}

		guestOrderIds = append(guestOrderIds, guestOrderId)
	}

	return guestOrderIds, nil
}
