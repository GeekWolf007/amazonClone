package main

type User struct {
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	IsLogged   bool       `json:"is_logged"`
	IsAdmin    bool       `json:"is_admin"`
	LoginToken string     `json:"login_token"`
	Cart       []CartItem `json:"cart"`
}

type CartItem struct {
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
}

type Product struct {
	ID                 string `json:"id" bson:"id"`
	ProductName        string `json:"product_name"`
	ProductPrice       int    `json:"product_price"`
	ProductDescription string `json:"product_description"`
}
