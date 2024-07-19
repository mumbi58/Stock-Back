package controllers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo/v4"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "stock-back/db"
    "stock-back/models"
    "stock-back/utils"
    "fmt"

)

func SuperAdminSignup(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    fmt.Printf("Received JSON: %+v\n", input) // Add this line to check the bound data

    fmt.Printf("Received RoleID: %d\n", input.RoleID)

    if input.RoleID != 1 {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Only super admin can sign up"})
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }
    input.Password = hashedPassword

    if err := db.GetDB().Create(&input).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    token, err := utils.GenerateJWT(input.ID, input.RoleID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func AddAdmin(c echo.Context) error {
    user, ok := c.Get("user").(models.User)
    if !ok {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "User not found in context"})
    }
    if user.RoleID != 1 { // Super Admin RoleID
        return c.JSON(http.StatusForbidden, echo.Map{"error": "Permission denied"})
    }

    var newAdmin models.User
    if err := c.Bind(&newAdmin); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    newAdmin.RoleID = 2 // Admin RoleID

    hashedPassword, err := utils.HashPassword(newAdmin.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }
    newAdmin.Password = hashedPassword

    if err := db.GetDB().Create(&newAdmin).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "Admin added successfully", "admin": newAdmin})
}

func SuperAdminLogin(c echo.Context) error {
    var input struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(user.ID, user.RoleID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    return c.JSON(http.StatusOK, echo.Map{"token": token})
}


func SuperAdminLogout(c echo.Context) error {
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func AdminLogin(c echo.Context) error {
    var input struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(user.ID, user.RoleID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func GetUserByID(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    var user models.User
    if err := db.GetDB().First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }
    return c.JSON(http.StatusOK, user)
}

func AdminAddUser(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    user, ok := c.Get("user").(models.User)
    if !ok {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    if user.RoleID != 2 { // Admin RoleID
        return c.JSON(http.StatusForbidden, echo.Map{"error": "Admins only"})
    }

    if input.RoleID != 3 && input.RoleID != 4 && input.RoleID != 2 { // Shopkeeper, Auditor, or Admin RoleIDs
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid role ID. Allowed roles: 3 (shopkeeper), 4 (auditor), 2 (admin)"})
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }
    input.Password = hashedPassword

    if err := db.GetDB().Create(&input).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "User created successfully"})
}
 
func EditUser(c echo.Context) error {
    id := c.Param("id")
    var user models.User
    if err := db.GetDB().First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }

    if err := c.Bind(&user); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    if err := db.GetDB().Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, user)
}


func SoftDeleteUser(c echo.Context) error {
    id := c.Param("id")
    var user models.User
    if err := db.GetDB().First(&user, id).Error; err != nil {
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }

    user.DeletedAt = gorm.DeletedAt{
        Time:  time.Now(),
        Valid: true,
    }

    if err := db.GetDB().Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, echo.Map{"message": "User soft deleted successfully"})
}

func AdminLogout(c echo.Context) error {
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func Login(c echo.Context) error {
    var loginData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&loginData); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID,
        "roleID": user.RoleID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(utils.JwtSecret) // Updated to use JwtSecret
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

func Logout(c echo.Context) error {
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func AuditorLogin(c echo.Context) error {
    var loginData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&loginData); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID,
        "roleID": user.RoleID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(utils.JwtSecret) // Updated to use JwtSecret
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

func AuditorLogout(c echo.Context) error {
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}
