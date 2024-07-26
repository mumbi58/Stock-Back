package middlewares

import (
    "net/http"
    "strings"
    "stock-back/models"
    "stock-back/utils"
    "github.com/labstack/echo/v4"
    "github.com/dgrijalva/jwt-go"
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

// AuthMiddleware checks if user has one of the allowed roles
func AuthMiddleware(allowedRoles ...int) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Allow access to the sign-up endpoint without authorization
            if c.Request().Method == http.MethodPost && c.Request().URL.Path == "/superadmin/signup" {
                return next(c)
            }

            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header is required"})
            }

            tokenParts := strings.Split(authHeader, " ")
            if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token format"})
            }

            tokenString := tokenParts[1]
            token, err := utils.ParseToken(tokenString) // Adjust this function as needed
            if err != nil {
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
            }

            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok || !token.Valid {
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"})
            }

            roleID, ok := claims["roleID"].(float64)
            if !ok {
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"})
            }

            // Set role ID in context
            c.Set("roleID", int(roleID))

            // Check if role ID is allowed
            if !contains(allowedRoles, int(roleID)) {
                return c.JSON(http.StatusForbidden, echo.Map{"error": "Access forbidden"})
            }

            return next(c) // Continue to the next middleware or handler
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
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.SuperAdminRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

// AdminOnly restricts access to admins
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.AdminRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

// ShopAttendantOnly restricts access to shop attendants
func ShopAttendantOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.ShopAttendantRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

// AuditorOnly restricts access to auditors
func AuditorOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.AuditorRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}
