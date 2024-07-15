// middlewares.go

package middlewares

import (
    "net/http"
    "strings"
    "stock-back/models"
    "stock-back/utils"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "gorm.io/gorm"
)

// AuthMiddleware ensures that the request is authenticated
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
        }

        token = strings.TrimPrefix(token, "Bearer ")
        _, err := utils.VerifyJWT(token)
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
        }

        return next(c)
    }
}

// AdminOnly ensures that only admins can access the route
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
        }

        token = strings.TrimPrefix(token, "Bearer ")
        claims, err := utils.VerifyJWT(token)
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
        }

        if claims.RoleID != models.AdminRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }

        return next(c)
    }
}

// SuperAdminOnly ensures that only super admins can access the route
func SuperAdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
        }

        token = strings.TrimPrefix(token, "Bearer ")
        claims, err := utils.VerifyJWT(token)
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
        }

        if claims.RoleID != models.SuperAdminRoleID {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
        }

        return next(c)
    }
}

// GetDBMiddleware attaches the database instance to the context
func GetDBMiddleware(db *gorm.DB) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            c.Set("db", db)
            return next(c)
        }
    }
}

// ExtractDB is a helper function to extract the database from the context
func ExtractDB(c echo.Context) *gorm.DB {
    db, ok := c.Get("db").(*gorm.DB)
    if !ok {
        panic("database connection not found in context")
    }
    return db
}

// JWTMiddleware is a middleware to check for JWT token in the header
func JWTMiddleware() echo.MiddlewareFunc {
    return middleware.JWTWithConfig(middleware.JWTConfig{
        SigningKey: []byte("secret"),
    })
}
