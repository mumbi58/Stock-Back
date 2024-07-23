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

    //public routes
    e.POST("/superadmin/login", controllers.SuperAdminLogin)
    e.POST("/superadmin/logout", controllers.SuperAdminLogout)
    e.POST("/login", controllers.AdminLogin)
    e.POST("/logout", controllers.AdminLogout)
    e.POST("/login", controllers.Login)
    e.POST("/logout", controllers.Logout)
    e.POST("/login", controllers.AuditorLogin)
    e.POST("/logout", controllers.AuditorLogout)



    // Super Admin routes
    e.POST("/superadmin/signup", controllers.SuperAdminSignup) // No middleware needed for signup

    e.POST("/superadmin/addadmin", middlewares.SuperAdminOnly(controllers.AddAdmin))
    
    // Admin routes
    adminGroup := e.Group("/admin")
    adminGroup.Use(middlewares.AuthMiddleware(models.AdminRoleID))
    adminGroup.POST("/adduser", controllers.AdminAddUser)
    adminGroup.GET("/user/:id", controllers.GetUserByID)
    adminGroup.PUT("/user/:id", controllers.EditUser)
    adminGroup.DELETE("/user/:id", controllers.SoftDeleteUser)

    
}
