package handlers

import (
	"awesomeProject9/database"
	model "awesomeProject9/models"
	"database/sql"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
) // SellProduct handles the sale of a product and updates the quantity in the database.

func SellProduct(c echo.Context) error {
	productIDStr := c.Param("product_id")
	quantitySoldStr := c.Param("quantity_sold")
	userID := c.QueryParam("user_id") // Assuming user_id is passed as a query parameter

	// Convert parameters to integer values
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", productIDStr)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	quantitySold, err := strconv.Atoi(quantitySoldStr)
	if err != nil {
		log.Printf("Invalid quantity sold: %s", quantitySoldStr)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity sold")
	}

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	defer tx.Rollback()

	// Check current quantity and retrieve category of the product
	var product model.Product
	err = tx.QueryRow(`
        SELECT product_id, category_name, product_name, product_code, date, quantity, reorder_level, price 
        FROM products 
        WHERE product_id = ?
    `, productID).Scan(
		&product.ProductID,
		&product.CategoryName,
		&product.ProductName,
		&product.ProductCode,
		&product.Date,
		&product.Quantity,
		&product.ReorderLevel,
		&product.Price,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Product not found with ID: %d", productID)
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		log.Printf("Error querying product details: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Check if enough quantity is available
	if product.Quantity < quantitySold {
		log.Printf("Insufficient quantity for product ID %d: Available %d, Requested %d", productID, product.Quantity, quantitySold)
		return echo.NewHTTPError(http.StatusBadRequest, "Insufficient quantity")
	}

	// Update the quantity of the product
	newQuantity := product.Quantity - quantitySold
	_, err = tx.Exec("UPDATE products SET quantity = ? WHERE product_id = ?", newQuantity, productID)
	if err != nil {
		log.Printf("Error updating product quantity: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Insert sale record into the sale table, including category_name
	_, err = tx.Exec(`
        INSERT INTO sale (name, price, quantity, user_id, date, category_name)
        VALUES (?, ?, ?, ?, ?, ?)
    `, product.ProductName, product.Price, quantitySold, userID, time.Now(), product.CategoryName)
	if err != nil {
		log.Printf("Error inserting sale record: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Log successful sale
	log.Printf("Sold %d units of product ID %d. Remaining quantity: %d", quantitySold, productID, newQuantity)

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{
		"message":       "Sale processed successfully",
		"product_id":    strconv.Itoa(productID),
		"quantity_sold": strconv.Itoa(quantitySold),
		"remaining_qty": strconv.Itoa(newQuantity),
	})
}

func SellProducttttttttyy(c echo.Context) error {
	productIDStr := c.Param("product_id")
	quantitySoldStr := c.Param("quantity_sold")
	userID := c.QueryParam("user_id") // Assuming user_id is passed as a query parameter

	// Convert parameters to integer values
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", productIDStr)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	quantitySold, err := strconv.Atoi(quantitySoldStr)
	if err != nil {
		log.Printf("Invalid quantity sold: %s", quantitySoldStr)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity sold")
	}

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	defer tx.Rollback()

	// Check current quantity of the product
	var product model.Product
	err = tx.QueryRow("SELECT product_id, category_name, product_name, product_code, date, quantity, reorder_level, price FROM products WHERE product_id = ?", productID).Scan(
		&product.ProductID,
		&product.CategoryName,
		&product.ProductName,
		&product.ProductCode,
		&product.Date,
		&product.Quantity,
		&product.ReorderLevel,
		&product.Price,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Product not found with ID: %d", productID)
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		log.Printf("Error querying product details: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Check if enough quantity is available
	if product.Quantity < quantitySold {
		log.Printf("Insufficient quantity for product ID %d: Available %d, Requested %d", productID, product.Quantity, quantitySold)
		return echo.NewHTTPError(http.StatusBadRequest, "Insufficient quantity")
	}

	// Update the quantity of the product
	newQuantity := product.Quantity - quantitySold
	_, err = tx.Exec("UPDATE products SET quantity = ? WHERE product_id = ?", newQuantity, productID)
	if err != nil {
		log.Printf("Error updating product quantity: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Insert sale record into the sale table
	_, err = tx.Exec(`
        INSERT INTO sale (name, price, quantity, user_id, date)
        VALUES (?, ?, ?, ?, ?)
    `, product.ProductName, product.Price, quantitySold, userID, time.Now())
	if err != nil {
		log.Printf("Error inserting sale record: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Log successful sale
	log.Printf("Sold %d units of product ID %d. Remaining quantity: %d", quantitySold, productID, newQuantity)

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{
		"message":       "Sale processed successfully",
		"product_id":    strconv.Itoa(productID),
		"quantity_sold": strconv.Itoa(quantitySold),
		"remaining_qty": strconv.Itoa(newQuantity),
	})
}

func SellProductt(c echo.Context) error {
	productIDStr := c.Param("product_id")
	quantitySoldStr := c.Param("quantity_sold")
	userID := c.QueryParam("user_id") // Assuming user_id is passed as a query parameter

	// Convert parameters to integer values
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		log.Printf("Invalid product ID: %s", productIDStr)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	quantitySold, err := strconv.Atoi(quantitySoldStr)
	if err != nil {
		log.Printf("Invalid quantity sold: %s", quantitySoldStr)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid quantity sold")
	}

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	defer tx.Rollback()

	// Check current quantity of the product
	var product model.Product
	err = tx.QueryRow("SELECT product_id, category_name, product_name, product_code, date, quantity, reorder_level, price FROM products WHERE product_id = ?", productID).Scan(
		&product.ProductID,
		&product.CategoryName,
		&product.ProductName,
		&product.ProductCode,
		&product.Date,
		&product.Quantity,
		&product.ReorderLevel,
		&product.Price,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Product not found with ID: %d", productID)
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		log.Printf("Error querying product details: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Check if enough quantity is available
	if product.Quantity < quantitySold {
		log.Printf("Insufficient quantity for product ID %d: Available %d, Requested %d", productID, product.Quantity, quantitySold)
		return echo.NewHTTPError(http.StatusBadRequest, "Insufficient quantity")
	}

	// Update the quantity of the product
	newQuantity := product.Quantity - quantitySold
	_, err = tx.Exec("UPDATE products SET quantity = ? WHERE product_id = ?", newQuantity, productID)
	if err != nil {
		log.Printf("Error updating product quantity: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Insert sale record into the correct table
	sale := model.Sale{
		Name:     product.ProductName,
		Price:    product.Price,
		Quantity: quantitySold,
		UserID:   userID,
		Date:     time.Now(),
	}

	_, err = tx.Exec(`
        INSERT INTO sale (name, price, quantity, user_id, date)
        VALUES (?, ?, ?, ?, ?)
    `, sale.Name, sale.Price, sale.Quantity, sale.UserID, sale.Date.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Printf("Error inserting sale record: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Log successful sale
	log.Printf("Sold %d units of product ID %d. Remaining quantity: %d", quantitySold, productID, newQuantity)

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{
		"message":       "Sale processed successfully",
		"product_id":    strconv.Itoa(productID),
		"quantity_sold": strconv.Itoa(quantitySold),
		"remaining_qty": strconv.Itoa(newQuantity),
	})
}

func GetSales(c echo.Context) error {
	log.Println("Received request to fetch sales")

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Query all sales from the Sales table
	rows, err := db.Query("SELECT sale_id, name, price, quantity, user_id FROM sale")
	if err != nil {
		log.Printf("Error querying sales from database: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()

	// Slice to hold the fetched sales
	var sales []model.Sale

	// Iterate over the query results
	for rows.Next() {
		var sale model.Sale

		// Scan each row into the Sale struct
		err := rows.Scan(&sale.SaleID, &sale.Name, &sale.Price, &sale.Quantity, &sale.UserID)
		if err != nil {
			log.Printf("Error scanning sale row: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
		}

		// Append the Sale to the slice
		sales = append(sales, sale)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over sale rows: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Log the number of Sales fetched
	log.Printf("Fetched %d sales", len(sales))

	// Return the fetched Sales as JSON
	return c.JSON(http.StatusOK, sales)
}

// GetSaleByID fetches a single sale by its ID from the database.
func GetSaleByID(c echo.Context) error {
	log.Println("Received request to fetch sale by ID")

	// Extract sale ID from path parameter
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Invalid sale ID: %s", c.Param("id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid sale ID")
	}

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Query the sale by sale ID
	row := db.QueryRow("SELECT sale_id, name, price, quantity, user_id FROM sale WHERE sale_id = ?", saleID)
	var sale model.Sale

	// Scan the row into the Sale struct
	err = row.Scan(&sale.SaleID, &sale.Name, &sale.Price, &sale.Quantity, &sale.UserID)
	if err != nil {
		// Check if the error is due to a "not found" condition
		if err == sql.ErrNoRows {
			log.Printf("Sale not found with ID: %d", saleID)
			return echo.NewHTTPError(http.StatusNotFound, "Sale not found")
		}

		// Handle other scanning errors
		log.Printf("Error scanning sale row: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch sale")
	}

	// Log the fetched sale
	log.Printf("Fetched sale: %+v", sale)

	// Return the fetched Sale as JSON
	return c.JSON(http.StatusOK, sale)
}

func AddSale(c echo.Context) error {
	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Parse JSON manually from request body into Sale struct
	var sale model.Sale
	if err := json.NewDecoder(c.Request().Body).Decode(&sale); err != nil {
		log.Printf("Error decoding JSON: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Error decoding JSON")
	}

	// Log the received sale details
	log.Printf("Received request to create a sale: %+v", sale)

	// Execute the SQL INSERT query to add the sale to the database
	result, err := db.Exec(`
        INSERT INTO sale (name, price, quantity, user_id)
        VALUES (?, ?, ?, ?)
    `, sale.Name, sale.Price, sale.Quantity, sale.UserID)
	if err != nil {
		log.Printf("Error inserting sale: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting sale")
	}

	// Get the ID of the newly inserted sale
	saleID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting sale ID: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting sale ID")
	}
	sale.SaleID = int(saleID)
	// Log the successful creation and the sale details
	log.Printf("Sale created successfully. Sale ID: %d, Name: %s, Price: %.2f, Quantity: %d, UserID: %s",
		sale.SaleID, sale.Name, sale.Price, sale.Quantity, sale.UserID)
	// Return the created sale as JSON with status 201 Created
	return c.JSON(http.StatusCreated, sale)
}

func DeleteSale(c echo.Context) error {
	// Extract sale_id from request parameters
	saleID := c.Param("sale_id")
	log.Printf("Received request to delete sale with ID: %s", saleID)

	// Convert sale_id to integer
	saleIDInt, err := strconv.Atoi(saleID)
	if err != nil {
		log.Printf("Invalid sale ID: %s", saleID)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid sale ID")
	}

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Prepare statement to delete a sale
	query := "DELETE FROM sales WHERE sale_id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing delete statement: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	defer stmt.Close()

	// Execute the delete operation
	result, err := stmt.Exec(saleIDInt)
	if err != nil {
		log.Printf("Error deleting sale with ID %s: %s", saleID, err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Check if any rows were affected
	if rowsAffected == 0 {
		log.Printf("Sale with ID %d not found", saleIDInt)
		return echo.NewHTTPError(http.StatusNotFound, "Sale not found")
	}

	// Log successful deletion
	log.Printf("Deleted sale with ID %d", saleIDInt)

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Sale deleted successfully"})
}
