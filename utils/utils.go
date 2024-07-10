package utils

import (
    "errors"
    "time"

    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

var DB *gorm.DB
var SecretKey = []byte("your-secret-key")

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateJWT(userID uint, roleID uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": userID,
        "roleID": roleID,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(SecretKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.StandardClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return SecretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

func ParseToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return SecretKey, nil
    })

    if err != nil {
        return nil, err
    }

    return token, nil
}
