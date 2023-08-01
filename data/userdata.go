package data

type UserLoginData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type BuyerSignUpData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type BuyerLoginResponseData struct {
	Email   string `json:"email" binding:"required"`
	BuyerId string `json:"buyer_id" binding:"required"`
}

type SellerSignUpData struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	SellerName string `json:"seller_name" binding:"required"`
}

type SellerResponseData struct {
	Email      string `json:"email" binding:"required"`
	SellerId   string `json:"seller_id" binding:"required"`
	SellerName string `json:"seller_name" binding:"required"`
	Followers  int    `json:"followers" binding:"required"`
}
