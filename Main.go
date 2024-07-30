package main

import (
	"awesomeProject9/database"
	"awesomeProject9/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize Echo
	e := echo.New()

	// Connect to the database
	db := database.InitDB()
	defer db.Close()

	// Middleware to enable CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("CORS_ALLOW_ORIGINS")},
		AllowHeaders: []string{os.Getenv("CORS_ALLOW_HEADERS")},
	}))

	// Define CRUD endpoints
	e.GET("/categories", handlers.GetCategories)
	e.GET("/categories/:category_id", handlers.GetCategoryByID)
	e.POST("/categories", handlers.CreateCategories)
	e.PUT("/categories/:category_id", handlers.UpdateCategory)
	e.DELETE("/categories/:id", handlers.DeleteCategoryByID)

	e.GET("/products", handlers.GetProducts)
	e.GET("/products/:product_id", handlers.GetProductByID)
	e.POST("/products", handlers.AddProduct)
	e.PUT("/products/:product_id", handlers.UpdateProduct)
	e.DELETE("/products/:product_id", handlers.DeleteProduct)

	e.GET("/sales", handlers.GetSales)
	e.GET("/sales/:sale_id", handlers.GetSaleByID)
	e.POST("/sales", handlers.AddSale)
	e.DELETE("/sales/:sale_id", handlers.DeleteSale)
	e.GET("/sales/:category_name", handlers.FetchSalesByCategory)
	e.GET("/sales/:date", handlers.FetchSalesByDate)
	e.GET("/sales/:user_id", handlers.FetchSalesByUserID)

	e.POST("/products/:product_id/sell/:quantity_sold", handlers.SellProduct)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port if not specified
	}
	e.Logger.Fatal(e.Start(":" + port))
}
