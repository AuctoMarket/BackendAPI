package main

import (
	"BackendAPI/api/product"
	"BackendAPI/data"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleGetProductById godoc
// @Summary      Gets a Product by its Product ID
// @Description  Checks to see if a product with a given id exists and returns its product information if it does.
// If not it returns a not found error (404)
// @Produce      json
// @Param id path string true "product_id"
// @Success      200  {object}  data.GetProductResponseData
// @Failure      400  {object}  data.Message
// @Failure      404  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /products/{id}  [get]
func handleGetProductById(c *gin.Context) {
	productId := c.Param("id")

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
// @Param		 language body string true "Language of the product, is either 'Eng' or 'Jap'"
// @Param 		 expansion body string true "Expansion of the product"
// @Param        product_quantity body int true "Quantity of product to be put for sale"
// @Success      201  {object}  data.CreateProductResponseData
// @Failure      400  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /products  [post]
func handleCreateProduct(c *gin.Context) {
	var createProduct data.CreateProductData
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

// handleCreateProductImages godoc
// @Summary      Adds images to products
// @Description  Adds images to an existing product with supplied product id. If product with product id does not exist returns a
// error (404), otherwise returns a 201.
// @Accept       mpfd
// @Produce      json
// @Param        id path string true "product_id"
// @Param 		 images formData file true "Array of image files to add to the product post"
// @Success      201  {object}  data.Message
// @Failure      415  {object}  data.Message
// @Failure      404  {object}  data.Message
// @Failure      500  {object}  data.Message
// @Router       /products/{id}/images  [post]
func handleCreateProductImages(c *gin.Context) {
	form, formErr := c.MultipartForm()
	productId := c.Param("id")

	if formErr != nil {
		r := data.Message{Message: "Bad Content-Type in Request"}
		c.JSON(http.StatusUnsupportedMediaType, r)
		return
	}
	imageFiles := form.File["images"]
	var images []io.Reader

	for i := 0; i < len(imageFiles); i++ {
		image, err := imageFiles[i].Open()

		if err != nil {
			r := data.Message{Message: "Bad Content-Type in Request"}
			c.JSON(http.StatusUnsupportedMediaType, r)
			return
		}

		images = append(images, image)
	}

	response, err := product.CreateProductImages(db, s3Client, productId, images)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusCreated, response)

}

// handleGetProducts godoc
// @Summary      Gets Products with given query parameters
// @Description  Gets product information of products given query parameters provided in the Request
// @Accept       json
// @Produce      json
// @Param        sort_by query string false  "Sort By a specific attribute of the product. Default is posted_date"
// @Param 		 product_types query []string false "Get products by a specific product type, the types are 'Pre-Order' or 'Buy-Now'. Default is both will be selected"
// @Param 		 languages query []string false "Get products filtered by the language of the expansion. The choices are 'Eng' or 'Jap' and default is both."
// @Param 		 expansions query []string false "Get products filtered by the expansion of the product. Default is all expansions"
// @Param		 prices  query []string false "Gets products filtered by prices ranges, the ranges are '0-20', '20-50', '50-100', '100-200', '200'"
// @Param 		 anchor query int true "Indicates the offset for the products"
// @Param 		 limit query int true "Indicates the number of products fetched"
// @Success      200  {object}  data.GetProductListResponseData
// @Failure      500  {object}  data.Message
// @Router       /products  [get]
func handleGetProductList(c *gin.Context) {
	var request data.GetProductListRequestData
	sortBy := c.DefaultQuery("sort_by", "None")
	prices := c.QueryArray("prices")
	productTypes := c.QueryArray("product_types")
	languages := c.QueryArray("languages")
	expansions := c.QueryArray("expansions")
	anchor := c.DefaultQuery("anchor", "None")
	limit := c.DefaultQuery("limit", "None")

	err := request.GetProductListDataRequestFromParams(sortBy, productTypes, languages, prices, expansions, anchor, limit)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	products, err := product.GetProductList(db, request)

	if err != nil {
		r := data.Message{Message: err.Error()}
		c.JSON(err.ErrorCode(), r)
		return
	}

	c.JSON(http.StatusOK, products)
}
