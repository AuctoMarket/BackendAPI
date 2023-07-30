package data

type ProductData struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"desc" binding:"required"`
	ProductType string `json:"product_type" binding:"required"`
	PostedData  string `json:"posted_date" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Condition   int8   `json:"condition" binding:"required"`
}
