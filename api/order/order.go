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
	"strconv"
	"time"
)

/*
Create a order for a specific product by a specific buyer from the order request and store it in the database. If the
buyer or product do not exist, return a BadRequestError (400).
*/
func CreateOrder(db *sql.DB, request data.CreateOrderRequestData) (data.CreateOrderResponseData, *utils.ErrorHandler) {
	var response data.CreateOrderResponseData
	orderDate := time.Now()

	//validate order details
	validErr := validateCreateOrderRequest(db, request)
	if validErr != nil {
		return response, validErr
	}

	//validate payment amount details
	amountErr := validatePaymentAmount(db, request.Products, request.Fees)
	if amountErr != nil {
		return response, amountErr
	}

	//SQL Query to insert new order
	query := `INSERT INTO orders(
		buyer_id, 
		delivery_type, 
		delivery_fee, 
		payment_type, 
		payment_fee, 
		small_order_fee, 
		total_paid,
		phone_number, 
		order_date, 
		address_line_1, 
		address_line_2, 
		postal_code, 
		telegram_handle) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) 
		RETURNING order_id;`

	err := db.QueryRowContext(
		context.Background(), query,
		request.BuyerId, request.Fees.DeliveryType, request.Fees.DeliveryFee,
		request.Fees.PaymentType, request.Fees.PaymentFee, request.Fees.SmallOrderFee, request.Fees.TotalPaid,
		request.PhoneNumber, orderDate, request.AddressLine1, utils.NewNullableString(request.AddressLine2),
		request.PostalCode, utils.NewNullableString(request.TelegramHandle)).Scan(&response.OrderId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in inserting Order rows")
		return response, errResp
	}

	query = `INSERT INTO order_products(product_id, order_id, quantity) VALUES`
	for i := 0; i < len(request.Products); i++ {
		query += `('` + request.Products[i].ProductId + `','` + response.OrderId + `',` + strconv.Itoa(request.Products[i].OrderQuantity) + `)`
		if i < len(request.Products)-1 {
			query += `,`
		}
	}
	_, err = db.ExecContext(context.Background(), query)

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

	//validate payment amount details
	amountErr := validatePaymentAmount(db, request.Products, request.Fees)
	if amountErr != nil {
		return response, amountErr
	}

	//SQL Query to insert new guest order
	query := `INSERT INTO guest_orders(
		email, 
		delivery_type, 
		delivery_fee, 
		payment_type, 
		payment_fee, 
		small_order_fee, 
		total_paid,
		phone_number, 
		order_date, 
		address_line_1, 
		address_line_2, 
		postal_code, 
		telegram_handle) 
		VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) 
		RETURNING guest_order_id;`

	err := db.QueryRowContext(
		context.Background(), query,
		request.Email, request.Fees.DeliveryType, request.Fees.DeliveryFee,
		request.Fees.PaymentType, request.Fees.PaymentFee, request.Fees.SmallOrderFee, request.Fees.TotalPaid,
		request.PhoneNumber, orderDate, request.AddressLine1, utils.NewNullableString(request.AddressLine2),
		request.PostalCode, utils.NewNullableString(request.TelegramHandle)).Scan(&response.GuestOrderId)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in inserting Order rows")
		return response, errResp
	}

	query = `INSERT INTO guest_order_products(product_id, guest_order_id, quantity) VALUES`
	for i := 0; i < len(request.Products); i++ {
		query += `('` + request.Products[i].ProductId + `','` + response.GuestOrderId + `',` + strconv.Itoa(request.Products[i].OrderQuantity) + `)`
		if i < len(request.Products)-1 {
			query += `,`
		}
	}
	_, err = db.ExecContext(context.Background(), query)

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
		buyer_id, 
		delivery_type, 
		delivery_fee, 
		payment_type, 
		payment_fee, 
		small_order_fee, 
		total_paid,
		phone_number, 
		order_date::TEXT, 
		address_line_1,
		COALESCE(address_line_2, ''), 
		postal_code, 
		payment_status,
		COALESCE(telegram_handle, ''),
		order_products.product_id,
		order_products.quantity
	FROM (orders INNER JOIN order_products ON orders.order_id = order_products.order_id) WHERE orders.order_id=$1;`

	rows, err := db.QueryContext(context.Background(), query, orderId)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting Order rows")
		return response, errResp
	}

	for rows.Next() {
		var productId string
		var quantity int
		err = rows.Scan(
			&response.BuyerId, &response.Fees.DeliveryType, &response.Fees.DeliveryFee,
			&response.Fees.PaymentType, &response.Fees.PaymentFee, &response.Fees.SmallOrderFee, &response.Fees.TotalPaid,
			&response.PhoneNumber, &response.OrderDate, &response.AddressLine1, &response.AddressLine2,
			&response.PostalCode, &response.PaymentStatus, &response.TelegramHandle, &productId, &quantity)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in selecting Order rows")
			return response, errResp
		}

		response.Products = append(response.Products, data.ProductOrder{ProductId: productId, OrderQuantity: quantity})
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
		email, 
		delivery_type, 
		delivery_fee, 
		payment_type, 
		payment_fee, 
		small_order_fee, 
		total_paid,
		phone_number, 
		order_date::TEXT, 
		address_line_1, 
		COALESCE(address_line_2, ''), 
		postal_code, 
		payment_status,
		COALESCE(telegram_handle, ''),
		guest_order_products.product_id,
		guest_order_products.quantity
	FROM (guest_orders INNER JOIN guest_order_products ON guest_orders.guest_order_id = guest_order_products.guest_order_id) WHERE guest_orders.guest_order_id=$1;`

	rows, err := db.QueryContext(context.Background(), query, guestOrderId)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting Guest Order rows")
		return response, errResp
	}

	for rows.Next() {
		var productId string
		var quantity int
		err = rows.Scan(
			&response.Email, &response.Fees.DeliveryType, &response.Fees.DeliveryFee,
			&response.Fees.PaymentType, &response.Fees.PaymentFee, &response.Fees.SmallOrderFee, &response.Fees.TotalPaid,
			&response.PhoneNumber, &response.OrderDate, &response.AddressLine1, &response.AddressLine2,
			&response.PostalCode, &response.PaymentStatus, &response.TelegramHandle, &productId, &quantity)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in selecting Order rows")
			return response, errResp
		}

		response.Products = append(response.Products, data.ProductOrder{ProductId: productId, OrderQuantity: quantity})
	}

	response.GuestOrderId = guestOrderId
	return response, nil
}

/*
Validate a create order request data
*/
func validateCreateOrderRequest(db *sql.DB, request data.CreateOrderRequestData) *utils.ErrorHandler {
	if len(request.Products) == 0 {
		utils.LogMessage("Order with no products selected")
		return utils.BadRequestError("Bad order products data")
	}

	for i := 0; i < len(request.Products); i++ {
		if request.Products[i].OrderQuantity <= 0 {
			utils.LogMessage("Invalid product amount selected")
			return utils.BadRequestError("Bad order products data")
		}

		if !product.DoesProductExist(db, request.Products[i].ProductId) {
			utils.LogMessage("Product with given id does not exist")
			return utils.BadRequestError("Bad product_id data")
		}
	}

	if !buyer.DoesBuyerExist(db, request.BuyerId) {
		utils.LogMessage("Buyer with given id does not exist")
		return utils.BadRequestError("Bad buyer_id data")
	}

	if len(request.PostalCode) != 6 {
		utils.LogMessage("Postal Code Data Incorrect")
		return utils.BadRequestError("Bad postal code data")
	}

	if request.Fees.PaymentType != "card" && request.Fees.PaymentType != "paynow_online" {
		utils.LogMessage("Payment Type is invalid")
		return utils.BadRequestError("Bad payment_type data")
	}

	if request.Fees.DeliveryType != "standard_delivery" && request.Fees.DeliveryType != "self_collection" {
		utils.LogMessage("Delivery Type is invalid")
		return utils.BadRequestError("Bad delivery_type data")
	}

	//Check to make sure all product quantities are within the stock limits
	query := `SELECT product_id, (product_quantity-sold_quantity) FROM products WHERE product_id IN (`
	for i := 0; i < len(request.Products); i++ {
		query += `'` + request.Products[i].ProductId + `'`
		if i < len(request.Products)-1 {
			query += `,`
		}
	}
	query += `)`

	var productMap map[string]int
	productMap = make(map[string]int)

	rows, err := db.QueryContext(context.Background(), query)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Selecting product rows")
		return errResp
	}

	for rows.Next() {
		var productId string
		var quantity int

		err = rows.Scan(&productId, &quantity)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in Selecting product rows")
			return errResp
		}

		productMap[productId] = quantity
	}

	for i := 0; i < len(request.Products); i++ {
		if productMap[request.Products[i].ProductId] < request.Products[i].OrderQuantity {
			utils.LogMessage("Quantity ordered is greater than available stock")
			return utils.BadRequestError("Bad quantity data")
		}
	}

	return nil
}

/*
Validate a create guest order request to ensure params are correct.
*/
func validateCreateGuestOrderRequest(db *sql.DB, request data.CreateGuestOrderRequestData) *utils.ErrorHandler {
	if len(request.Products) == 0 {
		utils.LogMessage("Guest order with no products selected")
		return utils.BadRequestError("Bad guest order products data")
	}

	for i := 0; i < len(request.Products); i++ {
		if request.Products[i].OrderQuantity <= 0 {
			utils.LogMessage("Invalid product amount selected")
			return utils.BadRequestError("Bad order products data")
		}

		if !product.DoesProductExist(db, request.Products[i].ProductId) {
			utils.LogMessage("Product with given id does not exist")
			return utils.BadRequestError("Bad product_id data")
		}
	}

	if len(request.PostalCode) != 6 {
		utils.LogMessage("Postal Code Data Incorrect")
		return utils.BadRequestError("Bad postal code data")
	}

	if request.Fees.PaymentType != "card" && request.Fees.PaymentType != "paynow_online" {
		utils.LogMessage("Payment Type is invalid")
		return utils.BadRequestError("Bad payment_type data")
	}

	if request.Fees.DeliveryType != "standard_delivery" && request.Fees.DeliveryType != "self_collection" {
		utils.LogMessage("Delivery Type is invalid")
		return utils.BadRequestError("Bad delivery_type data")
	}

	//Check to make sure all product quantities are within the stock limits
	query := `SELECT product_id, (product_quantity-sold_quantity) FROM products WHERE product_id IN (`
	for i := 0; i < len(request.Products); i++ {
		query += `'` + request.Products[i].ProductId + `'`
		if i < len(request.Products)-1 {
			query += `,`
		}
	}
	query += `)`

	var productMap map[string]int
	productMap = make(map[string]int)

	rows, err := db.QueryContext(context.Background(), query)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Selecting product rows")
		return errResp
	}

	for rows.Next() {
		var productId string
		var quantity int

		err = rows.Scan(&productId, &quantity)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in Selecting product rows")
			return errResp
		}

		productMap[productId] = quantity
	}

	for i := 0; i < len(request.Products); i++ {
		if productMap[request.Products[i].ProductId] < request.Products[i].OrderQuantity {
			utils.LogMessage("Quantity ordered is greater than available stock")
			return utils.BadRequestError("Bad quantity data")
		}
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

	query := `UPDATE orders SET payment_status = $2 WHERE order_id = $1;`
	_, err := db.ExecContext(context.Background(), query, orderId, req.Status)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Updating order rows")
		return errResp
	}

	var products []data.ProductOrder

	query = `SELECT product_id, quantity FROM order_products WHERE order_id = $1;`
	rows, err := db.QueryContext(context.Background(), query, orderId)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting product order rows")
		return errResp
	}

	for rows.Next() {
		var productId string
		var quantity int
		err = rows.Scan(&productId, &quantity)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in selecting product order rows")
			return errResp
		}

		prod := data.ProductOrder{ProductId: productId, OrderQuantity: quantity}
		products = append(products, prod)
	}

	if req.Status == "completed" {
		for i := 0; i < len(products); i++ {
			query = `UPDATE products SET sold_quantity = sold_quantity + $1 WHERE product_id = $2;`
			_, err = db.ExecContext(context.Background(), query, products[i].OrderQuantity, products[i].ProductId)

			if err != nil {
				errResp := utils.InternalServerError(nil)
				utils.LogError(err, "Error in Updating product rows")
				return errResp
			}
		}
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

	query := `UPDATE guest_orders SET payment_status = $2 WHERE guest_order_id = $1;`
	_, err := db.ExecContext(context.Background(), query, guestOrderId, req.Status)

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Updating guest order rows")
		return errResp
	}

	var products []data.ProductOrder

	query = `SELECT product_id, quantity FROM guest_order_products WHERE guest_order_id = $1;`
	rows, err := db.QueryContext(context.Background(), query, guestOrderId)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in selecting guest product order rows")
		return errResp
	}

	for rows.Next() {
		var productId string
		var quantity int
		err = rows.Scan(&productId, &quantity)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in selecting guest product order rows")
			return errResp
		}

		prod := data.ProductOrder{ProductId: productId, OrderQuantity: quantity}
		products = append(products, prod)
	}

	if req.Status == "completed" {
		for i := 0; i < len(products); i++ {
			query = `UPDATE products SET sold_quantity = sold_quantity + $1 WHERE product_id = $2;`
			_, err = db.ExecContext(context.Background(), query, products[i].OrderQuantity, products[i].ProductId)

			if err != nil {
				errResp := utils.InternalServerError(nil)
				utils.LogError(err, "Error in Updating product rows")
				return errResp
			}
		}
	}

	return nil
}

/*
Calculate the payment amount given an order request
*/
func validatePaymentAmount(db *sql.DB, products []data.ProductOrder, fees data.OrderFees) *utils.ErrorHandler {
	//calculate product costs
	query := `SELECT products.product_id, (price - COALESCE(discount, 0))
		FROM 
			(products LEFT OUTER JOIN product_discounts ON product_discounts.product_id = products.product_id)
	 	WHERE products.product_id IN (`
	for i := 0; i < len(products); i++ {
		query += `'` + products[i].ProductId + `'`
		if i < len(products)-1 {
			query += `,`
		}
	}
	query += `)`

	var productMap map[string]int
	productMap = make(map[string]int)

	rows, err := db.QueryContext(context.Background(), query)
	defer rows.Close()

	if err != nil {
		errResp := utils.InternalServerError(nil)
		utils.LogError(err, "Error in Selecting product rows")
		return errResp
	}

	for rows.Next() {
		var productId string
		var price int

		err = rows.Scan(&productId, &price)

		if err != nil {
			errResp := utils.InternalServerError(nil)
			utils.LogError(err, "Error in Selecting product rows")
			return errResp
		}

		productMap[productId] = price
	}

	var amountToBePaid int
	for i := 0; i < len(products); i++ {
		amountToBePaid += productMap[products[i].ProductId] * products[i].OrderQuantity
	}

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
