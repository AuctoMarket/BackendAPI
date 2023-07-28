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
<<<<<<< HEAD
<<<<<<< HEAD
	Email    string `json:"email"`
	Buyer_id string `json:"buyer_id"`
=======
	Email string `json:"email"`
	BUID  string `json:"buid"`
>>>>>>> 005bc68 (Add login and signup API)
=======
	Email    string `json:"email"`
	Buyer_id string `json:"buyer_id"`
>>>>>>> 1840890 (Update id tables)
}
