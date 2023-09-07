package order

import (
	"BackendAPI/api/buyer"
	"BackendAPI/api/product"
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"
)

/*
Create a order for a specific product by a specific buyer from the order request and store it in the database. If the
buyer or product do not exist, return a BadRequestError (400).
*/
func CreateOrder(db *sql.DB, request data.CreateOrderRequestData) (data.CreateOrderResponseData, *utils.ErrorHandler) {
	var response data.CreateOrderResponseData
	orderDate := time.Now()

	//validate input
	validErr := validateCreateOrderRequest(db, request)
	if validErr != nil {
		return response, validErr
	}

	//SQL Query to insert new order
	query := `INSERT INTO orders(
		product_id, buyer_id, delivery_type, delivery_fee, payment_type, payment_fee, small_order_fee, total_paid,
		order_quantity, phone_number, order_date, address_line_1, address_line_2, postal_code, telegram_handle, 
		product_price) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) 
		RETURNING order_id;`

	//SQL query execution dependant on if address_line_2 exists
	err := db.QueryRowContext(
		context.Background(), query,
		request.ProductId, request.BuyerId, request.Fees.DeliveryType, request.Fees.DeliveryFee,
		request.Fees.PaymentType, request.Fees.PaymentFee, request.Fees.SmallOrderFee, request.Fees.TotalPaid,
		request.OrderQuantity, request.PhoneNumber, orderDate, request.AddressLine1, utils.NewNullableString(request.AddressLine2),
		request.PostalCode, utils.NewNullableString(request.TelegramHandle), request.Fees.ProductPrice).Scan(&response.OrderId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in inserting Order rows")
		return response, errResp
	}

	paymentResponse, paymentErr := CreatePaymentRequest(float64(request.Fees.TotalPaid)/100, response.OrderId, request.Fees.PaymentType, false)
	response.RedirectUrl = paymentResponse.Url
	return response, paymentErr
}

/*
Create a order for a specific product by guest account from the order request and store it in the database. If the
product does not exist, return a BadRequestError (400).
*/
func CreateGuestOrder(db *sql.DB, request data.CreateGuestOrderRequestData) (data.CreateGuestOrderResponseData, *utils.ErrorHandler) {
	var response data.CreateGuestOrderResponseData
	orderDate := time.Now()

	//validate input
	validErr := validateCreateGuestOrderRequest(db, request)
	if validErr != nil {
		return response, validErr
	}

	//SQL Query to insert new guest order
	query := `INSERT INTO guest_orders(
		product_id, email, delivery_type, delivery_fee, payment_type, payment_fee, small_order_fee, total_paid,
		order_quantity, phone_number, order_date, address_line_1, address_line_2, postal_code, telegram_handle, 
		product_price) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) 
		RETURNING guest_order_id;`

	err := db.QueryRowContext(
		context.Background(), query,
		request.ProductId, request.Email, request.Fees.DeliveryType, request.Fees.DeliveryFee,
		request.Fees.PaymentType, request.Fees.PaymentFee, request.Fees.SmallOrderFee, request.Fees.TotalPaid,
		request.OrderQuantity, request.PhoneNumber, orderDate, request.AddressLine1, utils.NewNullableString(request.AddressLine2),
		request.PostalCode, utils.NewNullableString(request.TelegramHandle), request.Fees.ProductPrice).Scan(&response.GuestOrderId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in inserting Order rows")
		return response, errResp
	}

	paymentResponse, paymentErr := CreatePaymentRequest(float64(request.Fees.TotalPaid)/100, response.GuestOrderId, request.Fees.PaymentType, true)
	response.RedirectUrl = paymentResponse.Url
	return response, paymentErr
}

/*
Gets the order by its id, if orderid does not exist returns a 404 Error
*/
func GetOrderById(db *sql.DB, orderId string) (data.GetOrderByIdResponseData, *utils.ErrorHandler) {
	var response data.GetOrderByIdResponseData

	if !DoesOrderExist(db, orderId) {
		return response, utils.NotFoundError("Order with given id does not exist")
	}

	query := `SELECT
		product_id, buyer_id, delivery_type, delivery_fee, payment_type, payment_fee, small_order_fee, total_paid,
		order_quantity, phone_number, order_date::TEXT, address_line_1,COALESCE(address_line_2, ''), postal_code, payment_status,
		COALESCE(telegram_handle, ''), product_price
	FROM orders WHERE order_id=$1;`

	err := db.QueryRowContext(context.Background(), query, orderId).Scan(
		&response.ProductId, &response.BuyerId, &response.Fees.DeliveryType, &response.Fees.DeliveryFee,
		&response.Fees.PaymentType, &response.Fees.PaymentFee, &response.Fees.SmallOrderFee, &response.Fees.TotalPaid,
		&response.OrderQuantity, &response.PhoneNumber, &response.OrderDate, &response.AddressLine1, &response.AddressLine2,
		&response.PostalCode, &response.PaymentStatus, &response.TelegramHandle, &response.Fees.ProductPrice)

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

	query := `SELECT
		product_id, email, delivery_type, delivery_fee, payment_type, payment_fee, small_order_fee, total_paid,
		order_quantity, phone_number, order_date::TEXT, address_line_1, COALESCE(address_line_2, ''), postal_code, payment_status,
		COALESCE(telegram_handle, ''), product_price
	FROM guest_orders WHERE guest_order_id=$1;`

	err := db.QueryRowContext(context.Background(), query, guestOrderId).Scan(
		&response.ProductId, &response.Email, &response.Fees.DeliveryType, &response.Fees.DeliveryFee,
		&response.Fees.PaymentType, &response.Fees.PaymentFee, &response.Fees.SmallOrderFee, &response.Fees.TotalPaid,
		&response.OrderQuantity, &response.PhoneNumber, &response.OrderDate, &response.AddressLine1, &response.AddressLine2,
		&response.PostalCode, &response.PaymentStatus, &response.TelegramHandle, &response.Fees.ProductPrice)

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

	if request.Fees.PaymentType != "card" && request.Fees.PaymentType != "paynow_online" {
		utils.LogMessage("Payment Type is invalid")
		return utils.BadRequestError("Bad payment_type data")
	}

	if request.Fees.DeliveryType != "standard_delivery" && request.Fees.DeliveryType != "self_collection" {
		utils.LogMessage("Delivery Type is invalid")
		return utils.BadRequestError("Bad delivery_type data")
	}

	amountErr := validatePaymentAmount(request.OrderQuantity, request.Fees)

	if amountErr != nil {
		return amountErr
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

	if request.Fees.PaymentType != "card" && request.Fees.PaymentType != "paynow_online" {
		utils.LogMessage("Payment Type is invalid")
		return utils.BadRequestError("Bad payment_type data")
	}

	if request.Fees.DeliveryType != "standard_delivery" && request.Fees.DeliveryType != "self_collection" {
		utils.LogMessage("Delivery Type is invalid")
		return utils.BadRequestError("Bad delivery_type data")
	}

	amountErr := validatePaymentAmount(request.OrderQuantity, request.Fees)

	if amountErr != nil {
		return amountErr
	}

	return nil
}

/*
Updates the order payment status to either 'failed' or 'completed' based on the Hash sent by the Payment gateway
*/
func UpdateOrderPaymentStatus(db *sql.DB, orderId string, req data.PaymentValidationRequestData) *utils.ErrorHandler {
	if !DoesOrderExist(db, orderId) {
		return utils.NotFoundError("Order with given id does not exist")
	}

	var productId string
	var quantity int
	query := `UPDATE orders SET payment_status = $2 WHERE order_id = $1 RETURNING product_id, order_quantity;`
	err := db.QueryRowContext(context.Background(), query, orderId, req.Status).Scan(&productId, &quantity)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Updating order rows")
		return errResp
	}

	query = `UPDATE products SET sold_quantity = sold_quantity + $1 WHERE product_id = $2;`
	_, err = db.ExecContext(context.Background(), query, quantity, productId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Updating product rows")
		return errResp
	}

	return nil
}

/*
Updates the order payment status to either 'failed' or 'completed' based on the Hash sent by the Payment gateway
*/
func UpdateGuestOrderPaymentStatus(db *sql.DB, guestOrderId string, req data.PaymentValidationRequestData) *utils.ErrorHandler {
	if !DoesGuestOrderExist(db, guestOrderId) {
		return utils.NotFoundError("Guest order with given id does not exist")
	}

	query := `UPDATE guest_orders SET payment_status = $2 WHERE guest_order_id = $1`
	_, err := db.ExecContext(context.Background(), query, guestOrderId, req.Status)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Updating guest_order rows")
		return errResp
	}

	return nil
}

/*
Calculate the payment amount given an order request
*/
func validatePaymentAmount(quantity int, fees data.OrderFees) *utils.ErrorHandler {
	//calculate product cost
	amountToBePaid := quantity * fees.ProductPrice

	//calculate small order fee
	if amountToBePaid < 2500 {
		if fees.SmallOrderFee != 100 {
			utils.LogMessage("Small Order fee is incorrect")
			return utils.BadRequestError("Bad small_order_fee data")
		}
		amountToBePaid += 100
	}

	//calculate delivery fee
	if fees.DeliveryType == "standard_delivery" {
		if fees.DeliveryFee != 400 {
			utils.LogMessage("Delivery fee is incorrect")
			return utils.BadRequestError("Bad delivery_fee data")
		}
		amountToBePaid += 400
	}

	var paymentFee int
	//calculate payment fee
	if fees.PaymentType == "card" {
		paymentFee = int(math.Ceil((float64(amountToBePaid) * 2 / 100)))
		amountToBePaid += paymentFee
		if fees.PaymentFee != paymentFee {
			fmt.Println(paymentFee)
			utils.LogMessage("Payment fee is incorrect")
			return utils.BadRequestError("Bad payment_fee data")
		}
	}

	if amountToBePaid != fees.TotalPaid {
		utils.LogMessage("Total Paid amount is incorrect")
		return utils.BadRequestError("Bad total paid data")
	}

	return nil
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
