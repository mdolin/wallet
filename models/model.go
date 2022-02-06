package model

type Account struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}

type Response struct {
	Type     string    `json:"tpe"`
	Accounts []Account `json:"account"`
	Message  string    `json:"message"`
}
