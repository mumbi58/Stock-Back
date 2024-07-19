package model

type Category struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	ProductName  string `json:"product_name"`
}

type Product struct {
	ProductID    int     `json:"product_id"`
	CategoryName string  `json:"category_name"`
	ProductName  string  `json:"product_name"`
	ProductCode  string  `json:"product_code"`
	Date         string  `json:"date"` // Assuming date is a string in your database
	Quantity     int     `json:"quantity"`
	ReorderLevel int     `json:"reorder_level"`
	Price        float64 `json:"price"`
}

type Sale struct {
	SaleID   int     `json:"sale_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	UserID   string  `json:"user_id"`
}
