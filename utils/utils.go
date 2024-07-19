package utils

import (
    "errors"
    "time"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

var (
    JwtSecret = []byte("your-secret-key") // Replace with your actual secret key
    DB        *gorm.DB
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
    UserID uint `json:"user_id"`
    RoleID uint `json:"role_id"`
    jwt.StandardClaims
}

// ParseToken parses and validates the JWT token
func ParseToken(tokenString string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Make sure the token method conforms to your signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
        return JwtSecret, nil
    })
    return token, err
}

// GenerateJWT generates a JWT token
func GenerateJWT(userID, roleID uint) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID: userID,
        RoleID: roleID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(JwtSecret)
}

// VerifyJWT verifies the JWT token
func VerifyJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return JwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, ErrInvalidToken
    }

    return claims, nil
}

// HashPassword hashes a password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// CheckPasswordHash compares a password hash
func CheckPasswordHash(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
