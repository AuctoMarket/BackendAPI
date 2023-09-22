package data

import (
	"BackendAPI/utils"
	"strconv"
)

type CreateProductResponseData struct {
	ProductId    string `json:"product_id" binding:"required"`
	SellerId     string `json:"seller_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"desc" binding:"required"`
	ProductType  string `json:"product_type" binding:"required"`
	Language     string `json:"language" binding:"required"`
	Expansion    string `json:"expansion" binding:"required"`
	PostedDate   string `json:"posted_date" binding:"required" example:"2023-08-03 02:50:26.034552906 +0000 UTC m=+192.307467936"`
	Price        int    `json:"price" binding:"required"`
	Condition    int8   `json:"condition" binding:"required"`
	Quantity     int    `json:"product_quantity" binding:"required"`
	SoldQuantity int    `json:"sold_quantity" binding:"required"`
	OrderBy      string `json:"order_by"`
	ReleasesOn   string `json:"releases_on"`
	Discount     int    `json:"discount"`
}

type GetProductResponseData struct {
	ProductId     string                    `json:"product_id" binding:"required"`
	SellerInfo    GetSellerByIdResponseData `json:"seller_info" binding:"required"`
	Title         string                    `json:"title" binding:"required"`
	Description   string                    `json:"desc" binding:"required"`
	ProductType   string                    `json:"product_type" binding:"required"`
	Language      string                    `json:"language" binding:"required"`
	Expansion     string                    `json:"expansion" binding:"required"`
	PostedDate    string                    `json:"posted_date" binding:"required" example:"2023-08-03 02:50:26.034552906 +0000 UTC m=+192.307467936"`
	Price         int                       `json:"price" binding:"required"`
	Condition     int8                      `json:"condition" binding:"required"`
	Quantity      int                       `json:"product_quantity" binding:"required"`
	SoldQuantity  int                       `json:"sold_quantity" binding:"required"`
	OrderBy       string                    `json:"order_by"`
	ReleasesOn    string                    `json:"releases_on"`
	Discount      int                       `json:"discount"`
	ProductImages []ProductImageData        `json:"images" binding:"required"`
}

type ProductImageData struct {
	ProductImagePath string `json:"image_path" binding:"required"`
	ProductImageNo   int    `json:"image_no" binding:"required"`
}

type CreateProductData struct {
	Title       string `json:"title" binding:"required"`
	SellerId    string `json:"seller_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	ProductType string `json:"product_type" binding:"required"`
	Language    string `json:"language" binding:"required"`
	Expansion   string `json:"expansion" binding:"required"`
	Price       int    `json:"price"`
	Condition   int8   `json:"condition" `
	Quantity    int    `json:"product_quantity"`
	OrderBy     string `json:"order_by"`
	ReleasesOn  string `json:"releases_on"`
	Discount    int    `json:"discount"`
}

type CreateProductImageData struct {
	ProductId string   `json:"product_id" binding:"required"`
	Images    []string `json:"images" binding:"required"`
}

type GetProductListRequestData struct {
	SortBy       string   `json:"sort"`
	Prices       []string `json:"prices"`
	ProductTypes []string `json:"product_type"`
	Languages    []string `json:"language"`
	Expansions   []string `json:"expansion"`
	Anchor       int      `json:"anchor"`
	Limit        int      `json:"limit"`
}

type GetProductListResponseData struct {
	ProductCount int                      `json:"product_count" binding:"required"`
	Products     []GetProductResponseData `json:"products" binding:"required"`
}

/*
Takes a product request and creates the corresponding product response
*/
func (request *CreateProductData) ProductCreateResponseFromRequest(response *CreateProductResponseData) {
	response.Condition = request.Condition
	response.Description = request.Description
	response.Price = request.Price
	response.ProductType = request.ProductType
	response.SellerId = request.SellerId
	response.Title = request.Title
	response.Quantity = request.Quantity
	response.SoldQuantity = 0
	response.ReleasesOn = request.ReleasesOn
	response.OrderBy = request.OrderBy
	response.Discount = request.Discount
	response.Language = request.Language
	response.Expansion = request.Expansion
}

func (request *GetProductListRequestData) GetProductListDataRequestFromParams(sortBy string, productTypes []string, languages []string,
	prices []string, expansions []string, anchor string, limit string) *utils.ErrorHandler {
	request.SortBy = sortBy
	request.ProductTypes = productTypes
	request.Languages = languages
	request.Expansions = expansions
	request.Prices = prices

	if anchor != "None" {
		anch, err := strconv.Atoi(anchor)
		if err != nil || anch < 0 {
			return utils.BadRequestError("Bad anchor param")
		}

		request.Anchor = anch
	}

	if limit != "None" {
		lim, err := strconv.Atoi(limit)
		if err != nil || lim < 0 {
			return utils.BadRequestError("Bad limit param")
		}

		request.Limit = lim
	}

	return nil
}
