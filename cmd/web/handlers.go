package main

import (
	"BackendAPI/api/buyer"
	"BackendAPI/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleBuyerLogin godoc
// @Summary      Logs a buyer into their account
// @Description  Checks to see if a user email exists and if supplied password matches the stored password
// if not returns a unauthorized error (401).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "email"
// @Param 		 password body string true "password"
// @Success      200  {object}  data.LoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyer/login/ [post]
func handleBuyerLogin(c *gin.Context) {
	var loginData data.BuyerLoginData
	bindErr := c.ShouldBindJSON(&loginData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	loginResponse, err := buyer.BuyerLogin(db, loginData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &loginResponse)
}

// handleBuyerSignup godoc
// @Summary      Signs a new buyer up
// @Description  Checks to see if a user email exists and if not creates a new account with supplied email and password
// if not returns a bad request error (400).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "email"
// @Param 		 password body string true "password"
// @Success      200  {object}  data.LoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyer/signup/ [post]
func handleBuyerSignUp(c *gin.Context) {
	var signUpData data.BuyerSignUpData
	bindErr := c.ShouldBindJSON(&signUpData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	signUpResponse, err := buyer.BuyerSignUp(db, signUpData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &signUpResponse)
}

// handleGetProductById godoc
// @Summary      Gets a
// @Description  Checks to see if a product with a given id exists and returns its product information if it does.
// If not it returns a not found error (404)
// @Produce      json
// @Param 		 email body string true "email"
// @Param 		 password body string true "password"
// @Success      200  {object}  data.LoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      404  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyer/signup/ [post]
func handleGetProductById(c *gin.Context) {
	r := data.Message{Message: "TODO"}
	c.JSON(http.StatusOK, r)
}

/*
Test Ping as a sanity check
*/
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
