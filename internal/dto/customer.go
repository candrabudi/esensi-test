package dto

type (
	InsertCustomer struct {
		CustomerName    string `json:"customer_name"`
		CustomerAddress string `json:"customer_address" binding:"required"`
		CustomerPhone   string `json:"customer_phone" binding:"required"`
	}

	FindAllCustomer struct {
		ID              int    `json:"id"`
		CustomerName    string `json:"customer_name"`
		CustomerAddress string `json:"customer_address"`
		CustomerPhone   string `json:"customer_phone"`
	}
)
