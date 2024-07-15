package controllers

import (
    "net/http"
    "stock-back/models"
    "stock-back/utils"
    "stock-back/middlewares"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "strings"
)

// SuperAdminSignup handles the signup request for super admins
func SuperAdminSignup(c echo.Context) error {
    // Extract DB instance from context
    db := middlewares.ExtractDB(c)

    // Parse request body
    var newUser models.User
    if err := c.Bind(&newUser); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    // Validate and set default role for super admin
    if newUser.RoleID == 0 {
        // Set role ID for super admin (adjust according to your role definitions)
        newUser.RoleID = 1 // Assuming 1 is the ID for super admin
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(newUser.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
    }
    newUser.Password = hashedPassword

    // Save user to database
    if err := db.Create(&newUser).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
    }

    // Generate JWT token upon successful signup
    token, err := utils.GenerateJWT(newUser.ID, newUser.RoleID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
    }

    // Return token in response along with user details
    response := map[string]interface{}{
        "token": token,
        "user":  newUser, // Optionally return user details if needed
    }
    return c.JSON(http.StatusCreated, response)
}

// AddAdmin handles the request to add an admin
func AddAdmin(c echo.Context) error {
    // Extract JWT token from request header
    token := c.Request().Header.Get("Authorization")
    if token == "" {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
    }

    // Parse and validate JWT token
    claims, err := utils.VerifyJWT(strings.Replace(token, "Bearer ", "", 1)) // Strip "Bearer " prefix
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired JWT"})
    }

    // Only allow super admins to add admins
    if claims.RoleID != 1 { // Assuming 1 is the RoleID for super admin
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Not authorized to add admins"})
    }

    // Extract DB instance from context
    db := middlewares.ExtractDB(c)

    // Parse request body
    var newUser models.User
    if err := c.Bind(&newUser); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    // Validate and set default role for admin
    if newUser.RoleID == 0 {
        // Set role ID for admin (adjust according to your role definitions)
        newUser.RoleID = 2 // Assuming 2 is the ID for admin
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(newUser.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
    }
    newUser.Password = hashedPassword

    // Save user to database
    if err := db.Create(&newUser).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create admin"})
    }

    // Return success response
    return c.JSON(http.StatusCreated, newUser)
}


// SuperAdminAddOrganization handles the request to add an organization by superadmin
func SuperAdminAddOrganization(c echo.Context) error {
    // Extract DB instance from context
    db := middlewares.ExtractDB(c)

    // Parse request body
    var newOrganization models.Organization
    if err := c.Bind(&newOrganization); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
    }

    // Validate and save organization to database
    if err := db.Create(&newOrganization).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create organization"})
    }

    return c.JSON(http.StatusCreated, newOrganization)
}

// SuperAdminLogin handles the login process for super admins
// SuperAdminLogin handles the login process for super admins
func SuperAdminLogin(c echo.Context) error {
    db := c.Get("db").(*gorm.DB)

    var input struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required"`
    }

    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    // Validate input using Echo's validator
    if err := c.Validate(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    var user models.User
    if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
    }

    // Generate JWT token upon successful login
    token, err := utils.GenerateJWT(user.ID, user.RoleID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
    }

    // Return token in response
    return c.JSON(http.StatusOK, map[string]string{"token": token})
}


// SuperAdminLogout handles the logout process for super admins
func SuperAdminLogout(c echo.Context) error {
    // Implement logout logic if needed, e.g., invalidate JWT token, clear cookies, etc.
    return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
