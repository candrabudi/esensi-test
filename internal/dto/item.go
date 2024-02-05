package dto

type (
	InsertItem struct {
		ItemName        string `json:"item_name"`
		ItemDescription string `json:"item_description" binding:"required"`
		ItemType        string `json:"item_type" binding:"required"`
		ItemPrice       int    `json:"item_price" binding:"required"`
		ItemStock       int    `json:"item_stock" binding:"required"`
	}

	FindAllItem struct {
		ItemName        string `json:"item_name"`
		ItemDescription string `json:"item_description"`
		ItemType        string `json:"item_type"`
		ItemPrice       int    `json:"item_price"`
		ItemStock       int    `json:"item_stock"`
	}
)
