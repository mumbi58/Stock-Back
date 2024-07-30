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
)

func GetProducts(c echo.Context) error {
	log.Println("Received request to fetch products")

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Query all products from the products table, including product_description
	rows, err := db.Query("SELECT product_id, category_name, product_name, product_code, product_description, date, quantity, reorder_level, price FROM products")
	if err != nil {
		log.Printf("Error querying products from database: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()

	// Slice to hold the fetched products
	var products []model.Product

	// Iterate over the query results
	for rows.Next() {
		var prod model.Product

		// Scan each row into the Product struct
		err := rows.Scan(&prod.ProductID, &prod.CategoryName, &prod.ProductName, &prod.ProductCode, &prod.ProductDescription, &prod.Date, &prod.Quantity, &prod.ReorderLevel, &prod.Price)
		if err != nil {
			// Check if the error is due to NULL value in reorder_level or price
			if err.Error() == "sql: Scan error on column index 7, name \"reorder_level\": converting NULL to int is unsupported" {
				// Handle NULL value scenario for reorder_level
				log.Printf("Null value encountered in reorder_level column")
				prod.ReorderLevel = 0 // Or set it to any default value as per your application logic
			} else if err.Error() == "sql: Scan error on column index 8, name \"price\": converting NULL to float64 is unsupported" {
				// Handle NULL value scenario for price
				log.Printf("Null value encountered in price column")
				prod.Price = 0.0 // Or set it to any default value as per your application logic
			} else {
				log.Printf("Error scanning product row: %s", err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
			}
		}

		// Append the Product to the slice
		products = append(products, prod)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over product rows: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	// Log the number of Products fetched
	log.Printf("Fetched %d products", len(products))

	// Return the fetched Products as JSON
	return c.JSON(http.StatusOK, products)
}

// GetProductByID fetches a single product by its ID from the database.
func GetProductByID(c echo.Context) error {
	log.Println("Received request to fetch product by ID")

	// Extract product_id from path parameter
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		log.Printf("Invalid product ID: %s", c.Param("product_id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Query the product by product_id, including product_description
	row := db.QueryRow("SELECT product_id, category_name, product_name, product_code, product_description, date, quantity, reorder_level, price FROM products WHERE product_id = ?", productID)
	var prod model.Product

	// Scan the row into the Product struct
	err = row.Scan(&prod.ProductID, &prod.CategoryName, &prod.ProductName, &prod.ProductCode, &prod.ProductDescription, &prod.Date, &prod.Quantity, &prod.ReorderLevel, &prod.Price)
	if err != nil {
		// Check if the error is due to a "not found" condition
		if err == sql.ErrNoRows {
			log.Printf("Product not found with ID: %d", productID)
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}

		// Handle other scanning errors
		log.Printf("Error scanning product row: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch product")
	}

	// Log the fetched product
	log.Printf("Fetched product: %+v", prod)

	// Return the fetched Product as JSON
	return c.JSON(http.StatusOK, prod)
}
func AddProduct(c echo.Context) error {
	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Parse JSON manually from request body into Product struct
	var product model.Product
	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		log.Printf("Error decoding JSON: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Error decoding JSON")
	}

	// Log the received product details
	log.Printf("Received request to create a product: %+v", product)

	// Execute the SQL INSERT query to add the product to the database, including product_description
	result, err := db.Exec(`
        INSERT INTO products (category_name, product_name, product_code, product_description, date, quantity, reorder_level, price)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, product.CategoryName, product.ProductName, product.ProductCode, product.ProductDescription, product.Date, product.Quantity, product.ReorderLevel, product.Price)
	if err != nil {
		log.Printf("Error inserting product: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error inserting product")
	}

	// Get the ID of the newly inserted product
	productID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting product ID: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting product ID")
	}
	product.ProductID = int(productID)

	// Log the successful creation and the product details
	log.Printf("Product created successfully. Product ID: %d, Category: %s, Name: %s, Code: %s, Description: %s, Date: %s, Quantity: %d, Reorder Level: %d, Price: %.2f",
		product.ProductID, product.CategoryName, product.ProductName, product.ProductCode, product.ProductDescription, product.Date, product.Quantity, product.ReorderLevel, product.Price)

	// Return the created product as JSON with status 201 Created
	return c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c echo.Context) error {
	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Extract product_id from path parameter
	productID := c.Param("product_id")
	log.Printf("Received request to update product with ID: %s", productID)

	// Bind the request payload to a new Product struct
	var updatedProduct model.Product
	if err := c.Bind(&updatedProduct); err != nil {
		log.Printf("Error binding payload: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request body")
	}

	// Execute the SQL UPDATE query to update the product in the database, including product_description
	query := `
        UPDATE products 
        SET category_name = ?, 
            product_name = ?, 
            product_code = ?, 
            product_description = ?, 
            date = ?, 
            quantity = ?, 
            reorder_level = ?, 
            price = ?
        WHERE product_id = ?
    `
	_, err := db.Exec(query,
		updatedProduct.CategoryName,
		updatedProduct.ProductName,
		updatedProduct.ProductCode,
		updatedProduct.ProductDescription, // Added this line
		updatedProduct.Date,
		updatedProduct.Quantity,
		updatedProduct.ReorderLevel,
		updatedProduct.Price,
		productID,
	)
	if err != nil {
		log.Printf("Error updating product with ID %s: %s", productID, err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update product")
	}

	// Log successful update
	log.Printf("Product with ID %s updated successfully", productID)

	// Return success message
	return c.JSON(http.StatusOK, map[string]string{"message": "Product updated successfully"})
}

// DeleteProduct deletes a product from the database by its ID.
func DeleteProduct(c echo.Context) error {
	// Extract product_id from request parameters
	productID, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		log.Printf("Invalid product ID: %s", c.Param("product_id"))
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Prepare statement to delete a product
	query := "DELETE FROM products WHERE product_id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing delete statement: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	defer stmt.Close()

	// Execute the delete operation
	result, err := stmt.Exec(productID)
	if err != nil {
		log.Printf("Error deleting product: %s", err.Error())
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
		log.Printf("Product with ID %d not found", productID)
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	// Log successful deletion
	log.Printf("Deleted product with ID %d", productID)

	// Return success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
