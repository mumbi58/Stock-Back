package handlers

import (
	"awesomeProject9/database"
	models "awesomeProject9/models"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// FetchSalesByCategory fetches sales data filtered by category name
func FetchSalesByCategory(c echo.Context) error {
	// Extract category_name from request parameters
	categoryName := c.Param("category_name")
	if categoryName == "" {
		log.Printf("No category name provided in the request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category name is required"})
	}
	log.Printf("Received request to fetch sales for category name: %s", categoryName)

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Query sales data from the sale table filtered by category name
	query := `
        SELECT sale_id, name, price, quantity, user_id, date, category_name
        FROM sale
        WHERE category_name = ?`
	rows, err := db.Query(query, categoryName)
	if err != nil {
		log.Printf("Error querying sales from database: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()

	// Slice to hold the fetched sales data
	var sales []models.SaleByCategory

	// Iterate over the query results
	for rows.Next() {
		var sale models.SaleByCategory

		// Scan each row into the SaleByCategory struct
		if err := rows.Scan(
			&sale.SaleID,
			&sale.Name,
			&sale.Price,
			&sale.Quantity,
			&sale.UserID,
			&sale.Date,
			&sale.CategoryName,
		); err != nil {
			log.Printf("Error scanning sale row: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error processing data"})
		}

		// Append the SaleByCategory to the slice
		sales = append(sales, sale)
	}

	// Check if no sales were found
	if len(sales) == 0 {
		log.Printf("No sales found for category name %s", categoryName)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No sales found for this category"})
	}

	// Log the number of sales fetched
	log.Printf("Fetched %d sales for category name %s", len(sales), categoryName)

	// Return the fetched Sales as JSON
	return c.JSON(http.StatusOK, sales)
}

// FetchSalesByDate fetches sales data filtered by date
func FetchSalesByDate(c echo.Context) error {
	// Extract date from request parameters
	date := c.Param("date")
	if date == "" {
		log.Printf("No date provided in the request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Date is required"})
	}
	log.Printf("Received request to fetch sales for date: %s", date)

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Query sales data from the sale table filtered by date
	query := `
        SELECT sale_id, name, price, quantity, user_id, date, category_name
        FROM sale
        WHERE DATE(date) = ?`
	rows, err := db.Query(query, date)
	if err != nil {
		log.Printf("Error querying sales from database: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()

	// Slice to hold the fetched sales data
	var sales []models.SaleByCategory

	// Iterate over the query results
	for rows.Next() {
		var sale models.SaleByCategory

		// Scan each row into the SaleByCategory struct
		if err := rows.Scan(
			&sale.SaleID,
			&sale.Name,
			&sale.Price,
			&sale.Quantity,
			&sale.UserID,
			&sale.Date,
			&sale.CategoryName,
		); err != nil {
			log.Printf("Error scanning sale row: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error processing data"})
		}

		// Append the SaleByCategory to the slice
		sales = append(sales, sale)
	}

	// Check if no sales were found
	if len(sales) == 0 {
		log.Printf("No sales found for date %s", date)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No sales found for this date"})
	}

	// Log the number of sales fetched
	log.Printf("Fetched %d sales for date %s", len(sales), date)

	// Return the fetched Sales as JSON
	return c.JSON(http.StatusOK, sales)
}

// FetchSalesByUserID fetches sales data filtered by user ID
func FetchSalesByUserID(c echo.Context) error {
	// Extract user_id from request parameters
	userID := c.Param("user_id")
	if userID == "" {
		log.Printf("No user ID provided in the request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	log.Printf("Received request to fetch sales for user ID: %s", userID)

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Query sales data from the sale table filtered by user ID
	query := `
        SELECT sale_id, name, price, quantity, user_id, date, category_name
        FROM sale
        WHERE user_id = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error querying sales from database: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()

	// Slice to hold the fetched sales data
	var sales []models.SaleByCategory

	// Iterate over the query results
	for rows.Next() {
		var sale models.SaleByCategory

		// Scan each row into the SaleByCategory struct
		if err := rows.Scan(
			&sale.SaleID,
			&sale.Name,
			&sale.Price,
			&sale.Quantity,
			&sale.UserID,
			&sale.Date,
			&sale.CategoryName,
		); err != nil {
			log.Printf("Error scanning sale row: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error processing data"})
		}

		// Append the SaleByCategory to the slice
		sales = append(sales, sale)
	}

	// Check if no sales were found
	if len(sales) == 0 {
		log.Printf("No sales found for user ID %s", userID)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No sales found for this user"})
	}

	// Log the number of sales fetched
	log.Printf("Fetched %d sales for user ID %s", len(sales), userID)

	// Return the fetched Sales as JSON
	return c.JSON(http.StatusOK, sales)
}
