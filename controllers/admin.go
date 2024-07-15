package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	//"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"stock-back/db" // Import the db package where your db instance is defined
	"stock-back/models"
	"stock-back/utils"
)

// AdminLogin handles login for admins
func AdminLogin(c echo.Context) error {
	var loginData struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&loginData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&loginData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Fetch the admin user by username
	var adminUser models.User
	if err := db.GetDB().Where("username = ?", loginData.Username).First(&adminUser).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(adminUser.Password), []byte(loginData.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(adminUser.ID, adminUser.RoleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// Return token in response
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

// AdminLogout handles logout for admins
func AdminLogout(c echo.Context) error {
	// No action needed for logout since it's token-based authentication
	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}

// AdminAddUser allows admins to add new users
func AdminAddUser(c echo.Context) error {
	var newUser models.User
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not hash password"})
	}
	newUser.Password = hashedPassword

	// Create the user in the database
	if err := db.GetDB().Create(&newUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return success response
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "User added successfully", "user_id": strconv.Itoa(int(newUser.ID))})
}

// AdminGetUserByID fetches a user by ID for admins
func AdminGetUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var user models.User
	if err := db.GetDB().First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// AdminDeleteUser soft-deletes a user by ID for admins
func SoftDeleteUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	// Fetch the user
	var user models.User
	if err := db.GetDB().First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Check if the user to be deleted is not a super admin (RoleID 1)
	if user.RoleID == 1 {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot delete super admin"})
	}

	// Soft delete user
	if err := db.GetDB().Delete(&user, userID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// AdminUpdateUser updates user details for admins
func EditUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var updateUser models.User
	if err := c.Bind(&updateUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// Find the existing user
	var existingUser models.User
	if err := db.GetDB().First(&existingUser, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Update user fields
	existingUser.Username = updateUser.Username
	existingUser.Email = updateUser.Email
	existingUser.FirstName = updateUser.FirstName
	existingUser.LastName = updateUser.LastName

	// Save updated user details
	if err := db.GetDB().Save(&existingUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, existingUser)
}

// AdminListUsers returns a list of all users for admins
func AdminListUsers(c echo.Context) error {
	var users []models.User
	if err := db.GetDB().Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}

	return c.JSON(http.StatusOK, users)
}
