package routes

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "stock-back/controllers"
    "stock-back/middlewares"
    "gorm.io/gorm" // Import gorm package here
)

// SetupRouter initializes all routes for the application
func SetupRouter(e *echo.Echo, db *gorm.DB) {
    // Middleware to inject db instance into context
    e.Use(middlewares.GetDBMiddleware(db))

    // Public routes for login and logout
    e.POST("/login", func(c echo.Context) error { return controllers.Login(c, db) })
    e.POST("/logout", func(c echo.Context) error { return controllers.Logout(c, db) })
    e.POST("/auditor/login", func(c echo.Context) error { return controllers.AuditorLogin(c, db) })
    e.POST("/auditor/logout", func(c echo.Context) error { return controllers.AuditorLogout(c, db) })
    e.POST("/admin/login", func(c echo.Context) error { return controllers.AdminLogin(c) })
    e.POST("/admin/logout", func(c echo.Context) error { return controllers.AdminLogout(c) })

    // Routes specific to super admin
    superadmin := e.Group("/superadmin")
    {
        // Public routes for superadmin
        superadmin.POST("/signup", controllers.SuperAdminSignup)
        superadmin.POST("/login", controllers.SuperAdminLogin)

        // Protected routes for superadmin
        superadmin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
            SigningKey: []byte("secret"),
        }))
        superadmin.POST("/addadmin", controllers.AddAdmin)
        superadmin.POST("/addorganization", controllers.SuperAdminAddOrganization)
    }

    // Routes specific to admin
    admin := e.Group("/admin")
    {
        // Middleware to verify JWT token for admin routes and restrict to admins only
        admin.Use(
            middleware.JWTWithConfig(middleware.JWTConfig{
                SigningKey: []byte("secret"),
            }),
            middlewares.AdminOnly,
        )

        admin.POST("/adduser", controllers.AdminAddUser)
        admin.PUT("/users/:id", controllers.EditUser)
        admin.DELETE("/users/:id", controllers.SoftDeleteUser)
        admin.GET("/users/:id", controllers.AdminGetUserByID)
        admin.GET("/users", controllers.AdminListUsers)
    }
}
