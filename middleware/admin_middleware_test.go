package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdminMiddleware(t *testing.T) {
	e := echo.New()

	// Define a handler to test with
	testHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Success"})
	}

	// Apply middleware to the test handler
	middleware := AdminMiddleware(testHandler)

	// Create a test request with the role header set to '2' (admin)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Role", "2")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the middleware
	err := middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"message": "Success"}`, rec.Body.String())

	// Create a test request with a non-admin role header
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Role", "1") // Non-admin role
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Call the middleware
	err = middleware(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.JSONEq(t, `{"message": "Access denied"}`, rec.Body.String())
}
