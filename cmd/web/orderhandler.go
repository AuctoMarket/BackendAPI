package main

import (
	"BackendAPI/api/order"
	"BackendAPI/data"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleCreateOrder godoc
// @Summary      Creates a new order
// @Description  Creates a new order for a specific product. This order is created by an existing buyer with an account.
// @Accept       json
// @Produce      json
// @Param 		 product_id body string true "The product for which we are creating an order"
// @Param 		 buyer_id body string true "The id of the buyer who is creating the order"
// @Param 		 order_quantity body int true "Quantity of the product being ordered"
// @Param 		 phone_number body string true "Phone number of buyer"
// @Param        address_line_1 body string true "Delivery Address"
// @Param        address_line_2 body string false "Delivery Address 2"
// @Param        postal_code body string false "Postal code of address"
// @Param        fees body data.OrderFees false "Pricing Info, Delivery Type is either 'self_collection' or 'standard delivery', Payment type is 'card' or 'paynow_online'"
// @Success      201  {object}  data.CreateOrderResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders [post]
func handleCreateOrder(c *gin.Context) {
	var createOrderData data.CreateOrderRequestData
	bindErr := c.ShouldBindJSON(&createOrderData)

	if bindErr != nil {
		fmt.Println(bindErr)
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	response, err := order.CreateOrder(db, createOrderData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &response)
}

// handleCreateGuestOrder godoc
// @Summary      Creates a new guest order
// @Description  Creates a new order for a specific product. This order is created by a guest user.
// @Accept       json
// @Produce      json
// @Param 		 product_id body string true "The product for which we are creating an order"
// @Param 		 email body string true "The email of the guest user"
// @Param 		 order_quantity body int true "Quantity of the product being ordered"
// @Param 		 payment_type body string true "Payment method chosen for this order, can only be 'card' or 'paynow_online'"
// @Param 		 delivery_type body string true "Type of delivery method, can only be 'self_collection' or 'standard_delivery'"
// @Param 		 phone_number body string true "Phone number of buyer"
// @Param        address_line_1 body string true "Delivery Address"
// @Param        address_line_2 body string false "Delivery Address 2"
// @Param        postal_code body string false "Postal code of address"
// @Param		 amount body int true "Amount to be paid for the order in cents"
// @Success      201  {object}  data.CreateGuestOrderResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders/guest [post]
func handleCreateGuestOrder(c *gin.Context) {
	var createGuestOrderData data.CreateGuestOrderRequestData
	bindErr := c.ShouldBindJSON(&createGuestOrderData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	response, err := order.CreateGuestOrder(db, createGuestOrderData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &response)
}

// handleGetOrderById godoc
// @Summary      Fetched order details for an order with a specific order id
// @Description  Returns the order details of an order with a given order id. If the order id does not exists, returns a 404 error.
// @Accept       json
// @Produce      json
// @Success      200  {object}  data.GetOrderByIdResponseData
// @Failure      404  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders/{id} [get]
func handleGetOrderById(c *gin.Context) {
	productId := c.Param("id")

	product, err := order.GetOrderById(db, productId)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &product)
}

// handleGetGuestOrderById godoc
// @Summary      Fetched order details for an guest order with a specific guest order id
// @Description  Returns the order details of an guest order with a given guest order id. If the order id does not exists, returns a 404 error.
// @Accept       json
// @Produce      json
// @Success      200  {object}  data.GetGuestOrderByIdResponseData
// @Failure      404  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders/{id}/guest [get]
func handleGetGuestOrderById(c *gin.Context) {
	productId := c.Param("id")

	product, err := order.GetGuestOrderById(db, productId)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &product)
}

func handlePaymentComplete(c *gin.Context) {
	orderId := c.Param("id")
	var req data.PaymentValidationRequestData
	bindErr := c.ShouldBind(&req)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	err := order.UpdateOrderPaymentStatus(db, orderId, req)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.Status(200)
	return
}

func handleGuestPaymentComplete(c *gin.Context) {
	orderId := c.Param("id")
	var req data.PaymentValidationRequestData
	bindErr := c.ShouldBind(&req)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	err := order.UpdateGuestOrderPaymentStatus(db, orderId, req)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.Status(200)
	return
}
