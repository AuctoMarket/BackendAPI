package data

type CreateOrderRequestData struct {
	Products       []ProductOrder `json:"products" binding:"required"`
	BuyerId        string         `json:"buyer_id" binding:"required"`
	PhoneNumber    string         `json:"phone_number" binding:"required"`
	AddressLine1   string         `json:"address_line_1" binding:"required"`
	AddressLine2   string         `json:"address_line_2"`
	PostalCode     string         `json:"postal_code" binding:"required"`
	TelegramHandle string         `json:"telegram_handle"`
	Fees           OrderFees      `json:"fees" binding:"required"`
}

type ProductOrder struct {
	ProductId     string `json:"product_id" binding:"required"`
	OrderQuantity int    `json:"order_quantity" binding:"required"`
}

type CreateGuestOrderRequestData struct {
	Products       []ProductOrder `json:"products" binding:"required"`
	Email          string         `json:"email" binding:"required,email"`
	PhoneNumber    string         `json:"phone_number" binding:"required"`
	AddressLine1   string         `json:"address_line_1" binding:"required"`
	AddressLine2   string         `json:"address_line_2"`
	PostalCode     string         `json:"postal_code" binding:"required"`
	TelegramHandle string         `json:"telegram_handle"`
	Fees           OrderFees      `json:"fees" binding:"required"`
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
	OrderId        string         `json:"order_id" binding:"required"`
	Products       []ProductOrder `json:"products" binding:"required"`
	BuyerId        string         `json:"buyer_id" binding:"required"`
	PhoneNumber    string         `json:"phone_number" binding:"required"`
	AddressLine1   string         `json:"address_line_1" binding:"required"`
	AddressLine2   string         `json:"address_line_2"`
	PostalCode     string         `json:"postal_code" binding:"required"`
	TelegramHandle string         `json:"telegram_handle"`
	PaymentStatus  string         `json:"payment_status" binding:"required"`
	OrderDate      string         `json:"order_date" binding:"required"`
	Fees           OrderFees      `json:"fees" binding:"required"`
}

type GetGuestOrderByIdResponseData struct {
	GuestOrderId   string         `json:"guest_order_id" binding:"required"`
	Products       []ProductOrder `json:"products" binding:"required"`
	Email          string         `json:"email" binding:"required"`
	PhoneNumber    string         `json:"phone_number" binding:"required"`
	AddressLine1   string         `json:"address_line_1" binding:"required"`
	AddressLine2   string         `json:"address_line_2"`
	PostalCode     string         `json:"postal_code" binding:"required"`
	TelegramHandle string         `json:"telegram_handle"`
	PaymentStatus  string         `json:"payment_status" binding:"required"`
	OrderDate      string         `json:"order_date" binding:"required"`
	Fees           OrderFees      `json:"fees" binding:"required"`
}

type OrderFees struct {
	PaymentType   string `json:"payment_type" binding:"required"`
	PaymentFee    int    `json:"payment_fee"`
	DeliveryType  string `json:"delivery_type" binding:"required"`
	DeliveryFee   int    `json:"delivery_fee"`
	TotalPaid     int    `json:"total_paid" binding:"required"`
	SmallOrderFee int    `json:"small_order_fee"`
}
