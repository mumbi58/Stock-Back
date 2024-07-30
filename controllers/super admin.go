package controllers

import (
    "net/http"
    "strconv"
    "time"
    "log"

    "github.com/dgrijalva/jwt-go"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "stock-back/db"
    "stock-back/models"
    "stock-back/utils"
)

func SuperAdminSignup(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    log.Printf("Received JSON: %+v", input)
    log.Printf("Received RoleID: %d", input.RoleID)

    if input.RoleID != 1 {
        log.Println("Unauthorized signup attempt detected")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Only super admin can sign up"})
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        log.Printf("HashPassword error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }
    input.Password = hashedPassword

    if err := db.GetDB().Create(&input).Error; err != nil {
        log.Printf("Create error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    token, err := utils.GenerateJWT(input.ID, input.RoleID)
    if err != nil {
        log.Printf("GenerateJWT error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("Super admin signed up successfully")
    return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func AddAdmin(c echo.Context) error {
    log.Println("AddAdmin called")

    // Retrieve roleID and userID from context set by middleware
    roleID, ok := c.Get("roleID").(int)
    if !ok {
        log.Println("Failed to get roleID from context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    userID, ok := c.Get("userID").(int)
    if !ok {
        log.Println("Failed to get userID from context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    log.Printf("Received RoleID: %d, UserID: %d", roleID, userID)

    // Check if the roleID is 1 (Super Admin)
    if roleID != 1 {
        log.Println("Permission denied: non-super admin trying to add admin")
        return c.JSON(http.StatusForbidden, echo.Map{"error": "Permission denied"})
    }

    var newAdmin models.User
    if err := c.Bind(&newAdmin); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    log.Printf("New admin data: %+v", newAdmin)

    newAdmin.RoleID = 2 // Set roleID for new admin

    hashedPassword, err := utils.HashPassword(newAdmin.Password)
    if err != nil {
        log.Printf("HashPassword error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }
    newAdmin.Password = hashedPassword

    log.Printf("Saving new admin to database")

    if err := db.GetDB().Create(&newAdmin).Error; err != nil {
        log.Printf("Create error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("Admin added successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "Admin added successfully", "admin": newAdmin})
}

func SuperAdminAddOrganization(c echo.Context) error {
    userID, ok := c.Get("userID").(int)
    if !ok {
        log.Println("Failed to get userID from context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    roleID, ok := c.Get("roleID").(int)
    if !ok {
        log.Println("Failed to get roleID from context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    // Log the received userID and roleID
    log.Printf("Received UserID: %d, RoleID: %d", userID, roleID)

    // Check if the roleID is for a Super Admin (roleID = 1)
    if roleID != 1 {
        log.Println("Permission denied: non-super admin trying to add organization")
        return c.JSON(http.StatusForbidden, echo.Map{"error": "Permission denied"})
    }

    var newOrganization models.Organization
    if err := c.Bind(&newOrganization); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    newOrganization.RoleID = 5

    if err := db.GetDB().Create(&newOrganization).Error; err != nil {
        log.Printf("Create error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("Organization added successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "Organization added successfully", "organization": newOrganization})
}

func SuperAdminLogin(c echo.Context) error {
    var input struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&input); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
        log.Printf("Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        log.Printf("CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(user.ID, user.RoleID)
    if err != nil {
        log.Printf("GenerateJWT error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("Super admin logged in successfully")
    return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func SuperAdminLogout(c echo.Context) error {
    log.Println("Super admin logged out successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func AdminLogin(c echo.Context) error {
    var input struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&input); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
        log.Printf("Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        log.Printf("CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(user.ID, user.RoleID)
    if err != nil {
        log.Printf("GenerateJWT error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("Admin logged in successfully")
    return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func GetUserByID(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    log.Printf("GetUserByID called with ID: %d", id)
    
    var user models.User
    if err := db.GetDB().First(&user, id).Error; err != nil {
        log.Printf("First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }
    
    log.Printf("User found: %+v", user)
    return c.JSON(http.StatusOK, user)
}

func AdminAddUser(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    user, ok := c.Get("user").(models.User)
    if !ok {
        log.Println("Unauthorized: user not found in context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    if user.RoleID != 2 {
        log.Println("Admins only: unauthorized role")
        return c.JSON(http.StatusForbidden, echo.Map{"error": "Admins only"})
    }

    if input.RoleID != 3 && input.RoleID != 4 && input.RoleID != 2 {
        log.Println("Invalid role ID provided")
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid role ID. Allowed roles: 3 (shopkeeper), 4 (auditor), 2 (admin)"})
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        log.Printf("HashPassword error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }
    input.Password = hashedPassword

    if err := db.GetDB().Create(&input).Error; err != nil {
        log.Printf("Create error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("User created successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "User created successfully"})
}

func EditUser(c echo.Context) error {
    id := c.Param("id")
    log.Printf("EditUser called with ID: %s", id)
    
    var user models.User
    if err := db.GetDB().First(&user, id).Error; err != nil {
        log.Printf("First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }

    if err := c.Bind(&user); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    if err := db.GetDB().Save(&user).Error; err != nil {
        log.Printf("Save error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("User updated successfully")
    return c.JSON(http.StatusOK, user)
}

func SoftDeleteUser(c echo.Context) error {
    id := c.Param("id")
    log.Printf("SoftDeleteUser called with ID: %s", id)
    
    var user models.User

    if user.RoleID == 1 { // Superadmin
            return c.JSON(http.StatusForbidden, "Only superadmins can delete superadmin users")
        }
    
    if err := db.GetDB().First(&user, id).Error; err != nil {
        log.Printf("First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }

    user.DeletedAt = gorm.DeletedAt{
        Time:  time.Now(),
        Valid: true,
    }

    if err := db.GetDB().Save(&user).Error; err != nil {
        log.Printf("Save error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("User soft deleted successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "User soft deleted successfully"})
}

func AdminLogout(c echo.Context) error {
    log.Println("Admin logged out successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func Login(c echo.Context) error {
    var loginData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&loginData); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        log.Printf("Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(loginData.Password, user.Password); err != nil {
        log.Printf("CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID,
        "roleID": user.RoleID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(utils.JwtSecret)
    if err != nil {
        log.Printf("SignedString error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("User logged in successfully")
    return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

func Logout(c echo.Context) error {
    log.Println("User logged out successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func AuditorLogin(c echo.Context) error {
    var loginData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&loginData); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        log.Printf("Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(loginData.Password, user.Password); err != nil {
        log.Printf("CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID,
        "roleID": user.RoleID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(utils.JwtSecret)
    if err != nil {
        log.Printf("SignedString error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("Auditor logged in successfully")
    return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

func AuditorLogout(c echo.Context) error {
    log.Println("Auditor logged out successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}
