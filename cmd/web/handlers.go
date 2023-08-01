package main

import (
	"BackendAPI/api/buyer"
	"BackendAPI/api/product"
	"BackendAPI/api/seller"
	"BackendAPI/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleBuyerSignup godoc
// @Summary      Signs a new buyer up
// @Description  Checks to see if a buyer email exists and if not creates a new account with supplied email and password
// if not returns a bad request error (400).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Buyers email"
// @Param 		 password body string true "Buyers password as plaintext"
// @Success      200  {object}  data.BuyerLoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyers/signup [post]
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

// handleBuyerLogin godoc
// @Summary      Logs a buyer into their account
// @Description  Checks to see if a buyer email exists and if supplied password matches the stored password
// if not returns a unauthorized error (401).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Buyers email"
// @Param 		 password body string true "Buyers password as plaintext"
// @Success      200  {object}  data.BuyerLoginResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /buyers/login [post]
func handleBuyerLogin(c *gin.Context) {
	var loginData data.UserLoginData
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

// handleGetProductById godoc
// @Summary      Gets a Product by its Product ID
// @Description  Checks to see if a product with a given id exists and returns its product information if it does.
// If not it returns a not found error (404)
// @Produce      json
// @Param id path string true "product_id"
// @Success      200  {object}  data.ProductResponseData
// @Failure      400  {object}  data.Message
// @Failure      404  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /products/{id}  [get]
func handleGetProductById(c *gin.Context) {
	productId := c.Param("id")

	if productId == "" {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	product, err := product.GetProductById(db, productId)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &product)
}

// handleCreateProduct godoc
// @Summary      Creates a new product post
// @Description  Creates a new product post with the supplied data, if the data is not valid it throws and error
// @Produce      json
// @Param 		 seller_id body string true "The Seller who posted the product's seller_id"
// @Param 		 title body string true "Title of the product"
// @Param 		 description body string true "Short description of the product"
// @Param 		 price body int true "Price as an int of the product"
// @Param 		 condition body int true "Condition of the product from a scale of 0 to 5"
// @Param 		 product_type body string true "Type of product sale: Buy-Now or Pre-Order"
// @Success      201  {object}  data.ProductResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /products  [post]
func handleCreateProduct(c *gin.Context) {
	var createProduct data.ProductCreateData
	bindErr := c.ShouldBindJSON(&createProduct)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	product, err := product.CreateProduct(db, createProduct)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &product)
}

// handleSellerSignup godoc
// @Summary      Signs a new seller up
// @Description  Checks to see if a seller email does not already exists if so creates a new
// seller account with supplied email, password and seller_name
// if not returns a bad request error (400).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Sellers email"
// @Param 		 password body string true "Sellers password as plaintext"
// @Param 		 seller_name body string true "Sellers seller alias that is displayed as their seller name"
// @Success      200  {object}  data.SellerResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /sellers/signup [post]
func handleSellerSignUp(c *gin.Context) {
	var signUpData data.SellerSignUpData
	bindErr := c.ShouldBindJSON(&signUpData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	signUpResponse, err := seller.SellerSignUp(db, signUpData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, &signUpResponse)
}

// handleSellerLogin godoc
// @Summary      Logs a seller into their account
// @Description  Checks to see if a sellers email exists and if supplied password matches the stored password
// if not returns a unauthorized error (401).
// @Accept       json
// @Produce      json
// @Param 		 email body string true "Sellers email"
// @Param 		 password body string true "Sellers password as plaintext"
// @Success      200  {object}  data.SellerResponseData
// @Failure      400  {object}  data.Message
// @Failure      401  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /sellers/login [post]
func handleSellerLogin(c *gin.Context) {
	var loginData data.UserLoginData
	bindErr := c.ShouldBindJSON(&loginData)

	if bindErr != nil {
		r := data.Message{Message: "Bad Request Body"}
		c.JSON(http.StatusBadRequest, r)
		return
	}

	loginResponse, err := seller.SellerLogin(db, loginData)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, &loginResponse)
}

/*
Test Ping as a sanity check
*/
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
