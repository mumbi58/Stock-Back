package controllers

import (
    "net/http"
    "stock-back/models"
    "stock-back/utils"

    "github.com/labstack/echo/v4"
)

func SuperAdminSignup(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if input.RoleID != 1 {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Only super admin can sign up"})
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not hash password"})
    }
    input.Password = hashedPassword

    if err := utils.DB.Create(&input).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    token, err := utils.GenerateJWT(input.ID, input.RoleID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
    }

    return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func AddAdmin(c echo.Context) error {
    user := c.Get("user")
    currentUser, ok := user.(models.User)
    if !ok {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
    }

    if currentUser.RoleID != models.SuperAdminRoleID {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Permission denied"})
    }

    var newAdmin models.User
    if err := c.Bind(&newAdmin); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    newAdmin.RoleID = models.AdminRoleID

    hashedPassword, err := utils.HashPassword(newAdmin.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not hash password"})
    }
    newAdmin.Password = hashedPassword

    if err := utils.DB.Create(&newAdmin).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]interface{}{"message": "Admin added successfully", "admin": newAdmin})
}

func SuperAdminLogin(c echo.Context) error {
    var input struct {
        Email    string `json:"email" validate:"required"`
        Password string `json:"password" validate:"required"`
    }

    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := c.Validate(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    var user models.User
    if err := utils.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(user.ID, user.RoleID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
    }

    return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func SuperAdminLogout(c echo.Context) error {
    // Implement logout logic if needed, e.g., invalidate JWT token, clear cookies, etc.
    return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
