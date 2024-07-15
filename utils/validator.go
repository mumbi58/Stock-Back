package utils

import (
    "github.com/go-playground/validator/v10"
    "github.com/labstack/echo/v4"
)

// CustomValidator is a custom validator for Echo
type CustomValidator struct {
    validator *validator.Validate
}

// Validate performs validation on the provided interface
func (cv *CustomValidator) Validate(i interface{}) error {
    return cv.validator.Struct(i)
}

// NewValidator creates a new instance of CustomValidator
func NewValidator() *CustomValidator {
    return &CustomValidator{validator: validator.New()}
}

// SetupValidator sets up the custom validator for Echo
func SetupValidator(e *echo.Echo) {
    e.Validator = NewValidator()
}
