package handlers

import (
	"awesomeProject9/database"
	model "awesomeProject9/models"
	"database/sql"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// / Get all Categories
func GetCategories(c echo.Context) error {
	log.Println("Received request to fetch categories")
	// Initialize database connection
	db := database.InitDB()
	defer db.Close()
	// Query all categories from the Categories table
	rows, err := db.Query("SELECT category_id, category_name, product_name FROM Categories")
	if err != nil {
		log.Printf("Error querying categories from database: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer rows.Close()
	// Slice to hold the fetched categories
	var categories []model.Category
	// Iterate over the query results
	for rows.Next() {
		var cat model.Category

		// Scan each row into the Category struct
		if err := rows.Scan(&cat.CategoryID, &cat.CategoryName, &cat.ProductName); err != nil {
			log.Printf("Error scanning category row: %s", err.Error())
			return err
		}
		//Append the Category to the slice
		categories = append(categories, cat)
	}
	// Log the number of Categories fetched
	log.Printf("Fetched %d categories", len(categories))
	// Return the fetched Categories as JSON
	return c.JSON(http.StatusOK, categories)
}

func GetCategoryByID(c echo.Context) error {
	// Extract category_id from request parameters
	categoryID := c.Param("category_id")
	if categoryID == "" {
		log.Printf("No category ID provided in the request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category ID is required"})
	}
	log.Printf("Received request to fetch category with ID: %s", categoryID)

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()
	// Query the category from the Categories table
	query := "SELECT category_id, category_name, product_name FROM Categories WHERE category_id = ?"
	row := db.QueryRow(query, categoryID)

	// Initialize a Category struct to hold the fetched category
	var category model.Category

	// Scan the row into the Category struct
	err := row.Scan(&category.CategoryID, &category.CategoryName, &category.ProductName)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Category with ID %s not found", categoryID)
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
		}
		log.Printf("Error scanning category row: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}

	// Log the fetched Category
	log.Printf("Fetched category with"+
		" ID %s: %+v", categoryID, category)

	// Return the fetched Category as JSON
	return c.JSON(http.StatusOK, category)
}

func CreateCategories(c echo.Context) error {
	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Parse JSON manually from request body
	var category model.Category // Use full package path or imported name
	if err := json.NewDecoder(c.Request().Body).Decode(&category); err != nil {
		log.Printf("Error decoding JSON: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Error decoding JSON")
	}
	// Log the received category details
	log.Printf("Received request to create a category: %+v")

	// Execute the SQL INSERT query to add the category to the database
	result, err := db.Exec("INSERT INTO Categories (category_id, category_name, product_name) VALUES (?, ?, ?)",
		category.CategoryID,
		category.CategoryName,
		category.ProductName,
	)
	if err != nil {
		log.Printf("Error inserting a category: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Error inserting a category",
		})
	}

	// Get the ID of the newly inserted category
	categoryID, _ := result.LastInsertId()
	category.CategoryID = int(categoryID)

	// Log the successful creation and the category details
	log.Printf("Category created successfully. category_id: %d, category_name: %s, product_name: %s",
		category.CategoryID, category.CategoryName, category.ProductName)

	// Return the created category as JSON with status 201 Created
	return c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c echo.Context) error {
	// Initialize database connection
	db := database.InitDB()
	defer db.Close()
	// Extract the category ID from the request parameters
	categoryID := c.Param("category_id")

	// Bind the request payload to a new Category struct
	category := new(model.Category)
	if err := c.Bind(category); err != nil {
		log.Printf("Error binding payload: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Error binding payload")
	}
	// Execute the SQL UPDATE query to update the category in the database
	_, err := db.Exec("UPDATE Categories SET category_name = ?, product_name = ? WHERE category_id = ?",
		category.CategoryName,
		category.ProductName,
		categoryID,
	)
	if err != nil {
		log.Printf("Error updating a category: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error updating category")
	}

	// Return success message
	log.Printf("Category updated successfully")
	return c.JSON(http.StatusOK, "Category updated successfully")
}

func DeleteCategoryByID(c echo.Context) error {
	// Extract category_id from request parameters
	categoryID := c.Param("id")
	log.Printf("Received request to delete category with ID: %s", categoryID)

	// Initialize database connection
	db := database.InitDB()
	defer db.Close()

	// Prepare statement to delete a category
	query := "DELETE FROM Categories WHERE category_id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing delete statement: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	defer stmt.Close()

	// Execute the delete operation
	result, err := stmt.Exec(categoryID)
	if err != nil {
		log.Printf("Error deleting category: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
	// Check if any rows were affected
	if rowsAffected == 0 {
		log.Printf("Category with ID %s not found", categoryID)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Category not found"})
	}
	// Log successful deletion
	log.Printf("Deleted category with ID %s", categoryID)
	// Return success response
	return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
