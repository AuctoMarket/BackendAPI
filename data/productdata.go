package data

type ProductResponseData struct {
	ProductId   string `json:"product_id" binding:"required"`
	SellerId    string `json:"seller_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"desc" binding:"required"`
	ProductType string `json:"product_type" binding:"required"`
	PostedDate  string `json:"posted_date" binding:"required" example:""`
	Price       int    `json:"price" binding:"required"`
	Condition   int8   `json:"condition" binding:"required"`
}

type ProductCreateData struct {
	Title       string `json:"title" binding:"required"`
	SellerId    string `json:"seller_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	ProductType string `json:"product_type" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Condition   int8   `json:"condition" binding:"required"`
}

func (request *ProductCreateData) CreateResponseFromRequest(response *ProductResponseData) *ProductResponseData {
	response.Condition = request.Condition
	response.Description = request.Description
	response.Price = request.Price
	response.ProductType = request.ProductType
	response.SellerId = request.SellerId
	response.Title = request.Title

	return response
}
