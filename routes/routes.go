package routes

import (
	"awesomeProject9/handlers"
	"github.com/labstack/echo/v4"
)

// RegisterRoutes initializes all the routes for the Echo server
func RegisterRoutes(e *echo.Echo) {
	// Define CRUD endpoints for categories
	e.GET("/categories", handlers.GetCategories)
	e.GET("/categories/:category_id", handlers.GetCategoryByID)
	e.POST("/categories", handlers.CreateCategories)
	e.PUT("/categories/:category_id", handlers.UpdateCategory)
	e.DELETE("/categories/:id", handlers.DeleteCategoryByID)

	// Define CRUD endpoints for products
	e.GET("/products", handlers.GetProducts)
	e.GET("/products/:product_id", handlers.GetProductByID)
	e.POST("/products", handlers.AddProduct)
	e.PUT("/products/:product_id", handlers.UpdateProduct)
	e.DELETE("/products/:product_id", handlers.DeleteProduct)

	// Define CRUD endpoints for sales
	e.GET("/sales", handlers.GetSales)
	e.GET("/sales/:sale_id", handlers.GetSaleByID)
	e.POST("/sales", handlers.AddSale)
	e.DELETE("/sales/:sale_id", handlers.DeleteSale)
	e.GET("/sales/:category_name", handlers.FetchSalesByCategory)
	e.GET("/sales/:date", handlers.FetchSalesByDate)
	e.GET("/sales/:user_id", handlers.FetchSalesByUserID)

	// Endpoint for selling products
	e.POST("/products/:product_id/sell/:quantity_sold", handlers.SellProduct)
}
