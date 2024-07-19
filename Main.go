package main

import (
	"awesomeProject9/database"
	"awesomeProject9/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize Echo
	e := echo.New()

	// Connect to the database
	db := database.InitDB()
	defer db.Close()
	// Middleware to enable CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	// Define CRUD endpoints
	e.GET("/categories", handlers.GetCategories)
	e.GET("/categories/:category_id", handlers.GetCategoryByID)
	e.POST("/categories", handlers.CreateCategories)
	e.PUT("/categories/:category_id", handlers.UpdateCategory)
	e.DELETE("/categories/:id", handlers.DeleteCategoryByID)

	///Define CRUD endpoints
	e.GET("/products", handlers.GetProducts)
	e.GET("/products/:product_id", handlers.GetProductByID)
	e.POST("/products", handlers.AddProduct)
	e.PUT("/products/:product_id", handlers.UpdateProduct)
	//e.PUT("/productss/:product_id", handlers.UpdateProducteee)
	e.DELETE("/products/:product_id", handlers.DeleteProduct)
	//Start the server

	///Define CRUD endpoints
	e.GET("/sales", handlers.GetSales)
	e.GET("/sales/:sale_id", handlers.GetSaleByID)
	e.POST("/sales", handlers.AddSale)
	//e.PUT("/sales/:sale_id", handlers.UpdateSale)
	//e.PUT("/productss/:product_id", handlers.UpdateProducteee)
	e.DELETE("/sales/:sale_id", handlers.DeleteSale)
	//Start the server

	e.Logger.Fatal(e.Start(":8000"))
}
