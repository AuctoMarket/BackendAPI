package data

type BuyerLoginData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type BuyerSignUpData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type BuyerLoginResponseData struct {
	Email   string `json:"email"`
	BuyerId string `json:"buyer_id"`
}
