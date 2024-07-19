package routes

import (
    "stock-back/controllers"
    "stock-back/middlewares"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
    "stock-back/models"
)

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
    // Middleware
    e.Use(middlewares.GetDBMiddleware(db))
    e.Use(middlewares.JWTMiddleware)

    // Super Admin routes
    e.POST("/superadmin/signup", controllers.SuperAdminSignup) // No middleware needed for signup

    e.POST("/superadmin/addadmin", middlewares.SuperAdminOnly(controllers.AddAdmin))
    e.POST("/superadmin/login", controllers.SuperAdminLogin)
    e.POST("/superadmin/logout", controllers.SuperAdminLogout)

    // Admin routes
    adminGroup := e.Group("/admin")
    adminGroup.Use(middlewares.AuthMiddleware(models.AdminRoleID))
    adminGroup.POST("/login", controllers.AdminLogin)
    adminGroup.POST("/adduser", controllers.AdminAddUser)
    adminGroup.GET("/user/:id", controllers.GetUserByID)
    adminGroup.PUT("/user/:id", controllers.EditUser)
    adminGroup.DELETE("/user/:id", controllers.SoftDeleteUser)
    adminGroup.POST("/logout", controllers.AdminLogout)

    // Shop Attendant routes
    shopAttendantGroup := e.Group("/shopattendant")
    shopAttendantGroup.POST("/login", controllers.Login)
    shopAttendantGroup.POST("/logout", controllers.Logout)

    // Auditor routes
    auditorGroup := e.Group("/auditor")
    auditorGroup.POST("/login", controllers.AuditorLogin)
    auditorGroup.POST("/logout", controllers.AuditorLogout)
}
