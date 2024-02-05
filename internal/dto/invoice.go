package dto

type (
	ResultInvoice struct {
		Items    []FindAllInvoice `json:"items"`
		Metadata MetaData         `json:"metadata"`
	}

	FindAllInvoice struct {
		ID           int     `json:"id"`
		InvoiceNo    string  `json:"invoice_no"`
		Subject      string  `json:"subject"`
		IssueDate    string  `json:"issue_date"`
		CustomerName string  `json:"customer_name"`
		TotalItem    int     `json:"total_item"`
		SubTotal     float64 `json:"subtotal"`
		GrandTotal   float64 `json:"grandtotal"`
		DueDate      string  `json:"due_date"`
		Status       string  `json:"status"`
	}

	MetaData struct {
		TotalResults int `json:"total_results"`
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		Count        int `json:"count"`
	}

	InsertInvoiceDetail struct {
		ItemID    int     `json:"item_id"`
		ItemName  string  `json:"item_name"`
		Quantity  int     `json:"quantity"`
		DueDate   string  `json:"due_date"`
		UnitPrice float64 `json:"unit_price"`
		Amount    float64 `json:"amount"`
		IsDelete  int     `json:"is_delete"`
	}

	InsertInvoice struct {
		Subject         string                `json:"subject"`
		IssueDate       string                `json:"issue_date"`
		CustomerName    string                `json:"customer_name"`
		CustomerAddress string                `json:"customer_address"`
		SubTotal        float64               `json:"subtotal"`
		GrandTotal      float64               `json:"grandtotal"`
		Items           []InsertInvoiceDetail `json:"items"`
	}

	DetailInvoice struct {
		ID              int             `json:"id"`
		InvoiceNo       string          `json:"invoice_no"`
		CustomerName    string          `json:"customer_name"`
		CustomerAddress string          `json:"customer_address"`
		Subject         string          `json:"subject"`
		IssueDate       string          `json:"issue_date"`
		GrandTotal      float64         `json:"grand_total"`
		SubTotal        float64         `json:"sub_total"`
		TotalItem       int             `json:"total_item"`
		DueDate         string          `json:"due_date"`
		Status          string          `json:"status"`
		Items           []InvoiceDetail `json:"items"`
	}

	InvoiceDetail struct {
		ItemID    int     `json:"item_id"`
		ItemName  string  `json:"item_name"`
		Quantity  int     `json:"quantity"`
		UnitPrice float64 `json:"unit_price"`
		Amount    float64 `json:"amount"`
	}
)
