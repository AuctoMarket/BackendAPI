package order

import (
	"BackendAPI/api/buyer"
	"BackendAPI/api/product"
	"BackendAPI/data"
	"BackendAPI/utils"
	"context"
	"database/sql"
	"math"
	"time"
)

/*
Create a order for a specific product by a specific buyer from the order request and store it in the database. If the
buyer or product do not exist, return a BadRequestError (400).
*/
func CreateOrder(db *sql.DB, request data.CreateOrderDataRequest) (data.Message, *utils.ErrorHandler) {
	var response data.Message

	//validate input
	validErr := validateCreateOrderRequest(db, request)
	if validErr != nil {
		return response, validErr
	}

	cols := `product_id, buyer_id, delivery_type, order_quantity, payment_type, 
		payment_status, phone_number, order_date, address_line_1, postal_code`
	vals := `$1,$2,$3,$4,$5,$6,$7,$8,$9,$10`
	if request.AddressLine2 != "" {
		cols += `address_line_2`
		vals += `$11`
	}

	query := `INSERT INTO orders(` + cols + `) VALUES (` + vals + `) RETURNING order_id;`
	orderDate := time.Now()
	paymentStatus := `pending`
	var orderId string
	var err error

	if request.AddressLine2 != "" {
		err = db.QueryRowContext(
			context.Background(), query,
			request.ProductId, request.BuyerId, request.DeliveryType, request.OrderQuantity,
			request.PaymentType, paymentStatus, request.PhoneNumber, orderDate, request.AddressLine1,
			request.PostalCode, request.AddressLine2).Scan(&orderId)
	} else {
		err = db.QueryRowContext(
			context.Background(), query,
			request.ProductId, request.BuyerId, request.DeliveryType, request.OrderQuantity,
			request.PaymentType, paymentStatus, request.PhoneNumber, orderDate, request.AddressLine1,
			request.PostalCode).Scan(&orderId)
	}

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in inserting Order rows")
		return response, errResp
	}

	return response, nil
}

/*
Create a order for a specific product by guest account from the order request and store it in the database. If the
product does not exist, return a BadRequestError (400).
*/
func CreateGuestOrder(db *sql.DB, request data.CreateGuestOrderDataRequest) (data.Message, *utils.ErrorHandler) {
	var message data.Message
	return message, nil
}

/*
Validate a create order request to ensure params are correct.
*/
func validateCreateOrderRequest(db *sql.DB, request data.CreateOrderDataRequest) *utils.ErrorHandler {
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
	var paymentFee float64
	if paymentType == "card" {
		paymentFee = float64(amountToBePaid) * 0.02
	}

	amountToBePaid += int(math.Ceil(paymentFee))
	return amountToBePaid
}
