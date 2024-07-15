package utils

import (
    "errors"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/go-playground/validator/v10"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

var DB *gorm.DB
var SecretKey = []byte("your-secret-key")

// Validator holds the validator instance
var Validator = validator.New()

// CustomClaims struct to hold custom JWT claims
type CustomClaims struct {
    UserID uint `json:"userID"`
    RoleID uint `json:"roleID"`
    jwt.StandardClaims
}

// HashPassword hashes the given password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}

// CheckPasswordHash compares a hashed password with its possible plaintext equivalent
func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GenerateJWT generates a new JWT token with custom claims (userID and roleID)
func GenerateJWT(userID uint, roleID uint) (string, error) {
    // Set token claims
    claims := CustomClaims{
        UserID: userID,
        RoleID: roleID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
            IssuedAt:  time.Now().Unix(),
        },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Generate encoded token and return it
    tokenString, err := token.SignedString(SecretKey)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}

// VerifyJWT parses and verifies the JWT token string
func VerifyJWT(tokenString string) (*CustomClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return SecretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// ParseToken parses the JWT token string and returns the token object
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
