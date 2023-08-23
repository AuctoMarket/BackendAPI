package data

type CreateOrderRequestData struct {
	ProductId     string `json:"product_id" binding:"required"`
	BuyerId       string `json:"buyer_id" binding:"required"`
	OrderQuantity int    `json:"order_quantity" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	DeliveryType  string `json:"delivery_type" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	AddressLine1  string `json:"address_line_1" binding:"required"`
	AddressLine2  string `json:"address_line_2"`
	PostalCode    string `json:"postal_code" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
}

type CreateGuestOrderRequestData struct {
	ProductId     string `json:"product_id" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	OrderQuantity int    `json:"order_quantity" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	DeliveryType  string `json:"delivery_type" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	AddressLine1  string `json:"address_line_1" binding:"required"`
	AddressLine2  string `json:"address_line_2"`
	PostalCode    string `json:"postal_code" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
}

type CreateOrderResponseData struct {
	OrderId     string `json:"order_id" binding:"required"`
	RedirectUrl string `json:"redirect_url" binding:"required"`
}

type CreateGuestOrderResponseData struct {
	GuestOrderId string `json:"guest_order_id" binding:"required"`
	RedirectUrl  string `json:"redirect_url" binding:"required"`
}

type GetOrderByIdResponseData struct {
	OrderId       string `json:"order_id" binding:"required"`
	ProductId     string `json:"product_id" binding:"required"`
	BuyerId       string `json:"buyer_id" binding:"required"`
	OrderQuantity int    `json:"order_quantity" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	DeliveryType  string `json:"delivery_type" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	AddressLine1  string `json:"address_line_1" binding:"required"`
	AddressLine2  string `json:"address_line_2"`
	PostalCode    string `json:"postal_code" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
}

type GetGuestOrderByIdResponseData struct {
	GuestOrderId  string `json:"guest_order_id" binding:"required"`
	ProductId     string `json:"product_id" binding:"required"`
	Email         string `json:"email" binding:"required"`
	OrderQuantity int    `json:"order_quantity" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required"`
	DeliveryType  string `json:"delivery_type" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	AddressLine1  string `json:"address_line_1" binding:"required"`
	AddressLine2  string `json:"address_line_2"`
	PostalCode    string `json:"postal_code" binding:"required"`
	Amount        int    `json:"amount" binding:"required"`
}
