package middlewares

import (
    "net/http"
    "strings"

    "github.com/labstack/echo/v4"
    "github.com/dgrijalva/jwt-go"
    "stock-back/models"
    "stock-back/utils"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        if c.Request().Method == http.MethodPost && c.Request().URL.Path == "/superadmin/signup" {
            return next(c)
        }

        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Authorization header is required"})
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token format"})
        }

        tokenString := tokenParts[1]
        token, err := utils.ParseToken(tokenString)
        if err != nil {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token"})
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token claims"})
        }

        userID, ok := claims["userID"].(float64)
        if !ok {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token claims"})
        }

        roleID, ok := claims["roleID"].(float64)
        if !ok {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token claims"})
        }

        // Set user info in context
        c.Set("user", models.User{
            ID:     uint(userID),
            RoleID: uint(roleID),
        })

        return next(c)
    }
}

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := c.Get("user").(models.User)
        if user.RoleID != models.AdminRoleID {
            return c.JSON(http.StatusForbidden, map[string]interface{}{"error": "Forbidden"})
        }
        return next(c)
    }
}
