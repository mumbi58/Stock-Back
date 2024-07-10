package controllers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "stock-back/models"
    "stock-back/utils"
)

func AdminLogin(c echo.Context) error {
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

func AddUser(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := c.Validate(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    user := c.Get("user").(models.User)
    if user.RoleID != 2 {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Admins only"})
    }

    if input.RoleID != 3 && input.RoleID != 4 {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role ID. Allowed roles: 3 (shopkeeper), 4 (auditor)"})
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not hash password"})
    }
    input.Password = hashedPassword

    if err := utils.DB.Create(&input).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "User created successfully"})
}

func AdminGetUserByID(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    var user models.User
    if err := utils.DB.First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
    }
    return c.JSON(http.StatusOK, user)
}

func EditUser(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    var user models.User
    if err := utils.DB.First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
    }

    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := utils.DB.Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, user)
}

func SoftDeleteUser(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    var user models.User
    if err := utils.DB.First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
    }

    // Check if user is super admin
    if user.RoleID == 1 {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot delete super admin"})
    }

    user.DeletedAt = gorm.DeletedAt{
        Time:  time.Now(),
        Valid: true,
    }

    if err := utils.DB.Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "User soft deleted successfully"})
}

func GetAllUsers(c echo.Context) error {
    var users []models.User
    if err := utils.DB.Find(&users).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not fetch users"})
    }
    return c.JSON(http.StatusOK, users)
}

func AdminLogout(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
