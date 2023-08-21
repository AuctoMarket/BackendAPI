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
// @Success      201  {object}  data.Message
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders [post]
func handleCreateOrder(c *gin.Context) {
	var createOrderData data.CreateOrderDataRequest
	bindErr := c.ShouldBindJSON(&createOrderData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	message, err := order.CreateOrder(db, createOrderData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &message)
}

// handleCreateGuestOrder godoc
// @Summary      Creates a new order
// @Description  Creates a new order for a specific product. This order is created by a Guest account.
// @Accept       json
// @Produce      json
// @Success      201  {object}  data.Message
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /orders/guest [post]
func handleGuestCreateOrder(c *gin.Context) {
	var createGuestOrderData data.CreateGuestOrderDataRequest
	bindErr := c.ShouldBindJSON(&createGuestOrderData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	message, err := order.CreateGuestOrder(db, createGuestOrderData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &message)
}
