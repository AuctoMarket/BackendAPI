package data

type ProductCreateResponseData struct {
	ProductId    string `json:"product_id" binding:"required"`
	SellerId     string `json:"seller_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"desc" binding:"required"`
	ProductType  string `json:"product_type" binding:"required"`
	PostedDate   string `json:"posted_date" binding:"required" example:"2023-08-03 02:50:26.034552906 +0000 UTC m=+192.307467936"`
	Price        int    `json:"price" binding:"required"`
	Condition    int8   `json:"condition" binding:"required"`
	Quantity     int    `json:"product_quantity" binding:"required"`
	SoldQuantity int    `json:"sold_quantity" binding:"required"`
}

type GetProductResponseData struct {
	ProductId     string             `json:"product_id" binding:"required"`
	SellerId      string             `json:"seller_id" binding:"required"`
	Title         string             `json:"title" binding:"required"`
	Description   string             `json:"desc" binding:"required"`
	ProductType   string             `json:"product_type" binding:"required"`
	PostedDate    string             `json:"posted_date" binding:"required" example:"2023-08-03 02:50:26.034552906 +0000 UTC m=+192.307467936"`
	Price         int                `json:"price" binding:"required"`
	Condition     int8               `json:"condition" binding:"required"`
	Quantity      int                `json:"product_quantity" binding:"required"`
	SoldQuantity  int                `json:"sold_quantity" binding:"required"`
	ProductImages []ProductImageData `json:"images" binding:"required"`
}

type ProductImageData struct {
	ProductImagePath string `json:"image_path" binding:"required"`
	ProductImageNo   int    `json:"image_no" binding:"required"`
}

type ProductCreateData struct {
	Title       string `json:"title" binding:"required"`
	SellerId    string `json:"seller_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	ProductType string `json:"product_type" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Condition   int8   `json:"condition" binding:"required"`
	Quantity    int    `json:"product_quantity" binding:"required"`
}

type ProductImageCreateData struct {
	ProductId string   `json:"product_id" binding:"required"`
	Images    []string `json:"images" binding:"required"`
}

func (request *ProductCreateData) CreateResponseFromRequest(response *ProductCreateResponseData) *ProductCreateResponseData {
	response.Condition = request.Condition
	response.Description = request.Description
	response.Price = request.Price
	response.ProductType = request.ProductType
	response.SellerId = request.SellerId
	response.Title = request.Title
	response.Quantity = request.Quantity
	response.SoldQuantity = 0

	return response
}
