package data

type PaymentRequestData struct {
	Amount         float64  `json:"amount" binding:"required"`
	Currency       string   `json:"currency" binding:"required"`
	RedirectUrl    string   `json:"redirect_url" binding:"required"`
	Webhook        string   `json:"webhook" binding:"required"`
	PaymentMethods []string `json:"payment_methods" binding:"required"`
}

type PaymentRequestResponseData struct {
	Amount         float64  `json:"amount" binding:"required"`
	Currency       string   `json:"currency" binding:"required"`
	RedirectUrl    string   `json:"redirect_url" binding:"required"`
	Webhook        string   `json:"webhook" binding:"required"`
	PaymentMethods []string `json:"payment_methods" binding:"required"`
	Url            string   `json:"url" binding:"required"`
}

type PaymentValidationRequestData struct {
	Hmac   string `form:"hmac"`
	Status string `form:"status"`
}
