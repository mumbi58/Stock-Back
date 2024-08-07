package routes

import (
	"awesomeProject9/handlers"
	"awesomeProject9/middleware" // Import the middleware package
	"github.com/labstack/echo/v4"
)

// RegisterRoutes initializes all the routes for the Echo server
func RegisterRoutes(e *echo.Echo) {
	// Define CRUD endpoints for categories with admin middleware
	categoryGroup := e.Group("/categories")
	categoryGroup.Use(middleware.AdminMiddleware) // Apply middleware
	categoryGroup.GET("", handlers.GetCategories)
	categoryGroup.GET("/:category_id", handlers.GetCategoryByID)
	categoryGroup.POST("", handlers.CreateCategories)
	categoryGroup.PUT("/:category_id", handlers.UpdateCategory)
	categoryGroup.DELETE("/:id", handlers.DeleteCategoryByID)

	// Define CRUD endpoints for products with admin middleware
	productGroup := e.Group("/products")
	productGroup.Use(middleware.AdminMiddleware) // Apply middleware
	productGroup.GET("", handlers.GetProducts)
	productGroup.GET("/:product_id", handlers.GetProductByID)
	productGroup.POST("", handlers.AddProduct)
	productGroup.PUT("/:product_id", handlers.UpdateProduct)
	productGroup.DELETE("/:product_id", handlers.DeleteProduct)
	productGroup.DELETE("/:product_id/pending-deletion", handlers.MoveProductToPendingDeletion)
	productGroup.PUT("/:product_id/recover", handlers.MoveProductFromPendingDeletion)

	// Define CRUD endpoints for sales
	e.GET("/sales", handlers.GetSales)
	e.GET("/sales/:sale_id", handlers.GetSaleByID)
	e.POST("/sales", handlers.AddSale)
	e.DELETE("/sales/:sale_id", handlers.DeleteSale)
	e.GET("/salebycategory/:category_name", handlers.FetchSalesByCategory)
	e.GET("/salebycategory/:date", handlers.FetchSalesByDate)
	e.GET("/salebycategory/:user_id", handlers.FetchSalesByUserID)

	// Endpoint for selling products
	e.POST("/products/:product_id/sell/:quantity_sold", handlers.SellProduct)
}
