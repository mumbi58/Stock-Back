package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"stock-back/models"
	"stock-back/utils"
)


func Login(c echo.Context) error {
	var loginData struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&loginData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(&loginData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var user models.User
	if err := utils.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"roleID": user.RoleID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(utils.SecretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}

func Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
