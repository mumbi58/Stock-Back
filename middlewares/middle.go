package middlewares

import (
    "net/http"
    "strings"
    "stock-back/models"
    "stock-back/utils"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
)

// GetDBMiddleware injects the DB instance into the context
func GetDBMiddleware(db *gorm.DB) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            c.Set("db", db)
            return next(c)
        }
    }
}

// ExtractDB extracts the DB instance from the context
func ExtractDB(c echo.Context) *gorm.DB {
    db, ok := c.Get("db").(*gorm.DB)
    if !ok {
        return nil
    }
    return db
}

// JWTMiddleware is the middleware for JWT authentication
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Allow access to the sign-up endpoint without authorization
        if c.Request().Method == http.MethodPost && c.Request().URL.Path == "/superadmin/signup" {
            return next(c)
        }

        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authorization header is required"})
        }

        token = strings.Replace(token, "Bearer ", "", 1)
        claims, err := utils.VerifyJWT(token)
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
        }

        c.Set("userID", claims.UserID)
        c.Set("roleID", claims.RoleID)
        return next(c)
    }
}

// AuthMiddleware checks if user has one of the allowed roles
func AuthMiddleware(allowedRoles ...int) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            roleID, ok := c.Get("roleID").(int)  // Use int here to match expected roleID type
            if !ok {
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Role ID not found in context"})
            }

            if !contains(allowedRoles, roleID) {
                return c.JSON(http.StatusForbidden, echo.Map{"error": "Access forbidden"})
            }

            return next(c)
        }
    }
}

// Utility function to check if a slice contains a value
func contains(slice []int, value int) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
    }
    return false
}

// SuperAdminOnly restricts access to super admins
func SuperAdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)  // Use int here to match expected roleID type
        if !ok || roleID != models.SuperAdminRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

// AdminOnly restricts access to admins
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)  // Use int here to match expected roleID type
        if !ok || roleID != models.AdminRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

// ShopAttendantOnly restricts access to shop attendants
func ShopAttendantOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)  // Use int here to match expected roleID type
        if !ok || roleID != models.ShopAttendantRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

// AuditorOnly restricts access to auditors
func AuditorOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)  // Use int here to match expected roleID type
        if !ok || roleID != models.AuditorRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}
