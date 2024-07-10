package routes

import (
    "github.com/labstack/echo/v4"
    "stock-back/controllers"
    "stock-back/middlewares"

)

func Setup() *echo.Echo {
    e := echo.New()


    // Public routes
    e.POST("/login", controllers.Login)
    e.POST("/logout", controllers.Logout)
    e.POST("/auditor/login", controllers.AuditorLogin)
    e.POST("/auditor/logout", controllers.AuditorLogout)
    e.POST("/admin/login", controllers.AdminLogin)
    e.POST("/admin/logout", controllers.AdminLogout)
    e.POST("/superadmin/login", controllers.SuperAdminLogin)

    // Routes specific to super admin
    superadmin := e.Group("/superadmin")
    superadmin.POST("/signup", controllers.SuperAdminSignup)
    superadmin.Use(middlewares.AuthMiddleware)
    superadmin.POST("/addadmin", controllers.AddAdmin)

    // Routes specific to admin
    admin := e.Group("/admin")
    admin.Use(middlewares.AuthMiddleware, middlewares.AdminOnly)
    admin.POST("/adduser", controllers.AddUser)
    admin.PUT("/users/:id", controllers.EditUser)
    admin.DELETE("/users/:id", controllers.SoftDeleteUser)
    admin.DELETE("/users/:id", controllers.GetAllUsers)

    return e
}
