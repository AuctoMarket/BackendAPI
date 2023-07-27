package data

type LoginData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignUpData struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseData struct {
	Email string `json:"email"`
	BUID  string `json:"buid"`
}
