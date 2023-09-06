package data

type UserLoginData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type BuyerResendOtpData struct {
	BuyerId string `json:"buyer_id" binding:"required"`
}

type BuyerValidateOtpData struct {
	BuyerId string `json:"buyer_id" binding:"required"`
	Otp     string `json:"otp" binding:"required"`
}

type BuyerSignUpData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type BuyerLoginResponseData struct {
	Email        string `json:"email" binding:"required"`
	BuyerId      string `json:"buyer_id" binding:"required"`
	Verification string `json:"verification" binding:"required"`
}

type SellerSignUpData struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	SellerName string `json:"seller_name" binding:"required"`
}

type SellerLoginResponseData struct {
	Email      string `json:"email" binding:"required"`
	SellerId   string `json:"seller_id" binding:"required"`
	SellerName string `json:"seller_name" binding:"required"`
	Followers  int    `json:"followers" binding:"required"`
}

type GetSellerByIdResponseData struct {
	SellerId   string `json:"seller_id" binding:"required"`
	SellerName string `json:"seller_name" binding:"required"`
	Followers  int    `json:"followers" binding:"required"`
}
