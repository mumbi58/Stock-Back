package routes

import (
    "stock-back/controllers"
    "stock-back/middlewares"
    "github.com/labstack/echo/v4"
    "stock-back/models"
)

func SetupRoutes(e *echo.Echo) {
    // Public routes
    e.POST("/superadmin/login", controllers.SuperAdminLogin)
    e.POST("/superadmin/logout", controllers.SuperAdminLogout)
    e.POST("admin/login", controllers.AdminLogin)
    e.POST("admin/logout", controllers.AdminLogout)
    e.POST("/login", controllers.Login)
    e.POST("/logout", controllers.Logout)
    e.POST("auditor/login", controllers.AuditorLogin)
    e.POST("auditor/logout", controllers.AuditorLogout)

    // Super Admin routes
    superadmin := e.Group("/superadmin")
    superadmin.POST("/signup", controllers.SuperAdminSignup)
    superadmin.Use(middlewares.AuthMiddleware(models.SuperAdminRoleID)) // Ensure SuperAdmin is authorized
    superadmin.POST("/addadmin", controllers.AddAdmin)
    superadmin.POST("/addorganization", controllers.SuperAdminAddOrganization)

    // Admin routes
    adminGroup := e.Group("/admin")
    adminGroup.Use(middlewares.AuthMiddleware(models.AdminRoleID)) // Ensure Admin is authorized
    adminGroup.POST("/adduser", controllers.AdminAddUser)
    adminGroup.GET("/user/:id", controllers.GetUserByID)
    adminGroup.PUT("/user/:id", controllers.EditUser)
    adminGroup.DELETE("/user/:id", controllers.SoftDeleteUser)
    //adminGroup.DELETE("/organization:id", controllers.AdminDeleteOrganization)
}
