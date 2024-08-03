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
    "stock-back/validators"
)

func SuperAdminSignup(c echo.Context) error {
    var input models.User
    if err := c.Bind(&input); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    loginInput := validators.LoginInput{
        Username: input.Username,
        Password: input.Password,
    }
    if err := validators.ValidateLoginInput(loginInput); err != nil {
        log.Printf("Validation error: %v", err)
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

    loginInput := validators.LoginInput{
        Username: newAdmin.Email,
        Password: newAdmin.Password,
    }
    if err := validators.ValidateLoginInput(loginInput); err != nil {
        log.Printf("Validation error: %v", err)
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

/*func SoftDeleteOrganization(c echo.Context) error {
    id := c.Param("id")
    log.Printf("SoftDeleteOrganization called with ID: %s", id)

    // Retrieve roleID from the context
    roleID := c.Get("roleID").(int)
    if roleID != 2 {
        return c.JSON(http.StatusForbidden, "Only admins can delete organizations")
    }

    var organization models.Organization

    if err := db.GetDB().First(&organization, id).Error; err != nil {
        log.Printf("First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "Organization not found"})
    }

    if organization.RoleID != 5 {
        return c.JSON(http.StatusForbidden, "Unauthorized: Only organizations can be deleted")
    }

    // Set DeletedAt field
    organization.DeletedAt = &gorm.DeletedAt{
        Time:  time.Now(),
        Valid: true,
    }

    if err := db.GetDB().Save(&organization).Error; err != nil {
        log.Printf("Save error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("Organization soft deleted successfully")
    return c.JSON(http.StatusOK, echo.Map{"message": "Organization soft deleted successfully"})
}
am getting this error on postman  "error": "Access forbidden" and this error on 2024/07/30 15:10:22 Token parsed successfully. UserID: 30, RoleID: 1 Failed to get JWT token from context. help me fix it
*/

func SuperAdminLogin(c echo.Context) error {
    var input struct {
        Username    string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&input); err != nil {
        log.Printf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    var user models.User
    if err := db.GetDB().Where("username = ?", input.Username).First(&user).Error; err != nil {
        log.Printf("Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        log.Printf("CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid username or password"})
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
    log.Println("AdminLogin - Entry")

    var input struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&input); err != nil {
        log.Printf("AdminLogin - Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    log.Printf("AdminLogin - Received input: %+v", input)

    var user models.User
    if err := db.GetDB().Where("email = ?", input.Email).First(&user).Error; err != nil {
        log.Printf("AdminLogin - Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(input.Password, user.Password); err != nil {
        log.Printf("AdminLogin - CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token, err := utils.GenerateJWT(user.ID, user.RoleID)
    if err != nil {
        log.Printf("AdminLogin - GenerateJWT error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("AdminLogin - Admin logged in successfully")
    log.Println("AdminLogin - Exit")
    return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func GetUserByID(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    log.Printf("GetUserByID - Entry with ID: %d", id)
    
    var user models.User
    if err := db.GetDB().First(&user, id).Error; err != nil {
        log.Printf("GetUserByID - First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }
    
    log.Printf("GetUserByID - User found: %+v", user)
    log.Println("GetUserByID - Exit")
    return c.JSON(http.StatusOK, user)
}

func AdminAddUser(c echo.Context) error {
    log.Println("AdminAddUser - Entry")

    // Retrieve userID and roleID from context set by middleware
    userID, ok := c.Get("userID").(int)
    if !ok {
        log.Println("AdminAddUser - Unauthorized: userID not found in context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    roleID, ok := c.Get("roleID").(int)
    if !ok {
        log.Println("AdminAddUser - Unauthorized: roleID not found in context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    log.Printf("AdminAddUser - Received RoleID: %d, UserID: %d", roleID, userID)

    // Check if the roleID is 2 (Admin)
    if roleID != 2 {
        log.Println("AdminAddUser - Permission denied: non-admin trying to add user")
        return c.JSON(http.StatusForbidden, echo.Map{"error": "Permission denied"})
    }

    var input models.User
    if err := c.Bind(&input); err != nil {
        log.Printf("AdminAddUser - Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    log.Printf("AdminAddUser - New user data: %+v", input)

    // Validate roleID for new user
    if input.RoleID != 3 && input.RoleID != 4 && input.RoleID != 2 {
        log.Println("AdminAddUser - Invalid role ID provided")
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid role ID. Allowed roles: 3 (shopkeeper), 4 (auditor), 2 (admin)"})
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        log.Printf("AdminAddUser - HashPassword error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
    }
    input.Password = hashedPassword

    if err := db.GetDB().Create(&input).Error; err != nil {
        log.Printf("AdminAddUser - Create error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("AdminAddUser - User created successfully")
    log.Println("AdminAddUser - Exit")
    return c.JSON(http.StatusOK, echo.Map{"message": "User created successfully"})
}

func EditUser(c echo.Context) error {
    id := c.Param("id")
    log.Printf("EditUser - Entry with ID: %s", id)
    
    var user models.User
    if err := db.GetDB().First(&user, id).Error; err != nil {
        log.Printf("EditUser - First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }

    log.Printf("EditUser - Current user details: %+v", user)

    if err := c.Bind(&user); err != nil {
        log.Printf("EditUser - Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    if err := db.GetDB().Save(&user).Error; err != nil {
        log.Printf("EditUser - Save error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("EditUser - User updated successfully")
    log.Println("EditUser - Exit")
    return c.JSON(http.StatusOK, user)
}

func AdminViewAllUsers(c echo.Context) error {
    log.Println("AdminViewAllUsers - Entry")

    // Retrieve roleID from context set by middleware
    roleID, ok := c.Get("roleID").(int)
    if !ok {
        log.Println("AdminViewAllUsers - Failed to get roleID from context")
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
    }

    log.Printf("AdminViewAllUsers - Received RoleID: %d", roleID)

    // Check if the roleID is 2 (Admin)
    if roleID != 2 {
        log.Println("AdminViewAllUsers - Permission denied: non-admin trying to view users")
        return c.JSON(http.StatusForbidden, echo.Map{"error": "Permission denied"})
    }

    var users []models.User
    if err := db.GetDB().Find(&users).Error; err != nil {
        log.Printf("AdminViewAllUsers - Find error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not retrieve users"})
    }

    log.Printf("AdminViewAllUsers - Retrieved users: %+v", users)
    log.Println("AdminViewAllUsers - Exit")
    return c.JSON(http.StatusOK, users)
}

func SoftDeleteUser(c echo.Context) error {
    id := c.Param("id")
    log.Printf("SoftDeleteUser - Entry with ID: %s", id)
    
    var user models.User

    if user.RoleID == 1 { // Superadmin
        log.Println("SoftDeleteUser - Only superadmins can delete superadmin users")
        return c.JSON(http.StatusForbidden, "Only superadmins can delete superadmin users")
    }
    
    if err := db.GetDB().First(&user, id).Error; err != nil {
        log.Printf("SoftDeleteUser - First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
    }

    user.DeletedAt = gorm.DeletedAt{
        Time:  time.Now(),
        Valid: true,
    }

    if err := db.GetDB().Save(&user).Error; err != nil {
        log.Printf("SoftDeleteUser - Save error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
    }

    log.Println("SoftDeleteUser - User soft deleted successfully")
    log.Println("SoftDeleteUser - Exit")
    return c.JSON(http.StatusOK, echo.Map{"message": "User soft deleted successfully"})
}

func ActivateUser(c echo.Context) error {
    userID := c.Param("id")
    var user models.User

    if err := db.GetDB().First(&user, userID).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
    }

    user.IsActive = true
    if err := db.GetDB().Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving user"})
    }

    return c.JSON(http.StatusOK, user)
}

func DeactivateUser(c echo.Context) error {
    userID := c.Param("id")
    var user models.User

    if err := db.GetDB().First(&user, userID).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
    }

    user.IsActive = false
    if err := db.GetDB().Save(&user).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving user"})
    }

    return c.JSON(http.StatusOK, user)
}

func GetOrganizationByID(c echo.Context) error {
    id, _ := strconv.Atoi(c.Param("id"))
    log.Printf("GetOrganizationByID - Entry with ID: %d", id)

    var org models.Organization
    if err := db.GetDB().First(&org, id).Error; err != nil {
        log.Printf("GetOrganizationByID - First error: %v", err)
        return c.JSON(http.StatusNotFound, echo.Map{"error": "Organization not found"})
    }

    log.Printf("GetOrganizationByID - Organization found: %+v", org)
    log.Println("GetOrganizationByID - Exit")
    return c.JSON(http.StatusOK, org)
}

func GetAllOrganizations(c echo.Context) error {
    log.Println("GetAllOrganizations - Entry")

    var orgs []models.Organization
    if err := db.GetDB().Find(&orgs).Error; err != nil {
        log.Printf("GetAllOrganizations - Find error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error retrieving organizations"})
    }

    log.Printf("GetAllOrganizations - Organizations found: %+v", orgs)
    log.Println("GetAllOrganizations - Exit")
    return c.JSON(http.StatusOK, orgs)
}

// ActivateOrganization activates an organization
func ActivateOrganization(c echo.Context) error {
    orgID := c.Param("id")
    var org models.Organization

    if err := db.GetDB().First(&org, orgID).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Organization not found"})
    }

    org.IsActive = true
    if err := db.GetDB().Save(&org).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving organization"})
    }

    return c.JSON(http.StatusOK, org)
}

// DeactivateOrganization deactivates an organization
func DeactivateOrganization(c echo.Context) error {
    orgID := c.Param("id")
    var org models.Organization

    if err := db.GetDB().First(&org, orgID).Error; err != nil {
        return c.JSON(http.StatusNotFound, map[string]string{"message": "Organization not found"})
    }

    org.IsActive = false
    if err := db.GetDB().Save(&org).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error saving organization"})
    }

    return c.JSON(http.StatusOK, org)
}

func GetActiveUsers(c echo.Context) error {
    var users []models.User
    if err := db.GetDB().Where("is_active = ?", true).Find(&users).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error retrieving users"})
    }
    return c.JSON(http.StatusOK, users)
}

func GetInactiveUsers(c echo.Context) error {
    var users []models.User
    if err := db.GetDB().Where("is_active = ?", false).Find(&users).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error retrieving users"})
    }
    return c.JSON(http.StatusOK, users)
}

func GetActiveOrganizations(c echo.Context) error {
    var orgs []models.Organization
    if err := db.GetDB().Where("is_active = ?", true).Find(&orgs).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error retrieving organizations"})
    }
    return c.JSON(http.StatusOK, orgs)
}

func GetInactiveOrganizations(c echo.Context) error {
    var orgs []models.Organization
    if err := db.GetDB().Where("is_active = ?", false).Find(&orgs).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error retrieving organizations"})
    }
    return c.JSON(http.StatusOK, orgs)
}


func AdminLogout(c echo.Context) error {
    log.Println("AdminLogout - Entry")
    log.Println("AdminLogout - Admin logged out successfully")
    log.Println("AdminLogout - Exit")
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func Login(c echo.Context) error {
    log.Println("Login - Entry")

    var loginData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&loginData); err != nil {
        log.Printf("Login - Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    log.Printf("Login - Received data: %+v", loginData)

    var user models.User
    if err := db.GetDB().Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        log.Printf("Login - Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(loginData.Password, user.Password); err != nil {
        log.Printf("Login - CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID,
        "roleID": user.RoleID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(utils.JwtSecret)
    if err != nil {
        log.Printf("Login - SignedString error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("Login - User logged in successfully")
    log.Println("Login - Exit")
    return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

func Logout(c echo.Context) error {
    log.Println("Logout - Entry")
    log.Println("Logout - User logged out successfully")
    log.Println("Logout - Exit")
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}

func AuditorLogin(c echo.Context) error {
    log.Println("AuditorLogin - Entry")

    var loginData struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.Bind(&loginData); err != nil {
        log.Printf("AuditorLogin - Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
    }

    log.Printf("AuditorLogin - Received data: %+v", loginData)

    var user models.User
    if err := db.GetDB().Where("email = ?", loginData.Email).First(&user).Error; err != nil {
        log.Printf("AuditorLogin - Where error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    if err := utils.CheckPasswordHash(loginData.Password, user.Password); err != nil {
        log.Printf("AuditorLogin - CheckPasswordHash error: %v", err)
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": user.ID,
        "roleID": user.RoleID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(utils.JwtSecret)
    if err != nil {
        log.Printf("AuditorLogin - SignedString error: %v", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not generate token"})
    }

    log.Println("AuditorLogin - Auditor logged in successfully")
    log.Println("AuditorLogin - Exit")
    return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}

func AuditorLogout(c echo.Context) error {
    log.Println("AuditorLogout - Entry")
    log.Println("AuditorLogout - Auditor logged out successfully")
    log.Println("AuditorLogout - Exit")
    return c.JSON(http.StatusOK, echo.Map{"message": "Successfully logged out"})
}


