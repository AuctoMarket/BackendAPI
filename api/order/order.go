package order

import (
	"BackendAPI/api/buyer"
	"BackendAPI/api/product"
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"time"
)

/*
Create a order for a specific product by a specific buyer from the order request and store it in the database. If the
buyer or product do not exist, return a BadRequestError (400).
*/
func CreateOrder(db *sql.DB, request data.CreateOrderRequestData) (data.CreateOrderResponseData, *utils.ErrorHandler) {
	var response data.CreateOrderResponseData
	orderDate := time.Now()
	paymentStatus := `pending`

	//validate input
	validErr := validateCreateOrderRequest(db, request)
	if validErr != nil {
		return response, validErr
	}

	//SQL Query to insert new order
	query := `INSERT INTO orders(
		product_id, buyer_id, delivery_type, order_quantity, payment_type, 
		payment_status, phone_number, order_date, address_line_1, 
		address_line_2, postal_code) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) 
		RETURNING order_id;`

	//SQL query execution dependant on if address_line_2 exists
	err := db.QueryRowContext(
		context.Background(), query,
		request.ProductId, request.BuyerId, request.DeliveryType, request.OrderQuantity,
		request.PaymentType, paymentStatus, request.PhoneNumber, orderDate, request.AddressLine1,
		utils.NewNullableString(request.AddressLine2), request.PostalCode).Scan(&response.OrderId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in inserting Order rows")
		return response, errResp
	}

	response.RedirectUrl = "example.com"
	return response, nil
}

/*
Create a order for a specific product by guest account from the order request and store it in the database. If the
product does not exist, return a BadRequestError (400).
*/
func CreateGuestOrder(db *sql.DB, request data.CreateGuestOrderRequestData) (data.CreateGuestOrderResponseData, *utils.ErrorHandler) {
	var response data.CreateGuestOrderResponseData
	orderDate := time.Now()
	paymentStatus := `pending`

	//validate input
	validErr := validateCreateGuestOrderRequest(db, request)
	if validErr != nil {
		return response, validErr
	}

	//SQL Query to insert new guest order
	query := `INSERT INTO guest_orders(
		product_id, delivery_type, order_quantity, payment_type, 
		email, payment_status, phone_number, order_date, 
		address_line_1, address_line_2, postal_code) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) 
		RETURNING guest_order_id;`

	err := db.QueryRowContext(
		context.Background(), query,
		request.ProductId, request.DeliveryType, request.OrderQuantity,
		request.PaymentType, request.Email, paymentStatus, request.PhoneNumber, orderDate,
		request.AddressLine1, utils.NewNullableString(request.AddressLine2), request.PostalCode).Scan(&response.GuestOrderId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in inserting Order rows")
		return response, errResp
	}

	response.RedirectUrl = "example.com"
	return response, nil
}

/*
Gets the order by its id, if orderid does not exist returns a 404 Error
*/
func GetOrderById(db *sql.DB, orderId string) (data.GetOrderByIdResponseData, *utils.ErrorHandler) {
	var response data.GetOrderByIdResponseData

	if !DoesOrderExist(db, orderId) {
		return response, utils.NotFoundError("Order with given id does not exist")
	}

	query := `SELECT product_id, buyer_id, order_quantity, payment_type, delivery_type,
	phone_number, address_line_1, COALESCE(address_line_2, ''), postal_code
	FROM orders WHERE order_id=$1;`

	err := db.QueryRowContext(context.Background(), query, orderId).Scan(
		&response.ProductId, &response.BuyerId, &response.OrderQuantity, &response.PaymentType,
		&response.DeliveryType, &response.PhoneNumber, &response.AddressLine1, &response.AddressLine2,
		&response.PostalCode)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting Order rows")
		return response, errResp
	}

	response.OrderId = orderId
	return response, nil
}

/*
Gets the guest order by its id, if orderid does not exist returns a 404 Error
*/
func GetGuestOrderById(db *sql.DB, guestOrderId string) (data.GetGuestOrderByIdResponseData, *utils.ErrorHandler) {
	var response data.GetGuestOrderByIdResponseData

	if !DoesGuestOrderExist(db, guestOrderId) {
		return response, utils.NotFoundError("Guest order with given id does not exist")
	}

	query := `SELECT product_id, email, order_quantity, payment_type, delivery_type,
	phone_number, address_line_1, COALESCE(address_line_2, ''), postal_code
	FROM guest_orders WHERE guest_order_id=$1;`

	err := db.QueryRowContext(context.Background(), query, guestOrderId).Scan(
		&response.ProductId, &response.Email, &response.OrderQuantity, &response.PaymentType,
		&response.DeliveryType, &response.PhoneNumber, &response.AddressLine1, &response.AddressLine2,
		&response.PostalCode)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting Guest Order rows")
		return response, errResp
	}

	response.GuestOrderId = guestOrderId
	return response, nil
}

/*
Validate a create order request to ensure params are correct.
*/
func validateCreateOrderRequest(db *sql.DB, request data.CreateOrderRequestData) *utils.ErrorHandler {
	if !product.DoesProductExist(db, request.ProductId) {
		utils.LogMessage("Product with given id does not exist")
		return utils.BadRequestError("Bad product_id data")
	}

	if !buyer.DoesBuyerExist(db, request.BuyerId) {
		utils.LogMessage("Buyer with given id does not exist")
		return utils.BadRequestError("Bad buyer_id data")
	}

	var price int
	var availableStock int
	query := `SELECT price, (product_quantity - sold_quantity) FROM products WHERE product_id = $1;`
	err := db.QueryRowContext(context.Background(), query, request.ProductId).Scan(&price, &availableStock)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting product rows")
		return errResp
	}

	if request.OrderQuantity <= 0 || request.OrderQuantity > availableStock {
		utils.LogMessage("Order quanity is invalid")
		return utils.BadRequestError("Bad order_quantity data")
	}

	if request.PaymentType != "card" && request.PaymentType != "paynow_online" {
		utils.LogMessage("Payment Type is invalid")
		return utils.BadRequestError("Bad payment_type data")
	}

	if request.DeliveryType != "standard_delivery" && request.DeliveryType != "self_collection" {
		utils.LogMessage("Delivery Type is invalid")
		return utils.BadRequestError("Bad delivery_type data")
	}

	if request.Amount != calculatePaymentAmount(price, request.OrderQuantity, request.PaymentType, request.DeliveryType) {
		utils.LogMessage("Payment amount is incorrect")
		return utils.BadRequestError("Bad amount data")
	}

	return nil
}

/*
Validate a create order request to ensure params are correct.
*/
func validateCreateGuestOrderRequest(db *sql.DB, request data.CreateGuestOrderRequestData) *utils.ErrorHandler {
	if !product.DoesProductExist(db, request.ProductId) {
		utils.LogMessage("Product with given id does not exist")
		return utils.BadRequestError("Bad product_id data")
	}

	var price int
	var availableStock int
	query := `SELECT price, (product_quantity - sold_quantity) FROM products WHERE product_id = $1;`
	err := db.QueryRowContext(context.Background(), query, request.ProductId).Scan(&price, &availableStock)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting product rows")
		return errResp
	}

	if request.OrderQuantity <= 0 || request.OrderQuantity > availableStock {
		utils.LogMessage("Order quanity is invalid")
		return utils.BadRequestError("Bad order_quantity data")
	}

	if request.PaymentType != "card" && request.PaymentType != "paynow_online" {
		utils.LogMessage("Payment Type is invalid")
		return utils.BadRequestError("Bad payment_type data")
	}

	if request.DeliveryType != "standard_delivery" && request.DeliveryType != "self_collection" {
		utils.LogMessage("Delivery Type is invalid")
		return utils.BadRequestError("Bad delivery_type data")
	}

	if request.Amount != calculatePaymentAmount(price, request.OrderQuantity, request.PaymentType, request.DeliveryType) {
		utils.LogMessage("Payment amount is incorrect")
		return utils.BadRequestError("Bad amount data")
	}

	return nil
}

/*
Calculate the payment amount given an order request
*/
func calculatePaymentAmount(price int, quantity int, paymentType string, deliveryType string) int {
	//calculate product cost
	amountToBePaid := quantity * price

	//calculate small order fee
	if amountToBePaid < 2500 {
		amountToBePaid += 100
	}

	//calculate delivery fee
	if deliveryType == "standard_delivery" {
		amountToBePaid += 400
	}

	//calculate payment fee
	var paymentFee int
	if paymentType == "card" {
		paymentFee = (amountToBePaid * 2) / 100
	}

	amountToBePaid += paymentFee
	return amountToBePaid
}

/*
Checks wether a Order with a given product id already exists in the database
and returns true if it does false otherwise.
*/
func DoesOrderExist(db *sql.DB, orderId string) bool {
	var orderExists bool
	query := `SELECT EXISTS(SELECT * FROM orders WHERE order_id = $1);`
	err := db.QueryRowContext(context.Background(), query, orderId).Scan(&orderExists)

	if err != nil {
		return false
	}

	return orderExists
}

/*
Checks wether a Guest Order with a given product id already exists in the database
and returns true if it does false otherwise.
*/
func DoesGuestOrderExist(db *sql.DB, guestOrderId string) bool {
	var guestOrderExists bool
	query := `SELECT EXISTS(SELECT * FROM guest_orders WHERE guest_order_id = $1);`
	err := db.QueryRowContext(context.Background(), query, guestOrderId).Scan(&guestOrderExists)

	if err != nil {
		return false
	}

	return guestOrderExists
}
