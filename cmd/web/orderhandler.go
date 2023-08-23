package main

import (
	"BackendAPI/api/order"
	"BackendAPI/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleCreateOrder godoc
// @Summary      Creates a new order
// @Description  Creates a new order for a specific product. This order is created by an existing buyer with an account.
// @Accept       json
// @Produce      json
// @Success      201  {object}  data.CreateOrderResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders [post]
func handleCreateOrder(c *gin.Context) {
	var createOrderData data.CreateOrderRequestData
	bindErr := c.ShouldBindJSON(&createOrderData)

	if bindErr != nil {
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
// @Summary      Creates a new order
// @Description  Creates a new order for a specific product. This order is created by a Guest account.
// @Accept       json
// @Produce      json
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

// handleCreateGuestOrder godoc
// @Summary      Creates a new order
// @Description  Creates a new order for a specific product. This order is created by a Guest account.
// @Accept       json
// @Produce      json
// @Success      201  {object}  data.CreateGuestOrderResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders/guest [post]
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

// handleCreateGuestOrder godoc
// @Summary      Creates a new order
// @Description  Creates a new order for a specific product. This order is created by a Guest account.
// @Accept       json
// @Produce      json
// @Success      201  {object}  data.CreateGuestOrderResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders/guest [post]
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
