package dto

type (
	InsertInvoiceDetail struct {
		ItemID    int     `json:"item_id"`
		ItemName  string  `json:"item_name"`
		Quantity  int     `json:"quantity"`
		DueDate   string  `json:"due_date"`
		UnitPrice float64 `json:"unit_price"`
		Amount    float64 `json:"amount"`
	}

	InsertInvoice struct {
		Subject      string                `json:"subject"`
		IssueDate    string                `json:"issue_date"`
		CustomerName string                `json:"customer_name"`
		SubTotal     float64               `json:"subtotal"`
		GrandTotal   float64               `json:"grandtotal"`
		Items        []InsertInvoiceDetail `json:"items"`
	}
)
