package middlewares

import (
    "net/http"
    "strings"
    "stock-back/models"
    "stock-back/utils"
    "github.com/labstack/echo/v4"
    "log"
)

// AuthMiddleware validates the JWT token and checks if the user's role is allowed.
func AuthMiddleware(allowedRoles ...int) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // Skip authentication for login and logout routes
            if c.Path() == "/login" || c.Path() == "/logout" {
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
            token, err := utils.ParseToken(tokenString)
            if err != nil {
                log.Printf("Token parsing error: %v", err)
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
            }

            claims, ok := token.Claims.(*utils.Claims)
            if !ok || !token.Valid {
                log.Printf("Invalid token claims: %v", token)
                return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"})
            }

            userID := int(claims.UserID)
            roleID := int(claims.RoleID)
            log.Printf("Token parsed successfully. UserID: %d, RoleID: %d", userID, roleID)

            if roleID == 0 {
                log.Printf("RoleID is 0 for token: %s", tokenString)
                return c.JSON(http.StatusForbidden, echo.Map{"error": "Access forbidden"})
            }

            // Set context values
            c.Set("userID", userID)
            c.Set("roleID", roleID)

            if !contains(allowedRoles, roleID) {
                return c.JSON(http.StatusForbidden, echo.Map{"error": "Access forbidden"})
            }

            return next(c)
        }
    }
}

// Helper function to check if a slice contains a value.
func contains(slice []int, value int) bool {
    for _, v := range slice {
        if v == value {
            return true
        }
    }
    return false
}

// Role-based middlewares

func SuperAdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.SuperAdminRoleID {
            log.Printf("Access denied for RoleID %d", roleID)
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.AdminRoleID {
            log.Printf("Access denied for RoleID %d", roleID)
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

func ShopAttendantOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.ShopAttendantRoleID {
            log.Printf("Access denied for RoleID %d", roleID)
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}

func AuditorOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        roleID, ok := c.Get("roleID").(int)
        if !ok || roleID != models.AuditorRoleID {
            log.Printf("Access denied for RoleID %d", roleID)
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }
        return next(c)
    }
}
