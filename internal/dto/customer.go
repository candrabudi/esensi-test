package dto

type (
	FindAllCustomer struct {
		ID              int    `json:"id"`
		CustomerName    string `json:"customer_name"`
		CustomerAddress string `json:"customer_address"`
		CustomerPhone   string `json:"customer_phone"`
	}
)
